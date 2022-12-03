package qgenda

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
)

type XSchedule struct {
	RawMessage             *string         `json:"-" db:"_raw_message"`
	ProcessedMessage       *string         `json:"-" db:"_processed_message"` // RawMessage processed, with changing fields dropped
	SourceQuery            *string         `json:"_source_query" db:"_source_query"`
	ExtractDateTime        *Time           `json:"_extract_date_time" db:"_extract_date_time"`
	HashID                 *string         `json:"_hash_id" db:"_hash_id"` // hash of processed message
	ScheduleKey            *string         `json:"ScheduleKey,omitempty" primarykey:"true"`
	CallRole               *string         `json:"CallRole,omitempty"`
	CompKey                *string         `json:"CompKey,omitempty"`
	Credit                 *float64        `json:"Credit,omitempty"`
	Date                   *Date           `json:"Date,omitempty"`
	StartDateUTC           *Time           `json:"StartDateUTC,omitempty"`
	EndDateUTC             *Time           `json:"EndDateUTC,omitempty"`
	EndDate                *Date           `json:"EndDate,omitempty"`
	EndTime                *TimeOfDay      `json:"EndTime,omitempty"`
	IsCred                 *bool           `json:"IsCred,omitempty"`
	IsPublished            *bool           `json:"IsPublished,omitempty"`
	IsLocked               *bool           `json:"IsLocked,omitempty"`
	IsStruck               *bool           `json:"IsStruck,omitempty"`
	Notes                  *string         `json:"Notes,omitempty"`
	IsNotePrivate          *bool           `json:"IsNotePrivate,omitempty"`
	StaffAbbrev            *string         `json:"StaffAbbrev,omitempty"`
	StaffBillSysId         *string         `json:"StaffBillSysId,omitempty"`
	StaffEmail             *string         `json:"StaffEmail,omitempty"`
	StaffEmrId             *string         `json:"StaffEmrId,omitempty"`
	StaffErpId             *string         `json:"StaffErpId,omitempty"`
	StaffInternalId        *string         `json:"StaffInternalId,omitempty"`
	StaffExtCallSysId      *string         `json:"StaffExtCallSysId,omitempty"`
	StaffFName             *string         `json:"StaffFName,omitempty"`
	StaffId                *string         `json:"StaffId,omitempty"`
	StaffKey               *string         `json:"StaffKey,omitempty"`
	StaffLName             *string         `json:"StaffLName,omitempty"`
	StaffMobilePhone       *string         `json:"StaffMobilePhone,omitempty"`
	StaffNpi               *string         `json:"StaffNpi,omitempty"`
	StaffPager             *string         `json:"StaffPager,omitempty"`
	StaffPayrollId         *string         `json:"StaffPayrollId,omitempty"`
	StaffTags              []XScheduleTags `json:"StaffTags,omitempty"`
	StartDate              *Date           `json:"StartDate,omitempty"`
	StartTime              *TimeOfDay      `json:"StartTime,omitempty"`
	TaskAbbrev             *string         `json:"TaskAbbrev,omitempty"`
	TaskBillSysId          *string         `json:"TaskBillSysId,omitempty"`
	TaskContactInformation *string         `json:"TaskContactInformation,omitempty"`
	TaskExtCallSysId       *string         `json:"TaskExtCallSysId,omitempty"`
	TaskId                 *string         `json:"TaskId,omitempty"`
	TaskKey                *string         `json:"TaskKey,omitempty"`
	TaskName               *string         `json:"TaskName,omitempty"`
	TaskPayrollId          *string         `json:"TaskPayrollId,omitempty"`
	TaskEmrId              *string         `json:"TaskEmrId,omitempty"`
	TaskCallPriority       *string         `json:"TaskCallPriority,omitempty"`
	TaskDepartmentId       *string         `json:"TaskDepartmentId,omitempty"`
	TaskIsPrintEnd         *bool           `json:"TaskIsPrintEnd,omitempty"`
	TaskIsPrintStart       *bool           `json:"TaskIsPrintStart,omitempty"`
	TaskShiftKey           *string         `json:"TaskShiftKey,omitempty"`
	TaskType               *string         `json:"TaskType,omitempty"`
	TaskTags               []XScheduleTags `json:"TaskTags,omitempty"`
	LocationName           *string         `json:"LocationName,omitempty"`
	LocationAbbrev         *string         `json:"LocationAbbrev,omitempty"`
	LocationID             *string         `json:"LocationID,omitempty"`
	LocationAddress        *string         `json:"LocationAddress,omitempty"`
	TimeZone               *string         `json:"TimeZone,omitempty"`
	LastModifiedDateUTC    *Time           `json:"LastModifiedDateUTC,omitempty" primarykey:"true" querycondition:"ge" qf:"SinceModifiedTimestamp"`
	LocationTags           []XScheduleTags `json:"LocationTags,omitempty"`
	IsRotationTask         *bool           `json:"IsRotationTask"`
}

// UnmarshalJSON fulfils the json.Unmarshaler interface and
// assigns a compact json representation to .RawMessage
func (s *XSchedule) UnmarshalJSON(b []byte) error {
	// alias technique to avoid infinite recursion
	type Alias XSchedule
	var a Alias

	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}

	var aa Alias
	if err := json.Unmarshal(b, &aa); err != nil {
		return err
	}
	aa.RawMessage = nil
	aa.ProcessedMessage = nil
	aa.ExtractDateTime = nil
	aa.SourceQuery = nil
	aa.HashID = nil
	aab, err := json.Marshal(aa)
	if err != nil {
		return err
	}

	var bb bytes.Buffer
	if err := json.Compact(&bb, aab); err != nil {
		return err
	}
	rawMessage := bb.String()

	dest := XSchedule(a)
	dest.RawMessage = &rawMessage

	*s = dest
	return nil
}

// MarshalJSON satisfies the json.Marshaler interface
func (s XSchedule) MarshalJSON() ([]byte, error) {
	type Alias XSchedule
	a := Alias(s)

	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	var bb bytes.Buffer
	if err := json.Compact(&bb, b); err != nil {
		return nil, err
	}

	return bb.Bytes(), nil
}

// Process handles all the basic validating and processing of
// from the raw version of any values. It is idempotent.
func (s *XSchedule) Process() error {

	if err := ProcessStruct(s); err != nil {
		return fmt.Errorf("error processing %T:\t%q", s, err)
	}

	// process stafftags
	// if err := setScheduleTagsMetaData(s, s.StaffTags); err != nil {
	// 	return err
	// }

	// if err := sortScheduleTagsSlice(s.StaffTags); err != nil {
	// 	return err
	// }

	if len(s.StaffTags) > 0 {
		for i, _ := range s.StaffTags {
			s.StaffTags[i].ExtractDateTime = s.ExtractDateTime
			s.StaffTags[i].ScheduleKey = s.ScheduleKey
			s.StaffTags[i].LastModifiedDateUTC = s.LastModifiedDateUTC
			if err := s.StaffTags[i].Process(); err != nil {
				return err
			}
		}
		// sortScheduleTagsSlice(s.StaffTags)
		sort.SliceStable(s.StaffTags, func(i, j int) bool {
			return *(s.StaffTags[i].CategoryKey) < *(s.StaffTags[j].CategoryKey)
		})

	}

	// process TaskTags
	if len(s.TaskTags) > 0 {
		for i, _ := range s.TaskTags {
			s.TaskTags[i].ExtractDateTime = s.ExtractDateTime
			s.TaskTags[i].ScheduleKey = s.ScheduleKey
			s.TaskTags[i].LastModifiedDateUTC = s.LastModifiedDateUTC
			if err := s.TaskTags[i].Process(); err != nil {
				return err
			}
		}
		// sortScheduleTagsSlice(s.TaskTags)
		sort.SliceStable(s.TaskTags, func(i, j int) bool {
			return *(s.TaskTags[i].CategoryKey) < *(s.TaskTags[j].CategoryKey)
		})

	}

	// process LocationTags
	if len(s.LocationTags) > 0 {
		for i, _ := range s.LocationTags {
			s.LocationTags[i].ExtractDateTime = s.ExtractDateTime
			s.LocationTags[i].ScheduleKey = s.ScheduleKey
			s.LocationTags[i].LastModifiedDateUTC = s.LastModifiedDateUTC
			if err := s.LocationTags[i].Process(); err != nil {
				return err
			}
		}
		// sortScheduleTagsSlice(s.LocationTags)
		sort.SliceStable(s.LocationTags, func(i, j int) bool {
			return *(s.LocationTags[i].CategoryKey) < *(s.LocationTags[j].CategoryKey)
		})

	}
	if err := s.SetMessage(); err != nil {
		return err
	}
	if err := s.SetHashID(); err != nil {
		return err
	}
	return nil
}

// SetMessage uses a copy, strips metadata, remarshals to JSON and compacts it, and assigns the string to .ProcessedMessage
func (s *XSchedule) SetMessage() error {
	// take a copy and strip metadata, for good measure
	ss := *s
	ss.RawMessage = nil
	ss.ExtractDateTime = nil
	ss.ProcessedMessage = nil
	ss.SourceQuery = nil
	ss.HashID = nil

	b, err := json.Marshal(ss)
	if err != nil {
		return err
	}
	var bb bytes.Buffer
	if err := json.Compact(&bb, b); err != nil {
		return err
	}
	processedMessage := bb.String()
	s.ProcessedMessage = &processedMessage
	return nil
}

func (s *XSchedule) SetHashID() error {
	if s.ProcessedMessage == nil {
		return fmt.Errorf("ProcessedMessage is empty, cannot hash")
	}
	h := Hash(*s.ProcessedMessage)
	s.HashID = &h
	return nil
}

func DefaultXScheduleRequestQueryFields(rqf *RequestQueryFields) *RequestQueryFields {
	if rqf == nil {
		rqf = &RequestQueryFields{}
	}
	if rqf.StartDate == nil {
		rqf.SetStartDate(time.Now().UTC().Add(time.Hour * 24 * -15).Truncate(time.Hour * 24))
	}
	if rqf.EndDate == nil {
		rqf.SetEndDate(time.Now().UTC())
	}
	if rqf.EndDate.Sub(rqf.GetStartDate()) < time.Hour*0 {
		rqf.SetEndDate(rqf.GetStartDate().Add(time.Hour * 24 * 14))
	}
	if rqf.SinceModifiedTimestamp == nil {
		rqf.SetSinceModifiedTimestamp(rqf.GetStartDate())
	}
	if rqf.Includes == nil {
		rqf.SetIncludes("StaffTags,TaskTags,LocationTags")
	}
	return rqf
}

func NewXScheduleRequest(rqf *RequestQueryFields) *Request {
	requestPath := "schedule"
	queryFields := []string{
		"CompanyKey",
		"StartDate",
		"EndDate",
		"IncludeDeletes",
		"SinceModifiedTimestamp",
		"DateFormat",
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	rqf = DefaultScheduleRequestQueryFields(rqf)

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}

// func (s XSchedule) PGQuery(temporary bool, schema, table string) string {
// 	tbl := StructToTable(XSchedule{}, table, schema, temporary, "", nil, nil, "")
// 	tpl := `
// 	CREATE {{- if .Temporary }} TEMPORARY TABLE IF NOT EXISTS _tmp_{{- .UUID -}}_{{- .Name -}}
// 	{{ else }} TABLE IF NOT EXISTS {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }}
// 	{{- end }} (
// 	{{- $fields := pgincludefields .Fields -}}
// 	{{- range  $index, $field := $fields -}}
// 	{{- if ne $index 0 -}},{{- end }}
// 		{{ pgname $field }} {{ pgtype $field.Type }} {{ if $field.Unique }} unique {{ end -}} {{- if not $field.Nullable -}} not null {{- end }}
// 	{{- end -}}
// 	{{- if not .Temporary }}{{- if .PrimaryKey }},
// 		PRIMARY KEY ({{ .PrimaryKey }}){{- end -}}{{- end }}
// 	)
// 	`

// 	return PGTableStatement(tbl, tpl, nil)
// 	// return PGStatement(*new(XSchedule), schema, table, tpl)
// }

//	func (s XSchedule) ToTable(tablename, schema, id string, temporary bool, constraints map[string]string, tags map[string][]string, parent string) Table {
//		return StructToTable(XSchedule{}, tablename, schema, temporary, id, nil, nil, "")
//	}
func (s XSchedule) PGCreateTable(ctx context.Context, tx *sqlx.Tx, schema, tablename string, temporary bool, id string) (sql.Result, error) {

	basetable := "schedule"
	if tablename != "" {
		basetable = tablename
	}
	var res Result

	if !temporary && schema != "" {
		sqlResult, err := tx.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema))
		res = SQLResult(res, sqlResult)
		if err != nil {
			return res, err
		}
	}

	tablename = basetable
	tpl := `{{- $table := . -}}
	CREATE {{- if .Temporary }} TEMPORARY TABLE IF NOT EXISTS _tmp_{{- .UUID -}}_{{- .Name -}}
	{{ else }} TABLE IF NOT EXISTS {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }}
	{{- end }} (
		{{- $fields := pgincludefields .Fields -}}
	{{- range  $index, $field := $fields -}}
	{{- if ne $index 0 -}},{{- end }}
	{{ pgname $field }} {{ pgtype $field.Type }} {{ if not $table.Temporary }}{{ if $field.Unique }} unique {{ end -}} {{- if not $field.Nullable -}} not null {{- end }}{{- end }}
	{{- end -}}
	{{- if not .Temporary }}{{- if .PrimaryKey }}, 
		PRIMARY KEY ({{ .PrimaryKey }}){{- end -}}{{- end }}
	)	
	`

	table := StructToTable(XSchedule{}, tablename, schema, temporary, id, nil, nil, "")
	sqlStatement := PGTableStatement(table, tpl, nil)
	// fmt.Println(sqlStatement)

	sqlResult, err := tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprint(basetable, "stafftag")
	table = StructToTable(XScheduleTag{}, tablename, schema, temporary, id, nil, nil, basetable)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, tpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprint(basetable, "tasktag")
	table = StructToTable(XScheduleTag{}, tablename, schema, temporary, id, nil, nil, basetable)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, tpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprint(basetable, "locationtag")
	table = StructToTable(XScheduleTag{}, tablename, schema, temporary, id, nil, nil, basetable)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, tpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s XSchedule) PGQueryConstraints(ctx context.Context, db *sqlx.DB, schema, table string) (*RequestQueryFields, error) {
	rqf := RequestQueryFields{}
	tbl := StructToTable(XSchedule{}, table, schema, false, "", nil, nil, "")
	tpl := `
	SELECT
	{{- $fields := pgqueryfields .Fields -}}
	{{- range  $index, $field := $fields -}}
	{{- if ne $index 0 -}},{{- end }}
	MAX( {{ pgname $field }} ) AS {{ qfname $field }}
	{{- end }}
	FROM 
		{{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }}	
	`
	query := PGTableStatement(tbl, tpl, nil)
	// query := PGStatement(*new(XSchedule), schema, table, tpl)
	if err := db.GetContext(ctx, &rqf, query); err != nil {
		return nil, err
	}
	return &rqf, nil
}

// func (s XSchedules) PGInsertRows(ctx context.Context, tx *sqlx.Tx, schema, tablename string, temporary bool)

var XScehduleSQLTemplates = map[string]string{
	"postgres-createtable": `
CREATE {{- if .Temporary }} TEMPORARY TABLE IF NOT EXISTS _tmp_{{- .UUID -}}_{{- .Name -}}
{{ else }} TABLE IF NOT EXISTS {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }}
{{- end }} (
{{- $fields := pgincludefields .Fields -}}
{{- range  $index, $field := $fields -}}
{{- if ne $index 0 -}},{{- end }}
	{{ pgname $field }} {{ pgtype $field.Type }} {{ if $field.Unique }} unique {{ end -}} {{- if not $field.Nullable -}} not null {{- end }}
{{- end -}}
)`,
	"postgres-insertrows":       ``,
	"postgres-droptable":        ``,
	"postgres-queryconstraints": ``,
}

// extract
// process
// load
//   create table
//   create temp table
//   insert into temp table
//   insert into table from temp table
//
