package qgenda

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Staffs []Staff

func (s *Staffs) Get(ctx context.Context, c *Client, rc *RequestConfig) error {

	req := NewStaffRequest(rc)
	resp, err := c.Do(ctx, req)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	ss := []Staff{}
	if err := json.Unmarshal(data, &ss); err != nil {
		return err
	}

	sourceQuery := resp.Request.URL.String()
	extractDateTime, err := ParseTime(resp.Header.Get("date"))
	if err != nil {
		return err
	}

	if len(ss) > 0 {
		for i, _ := range ss {
			ss[i].SourceQuery = &sourceQuery
			ss[i].ExtractDateTime = &extractDateTime
		}
	}
	*s = ss
	return nil
}

func (s *Staffs) Process() error {
	ss := *s
	for i, _ := range ss {
		if err := ss[i].Process(); err != nil {
			return err
		}
	}
	sort.SliceStable(ss, func(i, j int) bool {
		return *(ss[i].StaffKey) < *(ss[j].StaffKey)
	})
	*s = ss
	return nil
}

func (s *Staffs) GetFromFile(filename string) error {
	ss := []Staff{}
	if err := GetFromFile(filename, &ss); err != nil {
		return err
	}
	*s = ss
	return nil
}

func (s Staffs) PGInsertRows(ctx context.Context, tx *sqlx.Tx, schema, tablename, id string) (sql.Result, error) {

	if len(s) < 1 {
		return nil, fmt.Errorf("%T.PGInsertRows: length of %T < 1, nothing to do", s, s)
	}
	var res Result
	id = strings.ReplaceAll(uuid.NewString(), "-", "")[0:16]
	if tablename == "" {
		tablename = "staff"
	}
	basetable := tablename

	sqlResult, err := s[0].PGCreateTable(ctx, tx, schema, basetable, false, id)
	res.AddResult(sqlResult)
	if err != nil {
		return res, err
	}
	tbl := StructToTable(s[0], tablename, schema, true, id, nil, nil, nil)
	// temp tables
	sqlResult, err = s[0].PGCreateTable(ctx, tx, "", basetable, true, id)
	res.AddResult(sqlResult)
	if err != nil {
		return res, err
	}
	insertTpl := `
	INSERT INTO 
	{{- if .Temporary }}  _tmp_{{- .UUID -}}_{{- .Name -}} 
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
	`

	fmt.Printf("Staff[0:%d]\tRows: %d\tFields: %d\tTotal Fields: %d\n", len(s), len(s), len(tbl.Fields), len(tbl.Fields)*len(s))
	sqlStatement := PGTableStatement(tbl, insertTpl, nil)
	// fmt.Sprintln(sqlStatement)
	sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, s)
	// fmt.Println("finished inserting to temp table?")
	res.AddResult(sqlResult)
	if err != nil {
		return res, err
	}

	// update table
	tbl.Temporary = false
	updateTpl := `
	with cte_most_recent as (
		select distinct on (s.staffkey)
		s._id_hash,
		s.staffkey
		from {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }} s
		order by s.staffkey, s._extract_date_time desc nulls last
	), cte_new as (
		select distinct * from _tmp_{{- .UUID -}}_{{- .Name }}
	), cte_updates as (
		select
		cn.*
		from cte_new cn
		where not exists (
			select 1
			from cte_most_recent cmr 
			where cmr._id_hash = cn._id_hash
			--and   cmr.staffkey     = cn.staffkey
		)
	) insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (
		select cu.* from cte_updates cu
	)
	`
	sqlStatement = PGTableStatement(tbl, updateTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res.AddResult(sqlResult)
	if err != nil {
		return res, err
	}
	// TODO: Rewrite childUpdateTpl to use _id_hash
	// childUpdateTpl := `
	// with cte_parent_ref as (
	// 	select distinct
	// 	p.staffkey,
	// 	p._id_hash,
	// 	p._extract_date_time
	// 	from {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Parent }} p
	// ), cte_current as (
	// 	select distinct
	// 	c.*
	// 	from {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }} c
	// 	order by c.staffkey, c._extract_date_time desc nulls last
	// ), cte_new as (
	// 	select distinct * from _tmp_{{- .UUID -}}_{{- .Name }}
	// ), cte_updates as (
	// 	select
	// 	cn.*
	// 	from cte_new cn
	// 	-- inner join on parent refs
	// 	inner join cte_parent_ref cpr on (
	// 		    cpr.staffkey = cn.staffkey
	// 		and cpr._extract_date_time = cn._extract_date_time
	// 	)
	// 	-- exclude on duplicates
	// 	where not exists (
	// 		select 1
	// 		from cte_current cc where
	// 		{{ $fields := pgincludefields .Fields -}}
	// 		{{- range  $index, $field := $fields -}}
	// 		{{ if ne $index 0 -}} and {{- end }} cc.{{ pgname $field }} = cn.{{ pgname $field }}
	// 		{{ end -}}
	// 	)
	// ) insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (
	// 	select cu.* from cte_updates cu
	// 	where not exists (
	// 		select 1 from
	// 		{{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} dst where
	// 		{{ $fields := pgincludefields .Fields -}}
	// 		{{- range  $index, $field := $fields -}}
	// 		{{ if ne $index 0 }}and {{ end -}} dst.{{ pgname $field }} = cu.{{ pgname $field }}
	// 		{{ end -}}
	// 	)
	// )
	// `

	childUpdateTpl := `
	with cte_new_data as (
		select distinct tmp.* from  _tmp_{{- .UUID -}}_{{- .Name }} tmp
		inner join {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Parent }} parent on (
			parent.staffkey 		  = tmp.staffkey
		and	parent._id_hash   		  = tmp._id_hash
		and parent._extract_date_time = tmp._extract_date_time
		)
	) insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (
		select distinct cnd.* from cte_new_data cnd
		where not exists (
			select 1 from 
			{{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} dst where
				dst._id_hash 		   = cnd._id_hash
			and dst._extract_date_time = cnd._extract_date_time
		)
	)
	`
	// tag
	stafftags := []XStaffTag{}
	for _, staff := range s {
		for _, tags := range staff.Tags {
			stafftags = append(stafftags, tags.Tags...)
		}
	}
	tablename = basetable + "tag"
	tbl = StructToTable(XStaffTag{}, tablename, schema, true, id, nil, nil, nil)
	sqlStatement = PGTableStatement(tbl, insertTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, stafftags)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}
	sqlStatement = PGTableStatement(tbl, childUpdateTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res.AddResult(sqlResult)
	if err != nil {
		return res, err
	}

	// ttcmtag
	staffttcmtags := []XStaffTag{}
	for _, staff := range s {
		for _, tags := range staff.TTCMTags {
			staffttcmtags = append(staffttcmtags, tags.Tags...)
		}
	}
	tablename = basetable + "ttcmtag"
	tbl = StructToTable(XStaffTag{}, tablename, schema, true, id, nil, nil, nil)
	sqlStatement = PGTableStatement(tbl, insertTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, staffttcmtags)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}
	sqlStatement = PGTableStatement(tbl, childUpdateTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res.AddResult(sqlResult)
	if err != nil {
		return res, err
	}

	// skillset
	staffskills := []StaffSkill{}
	for _, staff := range s {
		staffskills = append(staffskills, staff.Skillset...)
	}
	tablename = basetable + "skill"
	tbl = StructToTable(StaffSkill{}, tablename, schema, true, id, nil, nil, nil)
	sqlStatement = PGTableStatement(tbl, insertTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, staffskills)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}
	sqlStatement = PGTableStatement(tbl, childUpdateTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res.AddResult(sqlResult)
	if err != nil {
		return res, err
	}

	// profile
	staffprofiles := []XStaffProfile{}
	for _, staff := range s {
		staffprofiles = append(staffprofiles, staff.Profiles...)
	}
	tablename = basetable + "profile"
	tbl = StructToTable(XStaffProfile{}, tablename, schema, true, id, nil, nil, nil)
	sqlStatement = PGTableStatement(tbl, insertTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, staffprofiles)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}
	sqlStatement = PGTableStatement(tbl, childUpdateTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res.AddResult(sqlResult)
	if err != nil {
		return res, err
	}

	return nil, nil
}

func (s *Staffs) EPL(ctx context.Context, c *Client, rc *RequestConfig,
	db *sqlx.DB, schema, table string, newRowsOnly bool) (sql.Result, error) {
	return nil, nil
}

func PGChunks[SS []S, S []T, T any](data S) SS {
	maxPGParams := 65535
	// var tp T
	tp := *(new(T))
	numFields := reflect.ValueOf(tp).Type().NumField()
	chunkSize := maxPGParams / numFields

	out := SS{}
	for i := 0; i < len(data); i += chunkSize {
		j := i + chunkSize
		if j > len(data) {
			j = len(data)
		}
		out = append(out, data[i:j])
	}

	return out
}

func PGInsertChunks[SS []S, S []T, T any](ctx context.Context, tx sqlx.Tx, tpl string, data S) error {
	chunkSlice := PGChunks(data)
	for _, v := range chunkSlice {
		tx.NamedExecContext(ctx, "", v)
	}
	return nil
}
