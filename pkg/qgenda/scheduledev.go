package qgenda

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type SchedulesX struct {
	ExtractDateTime *Time `json:"extractDateTime"`
	Schedules       []ScheduleX
}

type ScheduleX struct {
	RawMessage             *string      `json:"-" db:"_raw_message"`
	ProcessedMessage       *string      `json:"-" db:"_processed_message"` // RawMessage processed, with changing fields dropped
	ExtractDateTime        *Time        `json:"-" db:"_extract_date_time"`
	ScheduleKey            *string      `json:"ScheduleKey,omitempty" primarykey:"true"`
	CallRole               *string      `json:"CallRole,omitempty"`
	CompKey                *string      `json:"CompKey,omitempty"`
	Credit                 *float64     `json:"Credit,omitempty"`
	Date                   *Date        `json:"Date,omitempty"`
	StartDateUTC           *Time        `json:"StartDateUTC,omitempty"`
	EndDateUTC             *Time        `json:"EndDateUTC,omitempty"`
	EndDate                *Date        `json:"EndDate,omitempty"`
	EndTime                *TimeOfDay   `json:"EndTime,omitempty"`
	IsCred                 *bool        `json:"IsCred,omitempty"`
	IsPublished            *bool        `json:"IsPublished,omitempty"`
	IsLocked               *bool        `json:"IsLocked,omitempty"`
	IsStruck               *bool        `json:"IsStruck,omitempty"`
	Notes                  *string      `json:"Notes,omitempty"`
	IsNotePrivate          *bool        `json:"IsNotePrivate,omitempty"`
	StaffAbbrev            *string      `json:"StaffAbbrev,omitempty"`
	StaffBillSysId         *string      `json:"StaffBillSysId,omitempty"`
	StaffEmail             *string      `json:"StaffEmail,omitempty"`
	StaffEmrId             *string      `json:"StaffEmrId,omitempty"`
	StaffErpId             *string      `json:"StaffErpId,omitempty"`
	StaffInternalId        *string      `json:"StaffInternalId,omitempty"`
	StaffExtCallSysId      *string      `json:"StaffExtCallSysId,omitempty"`
	StaffFName             *string      `json:"StaffFName,omitempty"`
	StaffId                *string      `json:"StaffId,omitempty"`
	StaffKey               *string      `json:"StaffKey,omitempty"`
	StaffLName             *string      `json:"StaffLName,omitempty"`
	StaffMobilePhone       *string      `json:"StaffMobilePhone,omitempty"`
	StaffNpi               *string      `json:"StaffNpi,omitempty"`
	StaffPager             *string      `json:"StaffPager,omitempty"`
	StaffPayrollId         *string      `json:"StaffPayrollId,omitempty"`
	StaffTags              ScheduleTags `json:"StaffTags,omitempty"`
	StartDate              *Date        `json:"StartDate,omitempty"`
	StartTime              *TimeOfDay   `json:"StartTime,omitempty"`
	TaskAbbrev             *string      `json:"TaskAbbrev,omitempty"`
	TaskBillSysId          *string      `json:"TaskBillSysId,omitempty"`
	TaskContactInformation *string      `json:"TaskContactInformation,omitempty"`
	TaskExtCallSysId       *string      `json:"TaskExtCallSysId,omitempty"`
	TaskId                 *string      `json:"TaskId,omitempty"`
	TaskKey                *string      `json:"TaskKey,omitempty"`
	TaskName               *string      `json:"TaskName,omitempty"`
	TaskPayrollId          *string      `json:"TaskPayrollId,omitempty"`
	TaskEmrId              *string      `json:"TaskEmrId,omitempty"`
	TaskCallPriority       *string      `json:"TaskCallPriority,omitempty"`
	TaskDepartmentId       *string      `json:"TaskDepartmentId,omitempty"`
	TaskIsPrintEnd         *bool        `json:"TaskIsPrintEnd,omitempty"`
	TaskIsPrintStart       *bool        `json:"TaskIsPrintStart,omitempty"`
	TaskShiftKey           *string      `json:"TaskShiftKey,omitempty"`
	TaskType               *string      `json:"TaskType,omitempty"`
	TaskTags               ScheduleTags `json:"TaskTags,omitempty"`
	LocationName           *string      `json:"LocationName,omitempty"`
	LocationAbbrev         *string      `json:"LocationAbbrev,omitempty"`
	LocationID             *string      `json:"LocationID,omitempty"`
	LocationAddress        *string      `json:"LocationAddress,omitempty"`
	TimeZone               *string      `json:"TimeZone,omitempty"`
	LastModifiedDateUTC    *Time        `json:"LastModifiedDateUTC,omitempty" primarykey:"true" querycondition:"ge" qf:"SinceModifiedTimestamp"`
	LocationTags           []Location   `json:"LocationTags,omitempty"`
	IsRotationTask         *bool        `json:"IsRotationTask"`
}

func (s *ScheduleX) UnmarshalJSON(b []byte) error {
	// alias technique to avoid infinite recursion
	type Alias ScheduleX
	var a Alias

	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}

	dest := ScheduleX(a)
	var bb bytes.Buffer
	if err := json.Compact(&bb, b); err != nil {
		return err
	}
	rawMessage := bb.String()
	dest.RawMessage = &rawMessage

	*s = dest
	return nil
}

// Process handles all the basic validating and processing of
// from the raw version of any values. It is idempotent.
func (s *ScheduleX) Process() error {

	if err := ProcessStruct(s); err != nil {
		return fmt.Errorf("error processing %T:\t%q", s, err)
	}

	// process stafftags
	// if StaffTags(s.StaffTags) != nil {

	// }

	// process tasktags

	// process locations
	return nil
}

// SetMessage processes the struct, remarshals and compacts it, and assigns the string to .ProcessedMessage
func (s *ScheduleX) SetMessage() error {
	if err := s.Process(); err != nil {
		return err
	}
	// take a copy and strip metadata, for good measure
	ss := *s
	ss.RawMessage = nil
	ss.ExtractDateTime = nil
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

type ScheduleTags struct {
	ExtractDateTime     *Time          `json:"-" db:"_extract_date_time"`
	ScheduleKey         *string        `json:"ScheduleKey,omitempty"`
	LastModifiedDateUTC *Time          `json:"LastModifiedDateUTC,omitempty"`
	CategoryKey         *int64         `json:"CategoryKey"`
	CategoryName        *string        `json:"CategoryName"`
	Tags                []ScheduleTagX `json:"Tags,omitempty"`
}

type ScheduleTagX struct {
	ExtractDateTime     *Time   `json:"-" db:"_extract_date_time"`
	ScheduleKey         *string `json:"ScheduleKey,omitempty"`
	LastModifiedDateUTC *Time   `json:"LastModifiedDateUTC,omitempty"`
	CategoryKey         *int64  `json:"CategoryKey,omitempty"`
	CategoryName        *string `json:"CategoryName,omitempty"`
	TagKey              *int64  `db:"tagkey"`
	TagName             *string `db:"tagname"`
}

func ProcessScheduleTags(st ScheduleTags) ScheduleTags {

	return st
}

type ScheduleLocations struct {
	ExtractDateTime     *Time   `json:"-" db:"_extract_date_time"`
	ScheduleKey         *string `json:"ScheduleKey,omitempty"`
	LastModifiedDateUTC *Time   `json:"LastModifiedDateUTC,omitempty"`
	Locations           []Location
}

type ScheduleLocation struct {
	CompanyKey  *string       `json:"CompanyKey,omitempty"`
	LocationKey *int64        `json:"LocationKey,omitempty"`
	ID          *string       `json:"Id,omitempty"`
	Name        *string       `json:"Name,omitempty"`
	Address     *string       `json:"Address,omitempty"`
	Abbrev      *string       `json:"Abbrev,omitempty"`
	Notes       *string       `json:"Notes,omitempty"`
	TimeZone    *string       `json:"TimeZone,omitempty"`
	Tags        []TagCategory `json:"Tags,omitempty"`
}

func ProcessScheduleLocations(sl ScheduleLocations) ScheduleLocations {
	return sl
}

// extract
// process
// load
//   create table
//   create temp table
//   insert into temp table
//   insert into table from temp table
//
