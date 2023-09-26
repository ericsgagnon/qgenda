package qgenda

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	// "github.com/ericsgagnon/qgenda/pkg/meta"
	"github.com/exiledavatar/gotoolkit/meta"
	"github.com/jmoiron/sqlx"
)

type Schedule struct {
	// ------- metadata ------------------- //
	RawMessage       *string `json:"-" db:"_raw_message" pgtype:"text"`
	ProcessedMessage *string `json:"-" db:"_processed_message" pgtype:"text"` // RawMessage processed, omits 'message' metadata and 'noisy' fields (eg lastlogin)
	SourceQuery      *string `json:"_source_query,omitempty" db:"_source_query" pgtype:"text"`
	ExtractDateTime  *Time   `json:"_extract_date_time,omitempty" db:"_extract_date_time" pgtype:"timestamp with time zone"`
	IDHash           *string `json:"_id_hash,omitempty" db:"_id_hash" pgtype:"text" primarykey:"true"` // hash of identifying fields: schedulekey-lastmodifieddateutc (rfc3339nano)
	// ------------------------------------ //
	ScheduleKey            *string      `json:"ScheduleKey,omitempty" db:"schedulekey" pgtype:"text" idhash:"true"`
	CallRole               *string      `json:"CallRole,omitempty" db:"callrole" pgtype:"text"`
	CompKey                *string      `json:"CompKey,omitempty" db:"compkey" pgtype:"text"`
	Credit                 *float64     `json:"Credit,omitempty" db:"credit" pgtype:"numeric"`
	Date                   *Date        `json:"Date,omitempty" db:"date" pgtype:"date"`
	StartDateUTC           *Time        `json:"StartDateUTC,omitempty" db:"startdateutc" pgtype:"timestamp with time zone"`
	EndDateUTC             *Time        `json:"EndDateUTC,omitempty" db:"enddateutc" pgtype:"timestamp with time zone"`
	EndDate                *Date        `json:"EndDate,omitempty" db:"enddate" pgtype:"date"`
	EndTime                *TimeOfDay   `json:"EndTime,omitempty" db:"endtime" pgtype:"time without time zone"`
	IsCred                 *bool        `json:"IsCred,omitempty" db:"iscred" pgtype:"boolean"`
	IsPublished            *bool        `json:"IsPublished,omitempty" db:"ispublished" pgtype:"boolean"`
	IsLocked               *bool        `json:"IsLocked,omitempty" db:"islocked" pgtype:"boolean"`
	IsStruck               *bool        `json:"IsStruck,omitempty" db:"isstruck" pgtype:"boolean"`
	Notes                  *string      `json:"Notes,omitempty" db:"notes" pgtype:"text"`
	IsNotePrivate          *bool        `json:"IsNotePrivate,omitempty" db:"isnoteprivate" pgtype:"boolean"`
	StaffAbbrev            *string      `json:"StaffAbbrev,omitempty" db:"staffabbrev" pgtype:"text"`
	StaffBillSysId         *string      `json:"StaffBillSysId,omitempty" db:"staffbillsysid" pgtype:"text"`
	StaffEmail             *string      `json:"StaffEmail,omitempty" db:"staffemail" pgtype:"text"`
	StaffEmrId             *string      `json:"StaffEmrId,omitempty" db:"staffemrid" pgtype:"text"`
	StaffErpId             *string      `json:"StaffErpId,omitempty" db:"stafferpid" pgtype:"text"`
	StaffInternalId        *string      `json:"StaffInternalId,omitempty" db:"staffinternalid" pgtype:"text"`
	StaffExtCallSysId      *string      `json:"StaffExtCallSysId,omitempty" db:"staffextcallsysid" pgtype:"text"`
	StaffFName             *string      `json:"StaffFName,omitempty" db:"stafffname" pgtype:"text"`
	StaffId                *string      `json:"StaffId,omitempty" db:"staffid" pgtype:"text"`
	StaffKey               *string      `json:"StaffKey,omitempty" db:"staffkey" pgtype:"text"`
	StaffLName             *string      `json:"StaffLName,omitempty" db:"stafflname" pgtype:"text"`
	StaffMobilePhone       *string      `json:"StaffMobilePhone,omitempty" db:"staffmobilephone" pgtype:"text"`
	StaffNpi               *string      `json:"StaffNpi,omitempty" db:"staffnpi" pgtype:"text"`
	StaffPager             *string      `json:"StaffPager,omitempty" db:"staffpager" pgtype:"text"`
	StaffPayrollId         *string      `json:"StaffPayrollId,omitempty" db:"staffpayrollid" pgtype:"text"`
	StaffTags              ScheduleTags `json:"StaffTags,omitempty" db:"-,stafftags" pgtype:"jsonb" table:"schedulestafftag"`
	StartDate              *Date        `json:"StartDate,omitempty" db:"startdate" pgtype:"date"`
	StartTime              *TimeOfDay   `json:"StartTime,omitempty" db:"starttime" pgtype:"time without time zone"`
	TaskAbbrev             *string      `json:"TaskAbbrev,omitempty" db:"taskabbrev" pgtype:"text"`
	TaskBillSysId          *string      `json:"TaskBillSysId,omitempty" db:"taskbillsysid" pgtype:"text"`
	TaskContactInformation *string      `json:"TaskContactInformation,omitempty" db:"taskcontactinformation" pgtype:"text"`
	TaskExtCallSysId       *string      `json:"TaskExtCallSysId,omitempty" db:"taskextcallsysid" pgtype:"text"`
	TaskId                 *string      `json:"TaskId,omitempty" db:"taskid" pgtype:"text"`
	TaskKey                *string      `json:"TaskKey,omitempty" db:"taskkey" pgtype:"text"`
	TaskName               *string      `json:"TaskName,omitempty" db:"taskname" pgtype:"text"`
	TaskPayrollId          *string      `json:"TaskPayrollId,omitempty" db:"taskpayrollid" pgtype:"text"`
	TaskEmrId              *string      `json:"TaskEmrId,omitempty" db:"taskemrid" pgtype:"text"`
	TaskCallPriority       *string      `json:"TaskCallPriority,omitempty" db:"taskcallpriority" pgtype:"text"`
	TaskDepartmentId       *string      `json:"TaskDepartmentId,omitempty" db:"taskdepartmentid" pgtype:"text"`
	TaskIsPrintEnd         *bool        `json:"TaskIsPrintEnd,omitempty" db:"taskisprintend" pgtype:"boolean"`
	TaskIsPrintStart       *bool        `json:"TaskIsPrintStart,omitempty" db:"taskisprintstart" pgtype:"boolean"`
	TaskShiftKey           *string      `json:"TaskShiftKey,omitempty" db:"taskshiftkey" pgtype:"text"`
	TaskType               *string      `json:"TaskType,omitempty" db:"tasktype" pgtype:"text"`
	TaskTags               ScheduleTags `json:"TaskTags,omitempty" db:"-,tasktags" pgtype:"jsonb" table:"scheduletasktag"`
	LocationName           *string      `json:"LocationName,omitempty" db:"locationname" pgtype:"text"`
	LocationAbbrev         *string      `json:"LocationAbbrev,omitempty" db:"locationabbrev" pgtype:"text"`
	LocationID             *string      `json:"LocationID,omitempty" db:"locationid" pgtype:"text"`
	LocationAddress        *string      `json:"LocationAddress,omitempty" db:"locationaddress" pgtype:"text"`
	TimeZone               *string      `json:"TimeZone,omitempty" db:"timezone" pgtype:"text"`
	LastModifiedDateUTC    *Time        `json:"LastModifiedDateUTC,omitempty" querycondition:"ge" qf:"SinceModifiedTimestamp" idhash:"true" db:"lastmodifieddateutc" pgtype:"timestamp with time zone" qgendarequestname:"SinceModifiedTimestamp"`
	LocationTags           ScheduleTags `json:"LocationTags,omitempty" db:"-,locationtags" pgtype:"text" table:"schedulelocationtag"`
	IsRotationTask         *bool        `json:"IsRotationTask" db:"isrotationtask" pgtype:"boolean"`
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

	if err := s.SetIDHash(); err != nil {
		return err
	}

	// fmt.Println("Length of StaffTags: ", len(s.StaffTags))
	if len(s.StaffTags) > 0 {
		for i, _ := range s.StaffTags {
			s.StaffTags[i].ExtractDateTime = s.ExtractDateTime
			s.StaffTags[i].ScheduleKey = s.ScheduleKey
			s.StaffTags[i].LastModifiedDateUTC = s.LastModifiedDateUTC
			s.StaffTags[i].ScheduleIDHash = s.IDHash
		}
		if err := s.StaffTags.Process(); err != nil {
			return err
		}
	}

	// process TaskTags
	// fmt.Println("Length of TaskTags: ", len(s.TaskTags))
	if len(s.TaskTags) > 0 {
		for i, _ := range s.TaskTags {
			s.TaskTags[i].ExtractDateTime = s.ExtractDateTime
			s.TaskTags[i].ScheduleKey = s.ScheduleKey
			s.TaskTags[i].LastModifiedDateUTC = s.LastModifiedDateUTC
			s.TaskTags[i].ScheduleIDHash = s.IDHash
		}
		if err := s.TaskTags.Process(); err != nil {
			return err
		}
	}

	// process LocationTags
	if len(s.LocationTags) > 0 {
		for i, _ := range s.LocationTags {
			s.LocationTags[i].ExtractDateTime = s.ExtractDateTime
			s.LocationTags[i].ScheduleKey = s.ScheduleKey
			s.LocationTags[i].LastModifiedDateUTC = s.LastModifiedDateUTC
			s.LocationTags[i].ScheduleIDHash = s.IDHash
		}
		if err := s.LocationTags.Process(); err != nil {
			return err
		}
	}
	if err := meta.ProcessStruct(s); err != nil {
		return fmt.Errorf("error processing %T:\t%q", s, err)
	}

	if err := s.SetMessage(); err != nil {
		return err
	}
	return nil
}

// SetMessage uses a copy, strips message metadata, remarshals to JSON and compacts it, and assigns the string to .ProcessedMessage
func (s *Schedule) SetMessage() error {
	// take a copy and strip message fields, for good measure
	ss := *s
	ss.RawMessage = nil
	ss.ProcessedMessage = nil

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

// SetIDHash takes the hash of the json encoded fields that (should) uniquely identify this instance
// for schedule, this is schedulekey, lastmodifieddateutc (in rfc3339 with nano precision)
func (s *Schedule) SetIDHash() error {
	idh := meta.ToValueMap(*s, "idhash").Hash()
	s.IDHash = &idh
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

func (s Schedule) CreatePGTable(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	return CreatePGTable(ctx, db, s, schema, table)

}

func (s Schedule) GetPGStatus(ctx context.Context, db *sqlx.DB, schema, table string) (*RequestConfig, error) {
	if table == "" {
		table = "schedule"
	}
	tpl := `
	{{- $field := .Struct.Fields.ByName "LastModifiedDateUTC" -}}
	select max ( {{ $field.TagName "db" }} ) {{ $field.TagName "qgendarequestname" }}
	from {{ .Struct.Identifier | tolower }}
	`
	return GetPGStatus(ctx, db, s, schema, table, tpl)
}
