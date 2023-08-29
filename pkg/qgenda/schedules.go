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

	"github.com/exiledavatar/gotoolkit/meta"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Schedules []Schedule

func GetSchedules(ctx context.Context, c *Client, rc *RequestConfig) (Schedules, error) {
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
		resp, err := c.Do(ctx, subreq)
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(data))
		// ss := []Schedule{}
		ss := Schedules{}
		if err := json.Unmarshal(data, &ss); err != nil {
			return nil, err
		}
		sourceQuery := resp.Request.URL.String()
		extractDateTime, err := ParseTime(resp.Header.Get("date"))
		if err != nil {
			return nil, err
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
	return schedules, nil
}

func GetSchedulesFromFiles(filenames ...string) (Schedules, error) {
	schedules := Schedules{}
	for _, filename := range filenames {
		fi, err := os.Stat(filename)
		if err != nil {
			return nil, err
		}
		modTime := fi.ModTime()

		b, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, &schedules); err != nil {
			return nil, err
		}
		for i, v := range schedules {
			if v.ExtractDateTime == nil {
				proxyExtractDateTime := NewTime(modTime)
				v.ExtractDateTime = &proxyExtractDateTime
			}
			schedules[i] = v
		}
	}
	return schedules, nil
}

func (s *Schedules) Get(ctx context.Context, c *Client, rc *RequestConfig) error {
	switch schedules, err := GetSchedules(ctx, c, rc); {
	case err != nil:
		return err
	default:
		*s = schedules
		return nil
	}
}

func (ss *Schedules) Process() error {
	sss := *ss
	for i, _ := range sss {
		if err := sss[i].Process(); err != nil {
			return err
		}
	}
	sort.SliceStable(sss, func(i, j int) bool {
		ikey := *(sss[i].ScheduleKey)
		jkey := *(sss[j].ScheduleKey)
		itime := *(sss[i].LastModifiedDateUTC)
		jtime := *(sss[j].LastModifiedDateUTC)
		switch {
		case ikey < jkey:
			return true
		case ikey == jkey && itime.Time.Before(jtime.Time):
			return true
		default:
			return false
		}
	})
	*ss = sss
	return nil
}

// // LoadFile is used to import any cached files
// func (s *Schedules) LoadFile(filename string) error {
// 	fi, err := os.Stat(filename)
// 	if err != nil {
// 		return err
// 	}
// 	modTime := fi.ModTime()

// 	b, err := os.ReadFile(filename)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	ss := []Schedule{}
// 	if err := json.Unmarshal(b, &ss); err != nil {
// 		log.Println(err)
// 	}
// 	for i, v := range ss {
// 		if v.ExtractDateTime == nil {
// 			proxyExtractDateTime := NewTime(modTime)
// 			v.ExtractDateTime = &proxyExtractDateTime
// 		}
// 		ss[i] = v
// 	}
// 	*s = ss
// 	return nil
// }

func (s Schedules) CreatePGTable(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	return Schedule{}.CreatePGTable(ctx, db, schema, table)
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

func (s Schedules) PutPG(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	if len(s) < 1 {
		return nil, fmt.Errorf("%T.PGInsertRows: length of %T < 1, nothing to do", s, s)
	}

	if schema == "" {
		schema = "qgenda"
	}
	if table == "" {
		table = "schedule"
	}

	var locationtags, stafftags, tasktags ScheduleTags
	for _, schedule := range s {
		locationtags = append(locationtags, schedule.LocationTags...)
		stafftags = append(stafftags, schedule.StaffTags...)
		tasktags = append(tasktags, schedule.TaskTags...)
	}

	// use meta.Struct to help with query templating
	str, err := meta.NewStruct(Schedule{}, meta.Structconfig{
		NameSpace: []string{schema},
		Name:      table,
	})
	if err != nil {
		return nil, err
	}

	data := map[string]meta.StructWithData{
		"schedule":    {Struct: str, Data: meta.ToData(s)},
		"locationtag": {Struct: str.Fields().ByName("LocationTags").ToStruct(), Data: meta.ToData(locationtags)},
		"stafftag":    {Struct: str.Fields().ByName("StaffTags").ToStruct(), Data: meta.ToData(stafftags)},
		"tasktag":     {Struct: str.Fields().ByName("TaskTags").ToStruct(), Data: meta.ToData(tasktags)},
	}
	// create target tables: schedule, schedulestafftag, scheduletasktag, schedulelocationtag
	result, err := Schedule{}.CreatePGTable(ctx, db, schema, table)
	if err != nil {
		return result, err
	}
	str.Value.Value.Interface()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// temporary tables
	tpl := `create temporary table _tmp_{{- .Struct.TagName "table" }} (like {{ .Struct.TagIdentifier "table" }})`
	for _, stri := range data {
		query, err := stri.Struct.ExecuteTemplate(tpl, nil, nil)
		if err != nil {
			return nil, err
		}
		if result, err := tx.ExecContext(ctx, query); err != nil {
			return result, err
		}
	}

	// batch insert
	// postgres has a 65535 'parameter' limit, there is an 'unnest' work around, but for now we're just going to chunk it
	tpl = `{{- "\n" -}}
	insert into _tmp_{{- .Struct.TagName "table" }} ( 
		{{- $fields := .Struct.Fields.WithTagTrue "db" -}}
		{{ $fields.TagNames "db" | join ", " }} )
	values (
		{{- $fields := .Struct.Fields.WithTagTrue "db" -}}
		:
		{{- $fields.TagNames "db" | join ", :" -}}
		{{- "\n)" }}
	`

	chunkSize := 65535 / len(str.Fields())
	for k, v := range data {
		query, err := v.Struct.ExecuteTemplate(tpl, nil, nil)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(v.Data); i += chunkSize {
			j := i + chunkSize
			if j > len(v.Data) {
				j = len(v.Data)
			}
			values := v.Data[i:j]
			result, err := tx.NamedExecContext(ctx, query, values)
			if err != nil {
				return result, err
			}
			log.Printf("%-25s[%10d:%10d] Rows: %10d Fields: %10d Total Parameters: %10d\n", k, i, j, len(values), len(v.Struct.Fields()), (len(values) * len(v.Fields())))
		}

	}
	// update from temp to permanent tables
	updateTpl := `
		insert into {{ .Struct.TagIdentifier "table" }} (
			select distinct tmp.*
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
	for k, v := range data {
		fmt.Println("updating ", k)
		query, err := v.Struct.ExecuteTemplate(updateTpl, nil, nil)
		if err != nil {
			return nil, err
		}
		fmt.Println(query)
		result, err := tx.ExecContext(ctx, query)
		if err != nil {
			return nil, err
		}
		log.Println(result)
	}
	fmt.Println(db.DriverName())
	return nil, tx.Commit()
}

func (s Schedules) GetPGStatus(ctx context.Context, db *sqlx.DB, schema, table string) (*RequestConfig, error) {
	return Schedule{}.GetPGStatus(ctx, db, schema, table)
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

func (s *Schedules) GetProcessPut(ctx context.Context, c *Client, rc *RequestConfig, db *sqlx.DB,
	schema, table string, newRowsOnly bool) (sql.Result, error) {

	rc = NewScheduleRequestConfig(rc)

	if result, err := (Schedule{}.CreatePGTable(ctx, db, schema, table)); err != nil {
		return result, err
	}
	// rc, err := Schedule{}.GetPGStatus(ctx, db, schema, table)
	return nil, nil
}

func GetProcessPutPGSchedules(ctx context.Context, c *Client, rc *RequestConfig,
	db *sqlx.DB, schema, table string, newRowsOnly bool, filenames ...string) (sql.Result, error) {

	schedules := Schedules{}
	var err error
	rc = NewScheduleRequestConfig(rc)
	if result, err := (Schedule{}.CreatePGTable(ctx, db, schema, table)); err != nil {
		return result, err
	}
	destRC, err := Schedule{}.GetPGStatus(ctx, db, schema, table)
	if err != nil {
		return nil, err
	}
	if destRC.SinceModifiedTimestamp != nil && newRowsOnly {
		rc.SetSinceModifiedTimestamp(destRC.GetSinceModifiedTimestamp())
	}
	// if len(filenames) > 0 {
	// 	schedules, err = GetSchedulesFromFiles(filenames...)
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// }

	schedules, err = GetSchedules(ctx, c, rc)
	if err != nil {
		return nil, err
	}

	// schedules = append(schedules, sch...)
	if err := schedules.Process(); err != nil {
		return nil, err
	}

	return schedules.PutPG(ctx, db, schema, table)

	// result, err := schedules.PutPG(ctx, db, schema, table)
	// if err != nil {
	// 	return nil, err
	// }

	// return result, nil
}
