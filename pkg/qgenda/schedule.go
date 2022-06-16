package qgenda

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Schedule struct {
	RawMessage             *string       `json:"-" db:"_raw_message"`
	ExtractDateTime        *Time         `json:"-" db:"_extract_date_time"`
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

type Schedules []Schedule

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

func (s Schedules) InsertRows(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	switch db.DriverName() {
	case "postgres":
		return s.InsertPGRows(ctx, db, schema, table)
	default:
		return nil, fmt.Errorf("%s not found in imported drivers: %s", db.DriverName(), sql.Drivers())
	}
}

// Load is just a wrapper for InsertRows
func (s Schedules) Load(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	return s.InsertRows(ctx, db, schema, table)
}

func (s Schedules) InsertPGRows(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	// fmt.Printf("InsertPGRows: length of schedules: %d\n", len(s))
	sch := []Schedule{}
	for _, v := range s {
		sch = append(sch, v)
	}
	return LoadSchedulesToPG(ctx, db, sch, schema, table)
}
func (s Schedules) QueryConstraints(ctx context.Context, db *sqlx.DB, schema, table string) (*RequestQueryFields, error) {
	switch db.DriverName() {
	case "postgres":
		return s.QueryPGConstraints(ctx, db, schema, table)
	default:
		return nil, fmt.Errorf("%s not found in imported drivers: %s", db.DriverName(), sql.Drivers())
	}
}

func (s Schedules) QueryPGConstraints(ctx context.Context, db *sqlx.DB, schema, table string) (*RequestQueryFields, error) {
	rqf := RequestQueryFields{}
	if err := db.GetContext(
		ctx,
		&rqf,
		PGQueryConstraintsStatement(Schedule{}, schema, table),
	); err != nil {
		return nil, err
	}
	return &rqf, nil
}

func (s Schedules) CreateTable(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	switch db.DriverName() {
	case "postgres":
		return s.CreatePGTable(ctx, db, schema, table)
	default:
		return nil, fmt.Errorf("%s not found in imported drivers: %s", db.DriverName(), sql.Drivers())
	}
}

func (s Schedules) CreatePGTable(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	cstResult, err := PGCreateTable(ctx, db, s, schema, table)
	if err != nil {
		return cstResult, err
	}

	sttablename := fmt.Sprintf("%sstafftag", table)
	stafftags := []pgScheduleTag{}
	csstResult, err := PGCreateTable(ctx, db, stafftags, schema, sttablename)
	if err != nil {
		return csstResult, err
	}

	tttablename := fmt.Sprintf("%stasktag", table)
	tasktags := []pgScheduleTag{}
	csttResult, err := PGCreateTable(ctx, db, tasktags, schema, tttablename)
	if err != nil {
		return csttResult, err
	}

	return cstResult, nil
}

func (s Schedules) DropTable(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	switch db.DriverName() {
	case "postgres":
		// fmt.Println("DropTable: DriverName ", db.DriverName())
		return s.DropPGTable(ctx, db, schema, table)
		// fmt.Println("----------------------------------------")
	default:
		return nil, fmt.Errorf("%s not found in imported drivers: %s", db.DriverName(), sql.Drivers())
	}
}

func (s Schedules) DropPGTable(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {

	var res Result
	sqlResult, err := PGDropTable(ctx, db, s, schema, table)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}

	sttablename := fmt.Sprintf("%sstafftag", table)
	stafftags := []pgScheduleTag{}
	sqlResult, err = PGDropTable(ctx, db, stafftags, schema, sttablename)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}

	tttablename := fmt.Sprintf("%stasktag", table)
	tasktags := []pgScheduleTag{}
	sqlResult, err = PGDropTable(ctx, db, tasktags, schema, tttablename)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}

	return res, nil
}

func LoadSchedulesToPG(ctx context.Context, db *sqlx.DB, value []Schedule, schema, table string) (sql.Result, error) {
	if len(value) < 1 {
		return nil, fmt.Errorf("LoadSchedulesToPG: length of %T < 1, nothing to do", value)
	}
	// fmt.Println("Creating Schedule Table")
	// rename - maybe pgloader or similar
	var res Result
	sqlResult, err := PGCreateTable(ctx, db, value, schema, table)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}
	// fmt.Println("Inserting Schedule Rows")
	sqlResult, err = PGInsertRows(ctx, db, value, schema, table)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
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
	// fmt.Println("Creating ScheduleStaffTag Table")
	sqlResult, err = PGCreateTable(ctx, db, stafftags, schema, sttablename)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}
	// fmt.Println("Inserting ScheduleStaffTag Rows")
	sqlResult, err = PGInsertRows(ctx, db, stafftags, schema, sttablename)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
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
	// fmt.Println("Creating ScheduleTaskTag Table")
	sqlResult, err = PGCreateTable(ctx, db, tasktags, schema, tttablename)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}
	// fmt.Println("Inserting ScheduleTaskTag Rows")
	sqlResult, err = PGInsertRows(ctx, db, tasktags, schema, fmt.Sprintf("%stasktag", table))
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}
	// fmt.Println("actually managed to insert all rows?")
	return res, err
}

func (s Schedules) Extract(ctx context.Context, c *Client, rqf *RequestQueryFields) (Schedules, error) {
	req := NewScheduleRequest(rqf)
	sch := Schedules{}
	// qgenda only supports 100 days of schedules per query
	// using 90 days in case there are other limits
	duration := time.Hour * 24 * 90
	// if req.GetEndDate().Sub(req.GetStartDate()) > duration {
	for t := req.GetStartDate(); t.Before(req.GetEndDate()); t = t.Add(duration) {
		vcpreq := *req
		subreq := &(vcpreq)
		subreq.SetStartDate(t)
		subreq.SetEndDate(t.Add(duration))
		// fmt.Printf("Original StartDate: %s\t Current StartDate: %s\n", req.GetStartDate(), subreq.GetStartDate())
		// fmt.Printf("Original SinceModified: %s\tCurrent SinceModified: %s\n", req.GetSinceModifiedTimestamp(), subreq.GetSinceModifiedTimestamp())
		resp, err := c.Do(ctx, subreq)
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		ts := Schedules{}
		if err := json.Unmarshal(data, &ts); err != nil {
			return nil, err
		}

		// user response header for extract time
		// fmt.Println(resp.Header.Get("date"))
		et, err := ParseTime(resp.Header.Get("date"))
		if err != nil {
			return nil, err
		}
		var messages []any
		if err := json.Unmarshal(data, &messages); err != nil {
			return nil, err
		}
		for i, v := range messages {
			rw, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			srw := string(rw)
			ts[i].RawMessage = &srw
			ts[i].ExtractDateTime = &et
			// *s[i].RawMessage = string(rw)
		}
		log.Printf("schedules: %s - %s (modTime>= %s)\ttotal: %d", subreq.GetStartDate(), subreq.GetEndDate(), subreq.GetSinceModifiedTimestamp(), len(ts))
		sch = append(sch, ts...)
	}
	// }
	return sch, nil
	// resp, err := c.Do(ctx, req)
	// if err != nil {
	// 	return nil, err
	// }
	// data, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := json.Unmarshal(data, &s); err != nil {
	// 	return nil, err
	// }

	// user response header for extract time
	// fmt.Println(resp.Header.Get("date"))
	// t, err := ParseTime(resp.Header.Get("date"))
	// if err != nil {
	// 	return nil, err
	// }

	// grab the 'raw message'
	// var messages []any
	// if err := json.Unmarshal(data, &messages); err != nil {
	// 	return nil, err
	// }
	// for i, v := range messages {
	// 	rw, err := json.Marshal(v)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	srw := string(rw)
	// 	s[i].RawMessage = &srw
	// 	s[i].ExtractDateTime = &t
	// 	// *s[i].RawMessage = string(rw)
	// }
	// return s, nil

}

func (s Schedules) Process() (Schedules, error) {
	out := []Schedule{}
	for _, v := range s {
		if err := Process(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func (s Schedules) EPL(ctx context.Context, c *Client, rqf *RequestQueryFields,
	db *sqlx.DB, schema, table string, newRowsOnly bool) (sql.Result, error) {
	rqf = DefaultScheduleRequestQueryFields(rqf)

	var res Result
	sqlResult, err := PGCreateSchema(ctx, db, s, schema, table)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}

	sqlResult, err = s.CreateTable(ctx, db, schema, table)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}

	qrqf, err := s.QueryConstraints(ctx, db, schema, table)
	if err != nil {
		return res, err
	}
	if qrqf.SinceModifiedTimestamp != nil && newRowsOnly {
		rqf.SetSinceModifiedTimestamp(qrqf.GetSinceModifiedTimestamp())
	}
	s, err = s.Extract(ctx, c, rqf)
	if err != nil {
		return res, err
	}
	s, err = s.Process()
	if err != nil {
		return res, err
	}
	sqlResult, err = s.Load(ctx, db, schema, table)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	return res, err
}
