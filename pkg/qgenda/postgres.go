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
	Tx     *sqlx.Tx // for use in transactions
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

// CreateSchema uses 'if not exists' to guarantee idompetency. Generally, owner is omitted, which will
// default to current user. If schema is omitted, it will default to PGClient.Config.Schema
func (c *PGClient) CreateSchema(ctx context.Context, schema string, owner string) (sql.Result, error) {
	if schema == "" {
		schema = c.Config.Schema
	}
	if owner != "" {
		owner = fmt.Sprintf("authorization %s ", pgx.Identifier{owner}.Sanitize())
	}
	return c.ExecContext(
		ctx,
		fmt.Sprintf("create schema if not exists %s %s", pgx.Identifier{schema}.Sanitize(), owner),
	)
}

// DropSchema uses 'if exists' to guarantee idempotency. Cascade will remove all objects that depend on
// an element in the schema.
func (c *PGClient) DropSchema(ctx context.Context, schema string, force bool) (sql.Result, error) {
	if schema == "" {
		schema = c.Config.Schema
	}
	var sql string
	if force {
		sql = fmt.Sprintf("drop schema if exists %s cascade", pgx.Identifier{schema}.Sanitize())
	} else {
		sql = fmt.Sprintf("drop schema if exists %s restrict", pgx.Identifier{schema}.Sanitize())
	}
	return c.ExecContext(ctx, sql)
}

// ForceDropSchema uses 'if exists' to guarantee idempotency and cascade to force the drop regardless of
// dependent objects in this or other schemas.
func (c *PGClient) ForceDropSchema(ctx context.Context, schema string) (sql.Result, error) {
	if schema == "" {
		schema = c.Config.Schema
	}
	return c.ExecContext(
		ctx,
		fmt.Sprintf("drop schema if exists %s cascade", pgx.Identifier{schema}.Sanitize()),
	)
}

// CreateTable uses 'if exists' to guarantee idempotency
func (c *PGClient) CreateTable(ctx context.Context, table Table, tx bool) (sql.Result, error) {
	result, err := c.CreateSchema(ctx, table.Schema, "")
	if err != nil {
		return result, err
	}
	// sql := PGTableStatement(table, PGCreateTableDevTpl, nil)
	if tx {
		if c.Tx == nil {
			c.Tx = c.MustBegin()
		}
		return c.Tx.ExecContext(ctx, PGTableStatement(table, PGCreateTableDevTpl, nil))
	}
	return c.ExecContext(ctx, PGTableStatement(table, PGCreateTableDevTpl, nil))
}

// DropTable uses 'if exists' to guarantee idempotency. Cascade will remove all objects that depend on
// the table.
func (c *PGClient) DropTable(ctx context.Context, table Table, force bool) (sql.Result, error) {
	if table.Name == "" {
		return nil, fmt.Errorf("drop table: table.Name appears empty: %v", table)
	}
	var sql string
	if force {
		sql = fmt.Sprintf("drop table if exists %s cascade", pgx.Identifier{table.Name}.Sanitize())
	} else {
		sql = fmt.Sprintf("drop table if exists %s restrict", pgx.Identifier{table.Name}.Sanitize())
	}
	return c.ExecContext(ctx, sql)
}

// ForceDropTable uses 'if exists' to guarantee idempotency and cascade to force the drop regardless of
// dependent objects in its own or other schemas.
func (c *PGClient) ForceDropTable(ctx context.Context, table Table) (sql.Result, error) {
	if table.Name == "" {
		return nil, fmt.Errorf("drop table: table.Name appears empty: %v", table)
	}
	sql := fmt.Sprintf("drop table if exists %s cascade", pgx.Identifier{table.Name}.Sanitize())
	return c.ExecContext(ctx, sql)
}

func (c *PGClient) InsertRows(ctx context.Context, table Table, value []any, tx bool) (sql.Result, error) {

	return PGInsertRowsDev(ctx, c.DB, table, value)

}

// DB is to enable using either (*sqlx.DB or *sqlx.TX), the compiler doesn't seem to really support
// our use case for generics yet, which seems silly, so we're just treating it like a regular interface...
type DB interface {
	// *sqlx.Tx | *sqlx.DB
	NamedExecContext(context.Context, string, interface{}) (sql.Result, error)
}

func DBTxx(ctx context.Context, db DB) (*sqlx.Tx, error) {
	var dbi interface{} = db
	var tx *sqlx.Tx
	var err error
	switch dbi.(type) {
	case *sqlx.Tx:
		tx = (dbi).(*sqlx.Tx)
	case *sqlx.DB:
		tx, err = (dbi).(*sqlx.DB).BeginTxx(ctx, nil)
		if err != nil {
			return nil, err
		}
	// case *sql.DB:
	// 	db, ok := (dbi).(*sql.DB)
	// 	if !ok {
	// 		return nil, fmt.Errorf("unable to assert %T to *sql.DB", db)
	// 	}
	// 	sqlx.NewDb(db, db.DriverName())
	default:
		return nil, fmt.Errorf("DBTxx received a %T and doesn't know what to do", db)
	}
	return tx, nil
}

// func PGInsertRowsTx[T any](ctx context.Context, db DB, table Table, value []T) (sql.Result, error) {

// 	db.NamedExecContext(ctx, PGStatement(*new(T), table.Schema, table.Name, pgInsertTpl), value[0:1])
// 	return nil, nil
// }

func PGInsertTx[T any](ctx context.Context, tx *sqlx.Tx, table Table, tpl string, value []T) (sql.Result, error) {

	// PGCreateTableDevTpl
	if tpl == "" {
		tpl = PGInsertRowsDevTpl
	}

	sqlStatement := PGStatement(*new(T), table.Schema, table.Name, pgInsertTpl)
	return tx.NamedExecContext(ctx, sqlStatement, value)
}

// func PGInsertRowsDev[T any](ctx context.Context, db *sqlx.DB, table Table, value []T) (sql.Result, error) {
func PGInsertRowsDev[T any](ctx context.Context, db DB, table Table, value []T) (sql.Result, error) {
	tx, err := DBTxx(ctx, db)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%T\n", db)
	if len(value) < 1 {
		return nil, fmt.Errorf("PGInsertRows: length of %T < 1, nothing to do", value)
	}
	// tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	// postgres has a 65535 'parameter' limit, there is an 'unnest' work around, but for now we're just going to chunk it
	chunkSize := 65535 / reflect.ValueOf(value[0]).NumField()

	var res Result
	for i := 0; i < len(value); i = i + chunkSize {
		j := i + chunkSize
		if j > len(value) {
			j = len(value)
		}
		sqlResult, err := db.NamedExecContext(
			ctx,
			// PGInsertStatement(value[0], schema, table),
			PGStatement(*new(T), table.Schema, table.Name, pgInsertTpl),
			value[i:j],
		)
		res = SQLResult(res, sqlResult)
		if err != nil {
			return res, err
		}

		// for _, field := range table.Fields {
		// 	switch field.Kind {
		// 	case "slice", "map":
		// 		continue
		// 		// v := reflect.ValueOf(value[i]).FieldByName(field.Name)
		// 		// PGInsertRowsDev(ctx, tx, table, v.InterfaceData())
		// 		// name := strings.ToLower(fmt.Sprintf("%s%s", table.Name, reflect.ValueOf()))
		// 		// table := StructToTable()
		// 		// PGInsertRowsDev(ctx, tx, table, )
		// 	default:
		// 		continue
		// 	}
		// }

	}

	tx.Commit()
	return res, nil
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
	nametags, ok := field.Tags()[PGTag]
	if !ok {
		nametags, ok = field.Tags()["db"]
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

func PGNames(fields []Field) []string {
	var fn []string
	for _, field := range fields {
		fn = append(fn, PGName(field))
	}
	return fn
}

func PGStatement[T any](value T, schema, table, tpl string) string {
	// var allfields []Field
	allfields := StructToFields(*new(T))
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
			"joinss":        JoinStringSlice,
			"pgtype":        GoToPGType,
			"pgname":        PGName,
			"pgnames":       PGNames,
			"pgqueryfields": PGQueryConditionFields,
			"pgomit":        PGOmit,
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
		if _, ok := field.Tags()["querycondition"]; ok {
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

func PGCreateTableDev(ctx context.Context, db *sqlx.DB, table Table) (sql.Result, error) {
	return db.ExecContext(ctx, PGTableStatement(table, PGCreateTableDevTpl, nil))

}

var PGCreateTableDevTpl = `
CREATE {{- if .Temporary }} TEMPORARY TABLE IF NOT EXISTS _tmp_{{- .Name -}}
{{ else }} TABLE IF NOT EXISTS {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }}
{{- end }} (
{{- $fields := pgincludefields .Fields -}}
{{- range  $index, $field := $fields -}}
{{- if ne $index 0 -}},{{- end }}
	{{ pgname $field }} {{ pgtype $field.Type }} {{ if $field.Unique }} unique {{ end -}} {{- if not $field.Nullable -}} not null {{- end }}
{{- end -}}
{{- if not .Temporary }}
{{- $pk := .Constraints.primarykey -}}
{{- $uf := .Constraints.unique -}}
{{ if $pk }}, 
PRIMARY KEY ( {{ $pk }} ) 
{{- end -}}
{{- if and $pk $uf -}},{{- end -}}
{{- if $uf }}
CONSTRAINT {{ .Name -}}_unique UNIQUE ( {{ $uf }} )
{{- end }}
{{- end }}
)
`

var PGInsertRowsDevTpl = `
INSERT INTO 
{{- if .Temporary }}  _tmp_{{- .Name -}} 
{{- else }} {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name -}}
{{- end }} (
	{{- $fields := pgincludefields .Fields -}}
	{{- range  $index, $field := $fields -}}
	{{- if ne $index 0 -}},{{- end }}
		{{ pgname $field }}
	{{- end }} 
	) VALUES (
	{{- range  $index, $field := pgincludefields .Fields -}}
	{{- if ne $index 0 -}},{{- end }}
		:{{ pgname $field }}
	{{- end }}	
)
{{- if not .Temporary }} ON CONFLICT (
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
{{ end }}
`

// {{ if .Constraints }}
// ON CONFLICT (
// 	{{- $primarykey := join .PrimaryKey  ", " -}}
// 	{{ if ne $primarykey "" }}
// 	{{ $primarykey }}
// 	{{ else }}
// 	{{- range  $index, $field := .Fields -}}
// 	{{- if ne $index 0 -}},{{- end }}
// 		{{ pgname $field }}
// 	{{- end -}}
// 	{{- end }}
// 	) DO NOTHING

// ) ON CONFLICT (
// 	{{- $pk := .Constraints.primarykey -}}r
// 	{{- $uf := .Constraints.unique -}}

// {{- $primarykey := join $pk  ", " -}}
// {{ if ne $primarykey "" }}
// {{ $primarykey }}
// {{ else }}
// {{- range  $index, $field := .Fields -}}
// {{- if ne $index 0 -}},{{- end }}
// 	{{ pgname $field }}
// {{- end -}}
// {{- end }}
// ) DO NOTHING

var PGInsertChangesOnlyDevTpl = `

with cte_partitioned_row_numbers as (
	select * , 
	row_number() over ( partition by {{ fieldswithtagvalue .Fields "idtype" "group" | pgnames | joinss " , " }} order by {{ fieldswithtagvalue .Fields "idtype" "order" | pgnames | joinss " desc , " }} desc ) rn
	FROM {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }}
), cte_last_inserted_rows as (
	select * from cte_partitioned_row_numbers cprn where cprn.rn = 1
), cte_new_rows as (
	select distinct * from _tmp_{{- .Name }}
), cte_anti_joined as (
	select
	cnr.*
	from cte_new_rows cnr
	where not exists (
		select 1
		from cte_last_inserted_rows clir where
		{{- $joinfields := .Fields.WithoutTagValue "idtype" "order" | pgincludefields | pgnames -}}
		{{- range  $index, $field := $joinfields  }}
		{{ if ne $index 0 }} and {{ end -}} clir.{{ $field }} = cnr.{{ $field }}
		{{- end }}	
	)
) insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (
	select caj.* from cte_anti_joined caj
)
`

// `
// insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (

// )
// `

var PGInsertSubtableRowsTpl = `
INSERT INTO 
{{- if .Temporary }}  _tmp_{{- .Name -}} 
{{- else }} {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name -}}
{{- end }} (
	{{- $fields := pgincludefields .Fields -}}
	{{- range  $index, $field := $fields -}}
	{{- if ne $index 0 -}},{{- end }}
		{{ pgname $field }}
	{{- end }} 
	) VALUES (
	{{- range  $index, $field := .Fields -}}
	{{- if ne $index 0 -}},{{- end }}
		:{{ pgname $field }}
	{{- end }}	
)	
`

var Xuseforreference = `
INSERT INTO 
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
	{{- $pk := fieldswithtagvalue .Fields "primarykey" "table" | pgnames | joinss ", " -}}
	{{- $uf := .Constraints.unique -}}		
	{{ if ne $pk "" }}
		{{ $pk }}
	{{ else }}
	{{- range  $index, $field := .Fields -}}
	{{- if ne $index 0 -}},{{- end }}
		{{ pgname $field }}
	{{- end -}}	
	{{- end }}
	) DO NOTHING	
`

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

// PGIncludeFields is intended to be used in templates and is included in the default
// PGTableStatement funcmap as pgincludefields
func PGIncludeFields(fields []Field) []Field {
	f := []Field{}
	for _, field := range fields {
		if !PGOmit(field) {
			f = append(f, field)
		}
	}
	return f
}

// PGTableStatement wraps TableStatement and adds several pg specific funcs
func PGTableStatement(table Table, tpl string, funcs template.FuncMap) string {
	fm := template.FuncMap{
		"join":               strings.Join,
		"joinss":             JoinStringSlice,
		"pgtype":             GoToPGType,
		"pgname":             PGName,
		"pgnames":            PGNames,
		"pgqueryfields":      PGQueryConditionFields,
		"qfname":             QueryFieldName,
		"pgomit":             PGOmit,
		"pgincludefields":    PGIncludeFields,
		"uniquefields":       UniqueFields,
		"fieldswithtagvalue": FieldsWithTagValue,
		"fieldnames":         FieldNames,
	}
	for k, v := range funcs {
		fm[k] = v
	}
	return TableStatement(table, tpl, fm)

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
