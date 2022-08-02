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

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBClient interface {
	// Open(cfg *DBClientConfig) (*sqlx.DB, error)
	// Ping(ctx context.Context) (bool, error)
	// DB() *sqlx.DB
	CreateSchema(ctx context.Context, schema string) (sql.Result, error)
	CreateTable(ctx context.Context, table Table) (sql.Result, error)
	DropTable(ctx context.Context, table Table) (sql.Result, error)
	InsertRows(ctx context.Context, data Dataset) (sql.Result, error)
	QueryConstraints(ctx context.Context, data Dataset) error
}

// type DClient struct {
// 	*sqlx.DB
// }

// func (c *DClient) CreateTable() {
// 	switch c.DriverName() {
// 	case "postgres":
// 	}
// }

type DBClientConfig struct {
	Name               string // descriptive name - used only for logs and reference
	Type               string // descriptive type - used only for logs and reference
	Driver             string // driver name - will be passed to sqlx.DB
	DataSourceName     string // only applicable if using DSN's
	ConnectionString   string // will be parsed to url
	Schema             string // schema to use for this client
	ExpandEnvVars      bool   // whether or not to interpolate env vars of the form ${ENV_VAR} in connection string and dsn
	ExpandFileContents bool   // whether or not to interpolate file contents of the form {file:/path/to/file} in connection string and dsn
	// User             string   // prefer to reference env var or file contents by ${ENV_VAR_NAME} or {file:/path/to/file}
	// Password         string   // prefer to reference env var or file contents by ${ENV_VAR_NAME} or {file:/path/to/file}
	// url *url.URL // let the program handle this
}

func ExampleDBClientConfig() DBClientConfig {
	cfg := DBClientConfig{
		Name:               "database",
		Type:               "odbc",
		Driver:             "odbc",
		ConnectionString:   "${DB_CONN_SCHEME}://${DB_USER}:${DB_PASSWORD}@${DB_HOSTNAME}:${DB_PORT}/${DB_DATABASE}?${DB_ARGUMENTS}",
		Schema:             "qgenda",
		ExpandEnvVars:      true,
		ExpandFileContents: true,
		// User:             "${DB_USER}",
		// Password:         "${DB_PASSWORD}",
	}
	return cfg
}

func (cfg DBClientConfig) String() string {
	s := cfg.ConnectionString
	if cfg.ExpandEnvVars {
		s = ExpandEnvVars(s)
	}
	if cfg.ExpandFileContents {
		s = ExpandFileContents(s)
	}
	return s
}

// OpenDBConnection doesn't technically 'open' a real connection, it follows
// the go default of creating a DB struct that manages connections as needed
func OpenDBClient(cfg *DBClientConfig) (DBClient, error) {
	return nil, nil
}

func NewDBClient(cfg *DBClientConfig) (*sqlx.DB, error) {
	connString := ExpandEnvVars(cfg.ConnectionString)
	// fmt.Printf("Driver: %s\t ConnString: %s\n", cfg.Driver, connString)
	return sqlx.Open(cfg.Driver, connString)
}

type Table struct {
	Name            string
	Schema          string
	Temporary       bool
	UUID            string
	Constraints     map[string]string
	Fields          Fields
	FlattenChildren bool // by default, slices and maps will be handled by creating a child table for each and 'flattening' any nested slices or maps
	Tags            map[string][]string
	UpdateStrategy  string
	Parent          string // parent table, for 'child' tables
}

type Fields []Field

// type UpdateStrategy int8

// func (u UpdateStrategy) String() string {
// 	return ""
// }

// const (
// 	AppendNew UpdateStrategy = iota
// 	FullReplace
// 	AppendAll

// )

func StructToTable[T any](value T, name, schema string, temporary bool, id string, constraints map[string]string, tags map[string][]string, parent string) Table {

	fields := StructToFields(value)

	if name == "" {
		name = strings.ToLower(reflect.ValueOf(value).Type().Name())
	}

	if len(constraints) == 0 {
		constraints = map[string]string{}
		pk := strings.Join(PrimaryKey(fields), ", ")
		if pk != "" {
			constraints["primarykey"] = pk

		}
		uf := []string{}
		for _, v := range UniqueFields(fields) {

			uf = append(uf, PGName(v))
		}
		if len(uf) > 0 {
			constraints["unique"] = strings.Join(uf, ", ")
		}
	}

	if id == "" {
		id = strings.ReplaceAll(uuid.NewString(), "-", "")
	}
	// check https://www.postgresql.org/docs/current/limits.html for current
	// identifier limits. Limit is 63 at time of coding.
	pgIDLimit := 63
	idLength := len(id)
	permIDLength := len("_tmp_") + len(name)
	maxIDLength := pgIDLimit - permIDLength
	if idLength > maxIDLength {
		id = id[0:maxIDLength]
		log.Printf("length of _tmp_[id]_[name] is %d, exceeding postgres identifier limit of %d, truncating [id] to %d characters: %s", (permIDLength + idLength), pgIDLimit, maxIDLength, id)
	}
	return Table{
		Name:            name,
		Schema:          schema,
		Temporary:       temporary,
		UUID:            id,
		Constraints:     constraints,
		Fields:          fields,
		FlattenChildren: true,
		Tags:            tags,
		Parent:          parent,
	}

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

func FieldNames(fields []Field) []string {
	var fn []string
	for _, field := range fields {
		fn = append(fn, field.Name)
	}
	return fn
}

func JoinStringSlice(sep string, s []string) string {
	return strings.Join(s, sep)
}

// SQLResult combines any number of sql.Result's
func SQLResult(res ...sql.Result) Result {
	var lis, ras int64
	var lies, raes error
	for _, r := range res {
		if r == nil {
			continue
		}
		li, lie := r.LastInsertId()
		ra, rae := r.RowsAffected()
		lis = li
		ras = ras + ra
		lies = fmt.Errorf("[%v]: [%w]", lie, lies)
		raes = fmt.Errorf("[%v]: [%w]", rae, raes)
	}
	return Result{
		lastInsertID:      lis,
		lastInsertIDError: lies,
		rowsAffected:      ras,
		rowsAffectedError: raes,
	}
}

// Result is used to satisfy the sql.Result interface and enable aggregating multiple sql.Results
type Result struct {
	lastInsertID      int64
	lastInsertIDError error
	rowsAffected      int64
	rowsAffectedError error
}

func (r Result) LastInsertId() (int64, error) {
	return r.lastInsertID, r.lastInsertIDError
}

func (r Result) RowsAffected() (int64, error) {
	return r.rowsAffected, r.rowsAffectedError
}

// TableStatement is primarily intended for creating SQL statements from a Table
// and a template. It includes a handful of funcs, but accepts a funcmap that can
// overwrite the included funcs
func TableStatement(table Table, tpl string, funcs template.FuncMap) string {

	var buf bytes.Buffer

	if err := template.Must(template.
		New("").
		Option("missingkey=zero").
		Funcs(template.FuncMap{
			"join":               strings.Join,
			"joinss":             JoinStringSlice,
			"qfname":             QueryFieldName,
			"uniquefields":       UniqueFields,
			"fieldswithtagvalue": FieldsWithTagValue,
			"fieldnames":         FieldNames,
		}).
		Funcs(funcs).
		Parse(tpl)).
		Execute(&buf, table); err != nil {
		log.Println(err)
		panic(err)
	}
	return buf.String()
}

// UniqueFields is intended to be used in templates and is included in the
// default TableStatement funcmap as uniquefields
func UniqueFields(fields []Field) []Field {
	f := []Field{}
	for _, field := range fields {
		if field.Unique {
			f = append(f, field)
		}
	}
	return f
}

// FieldsWithTagValue returns only those fields with the given key-value pair
// it is included in the TableStatement funcmap as fieldswithtagvalue
func FieldsWithTagValue(fields []Field, key, value string) []Field {
	f := []Field{}
	for _, field := range fields {
		tagSlice, ok := field.Tags[key]
		if ok && len(tagSlice) > 0 {
			for _, tagi := range tagSlice {
				if tagi == value {
					f = append(f, field)
				}
			}
		}
	}
	return f
}

// FieldsWithoutTagValue returns only those fields with the given key-value pair
// it is included in the TableStatement funcmap as fieldswithtagvalue
func FieldsWithoutTagValue(fields []Field, key, value string) []Field {
	f := []Field{}

	for _, field := range fields {
		tagSlice, ok := field.Tags[key]
		if ok && len(tagSlice) > 0 {
			for _, tagi := range tagSlice {
				if tagi == value {
					f = append(f, field)
				}
			}
		}
	}
	return f
}

// FieldHasTagValue returns true if the key: value exists in the given tag
// it is included in the TableStatement funcmap as fieldhastagvalue
func FieldHasTagValue(field Field, key, value string) bool {

	tag, ok := field.Tags[key]
	if !ok {
		return false
	}
	for _, tv := range tag {
		if tv == value {
			return true
		}
	}
	return false
}

// Field.HasTagValue is a direct wrap of FieldHasTagValue
func (f Field) HasTagValue(key, value string) bool {
	return FieldHasTagValue(f, key, value)
}

// WithTagValue returns fields that test true for Field.HasTagValue
func (f Fields) WithTagValue(key, value string) Fields {
	ff := Fields{}
	for _, field := range f {
		if field.HasTagValue(key, value) {
			ff = append(ff, field)
		}
	}
	return ff
}

// WithoutTagValue returns fields that test false for Field.HasTagValue
func (f Fields) WithoutTagValue(key, value string) Fields {
	ff := Fields{}
	for _, field := range f {
		if !field.HasTagValue(key, value) {
			ff = append(ff, field)
		}
	}

	return ff
}
