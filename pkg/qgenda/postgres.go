package qgenda

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"text/template"

	"github.com/jackc/pgx/v4"
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
		Schema:           "qgenda",
	}

}

type PGClient struct {
	*sqlx.DB
	Config DBClientConfig
}

// copied here just to know the methods I need to implement to satisfy the interface
// type DBClient interface {
// 	// Open(cfg *DBClientConfig) (*sqlx.DB, error)
// 	// Ping(ctx context.Context) (bool, error)
// 	DB() *sqlx.DB
// 	CreateSchema(ctx context.Context, schema string) (sql.Result, error)
// 	CreateTable(ctx context.Context, table Table) (sql.Result, error)
// 	DropTable(ctx context.Context, table Table) (sql.Result, error)
// 	InsertRows(ctx context.Context, data Dataset) (sql.Result, error)
// 	QueryConstraints(ctx context.Context, data Dataset) error
// }

func (c *PGClient) CreateSchema(ctx context.Context, schema string) (sql.Result, error) {
	return c.ExecContext(
		ctx,
		fmt.Sprintf("create schema if not exists %s", pgx.Identifier{schema}.Sanitize()),
	)
}

func CreateTable(ctx context.Context, table Table) (sql.Result, error) {
	// return pgCreateTable[]()

	// return c.ExecContext(
	// 	ctx,
	// 	// PGCreateTableStatement(value[0], schema, table),
	// 	PGStatement(*new(T), schema, table, pgCreateNewTableTpl),
	// )
	return nil, nil
}

func (c *PGClient) CreateTable(ctx context.Context, db *sqlx.DB, value []any, schema, table string) (sql.Result, error) {
	return PGCreateTable(ctx, c.DB, value, schema, table)
}

func (c *PGClient) DropTable(ctx context.Context, db *sqlx.DB, value []any, schema, table string) (sql.Result, error) {
	return PGDropTable(ctx, c.DB, value, schema, table)
}

func (c *PGClient) InsertRows(ctx context.Context, db *sqlx.DB, value []any, schema, table string) (sql.Result, error) {
	return PGInsertRows(ctx, c.DB, value, schema, table)
}

func (c *PGClient) QueryConstraints(ctx context.Context, db *sqlx.DB, value []any, schema, table string) (sql.Result, error) {
	return PGQueryConstraint(ctx, c.DB, value, schema, table)
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
	var allfields []Field
	allfields = StructToFields(*new(T))
	// fmt.Println("PGStatement---------------------------------------")
	var fields []Field
	for _, field := range allfields {
		if PGOmit(field) {
			continue
		}
		fields = append(fields, field)
	}
	// if schema != "" {
	// 	schema = fmt.Sprintf("%s.", schema)
	// }
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

var pgCreateSchemaTpl = `
CREATE SCHEMA IF NOT EXISTS {{ .Schema }}
`

var pgCreateNewTableTpl = `
CREATE TABLE IF NOT EXISTS {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Table }} (
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

var pgDropTableTpl = `
DROP TABLE IF EXISTS {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Table }}
`

//create unique index if not exists schedulestafftag_all_columns_unique on schedulestafftag (schedulekey, lastmodifieddateutc, categorykey, categoryname, tagkey, tagname)

var pgInsertTpl = `
INSERT INTO {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Table }} (
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
	{{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Table }}
`

func PGCreateTableStatement[T any](value T, schema, table string) string {
	return PGStatement(value, schema, table, pgCreateNewTableTpl)
}

func PGDropTableStatement[T any](value T, schema, table string) string {
	return PGStatement(value, schema, table, pgDropTableTpl)
}

func PGInsertStatement[T any](value T, schema, table string) string {
	return PGStatement(value, schema, table, pgInsertTpl)
}

func PGQueryConstraintsStatement[T any](value T, schema, table string) string {
	return PGStatement(value, schema, table, pgSelectMaxConstraintsTpl)
}

func PGCreateSchema[T any](ctx context.Context, db *sqlx.DB, value []T, schema, table string) (sql.Result, error) {
	// fmt.Printf("PGCreateSchema: %T\n", *new(T))
	return db.ExecContext(
		ctx,
		PGStatement(*new(T), schema, table, pgCreateSchemaTpl),
	)
	// return db.NamedExecContext(
	// 	ctx,
	// 	PGStatement(*new(T), schema, table, pgCreateSchemaTpl),
	// 	value,
	// )

}

func pgCreateTable[T any](ctx context.Context, db *sqlx.DB, value T, schema, table string) (sql.Result, error) {
	// fmt.Println(PGCreateTableStatement(value[0], schema, table))
	// fmt.Println(PGStatement(*new(T), schema, table, pgCreateNewTableTpl))
	// if IsSlice(value) || IsMap(value){
	// 	return db.ExecContext(
	// 		ctx,
	// 		PGStatement(*new(T), schema, table, pgCreateNewTableTpl),
	// 	)

	// }
	return db.ExecContext(
		ctx,
		// PGCreateTableStatement(value[0], schema, table),
		PGStatement(*new(T), schema, table, pgCreateNewTableTpl),
	)
}

func PGCreateTable[T any](ctx context.Context, db *sqlx.DB, value []T, schema, table string) (sql.Result, error) {
	// fmt.Println(PGCreateTableStatement(value[0], schema, table))
	// fmt.Println(PGStatement(*new(T), schema, table, pgCreateNewTableTpl))

	return db.ExecContext(
		ctx,
		// PGCreateTableStatement(value[0], schema, table),
		PGStatement(*new(T), schema, table, pgCreateNewTableTpl),
	)
}

func PGDropTable[T any](ctx context.Context, db *sqlx.DB, value []T, schema, table string) (sql.Result, error) {
	// fmt.Println(PGCreateTableStatement(value[0], schema, table))
	// fmt.Println(PGStatement(*new(T), schema, table, pgDropTableTpl))
	return db.ExecContext(
		ctx,
		// PGDropTableStatement(value[0], schema, table),
		PGStatement(*new(T), schema, table, pgDropTableTpl),
		// PGStatement(*new(T), schema, table, pgDropTableTpl),
	)
}

func PGInsertRows[T any](ctx context.Context, db *sqlx.DB, value []T, schema, table string) (sql.Result, error) {
	if len(value) < 1 {
		return nil, fmt.Errorf("PGInsertRows: length of %T < 1, nothing to do", value)
	}
	// postgres has a 65535 'parameter' limit, there is an 'unnest' work around, but for now we're just going to chunk it
	chunkSize := 65535 / reflect.ValueOf(value[0]).NumField()
	// var rowsAffected int64
	var res Result
	for i := 0; i < len(value); i = i + chunkSize {
		j := i + chunkSize
		if j > len(value) {
			j = len(value)
		}
		sqlResult, err := db.NamedExecContext(
			ctx,
			// PGInsertStatement(value[0], schema, table),
			PGStatement(*new(T), schema, table, pgInsertTpl),
			value[i:j],
		)
		if sqlResult != nil {
			// ra, _ := sqlResult.RowsAffected()

			// fmt.Printf("Insert %T[%d:%d]: RowsAffected: %d\n", value[i], i, j, ra)
			res = SQLResult(res, sqlResult)

			// rowsAffected = rowsAffected + ra

		}
		if err != nil {
			return res, err
		}

	}

	return res, nil
}

func PGQueryConstraint[T any](ctx context.Context, db *sqlx.DB, value []T, schema, table string) (sql.Result, error) {
	// fmt.Println(PGQueryConstraintsStatement(value[0], schema, table))
	result, err := db.NamedExecContext(
		ctx,
		// PGQueryConstraintsStatement(value[0], schema, table),
		PGStatement(*new(T), schema, table, pgSelectMaxConstraintsTpl),
		value,
	)
	return result, err
}

func PGStatementDev[T any](value T, tpl string) string {
	
	return ""
}

func pgStatement[T any](value T, schema, table, tpl string) string {
	var allfields []Field
	allfields = StructToFields(*new(T))
	// fmt.Println("PGStatement---------------------------------------")
	var fields []Field
	for _, field := range allfields {
		if PGOmit(field) {
			continue
		}
		fields = append(fields, field)
	}
	// if schema != "" {
	// 	schema = fmt.Sprintf("%s.", schema)
	// }
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
