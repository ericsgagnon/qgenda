package qgenda

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type XSchedules []XSchedule

func (s *XSchedules) Get(ctx context.Context, c *Client, rqf *RequestQueryFields) error {
	req := NewScheduleRequest(rqf)
	// qgenda only supports 100 days of schedules per query
	// using 90 days in case there are other limits
	duration := time.Hour * 24 * 90
	for t := req.GetStartDate(); t.Before(req.GetEndDate()); t = t.Add(duration) {
		vcpreq := *req
		subreq := &(vcpreq)
		subreq.SetStartDate(t)
		subreq.SetEndDate(t.Add(duration))
		resp, err := c.Do(ctx, subreq)
		if err != nil {
			return err
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		ss := []XSchedule{}
		if err := json.Unmarshal(data, &ss); err != nil {
			return err
		}
		fmt.Println(resp.Request.URL.String())
	}
	return nil
}

func (ss *XSchedules) Process() error {
	sss := *ss
	for i, _ := range sss {
		if err := sss[i].Process(); err != nil {
			return err
		}
	}
	sort.SliceStable(sss, func(i, j int) bool {
		return *(sss[i].ScheduleKey) < *(sss[j].ScheduleKey)
	})
	*ss = sss
	return nil
}

// LoadFile is used to import any cached files
func (s *XSchedules) LoadFile(filename string) error {
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}
	modTime := fi.ModTime()

	b, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	ss := []XSchedule{}
	if err := json.Unmarshal(b, &ss); err != nil {
		log.Println(err)
	}
	for i, v := range ss {
		if v.ExtractDateTime == nil {
			proxyExtractDateTime := NewTime(modTime)
			v.ExtractDateTime = &proxyExtractDateTime
		}
		ss[i] = v
	}
	*s = ss
	return nil
}

func (s XSchedules) PGInsertRows(ctx context.Context, tx *sqlx.Tx, schema, tablename, id string) (sql.Result, error) {

	if len(s) < 1 {
		return nil, fmt.Errorf("%T.PGInsertRows: length of %T < 1, nothing to do", s)
	}
	var res Result
	id = strings.ReplaceAll(uuid.NewString(), "-", "")[0:16]
	// fmt.Println("Length of these XSchedules: ", len(s))
	if tablename == "" {
		tablename = "schedule"
	}
	basetable := tablename

	sqlResult, err := s[0].PGCreateTable(ctx, tx, schema, basetable, false, id)
	res.AddResult(sqlResult)
	if err != nil {
		return res, err
	}
	tbl := StructToTable(s[0], tablename, schema, true, id, nil, nil, "")
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
	// postgres has a 65535 'parameter' limit, there is an 'unnest' work around, but for now we're just going to chunk it
	chunkSize := 65535 / len(tbl.Fields) // reflect.ValueOf(XSchedule{}).NumField()
	// this shouldn't be an issue for StaffMember, only implement if you need to
	for i := 0; i < len(s); i += chunkSize {
		j := i + chunkSize
		if j > len(s) {
			j = len(s)
		}
		v := s[i:j]
		fmt.Printf("Schedules[%d:%d]\tRows: %d\tFields: %d\tTotal Fields: %d\n", i, j, (j - i), len(tbl.Fields), len(tbl.Fields)*(j-i))
		// sqlResult, err := s[i:j].PGInsertRows(ctx, tx, schema, tablename, id)
		// res = SQLResult(res, sqlResult)
		// if err != nil {
		// 	return res, err
		// }
		// id = tbl.UUID

		// fmt.Println("did I get this far?")

		// insert to temp tables

		sqlStatement := PGTableStatement(tbl, insertTpl, nil)
		// fmt.Sprintln(sqlStatement)
		sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, v)
		// fmt.Println("finished inserting to temp table?")
		res.AddResult(sqlResult)
		if err != nil {
			return res, err
		}
	}
	// update table
	tbl.Temporary = false
	updateTpl := `
	with cte_most_recent as (
		select distinct on (s.schedulekey)		
		s.schedulekey, 
		s.lastmodifieddateutc
		from {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }} s
		order by s.schedulekey, s.lastmodifieddateutc desc nulls last
	), cte_new as (
		select distinct * from _tmp_{{- .UUID -}}_{{- .Name }}
	), cte_updates as (
		select
		cn.*
		from cte_new cn
		where not exists (
			select 1
			from cte_most_recent cmr 
			where cmr.lastmodifieddateutc = cn.lastmodifieddateutc
			and   cmr.schedulekey     = cn.schedulekey
		)
	) insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (
	select cu.* from cte_updates cu
	)
	`
	sqlStatement := PGTableStatement(tbl, updateTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res.AddResult(sqlResult)
	if err != nil {
		return res, err
	}

	// child tables
	childUpdateTpl := `
	with cte_updates as (
		select distinct
		tmp.*
		from _tmp_{{- .UUID -}}_{{- .Name }} tmp
		-- inner join on parent refs
		inner join {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Parent }} prnt on (
			    prnt.schedulekey         = tmp.schedulekey
			and prnt.lastmodifieddateutc = tmp.lastmodifieddateutc
		) 
		-- exclude on duplicates
		where not exists (
			select 1
			from {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} dst 
			where
			    dst.lastmodifieddateutc = tmp.lastmodifieddateutc
			and dst.schedulekey 		= tmp.schedulekey
			and dst.categorykey 		= tmp.categorykey
			and dst.categoryname 		= tmp.categoryname
			and dst.tagkey 				= tmp.tagkey
			and dst.tagname 			= tmp.tagname
		)
	) insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (
		select cu.* from cte_updates cu
	)
	`

	// stafftags
	stafftags := []XScheduleTag{}
	for _, schedule := range s {
		for _, tags := range schedule.StaffTags {
			stafftags = append(stafftags, tags.Tags...)
		}
	}
	tablename = basetable + "stafftag"
	stafftagTbl := StructToTable(stafftags[0], tablename, schema, true, id, nil, nil, "")
	stafftagTbl.Parent = basetable
	chunkSize = 65535 / len(stafftagTbl.Fields)
	// sqlStatement = PGTableStatement(stafftagTbl, childUpdateTpl, nil)
	// fmt.Println(sqlStatement)

	if len(stafftags) > 0 {

		for i := 0; i < len(stafftags); i += chunkSize {
			j := i + chunkSize
			if j > len(stafftags) {
				j = len(stafftags)
			}
			v := stafftags[i:j]
			sqlStatement := PGTableStatement(stafftagTbl, insertTpl, nil)
			// fmt.Sprintln(sqlStatement)
			sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, v)
			// fmt.Println("finished inserting to temp table?")
			res.AddResult(sqlResult)
			if err != nil {
				return res, err
			}
		}
		sqlStatement := PGTableStatement(stafftagTbl, childUpdateTpl, nil)
		// fmt.Println(sqlStatement)
		sqlResult, err = tx.ExecContext(ctx, sqlStatement)
		res.AddResult(sqlResult)
		if err != nil {
			return res, err
		}

	}

	// tasktags
	tasktags := []XScheduleTag{}
	for _, schedule := range s {
		for _, tags := range schedule.TaskTags {
			tasktags = append(tasktags, tags.Tags...)
		}
	}
	tablename = basetable + "tasktag"
	tasktagTbl := StructToTable(tasktags[0], tablename, schema, true, id, nil, nil, "")
	tasktagTbl.Parent = basetable
	chunkSize = 65535 / len(tasktagTbl.Fields)
	// sqlStatement = PGTableStatement(tasktagTbl, childUpdateTpl, nil)
	// fmt.Println(sqlStatement)

	if len(tasktags) > 0 {

		for i := 0; i < len(tasktags); i += chunkSize {
			j := i + chunkSize
			if j > len(tasktags) {
				j = len(tasktags)
			}
			v := tasktags[i:j]
			sqlStatement := PGTableStatement(tasktagTbl, insertTpl, nil)
			// fmt.Sprintln(sqlStatement)
			sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, v)
			// fmt.Println("finished inserting to temp table?")
			res.AddResult(sqlResult)
			if err != nil {
				return res, err
			}
		}
		sqlStatement := PGTableStatement(tasktagTbl, childUpdateTpl, nil)
		// fmt.Println(sqlStatement)
		sqlResult, err = tx.ExecContext(ctx, sqlStatement)
		res.AddResult(sqlResult)
		if err != nil {
			return res, err
		}

	}

	// locationtags
	locationtags := []XScheduleTag{}
	for _, schedule := range s {
		for _, tags := range schedule.LocationTags {
			locationtags = append(locationtags, tags.Tags...)
		}
	}
	tablename = basetable + "locationtag"
	locationtagTbl := StructToTable(XScheduleTag{}, tablename, schema, true, id, nil, nil, "")
	locationtagTbl.Parent = basetable
	chunkSize = 65535 / len(locationtagTbl.Fields)
	// sqlStatement = PGTableStatement(locationtagTbl, childUpdateTpl, nil)
	// fmt.Println(sqlStatement)

	if len(locationtags) > 0 {

		for i := 0; i < len(locationtags); i += chunkSize {
			j := i + chunkSize
			if j > len(locationtags) {
				j = len(locationtags)
			}
			v := locationtags[i:j]
			sqlStatement := PGTableStatement(locationtagTbl, insertTpl, nil)
			// fmt.Sprintln(sqlStatement)
			sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, v)
			// fmt.Println("finished inserting to temp table?")
			res.AddResult(sqlResult)
			if err != nil {
				return res, err
			}
		}
		sqlStatement := PGTableStatement(locationtagTbl, childUpdateTpl, nil)
		// fmt.Println(sqlStatement)
		sqlResult, err = tx.ExecContext(ctx, sqlStatement)
		res.AddResult(sqlResult)
		if err != nil {
			return res, err
		}

	}

	// fmt.Println("finished inserting to permanent table?")
	return res, nil
}
