package qgenda

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"os"
	"sort"
	"time"

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
		// fmt.Println(string(data))
		// ss := []Schedule{}
		ss := Schedules{}
		// fmt.Println(string(data))
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

func (s Schedules) CreatePGTable(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	return Schedule{}.CreatePGTable(ctx, db, schema, table)
}

func (s Schedules) PutPG(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	return PutPG(ctx, db, s, schema, table)
}

// func (s Schedules) PutPG(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
// 	if len(s) < 1 {
// 		return nil, fmt.Errorf("%T.PGInsertRows: length of %T < 1, nothing to do", s, s)
// 	}

// 	if schema == "" {
// 		schema = "qgenda"
// 	}
// 	if table == "" {
// 		table = "schedule"
// 	}

// 	var locationtags, stafftags, tasktags ScheduleTags
// 	for _, schedule := range s {
// 		locationtags = append(locationtags, schedule.LocationTags...)
// 		stafftags = append(stafftags, schedule.StaffTags...)
// 		tasktags = append(tasktags, schedule.TaskTags...)
// 	}

// 	// use meta.Struct to help with query templating
// 	str, err := meta.NewStruct(Schedule{}, meta.Structconfig{
// 		NameSpace: []string{schema},
// 		Name:      table,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	data := map[string]meta.StructWithData{
// 		"schedule":    {Struct: str, Data: meta.ToData(s)},
// 		"locationtag": {Struct: str.Fields().ByName("LocationTags").ToStruct(), Data: meta.ToData(locationtags)},
// 		"stafftag":    {Struct: str.Fields().ByName("StaffTags").ToStruct(), Data: meta.ToData(stafftags)},
// 		"tasktag":     {Struct: str.Fields().ByName("TaskTags").ToStruct(), Data: meta.ToData(tasktags)},
// 	}
// 	// create target tables: schedule, schedulestafftag, scheduletasktag, schedulelocationtag
// 	result, err := Schedule{}.CreatePGTable(ctx, db, schema, table)
// 	if err != nil {
// 		return result, err
// 	}
// 	str.Value.Value.Interface()
// 	tx, err := db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// temporary tables
// 	tpl := `create temporary table _tmp_{{- .Struct.TagName "table" }} (like {{ .Struct.TagIdentifier "table" }})`
// 	for _, stri := range data {
// 		query, err := stri.Struct.ExecuteTemplate(tpl, nil, nil)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if result, err := tx.ExecContext(ctx, query); err != nil {
// 			return result, err
// 		}
// 	}

// 	// batch insert
// 	// postgres has a 65535 'parameter' limit, there is an 'unnest' work around, but for now we're just going to chunk it
// 	tpl = `{{- "\n" -}}
// 	insert into _tmp_{{- .Struct.TagName "table" }} (
// 		{{- $fields := .Struct.Fields.WithTagTrue "db" -}}
// 		{{ $fields.TagNames "db" | join ", " }} )
// 	values (
// 		{{- $fields := .Struct.Fields.WithTagTrue "db" -}}
// 		:
// 		{{- $fields.TagNames "db" | join ", :" -}}
// 		{{- "\n)" }}
// 	`

// 	chunkSize := 65535 / len(str.Fields())
// 	for k, v := range data {
// 		query, err := v.Struct.ExecuteTemplate(tpl, nil, nil)
// 		if err != nil {
// 			return nil, err
// 		}
// 		for i := 0; i < len(v.Data); i += chunkSize {
// 			j := i + chunkSize
// 			if j > len(v.Data) {
// 				j = len(v.Data)
// 			}
// 			values := v.Data[i:j]
// 			result, err := tx.NamedExecContext(ctx, query, values)
// 			if err != nil {
// 				return result, err
// 			}
// 			log.Printf("%-25s[%10d:%10d] Rows: %10d Fields: %10d Total Parameters: %10d\n", k, i, j, len(values), len(v.Struct.Fields()), (len(values) * len(v.Fields())))
// 		}

// 	}
// 	// update from temp to permanent tables
// 	updateTpl := `
// 		insert into {{ .Struct.TagIdentifier "table" }} (
// 			select distinct tmp.*
// 			from _tmp_{{- .Struct.TagName "table" }} tmp
// 			{{- if ne .Struct.Parent nil }}
// 			inner join _tmp_{{ .Struct.Parent.TagName "table" }} ptmp
// 			{{- $parentprimarykey := (index ( .Struct.Parent.Fields.WithTagTrue "primarykey" ) 0 ).TagName "db" -}}
// 			{{- $parentpkey := ( index ( .Struct.Fields.WithTagTrue "parentprimarykey" ) 0  ).TagName "db" }}
// 			on tmp.{{ $parentpkey }} = ptmp.{{ $parentprimarykey -}}
// 			{{ end }}
// 			where not exists (
// 				select 1
// 				from {{ .Struct.TagIdentifier "table" }} dst
// 				{{- $pkey := index ( ( .Struct.Fields.WithTagTrue "primarykey" ).TagNames "db" ) 0 }}
// 				where dst.{{ $pkey }} = tmp.{{ $pkey }}
//  			)
// 		)
// 		`

// 	// update permanent tables from temp tables
// 	for k, v := range data {
// 		fmt.Println("updating ", k)
// 		query, err := v.Struct.ExecuteTemplate(updateTpl, nil, nil)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fmt.Println(query)
// 		result, err := tx.ExecContext(ctx, query)
// 		if err != nil {
// 			return nil, err
// 		}
// 		log.Println(result)
// 	}
// 	fmt.Println(db.DriverName())
// 	return nil, tx.Commit()
// }

func (s Schedules) GetPGStatus(ctx context.Context, db *sqlx.DB, schema, table string) (*RequestConfig, error) {
	return Schedule{}.GetPGStatus(ctx, db, schema, table)
}

func (s *Schedules) GetProcessPut(ctx context.Context, c *Client, rc *RequestConfig, db *sqlx.DB,
	schema, table string, newRowsOnly bool) (sql.Result, error) {
	if result, err := (Schedule{}.CreatePGTable(ctx, db, schema, table)); err != nil {
		return result, err
	}
	return nil, nil
}

// func GetProcessPutPGSchedules(ctx context.Context, c *Client, rc *RequestConfig,
// 	db *sqlx.DB, schema, table string, newRowsOnly bool, filenames ...string) (sql.Result, error) {

// 	schedules := Schedules{}
// 	var err error
// 	rc = NewScheduleRequestConfig(rc)
// 	if result, err := CreatePGTable(ctx, db, schedules, schema, table); err != nil {
// 		return result, err
// 	}
// 	// if result, err := (Schedule{}.CreatePGTable(ctx, db, schema, table)); err != nil {
// 	// 	return result, err
// 	// }
// 	destRC, err := Schedule{}.GetPGStatus(ctx, db, schema, table)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if destRC.SinceModifiedTimestamp != nil && newRowsOnly {
// 		rc.SetSinceModifiedTimestamp(destRC.GetSinceModifiedTimestamp())
// 	}

// 	schedules, err = GetSchedules(ctx, c, rc)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := schedules.Process(); err != nil {
// 		return nil, err
// 	}

// 	// return schedules.PutPG(ctx, db, schema, table)
// 	return nil, nil
// }
