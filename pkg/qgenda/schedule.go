package qgenda

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Schedule struct {
	// RawMessage        *string    `json:"-" db:"_raw_message"`
	// ExtractDateTime   *Time      `json:"-" db:"_extract_date_time"`
	ScheduleKey            *string       `json:"ScheduleKey,omitempty" primarykey:"true"`
	CallRole               *string       `json:"CallRole,omitempty"`
	CompKey                *string       `json:"CompKey,omitempty"`
	Credit                 *float64      `json:"Credit,omitempty"`
	Date                   *Date         `json:"Date,omitempty"`
	StartDateUTC           *Time         `json:"StartDateUTC,omitempty"`
	EndDateUTC             *Time         `json:"EndDateUTC,omitempty"`
	EndDate                *Date         `json:"EndDate,omitempty"`
	EndTime                *TimeOfDay    `json:"EndTime,omitempty"`
	IsCred                 *bool         `json:"IsCred,omitempty"`
	IsPublished            *bool         `json:"IsPublished,omitempty"`
	IsLocked               *bool         `json:"IsLocked,omitempty"`
	IsStruck               *bool         `json:"IsStruck,omitempty"`
	Notes                  *string       `json:"Notes,omitempty"`
	IsNotePrivate          *bool         `json:"IsNotePrivate,omitempty"`
	StaffAbbrev            *string       `json:"StaffAbbrev,omitempty"`
	StaffBillSysId         *string       `json:"StaffBillSysId,omitempty"`
	StaffEmail             *string       `json:"StaffEmail,omitempty"`
	StaffEmrId             *string       `json:"StaffEmrId,omitempty"`
	StaffErpId             *string       `json:"StaffErpId,omitempty"`
	StaffInternalId        *string       `json:"StaffInternalId,omitempty"`
	StaffExtCallSysId      *string       `json:"StaffExtCallSysId,omitempty"`
	StaffFName             *string       `json:"StaffFName,omitempty"`
	StaffId                *string       `json:"StaffId,omitempty"`
	StaffKey               *string       `json:"StaffKey,omitempty"`
	StaffLName             *string       `json:"StaffLName,omitempty"`
	StaffMobilePhone       *string       `json:"StaffMobilePhone,omitempty"`
	StaffNpi               *string       `json:"StaffNpi,omitempty"`
	StaffPager             *string       `json:"StaffPager,omitempty"`
	StaffPayrollId         *string       `json:"StaffPayrollId,omitempty"`
	StaffTags              []ScheduleTag `json:"StaffTags,omitempty"`
	StartDate              *Date         `json:"StartDate,omitempty"`
	StartTime              *TimeOfDay    `json:"StartTime,omitempty"`
	TaskAbbrev             *string       `json:"TaskAbbrev,omitempty"`
	TaskBillSysId          *string       `json:"TaskBillSysId,omitempty"`
	TaskContactInformation *string       `json:"TaskContactInformation,omitempty"`
	TaskExtCallSysId       *string       `json:"TaskExtCallSysId,omitempty"`
	TaskId                 *string       `json:"TaskId,omitempty"`
	TaskKey                *string       `json:"TaskKey,omitempty"`
	TaskName               *string       `json:"TaskName,omitempty"`
	TaskPayrollId          *string       `json:"TaskPayrollId,omitempty"`
	TaskEmrId              *string       `json:"TaskEmrId,omitempty"`
	TaskCallPriority       *string       `json:"TaskCallPriority,omitempty"`
	TaskDepartmentId       *string       `json:"TaskDepartmentId,omitempty"`
	TaskIsPrintEnd         *bool         `json:"TaskIsPrintEnd,omitempty"`
	TaskIsPrintStart       *bool         `json:"TaskIsPrintStart,omitempty"`
	TaskShiftKey           *string       `json:"TaskShiftKey,omitempty"`
	TaskType               *string       `json:"TaskType,omitempty"`
	TaskTags               []ScheduleTag `json:"TaskTags,omitempty"`
	LocationName           *string       `json:"LocationName,omitempty"`
	LocationAbbrev         *string       `json:"LocationAbbrev,omitempty"`
	LocationID             *string       `json:"LocationID,omitempty"`
	LocationAddress        *string       `json:"LocationAddress,omitempty"`
	TimeZone               *string       `json:"TimeZone,omitempty"`
	LastModifiedDateUTC    *Time         `json:"LastModifiedDateUTC,omitempty" primarykey:"true" querycondition:"ge" qf:"SinceModifiedTimestamp"`
	LocationTags           []Location    `json:"LocationTags,omitempty"`
	IsRotationTask         *bool         `json:"IsRotationTask"`
}

type ScheduleTag struct {
	ScheduleKey         *string `json:"ScheduleKey,omitempty" nullable:"false"`
	LastModifiedDateUTC *Time   `json:"LastModifiedDateUTC,omitempty" nullable:"false"`
	CategoryKey         *int64  `json:"CategoryKey" nullable:"false"`
	CategoryName        *string `json:"CategoryName" nullable:"false"`
	Tags                []struct {
		Key  *int64  `json:"Key" db:"tagkey" nullable:"false"`
		Name *string `json:"Name" db:"tagname" nullable:"false"`
	}
}

func DefaultScheduleRequestQueryFields(rqf *RequestQueryFields) *RequestQueryFields {
	if rqf == nil {
		rqf = &RequestQueryFields{}
	}
	if rqf.StartDate == nil {
		rqf.SetStartDate(time.Now().UTC().Add(time.Hour * 24 * -1).Truncate(time.Hour * 24))
	}
	if rqf.EndDate == nil || rqf.EndDate.Sub(rqf.GetStartDate()) > time.Hour*24*100 {
		rqf.SetEndDate(rqf.GetStartDate().Add(time.Hour * 24 * 100))
	}
	if rqf.SinceModifiedTimestamp == nil {
		rqf.SetSinceModifiedTimestamp(rqf.GetStartDate())
	}
	if rqf.Includes == nil {
		rqf.SetIncludes("StaffTags,TaskTags,LocationTags")
	}
	return rqf
}

func NewScheduleRequest(rqf *RequestQueryFields) *Request {
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

type pgScheduleTag struct {
	ScheduleKey         *string `json:"ScheduleKey,omitempty" nullable:"false"`
	LastModifiedDateUTC *Time   `json:"LastModifiedDateUTC,omitempty" nullable:"false"`
	CategoryKey         *int64  `json:"CategoryKey" nullable:"false"`
	CategoryName        *string `json:"CategoryName" nullable:"false"`
	TagKey              *int64  `db:"tagkey" nullable:"false"`
	TagName             *string `db:"tagname" nullable:"false"`
}

func ExecSchedulePipeline(ctx context.Context, db *sqlx.DB, value []Schedule, schema, table string) (sql.Result, error) {
	// rename - maybe pgloader or similar
	cstResult, err := PGCreateTable(ctx, db, value, schema, table)
	if err != nil {
		return cstResult, err
	}
	isrResult, err := PGInsertRows(ctx, db, value, schema, table)
	if err != nil {
		return isrResult, err
	}

	sttablename := fmt.Sprintf("%sstafftag", table)
	stafftags := []pgScheduleTag{}
	// csstResult, err := PGCreateTable(ctx, db, value[0].StaffTags, schema, sttablename)
	for _, sch := range value {
		for _, cat := range sch.StaffTags {
			for _, tag := range cat.Tags {
				stafftag := pgScheduleTag{
					ScheduleKey:         sch.ScheduleKey,
					LastModifiedDateUTC: sch.LastModifiedDateUTC,
					CategoryKey:         cat.CategoryKey,
					CategoryName:        cat.CategoryName,
					TagKey:              tag.Key,
					TagName:             tag.Name,
				}
				stafftags = append(stafftags, stafftag)
			}
		}
	}
	csstResult, err := PGCreateTable(ctx, db, stafftags, schema, sttablename)
	if err != nil {
		return csstResult, err
	}
	istrResult, err := PGInsertRows(ctx, db, stafftags, schema, sttablename)
	if err != nil {
		return istrResult, err
	}

	tttablename := fmt.Sprintf("%stasktag", table)
	tasktags := []pgScheduleTag{}
	for _, sch := range value {
		for _, cat := range sch.TaskTags {
			for _, tag := range cat.Tags {
				tasktag := pgScheduleTag{
					ScheduleKey:         sch.ScheduleKey,
					LastModifiedDateUTC: sch.LastModifiedDateUTC,
					CategoryKey:         cat.CategoryKey,
					CategoryName:        cat.CategoryName,
					TagKey:              tag.Key,
					TagName:             tag.Name,
				}
				tasktags = append(tasktags, tasktag)
			}
		}
	}
	csttResult, err := PGCreateTable(ctx, db, tasktags, schema, tttablename)
	if err != nil {
		return csttResult, err
	}
	ittrResult, err := PGInsertRows(ctx, db, tasktags, schema, fmt.Sprintf("%stasktag", table))
	if err != nil {
		return ittrResult, err
	}

	return isrResult, err
}

// func PGCreateScheduleTableStatement(schema, table string) string {
// 	return PGCreateTableStatement(Schedule{}, schema, table)
// }

// func PGCreateScheduleTagStatement(schema, table string) string {
// 	return PGCreateTableStatement(ScheduleTag{}, schema, table)
// }

// func PGInsertScheduleStatement(schema, table string) string {
// 	return PGInsertStatement(Schedule{}, schema, table)
// }

// func PGInsertScheduleTagStatement(schema, table string) string {
// 	return PGInsertStatement(ScheduleTag{}, schema, table)
// }

// func PGQueryScheduleConstraintStatement(schema, table string) string {
// 	return PGQueryConstraintsStatement(Schedule{}, schema, table)
// }

// func PGQueryScheduleTagConstraintStatement(schema, table string) string {
// 	return PGQueryConstraintsStatement(ScheduleTag{}, schema, table)
// }

// func (sch Schedule) PGCreateTable() {

// }
// if rqf == nil {
// 	rqf = &RequestQueryFields{}
// }
// if rqf.StartDate == nil {
// 	rqf.SetStartDate(time.Now().UTC().Add(time.Hour * 24 * -1).Truncate(time.Hour * 24))
// }
// if rqf.EndDate == nil || rqf.EndDate.Sub(rqf.GetStartDate()) > time.Hour*24*100 {
// 	rqf.SetEndDate(rqf.GetStartDate().Add(time.Hour * 24 * 100))
// }
// if rqf.SinceModifiedTimestamp == nil {
// 	rqf.SetSinceModifiedTimestamp(rqf.GetStartDate())
// }
// if rqf.Includes == nil {
// 	rqf.SetIncludes("StaffTags,TaskTags,LocationTags")
// }
