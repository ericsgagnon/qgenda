package qgenda

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var PGTag = "pg"

func NewPGClientConfig(connString string) *DBClientConfig {
	return &DBClientConfig{
		Name:             "postgres",
		Type:             "postgres",
		Driver:           "postgres",
		ConnectionString: connString,
	}

}

// PGToGoTypeMap represents the default type mapping
// we expect to use when retrieving data from postgres
var PGToGoTypeMap = map[string]string{
	"text": "string",
}

// GoToPGTypeMap represents the default type mapping
// we expect to use when sending data to postgres
var GoToPGTypeMap = map[string]string{
	"default":          "bytea[]",
	"string":           "text",
	"Date":             "date",
	"qgenda.Date":      "date",
	"Time":             "timestamp with time zone",
	"qgenda.Time":      "timestamp with time zone",
	"TimeOfDay":        "time without time zone",
	"qgenda.TimeOfDay": "time without time zone",
	"time.Time":        "timestamp with time zone",
	"bool":             "boolean",

	"int":     "bigint",
	"int8":    "smallint",
	"int16":   "smallint",
	"int32":   "bigint",
	"int64":   "bigint",
	"float32": "double precision",
	"float64": "double precision",
}

func GoToPGType(gotype string) string {
	pgtype, ok := GoToPGTypeMap[gotype]
	if !ok {
		pgtype = "text"
	}
	return pgtype
}

type Field struct {
	Name        string
	Kind        string
	Type        string
	Pointer     bool
	PrimaryKey  bool
	Unique      bool
	Nullable    bool // nullable follows the sql standard of defaulting to true
	Constraints []string
	Tags        map[string][]string
	StructField reflect.StructField
}

func StructToFields[T any](value T) []Field {

	v := reflect.ValueOf(*new(T))
	iv := reflect.Indirect(v)
	// fmt.Printf("%+v\n", iv)
	// handle zero pointers - look up code above
	structfields := StructFields(iv)
	fields := []Field{}
	for i := 0; i < iv.NumField(); i++ {
		sf := structfields[i]
		ivField := iv.Field(i)

		fieldType := reflect.TypeOf(ivField.Interface())
		if ivField.Kind() == reflect.Pointer {
			fieldType = reflect.TypeOf(ivField.Interface()).Elem()
		}
		fieldKind := fieldType.Kind()
		pointer := ivField.Kind() == reflect.Pointer
		tags := TagKeyValues(fmt.Sprint(sf.Tag))
		val, ok := tags["primarykey"]
		primarykey := ok && strings.ToLower(val[0]) != "false"
		val, ok = tags["unique"]
		unique := ok && strings.ToLower(val[0]) != "false"
		val, ok = tags["nullable"]
		nullable := !ok || strings.ToLower(val[0]) != "false"

		field := Field{
			Name:        sf.Name,
			Kind:        fieldKind.String(),
			Type:        fieldType.String(),
			Pointer:     pointer,
			PrimaryKey:  primarykey,
			Unique:      unique,
			Nullable:    nullable,
			Constraints: tags["constraints"],
			Tags:        tags,
			StructField: sf,
		}

		field.StructField = sf
		fields = append(fields, field)
	}

	return fields
}

func TagKeyValues(s string) map[string][]string {

	pattern := regexp.MustCompile(`(?m)(?P<key>\w+):\"(?P<value>[^"]+)\"`)
	matches := pattern.FindAllStringSubmatch(s, -1)
	var out = map[string][]string{}
	for _, match := range matches {
		out[match[1]] = strings.Split(match[2], ",")
	}
	return out
}

func PrimaryKey(fields []Field) []string {
	pk := []string{}
	for _, field := range fields {
		if field.PrimaryKey {
			pk = append(pk, PGName(field))
		}
	}
	return pk
}

func QueryFieldName(field Field) string {
	if nametags, ok := field.Tags["qf"]; ok {
		return nametags[0]
	}
	return ""
}

// PGName attempts to return the field name to be used for postgres endpoints
// it does not validate anything. Rather, it check the pg struct tag, if that
// is empty it checks the db struct tag. If that is empty it takes a lowercase
// of the name.
func PGName(field Field) string {
	var name string
	nametags, ok := field.Tags[PGTag]
	if !ok {
		nametags, ok = field.Tags["db"]
	}
	if ok {
		name = nametags[0]
	} else {
		name = strings.ToLower(field.Name)
	}
	return name
}

func PGOmit(field Field) bool {
	if !field.StructField.IsExported() {
		return true
	}
	// slices, maps, arrays, and channel's must be handled separately for now
	switch field.StructField.Type.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		return true
	}
	return PGName(field) == "-"
}

func PGStatement[T any](value T, schema, table, tpl string) string {
	allfields := StructToFields(*new(T))

	var fields []Field
	for _, field := range allfields {
		if PGOmit(field) {
			continue
		}
		fields = append(fields, field)
	}

	tplValues := struct {
		Schema     string
		Table      string
		Fields     []Field
		PrimaryKey []string
	}{
		Schema:     schema,
		Table:      table,
		Fields:     fields,
		PrimaryKey: PrimaryKey(fields),
	}

	var buf bytes.Buffer

	if err := template.Must(template.
		New("").
		Funcs(template.FuncMap{
			"join":          strings.Join,
			"pgtype":        GoToPGType,
			"pgname":        PGName,
			"pgqueryfields": PGQueryConditionFields,
			"qfname":        QueryFieldName,
		}).
		Parse(tpl)).
		Execute(&buf, tplValues); err != nil {
		log.Println(err)
		panic(err)
	}
	return buf.String()
}

func PGQueryConditionFields(fields []Field) []Field {
	qfields := []Field{}
	for _, field := range fields {
		if _, ok := field.Tags["querycondition"]; ok {
			qfields = append(qfields, field)
		}
	}
	return qfields
}

var pgCreateNewTableTpl = `
CREATE TABLE IF NOT EXISTS {{ .Schema -}}{{- .Table }} (
{{- range  $index, $field := .Fields -}}
{{- if ne $index 0 -}},{{- end }}
	{{ pgname $field }} {{ pgtype $field.Type }} {{ if $field.Unique }} unique {{ end -}} {{- if not $field.Nullable -}} not null {{- end }}
{{- end -}}
{{- $primarykey := join .PrimaryKey  ", " -}}
{{ if ne $primarykey "" }}, 
PRIMARY KEY ( {{ $primarykey }} ) 
{{ else }},
CONSTRAINT {{ .Table -}}_all_columns_unique UNIQUE (
{{- range  $index, $field := .Fields -}}
{{- if ne $index 0 -}},{{- end }}
	{{ pgname $field }} 
{{- end -}} )
{{ end }} 
)
`

//create unique index if not exists schedulestafftag_all_columns_unique on schedulestafftag (schedulekey, lastmodifieddateutc, categorykey, categoryname, tagkey, tagname)

var pgInsertTpl = `
INSERT INTO {{ .Schema -}}{{- .Table }} (
{{- range  $index, $field := .Fields -}}
{{- if ne $index 0 -}},{{- end }}
	{{ pgname $field }}
{{- end }} 
) VALUES (
{{- range  $index, $field := .Fields -}}
{{- if ne $index 0 -}},{{- end }}
	:{{ pgname $field }}
{{- end }}	
) ON CONFLICT (
{{- $primarykey := join .PrimaryKey  ", " -}}
{{ if ne $primarykey "" }}
{{ $primarykey }}
{{ else }}
{{- range  $index, $field := .Fields -}}
{{- if ne $index 0 -}},{{- end }}
	{{ pgname $field }}
{{- end -}}	
{{- end }}
) DO NOTHING
`

var pgSelectMaxConstraintsTpl = `
SELECT
{{- $fields := pgqueryfields .Fields -}}
{{- range  $index, $field := $fields -}}
{{- if ne $index 0 -}},{{- end }}
	MAX( {{ pgname $field }} ) AS {{ qfname $field }}
{{- end }}
FROM 
	{{ .Schema -}}{{- .Table }}
`

func PGCreateTableStatement[T any](value T, schema, table string) string {
	return PGStatement(value, schema, table, pgCreateNewTableTpl)
}

func PGInsertStatement[T any](value T, schema, table string) string {
	return PGStatement(value, schema, table, pgInsertTpl)
}

func PGQueryConstraintsStatement[T any](value T, schema, table string) string {
	return PGStatement(value, schema, table, pgSelectMaxConstraintsTpl)
}

func PGCreateTable[T any](ctx context.Context, db *sqlx.DB, value []T, schema, table string) (sql.Result, error) {
	// fmt.Println(PGCreateTableStatement(value[0], schema, table))
	return db.NamedExecContext(
		ctx,
		PGCreateTableStatement(value[0], schema, table),
		value,
	)
}

func PGInsertRows[T any](ctx context.Context, db *sqlx.DB, value []T, schema, table string) (sql.Result, error) {
	// fmt.Println(PGInsertStatement(value[0], schema, table))
	return db.NamedExecContext(
		ctx,
		PGInsertStatement(value[0], schema, table),
		value,
	)
}

func PGQueryConstraint[T any](ctx context.Context, db *sqlx.DB, value []T, schema, table string) (sql.Result, error) {
	// fmt.Println(PGQueryConstraintsStatement(value[0], schema, table))
	result, err := db.NamedExecContext(
		ctx,
		PGQueryConstraintsStatement(value[0], schema, table),
		value,
	)
	return result, err
}

// func CreateTable[T any](a T) (bool, error) {
// 	if IsStruct(a) {
// 		// sf := StructFields(a)
// 		// sf[0].Anonymous
// 	}
// 	return false, nil
// }

// func CreateTableSQL[T any](schema, table string) string {

// 	v := reflect.ValueOf(*new(T))
// 	iv := reflect.Indirect(v)
// 	fields := StructFields(iv)
// 	sqlFieldDefs := []string{}
// 	primaryKey := []string{}
// 	for i := 0; i < iv.NumField(); i++ {
// 		sf := fields[i]
// 		tags := strings.Split(sf.Tag.Get("sql"), ",")
// 		tagMap := map[string]bool{}
// 		for i, v := range tags {
// 			tag := strings.ToLower(v)
// 			tags[i] = tag
// 			if i != 0 {
// 				// skip tag[0] - it's the sql name
// 				tagMap[tag] = true
// 			}
// 		}
// 		if !sf.IsExported() || tags[0] == "-" {
// 			continue
// 		}
// 		field := iv.Field(i)
// 		fieldType := reflect.TypeOf(field.Interface())
// 		if field.Kind() == reflect.Pointer {
// 			fieldType = reflect.TypeOf(field.Interface()).Elem()
// 		}
// 		switch fieldType.Kind() {
// 		case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
// 			continue
// 		}

// 		sqlName := strings.ToLower(sf.Name)
// 		if tags[0] != "" {
// 			sqlName = tags[0]
// 		}

// 		sqlType := GoToPGTypeMap[fieldType.Name()]
// 		if tagMap["primarykey"] || tagMap["primary key"] {
// 			primaryKey = append(primaryKey, sqlName)
// 			delete(tagMap, "primarykey")
// 			delete(tagMap, "primary key")
// 		}
// 		constraints := ""
// 		for _, v := range tags {
// 			if tagMap[v] {
// 				constraints = fmt.Sprintf("%s %s", constraints, v)
// 			}
// 		}
// 		sqlFieldDefs = append(sqlFieldDefs, fmt.Sprintf("%s %s %s", sqlName, sqlType, constraints))

// 	}
// 	// fmt.Println(strings.Join(sqlFieldDefs, ",\n"))

// 	sqlTpl := `
// CREATE TABLE IF NOT EXISTS {{ .Schema -}}{{- .Table }} (
// {{- range  $index, $field := .Fields -}}
// {{- if ne $index 0 -}},{{- end }}
// 	{{ $field -}}
// {{- end -}}
// {{- $primarykey := join .PrimaryKey  ", " -}}
// {{ if ne $primarykey "" }},
// 	PRIMARY KEY ( {{ $primarykey }} ) {{ end }}
// )
// `
// 	tplValues := struct {
// 		Schema     string
// 		Table      string
// 		Fields     []string
// 		PrimaryKey []string
// 	}{
// 		Schema:     schema,
// 		Table:      table,
// 		Fields:     sqlFieldDefs,
// 		PrimaryKey: primaryKey,
// 	}
// 	var buf bytes.Buffer

// 	if err := template.Must(template.
// 		New("createScheduleTable").
// 		Funcs(template.FuncMap{"join": strings.Join}).
// 		Parse(sqlTpl)).
// 		Execute(&buf, tplValues); err != nil {
// 		log.Println(err)
// 		panic(err)
// 	}
// 	return buf.String()
// 	// return ""
// }

// func createSQLStatement[T any](schema, table, tpl string) string {

// 	v := reflect.ValueOf(*new(T))
// 	iv := reflect.Indirect(v)
// 	fields := StructFields(iv)
// 	sqlFieldDefs := []string{}
// 	primaryKey := []string{}
// 	for i := 0; i < iv.NumField(); i++ {
// 		sf := fields[i]
// 		tags := strings.Split(sf.Tag.Get("sql"), ",")
// 		tagMap := map[string]bool{}
// 		for i, v := range tags {
// 			tag := strings.ToLower(v)
// 			tags[i] = tag
// 			if i != 0 {
// 				// skip tag[0] - it's the sql name
// 				tagMap[tag] = true
// 			}
// 		}
// 		if !sf.IsExported() || tags[0] == "-" {
// 			continue
// 		}
// 		field := iv.Field(i)
// 		fieldType := reflect.TypeOf(field.Interface())
// 		if field.Kind() == reflect.Pointer {
// 			fieldType = reflect.TypeOf(field.Interface()).Elem()
// 		}
// 		switch fieldType.Kind() {
// 		case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
// 			continue
// 		}

// 		sqlName := strings.ToLower(sf.Name)
// 		if tags[0] != "" {
// 			sqlName = tags[0]
// 		}

// 		sqlType := GoToPGTypeMap[fieldType.Name()]
// 		if tagMap["primarykey"] || tagMap["primary key"] {
// 			primaryKey = append(primaryKey, sqlName)
// 			delete(tagMap, "primarykey")
// 			delete(tagMap, "primary key")
// 		}
// 		constraints := ""
// 		for _, v := range tags {
// 			if tagMap[v] {
// 				constraints = fmt.Sprintf("%s %s", constraints, v)
// 			}
// 		}
// 		sqlFieldDefs = append(sqlFieldDefs, fmt.Sprintf("%s %s %s", sqlName, sqlType, constraints))

// 	}
// 	fmt.Println(strings.Join(sqlFieldDefs, ",\n"))

// 	tplValues := struct {
// 		Schema     string
// 		Table      string
// 		Fields     []string
// 		PrimaryKey []string
// 	}{
// 		Schema:     schema,
// 		Table:      table,
// 		Fields:     sqlFieldDefs,
// 		PrimaryKey: primaryKey,
// 	}
// 	var buf bytes.Buffer

// 	if err := template.Must(template.
// 		New("").
// 		Funcs(template.FuncMap{"join": strings.Join}).
// 		Parse(tpl)).
// 		Execute(&buf, tplValues); err != nil {
// 		log.Println(err)
// 		panic(err)
// 	}
// 	return buf.String()
// }
// var pgCreateTableTpl = `
// CREATE TABLE IF NOT EXISTS {{ .Schema -}}{{- .Table }} (
// {{- range  $index, $field := .Fields -}}
// {{- if ne $index 0 -}},{{- end }}
// 	{{ $field -}}
// {{- end -}}
// {{- $primarykey := join .PrimaryKey  ", " -}}
// {{ if ne $primarykey "" }},
// 	PRIMARY KEY ( {{ $primarykey }} ) {{ end }}
// )
// `

// // func PGCreateTableStatement[T any](schema, table string) string {
// // 	return createSQLStatement[T](schema, table, pgCreateTableTpl)
// // }

// var pgInsertRowsTpl = `
// INSERT INTO {{ .Schema -}}{{- .Table }} (
// {{- range  $index, $field := .Fields -}}
// {{- if ne $index 0 -}},{{- end }}
// 	{{ $field -}}
// {{- end }}
// ) VALUES (
// {{- range  $index, $field := .Fields -}}
// {{- if ne $index 0 -}},{{- end }}
// 	:{{ $field -}}
// {{- end }}
// )
// `

// {{- range $dataindex, $row := .Data -}}
// {{- if ne $dataindex 0 -}},{{- end }} (
// {{- range $index, $field := $row -}}
// 	{{ $field -}}
// {{- end -}} )
// {{- end -}}
// ('pigeon', 'common in cities')
// ('eagle', 'bird of prey');

// func PGInsertRowsStatement[T any](schema, table string) string {
// 	return createSQLStatement[T](schema, table, pgInsertRowsTpl)
// }

// func LetsInsertRows[T any](schema, table string, data []T) string {
// 	tpl := pgInsertRowsTpl

// 	v := reflect.ValueOf(*new(T))
// 	iv := reflect.Indirect(v)
// 	fields := StructFields(iv)
// 	sqlFieldDefs := []string{}
// 	primaryKey := []string{}
// 	for i := 0; i < iv.NumField(); i++ {
// 		sf := fields[i]
// 		tags := strings.Split(sf.Tag.Get("sql"), ",")
// 		tagMap := map[string]bool{}
// 		for i, v := range tags {
// 			tag := strings.ToLower(v)
// 			tags[i] = tag
// 			if i != 0 {
// 				// skip tag[0] - it's the sql name
// 				tagMap[tag] = true
// 			}
// 		}
// 		if !sf.IsExported() || tags[0] == "-" {
// 			continue
// 		}
// 		field := iv.Field(i)
// 		fieldType := reflect.TypeOf(field.Interface())
// 		if field.Kind() == reflect.Pointer {
// 			fieldType = reflect.TypeOf(field.Interface()).Elem()
// 		}
// 		switch fieldType.Kind() {
// 		case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
// 			continue
// 		}

// 		sqlName := strings.ToLower(sf.Name)
// 		if tags[0] != "" {
// 			sqlName = tags[0]
// 		}

// 		// sqlType := GoToPGTypeMap[fieldType.Name()]
// 		if tagMap["primarykey"] || tagMap["primary key"] {
// 			primaryKey = append(primaryKey, sqlName)
// 			delete(tagMap, "primarykey")
// 			delete(tagMap, "primary key")
// 		}
// 		constraints := ""
// 		for _, v := range tags {
// 			if tagMap[v] {
// 				constraints = fmt.Sprintf("%s %s", constraints, v)
// 			}
// 		}
// 		sqlFieldDefs = append(sqlFieldDefs, sqlName)

// 	}
// 	// fmt.Println(strings.Join(sqlFieldDefs, ",\n"))

// 	tplValues := struct {
// 		Schema     string
// 		Table      string
// 		Fields     []string
// 		PrimaryKey []string
// 	}{
// 		Schema:     schema,
// 		Table:      table,
// 		Fields:     sqlFieldDefs,
// 		PrimaryKey: primaryKey,
// 	}
// 	var buf bytes.Buffer

// 	if err := template.Must(template.
// 		New("").
// 		Funcs(template.FuncMap{"join": strings.Join}).
// 		Parse(tpl)).
// 		Execute(&buf, tplValues); err != nil {
// 		log.Println(err)
// 		panic(err)
// 	}
// 	return buf.String()
// 	// return ""
// }
