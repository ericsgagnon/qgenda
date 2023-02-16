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

type Schedules []Schedule

func (s *Schedules) Get(ctx context.Context, c *Client, rc *RequestConfig) error {
	req := NewScheduleRequest(rc)

	// qgenda only supports 100 days of schedules per query
	// using 90 days in case there are other limits
	schedules := Schedules{}
	duration := time.Hour * 24 * 90
	for t := req.GetStartDate(); t.Before(req.GetEndDate()); t = t.Add(duration) {
		vcpreq := *req
		subreq := &(vcpreq)
		subreq.SetStartDate(t)
		endDate := subreq.GetEndDate()
		if endDate.After(t.Add(duration)) {
			endDate = t.Add(duration)
		}
		subreq.SetEndDate(endDate)
		// fmt.Println(subreq.ToHTTPRequest().URL.String())
		resp, err := c.Do(ctx, subreq)
		if err != nil {
			return err
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		ss := []Schedule{}
		// fmt.Printf("response header:\t%s\n", resp.Header)
		// fmt.Printf("response body:\t%s\n", data)
		if err := json.Unmarshal(data, &ss); err != nil {
			return err
		}
		sourceQuery := resp.Request.URL.String()
		// user response header for extract time
		// fmt.Println(resp.Header.Get("date"))
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
		log.Printf("schedules: %s - %s (modTime>= %s)\ttotal: %d\textractDateTime:%s\n", subreq.GetStartDate(), subreq.GetEndDate(), subreq.GetSinceModifiedTimestamp(), len(ss), extractDateTime)

		schedules = append(schedules, ss...)

	}
	*s = schedules
	return nil
}

func (ss *Schedules) Process() error {
	sss := *ss
	// ts := map[string]bool{}
	for i, _ := range sss {
		if err := sss[i].Process(); err != nil {
			return err
		}
		// ts[fmt.Sprint(*(sss[i].ExtractDateTime))] = true
	}
	sort.SliceStable(sss, func(i, j int) bool {
		return *(sss[i].ScheduleKey) < *(sss[j].ScheduleKey)
	})
	// fmt.Println(ts)
	// fmt.Println(*(sss[0].ExtractDateTime))
	*ss = sss
	return nil
}

// LoadFile is used to import any cached files
func (s *Schedules) LoadFile(filename string) error {
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}
	modTime := fi.ModTime()

	b, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	ss := []Schedule{}
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

func (s Schedules) PGInsertRows(ctx context.Context, tx *sqlx.Tx, schema, tablename, id string) (sql.Result, error) {

	if len(s) < 1 {
		return nil, fmt.Errorf("%T.PGInsertRows: length of %T < 1, nothing to do", s, s)
	}
	var res Result
	id = strings.ReplaceAll(uuid.NewString(), "-", "")[0:16]
	// fmt.Println("Length of these Schedules: ", len(s))
	if tablename == "" {
		tablename = "schedule"
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
	// postgres has a 65535 'parameter' limit, there is an 'unnest' work around, but for now we're just going to chunk it
	chunkSize := 65535 / len(tbl.Fields) // reflect.ValueOf(Schedule{}).NumField()
	// this shouldn't be an issue for StaffMember, only implement if you need to
	for i := 0; i < len(s); i += chunkSize {
		j := i + chunkSize
		if j > len(s) {
			j = len(s)
		}
		v := s[i:j]
		fmt.Printf("Schedules[%d:%d]\tRows: %d\tFields: %d\tTotal Fields: %d\n", i, j, (j - i), len(tbl.Fields), len(tbl.Fields)*(j-i))

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
	// updateTpl := `
	// with cte_most_recent as (
	// 	select distinct on (s.schedulekey)
	// 	s.schedulekey,
	// 	s.lastmodifieddateutc
	// 	from {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }} s
	// 	order by s.schedulekey, s.lastmodifieddateutc desc nulls last
	// ), cte_new as (
	// 	select distinct * from _tmp_{{- .UUID -}}_{{- .Name }}
	// ), cte_updates as (
	// 	select
	// 	cn.*
	// 	from cte_new cn
	// 	where not exists (
	// 		select 1
	// 		from cte_most_recent cmr
	// 		where cmr.lastmodifieddateutc = cn.lastmodifieddateutc
	// 		and   cmr.schedulekey     = cn.schedulekey
	// 	)
	// ) insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (
	// select cu.* from cte_updates cu
	// )
	// `

	updateTpl := `
	insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (
		select distinct tmp.* 
		from _tmp_{{- .UUID -}}_{{- .Name }} tmp
		where not exists (
			select 1
			from {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} dst
			where dst.lastmodifieddateutc	=	tmp.lastmodifieddateutc
			and   dst.schedulekey			=	tmp.schedulekey
		)
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
		inner join {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Parent.Name }} prnt on (
				prnt._extract_date_time = tmp._extract_date_time
			and prnt.schedulekey         = tmp.schedulekey
			and prnt.lastmodifieddateutc = tmp.lastmodifieddateutc
		) 
		-- exclude on duplicates
		where not exists (
			select 1
			from {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} dst 
			where
				dst._extract_date_time  = tmp._extract_date_time
			and dst.lastmodifieddateutc = tmp.lastmodifieddateutc
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

	// childUpdateTpl = `
	// with cte_updates as (
	// 	select
	// 	tmp.*
	// 	from _tmp_{{- .UUID -}}_{{- .Name }} tmp
	// 	-- inner join on parent refs
	// 	inner join {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Parent }} prnt on (
	// 			prnt._extract_date_time = tmp._extract_date_time
	// 		and prnt.schedulekey         = tmp.schedulekey
	// 		and prnt.lastmodifieddateutc = tmp.lastmodifieddateutc
	// 	)
	// ) insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (
	// 	select cu.* from cte_updates cu
	// )
	// `

	// stafftags
	// stafftags := []ScheduleTag{}
	stafftags := ScheduleTags{}
	for _, schedule := range s {
		stafftags = append(stafftags, schedule.StaffTags...)
		// for _, tags := range schedule.StaffTags {
		// 	stafftags = append(stafftags, tags)
		// }
	}

	tablename = basetable + "stafftag"
	stafftagTbl := StructToTable(ScheduleTag{}, tablename, schema, true, id, nil, nil, nil)
	stafftagTbl.Parent = &tbl
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
	// tasktags := []ScheduleTag{}
	tasktags := ScheduleTags{}
	for _, schedule := range s {
		tasktags = append(tasktags, schedule.TaskTags...)
		// for _, tags := range schedule.TaskTags {
		// 	tasktags = append(tasktags, tags)
		// }
	}
	tablename = basetable + "tasktag"
	tasktagTbl := StructToTable(ScheduleTag{}, tablename, schema, true, id, nil, nil, nil)
	tasktagTbl.Parent = &tbl
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
	// locationtags := []ScheduleTag{}
	locationtags := ScheduleTags{}
	for _, schedule := range s {
		locationtags = append(locationtags, schedule.LocationTags...)
		// for _, tags := range schedule.LocationTags {
		// 	locationtags = append(locationtags, tags)
		// }
	}
	tablename = basetable + "locationtag"
	locationtagTbl := StructToTable(ScheduleTag{}, tablename, schema, true, id, nil, nil, nil)
	locationtagTbl.Parent = &tbl
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

func (s *Schedules) EPL(ctx context.Context, c *Client, rc *RequestConfig,
	db *sqlx.DB, schema, table string, newRowsOnly bool) (sql.Result, error) {

	rc = NewScheduleRequestConfig(rc)

	var res Result

	tx := db.MustBeginTx(ctx, nil)
	sqlResult, err := Schedule{}.PGCreateTable(ctx, tx, schema, table, false, "")
	res.AddResult(sqlResult)
	if err != nil {
		return res, err
	}

	qrqf, err := Schedule{}.PGQueryConstraints(ctx, db, schema, table)
	if err != nil {
		return res, err
	}
	if qrqf.SinceModifiedTimestamp != nil && newRowsOnly {
		rc.SetSinceModifiedTimestamp(qrqf.GetSinceModifiedTimestamp())
	}
	if err := s.Get(ctx, c, rc); err != nil {
		return res, err
	}
	if err := s.Process(); err != nil {
		return res, err
	}
	sqlResult, err = s.PGInsertRows(ctx, tx, schema, table, "")
	res.AddResult(sqlResult)
	if err != nil {
		return res, err

	}
	err = tx.Commit()
	if err != nil {
		return res, err
	}
	return res, nil
}

func GetFromFile[P *[]T, T any](filename string, dst P) error {

	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	t := []T{}
	if err := json.Unmarshal(b, &t); err != nil {
		log.Println(err)
	}
	*dst = t
	return nil

}
