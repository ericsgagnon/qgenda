package qgenda

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"text/template"

	"github.com/exiledavatar/gotoolkit/meta"
	"github.com/exiledavatar/gotoolkit/typemap"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var PGTag = "pg"

func CreatePGTable[T any](ctx context.Context, db *sqlx.DB, value T, schema, table string) (sql.Result, error) {

	if schema == "" {
		schema = "qgenda"
	}
	str, err := meta.NewStruct(value, meta.Structconfig{
		NameSpace: []string{schema},
		Name:      table,
		Tags:      meta.ToTags(fmt.Sprintf(`table:"%s"`, table)),
	})
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	result, err := tx.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", str.NameSpace[0]))
	if err != nil {
		return result, err
	}

	tpl := `{{- "\n" -}}
	CREATE TABLE IF NOT EXISTS {{ .Struct.TagIdentifier "table" | tolower }} (
		{{- $fields := .Struct.Fields.WithTagTrue "db" -}}
		{{- $names := $fields.TagNames "db" -}}
		{{- $types := $fields.NonEmptyTagValues "pgtype" -}}
		{{- $columnDefs := joinslices "\t" ",\n\t" $names $types -}}
		{{- print "\n\t" $columnDefs -}}
		{{- $primarykeyfields := .Struct.Fields.WithTagTrue "primarykey" -}}
		{{- $primarykey := $primarykeyfields.TagNames "db" | join ", " -}}
		{{- if ne $primarykey "" -}}{{- printf ",\n\tPRIMARY KEY ( %s )" $primarykey -}}{{- end -}}
		{{- "\n)" -}}
	`
	funcs := template.FuncMap{
		// "gotopgtype": qgenda.GoToPGType,
		// "joinslices": meta.JoinSlices,
	}

	data := map[string]any{
		"postgres": typemap.TypeMaps["postgres"].ToType,
		// "Schema":   schema,
		// "Table":    table,
	}

	query, err := str.ExecuteTemplate(tpl, funcs, data)
	// fmt.Println(query)
	if err != nil {
		return nil, err
	}

	if result, err := tx.ExecContext(ctx, query); err != nil {
		return result, err
	}

	children := str.Fields().WithTagTrue("table").ToStructs()
	for _, child := range children {
		query, err := child.ExecuteTemplate(tpl, funcs, data)
		if err != nil {
			return nil, err
		}
		if result, err := tx.ExecContext(ctx, query); err != nil {
			return result, err
		}

	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func GetPGStatus[T any](ctx context.Context, db *sqlx.DB, value T, schema, table, tpl string) (*RequestConfig, error) {
	if schema == "" {
		schema = "qgenda"
	}

	str, err := meta.NewStruct(value, meta.Structconfig{
		NameSpace: []string{schema},
		Name:      table,
	})
	if err != nil {
		return nil, err
	}

	if tpl == "" {
		tpl = `
		{{- $field := .Struct.Fields.ByName "LastModifiedDateUTC" -}}
		select max ( {{ $field.TagName "db" }} ) {{ $field.TagName "qgendarequestname" }}
		from {{ .Struct.TagIdentifier "table" | tolower }}
		`
	}

	query, err := str.ExecuteTemplate(tpl, nil, nil)
	if err != nil {
		return nil, err
	}

	rc := RequestConfig{}
	if err := db.GetContext(ctx, &rc, query); err != nil {
		return nil, err
	}

	return &rc, nil

}

func PutPG[S ~[]T, T any](ctx context.Context, db *sqlx.DB, value S, schema, table string) (sql.Result, error) {
	if len(value) < 1 {
		return nil, fmt.Errorf("%T.PGInsertRows: length of %T < 1, nothing to do", value, value)
	}

	if schema == "" {
		schema = "qgenda"
	}

	st, err := meta.NewStruct(value, meta.Structconfig{
		NameSpace: []string{schema},
		Name:      table,
		Tags:      meta.ToTags(fmt.Sprintf(`table:"%s"`, table)),
	})
	if err != nil {
		return nil, err
	}
	structs := st.Extract(st.Fields().WithTagTrue("table").Names()...)
	// fmt.Println("Prior to Temp Tables: ----------------------------------------------------------")
	// for k, str := range structs {
	// 	fmt.Printf("%20s%30s%30s\n", k, str.TagName("table"), str.TagIdentifier("table"))
	// }
	// fmt.Println("--------------------------------------------------------------------------------")
	// for _, str := range structs {
	// 	switch {
	// 	case str.Tags.Exists("table"):
	// 		// do nothing
	// 	case str.Parent != nil:
	// 		// parentTableName := strings.ToLower(str.Parent.TagName("table"))

	// 	case str.Parent == nil:

	// 	}
	// }
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	for _, str := range structs {
		result, err := CreatePGTable(ctx, db, str, schema, str.TagName("table"))
		if err != nil {
			return result, err
		}
	}
	// fmt.Println("Prior to Temp Tables: ----------------------------------------------------------")
	// for k, str := range structs {
	// 	fmt.Printf("%20s%30s%30s\n", k, str.TagName("table"), str.TagIdentifier("table"))
	// }
	// fmt.Println("--------------------------------------------------------------------------------")

	// temporary tables
	tpl := `create temporary table if not exists _tmp_{{- .Struct.TagName "table" | tolower }} (like {{ .Struct.TagIdentifier "table" | tolower }} excluding constraints)`
	for _, str := range structs {
		query, err := str.ExecuteTemplate(tpl, nil, nil)
		if err != nil {
			return nil, err
		}
		if result, err := tx.ExecContext(ctx, query); err != nil {
			return result, err
		}
	}

	// insert to temp tables
	tpl = `{{- "\n" -}}
	insert into _tmp_{{- .Struct.TagName "table" | tolower }} ( 
		{{- $fields := .Struct.Fields.WithTagTrue "db" -}}
		{{ $fields.TagNames "db" | join ", " }} )
		values (
			{{- $fields := .Struct.Fields.WithTagTrue "db" -}}
			:
			{{- $fields.TagNames "db" | join ", :" -}}
			{{- "\n)" }}
			`
	// postgres has a 65535 'parameter' limit, there is an 'unnest' work around, but for now we're just going to chunk it
	for k, str := range structs {
		chunkSize := 65535 / len(str.Fields())
		query, err := str.ExecuteTemplate(tpl, nil, nil)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(str.Data); i += chunkSize {
			j := i + chunkSize
			if j > len(str.Data) {
				j = len(str.Data)
			}
			values := str.Data[i:j]
			result, err := tx.NamedExecContext(ctx, query, values)
			if err != nil {
				return result, err
			}
			log.Printf("%-25s[%10d:%10d] Rows: %10d Fields: %10d Total Parameters: %10d\n", k, i, j, len(values), len(str.Fields()), (len(values) * len(str.Fields())))
		}

	}
	// fmt.Println("Prior to Update: -----------------------------------------------------------------------------------------")
	// for k, str := range structs {
	// 	fmt.Println(k, " ", str.TagName("table"))
	// 	fmt.Println(k, " ", str.TagIdentifier("table"))
	// }

	// fmt.Println("----------------------------------------------------------------------------------------------------------")
	// update from temp to permanent tables
	updateTpl := `{{- "\n" -}}
	insert into {{ .Struct.TagIdentifier "table" | tolower }} (
		select distinct on (tmp._id_hash) tmp.*
		from _tmp_{{- .Struct.TagName "table" }} tmp 
		{{- if ne .Struct.Parent nil }}
		inner join _tmp_{{ .Struct.Parent.TagName "table" }} ptmp
		{{- $parentprimarykey := (index ( .Struct.Parent.Fields.WithTagTrue "primarykey" ) 0 ).TagName "db" -}}
		{{- $parentpkey := ( index ( .Struct.Fields.WithTagTrue "parentprimarykey" ) 0  ).TagName "db" }} 
		on tmp.{{ $parentpkey }} = ptmp.{{ $parentprimarykey -}}
		{{ end }}
		where not exists (
			select 1
			from {{ .Struct.TagIdentifier "table" }} dst
			{{- $pkey := index ( ( .Struct.Fields.WithTagTrue "primarykey" ).TagNames "db" ) 0 }}
			where dst.{{ $pkey }} = tmp.{{ $pkey }}
		) 
	)
	`

	// update permanent tables from temp tables
	for _, str := range structs {
		// fmt.Println("updating ", k)
		// fmt.Println(str.TagIdentifier("table"))
		query, err := str.ExecuteTemplate(updateTpl, nil, nil)
		if err != nil {
			return nil, err
		}
		// fmt.Println(query)
		result, err := tx.ExecContext(ctx, query)
		if err != nil {
			return result, err
		}
		// log.Println(result)
	}
	// fmt.Println(db.DriverName())
	return nil, tx.Commit()
}

func DropPGTable[T any](ctx context.Context, db *sqlx.DB, value T, force bool, schema, table string) (sql.Result, error) {
	if schema == "" {
		schema = "qgenda"
	}
	str, err := meta.NewStruct(value, meta.Structconfig{
		NameSpace: []string{schema},
		Name:      table,
		Tags:      meta.ToTags(fmt.Sprintf(`table:"%s"`, table)),
	})
	if err != nil {
		return nil, err
	}

	tpl := `drop table if exists {{ .Struct.TagIdentifier "table" | tolower }}`
	if force {
		tpl = tpl + ` cascade `
	}
	query, err := str.ExecuteTemplate(tpl, nil, nil)
	if err != nil {
		return nil, err
	}
	return db.ExecContext(ctx, query)
}

func DropPGSchema(ctx context.Context, db *sqlx.DB, force bool, schema string) (sql.Result, error) {
	if schema == "" {
		schema = "qgenda"
	}
	var cascade string
	if force {
		cascade = "cascade"
	}
	query := fmt.Sprintf(`drop schema if exists %s %s`, schema, cascade)
	return db.ExecContext(ctx, query)
}
