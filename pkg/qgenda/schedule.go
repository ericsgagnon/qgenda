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

type Schedule struct {
	// ------- metadata ------------------- //
	RawMessage       *string `json:"-" db:"_raw_message"`
	ProcessedMessage *string `json:"-" db:"_processed_message"` // RawMessage processed, omits 'message' metadata and 'noisy' fields (eg lastlogin)
	SourceQuery      *string `json:"_source_query,omitempty" db:"_source_query"`
	ExtractDateTime  *Time   `json:"_extract_date_time,omitempty" db:"_extract_date_time"`
	IDHash           *string `json:"_id_hash,omitempty" db:"_id_hash"` // hash of identifying fields: schedulekey-lastmodifieddateutc (rfc3339nano)
	// ------------------------------------ //
	ScheduleKey            *string        `json:"ScheduleKey,omitempty" primarykey:"true"`
	CallRole               *string        `json:"CallRole,omitempty"`
	CompKey                *string        `json:"CompKey,omitempty"`
	Credit                 *float64       `json:"Credit,omitempty"`
	Date                   *Date          `json:"Date,omitempty"`
	StartDateUTC           *Time          `json:"StartDateUTC,omitempty"`
	EndDateUTC             *Time          `json:"EndDateUTC,omitempty"`
	EndDate                *Date          `json:"EndDate,omitempty"`
	EndTime                *TimeOfDay     `json:"EndTime,omitempty"`
	IsCred                 *bool          `json:"IsCred,omitempty"`
	IsPublished            *bool          `json:"IsPublished,omitempty"`
	IsLocked               *bool          `json:"IsLocked,omitempty"`
	IsStruck               *bool          `json:"IsStruck,omitempty"`
	Notes                  *string        `json:"Notes,omitempty"`
	IsNotePrivate          *bool          `json:"IsNotePrivate,omitempty"`
	StaffAbbrev            *string        `json:"StaffAbbrev,omitempty"`
	StaffBillSysId         *string        `json:"StaffBillSysId,omitempty"`
	StaffEmail             *string        `json:"StaffEmail,omitempty"`
	StaffEmrId             *string        `json:"StaffEmrId,omitempty"`
	StaffErpId             *string        `json:"StaffErpId,omitempty"`
	StaffInternalId        *string        `json:"StaffInternalId,omitempty"`
	StaffExtCallSysId      *string        `json:"StaffExtCallSysId,omitempty"`
	StaffFName             *string        `json:"StaffFName,omitempty"`
	StaffId                *string        `json:"StaffId,omitempty"`
	StaffKey               *string        `json:"StaffKey,omitempty"`
	StaffLName             *string        `json:"StaffLName,omitempty"`
	StaffMobilePhone       *string        `json:"StaffMobilePhone,omitempty"`
	StaffNpi               *string        `json:"StaffNpi,omitempty"`
	StaffPager             *string        `json:"StaffPager,omitempty"`
	StaffPayrollId         *string        `json:"StaffPayrollId,omitempty"`
	StaffTags              []ScheduleTags `json:"StaffTags,omitempty"`
	StartDate              *Date          `json:"StartDate,omitempty"`
	StartTime              *TimeOfDay     `json:"StartTime,omitempty"`
	TaskAbbrev             *string        `json:"TaskAbbrev,omitempty"`
	TaskBillSysId          *string        `json:"TaskBillSysId,omitempty"`
	TaskContactInformation *string        `json:"TaskContactInformation,omitempty"`
	TaskExtCallSysId       *string        `json:"TaskExtCallSysId,omitempty"`
	TaskId                 *string        `json:"TaskId,omitempty"`
	TaskKey                *string        `json:"TaskKey,omitempty"`
	TaskName               *string        `json:"TaskName,omitempty"`
	TaskPayrollId          *string        `json:"TaskPayrollId,omitempty"`
	TaskEmrId              *string        `json:"TaskEmrId,omitempty"`
	TaskCallPriority       *string        `json:"TaskCallPriority,omitempty"`
	TaskDepartmentId       *string        `json:"TaskDepartmentId,omitempty"`
	TaskIsPrintEnd         *bool          `json:"TaskIsPrintEnd,omitempty"`
	TaskIsPrintStart       *bool          `json:"TaskIsPrintStart,omitempty"`
	TaskShiftKey           *string        `json:"TaskShiftKey,omitempty"`
	TaskType               *string        `json:"TaskType,omitempty"`
	TaskTags               []ScheduleTags `json:"TaskTags,omitempty"`
	LocationName           *string        `json:"LocationName,omitempty"`
	LocationAbbrev         *string        `json:"LocationAbbrev,omitempty"`
	LocationID             *string        `json:"LocationID,omitempty"`
	LocationAddress        *string        `json:"LocationAddress,omitempty"`
	TimeZone               *string        `json:"TimeZone,omitempty"`
	LastModifiedDateUTC    *Time          `json:"LastModifiedDateUTC,omitempty" primarykey:"true" querycondition:"ge" qf:"SinceModifiedTimestamp"`
	LocationTags           []ScheduleTags `json:"LocationTags,omitempty"`
	IsRotationTask         *bool          `json:"IsRotationTask"`
}

// UnmarshalJSON fulfils the json.Unmarshaler interface and
// assigns a compact json representation to .RawMessage
func (s *Schedule) UnmarshalJSON(b []byte) error {
	// alias technique to avoid infinite recursion
	type Alias Schedule
	var a Alias

	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}

	var bb bytes.Buffer
	if err := json.Compact(&bb, b); err != nil {
		return err
	}
	rawMessage := bb.String()

	dest := Schedule(a)
	dest.RawMessage = &rawMessage

	*s = dest
	return nil
}

// MarshalJSON satisfies the json.Marshaler interface
func (s Schedule) MarshalJSON() ([]byte, error) {
	type Alias Schedule
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
func (s *Schedule) Process() error {

	if err := ProcessStruct(s); err != nil {
		return fmt.Errorf("error processing %T:\t%q", s, err)
	}

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
	if err := s.SetIDHash(); err != nil {
		return err
	}
	return nil
}

// SetMessage uses a copy, strips message metadata, remarshals to JSON and compacts it, and assigns the string to .ProcessedMessage
func (s *Schedule) SetMessage() error {
	// take a copy and strip metadata, for good measure
	ss := *s
	ss.RawMessage = nil
	ss.ProcessedMessage = nil
	// ss.MessageHash = nil

	b, err := json.Marshal(ss)
	// b, err := json.Marshal(s)
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

// SetIDHash takes the hash of the json encoded fields that (should) uniquely identify this instance
// for schedule, this is schedulekey, lastmodifieddateutc (in rfc3339 with nano precision)
func (s *Schedule) SetIDHash() error {
	if s.ProcessedMessage == nil {
		return fmt.Errorf("ProcessedMessage is empty, cannot hash")
	}
	id := map[string]any{
		"ScheduleKey":         (*s.ScheduleKey),
		"LastModifiedDateUTC": (*s.LastModifiedDateUTC).Time.Format(time.RFC3339Nano),
	}

	b, err := json.Marshal(id)
	if err != nil {
		return err
	}
	var bb bytes.Buffer
	if err := json.Compact(&bb, b); err != nil {
		return err
	}

	// id := fmt.Sprint(*s.ScheduleKey, "-", (*s.LastModifiedDateUTC).Time.Format(time.RFC3339Nano))
	h := Hash(bb.String())
	s.IDHash = &h
	return nil
}

func DefaultScheduleRequestConfig() *RequestConfig {
	requestPath := "schedule"
	allowedFields := []string{
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
	rc := NewRequestConfig(requestPath, allowedFields)
	rc.SetStartDate(time.Now().UTC().Add(time.Hour * 24 * -15).Truncate(time.Hour * 24))
	rc.SetEndDate(time.Now().UTC())
	rc.SetEndDate(rc.GetStartDate().Add(time.Hour * 24 * 14))
	rc.SetSinceModifiedTimestamp(rc.GetStartDate())
	rc.SetIncludes("StaffTags,TaskTags,LocationTags")
	return rc
}

func NewScheduleRequestConfig(rc *RequestConfig) *RequestConfig {
	return DefaultScheduleRequestConfig().Merge(rc)
}

func NewScheduleRequest(rc *RequestConfig) *Request {
	rc = NewScheduleRequestConfig(rc)
	return NewRequest(rc)
}

func (s Schedule) PGCreateTable(ctx context.Context, tx *sqlx.Tx, schema, tablename string, temporary bool, id string) (sql.Result, error) {

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

	table := StructToTable(Schedule{}, tablename, schema, temporary, id, nil, nil, "")
	sqlStatement := PGTableStatement(table, tpl, nil)
	// fmt.Println(sqlStatement)

	sqlResult, err := tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprint(basetable, "stafftag")
	table = StructToTable(ScheduleTag{}, tablename, schema, temporary, id, nil, nil, basetable)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, tpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprint(basetable, "tasktag")
	table = StructToTable(ScheduleTag{}, tablename, schema, temporary, id, nil, nil, basetable)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, tpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprint(basetable, "locationtag")
	table = StructToTable(ScheduleTag{}, tablename, schema, temporary, id, nil, nil, basetable)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, tpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s Schedule) PGQueryConstraints(ctx context.Context, db *sqlx.DB, schema, table string) (*RequestConfig, error) {
	rc := RequestConfig{}
	tbl := StructToTable(Schedule{}, table, schema, false, "", nil, nil, "")
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
	// query := PGStatement(*new(Schedule), schema, table, tpl)
	if err := db.GetContext(ctx, &rc, query); err != nil {
		return nil, err
	}
	return &rc, nil
}

// PGGetCDC returns a single row from the destination that (should) only have relevant CDC data
// for conversion to a RequestConfig using ToRequestConfig
func (s *Schedule) PGGetCDC(ctx context.Context, db *sqlx.DB, schema, table string) (*Schedule, error) {
	var ss *Schedule
	if schema != "" {
		table = schema + "." + table
	}
	query := fmt.Sprintf("SELECT MAX ( lastmodifieddateutc ) FROM %s ", table)
	if err := db.GetContext(ctx, ss, query); err != nil {
		return nil, err
	}
	return ss, nil

}

func (s Schedule) ToRequestConfig() (*RequestConfig, error) {
	// if s.LastModifiedDateUTC == nil {
	// 	return nil, fmt.Errorf("*LastModifiedDateUTC is nil")
	// }
	// rc := RequestConfig{}

	rc := DefaultScheduleRequestConfig()
	if s.LastModifiedDateUTC != nil {
		rc.SetSinceModifiedTimestamp(s.LastModifiedDateUTC.Time)
	}

	return rc, nil
}
