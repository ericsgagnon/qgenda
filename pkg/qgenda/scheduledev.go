package qgenda

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

// type ScheduleDataset struct {
// 	MetaData
// 	Schedules []XSchedule
// }

type XSchedules struct {
	ExtractDateTime *Time `json:"extractDateTime"`
	Schedules       []XSchedule
}

func (ss *XSchedules) Process() error {
	scs := []XSchedule{}
	for _, s := range ss.Schedules {
		sp := &s
		if err := sp.Process(); err != nil {
			return err
		}
		scs = append(scs, s)

	}
	sort.SliceStable(scs, func(i, j int) bool {
		return *(scs[i].ScheduleKey) < *(scs[j].ScheduleKey)
	})
	ss.Schedules = scs

	return nil
}

func (ss XSchedules) MarshalJSON() ([]byte, error) {
	return json.Marshal(ss)
}

type XSchedule struct {
	RawMessage             *string         `json:"-" db:"_raw_message"`
	ProcessedMessage       *string         `json:"-" db:"_processed_message"` // RawMessage processed, with changing fields dropped
	ExtractDateTime        *Time           `json:"-" db:"_extract_date_time"`
	HashID                 *string         `json:"-" db:"_hash_id"`
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
// assigns a compact json representation to the RawMessage member
func (s *XSchedule) UnmarshalJSON(b []byte) error {
	// alias technique to avoid infinite recursion
	type Alias XSchedule
	var a Alias

	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}

	dest := XSchedule(a)
	var bb bytes.Buffer
	if err := json.Compact(&bb, b); err != nil {
		return err
	}
	rawMessage := bb.String()
	dest.RawMessage = &rawMessage

	*s = dest
	return nil
}

// MarshalJSON satisfies the json.Marshaler interface
func (s *XSchedule) MarshalJSON() ([]byte, error) {
	type Alias XSchedule
	a := Alias(*s)

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
	// s.StaffTags.ExtractDateTime = s.ExtractDateTime
	// s.StaffTags.ScheduleKey = s.ScheduleKey
	// s.StaffTags.LastModifiedDateUTC = s.LastModifiedDateUTC
	// sp := s.StaffTags
	// if err := sp.Process(); err != nil {
	// 	return err
	// }
	// s.StaffTags = sp

	// process stafftags
	if len(s.StaffTags) > 0 {
		for _, v := range s.StaffTags {
			v.ExtractDateTime = s.ExtractDateTime
			v.ScheduleKey = s.ScheduleKey
			v.LastModifiedDateUTC = s.LastModifiedDateUTC
			if err := v.Process(); err != nil {
				return err
			}
		}

	}

	// process TaskTags
	// ss := s.TaskTags
	if len(s.TaskTags) > 0 {
		for _, v := range s.TaskTags {
			v.ExtractDateTime = s.ExtractDateTime
			v.ScheduleKey = s.ScheduleKey
			v.LastModifiedDateUTC = s.LastModifiedDateUTC
			if err := v.Process(); err != nil {
				return err
			}
		}

	}

	// process LocationTags
	if len(s.LocationTags) > 0 {
		for _, v := range s.LocationTags {
			v.ExtractDateTime = s.ExtractDateTime
			v.ScheduleKey = s.ScheduleKey
			v.LastModifiedDateUTC = s.LastModifiedDateUTC
			if err := v.Process(); err != nil {
				return err
			}
		}

	}

	return nil
}

// SetMessage processes the struct, remarshals and compacts it, and assigns the string to .ProcessedMessage
func (s *XSchedule) SetMessage() error {
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
	s.Schedules = ss
	return nil
}

// extract
// process
// load
//   create table
//   create temp table
//   insert into temp table
//   insert into table from temp table
//
