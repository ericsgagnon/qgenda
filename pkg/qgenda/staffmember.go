package qgenda

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// func NewStaffMemberStaffIdRequest(rqf *RequestQueryFields) *Request {
// 	requestPath := "staffmember/:staffId"
// 	queryFields := []string{
// 		"CompanyKey",
// 		"Includes",
// 		"Select",
// 		"Filter",
// 		"OrderBy",
// 		"Expand",
// 	}
// 	if rqf != nil {
// 		if rqf.Includes == nil {
// 			rqf.SetIncludes("Skillset,Tags,Profiles")
// 		}
// 	}

// 	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
// 	return r
// }
// func NewStaffMemberLocationRequest(rqf *RequestQueryFields) *Request {
// 	requestPath := "staffmember/:staffId/location"
// 	queryFields := []string{
// 		"CompanyKey",
// 		"Select",
// 		"Filter",
// 		"OrderBy",
// 		"Expand",
// 	}

// 	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
// 	return r
// }
// func NewStaffMemberRequestLimitRequest(rqf *RequestQueryFields) *Request {
// 	requestPath := "staffmember/:staffId/requestlimit"
// 	queryFields := []string{
// 		"CompanyKey",
// 		"Includes",
// 		"Select",
// 		"Filter",
// 		"OrderBy",
// 		"Expand",
// 	}
// 	if rqf != nil {
// 		if rqf.Includes == nil {
// 			rqf.SetIncludes("ShiftsCredit")
// 		}
// 	}

// 	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
// 	return r
// }

type StaffMember struct {
	RawMessage               *string         `json:"-" db:"_raw_message" primarykey:"true"`
	ExtractDateTime          *Time           `json:"-" db:"_extract_date_time"`
	Abbrev                   *string         `json:"Abbrev,omitempty"`
	BgColor                  *string         `json:"BgColor,omitempty"`
	BillSysID                *string         `json:"BillSysId,omitempty"`
	CompKey                  *string         `json:"CompKey,omitempty"`
	ContactInstructions      *string         `json:"Contact Instructions,omitempty"`
	Email                    *string         `json:"Email,omitempty"`
	SsoID                    *string         `json:"SsoId,omitempty"`
	EmrID                    *string         `json:"EmrId,omitempty"`
	ErpID                    *string         `json:"ErpId,omitempty"`
	EndDate                  *Date           `json:"EndDate,omitempty"`
	ExtCallSysId             *string         `json:"ExtCallSysId,omitempty"`
	FirstName                *string         `json:"FirstName,omitempty"`
	LastName                 *string         `json:"LastName,omitempty"`
	HomePhone                *string         `json:"HomePhone,omitempty"`
	MobilePhone              *string         `json:"MobilePhone,omitempty"`
	Npi                      *string         `json:"Npi,omitempty"`
	OtherNumber1             *string         `json:"OtherNumber1,omitempty"`
	OtherNumber2             *string         `json:"OtherNumber2,omitempty"`
	OtherNumber3             *string         `json:"OtherNumber3,omitempty"`
	OtherNumberType1         *string         `json:"OtherNumberType1,omitempty"`
	OtherNumberType2         *string         `json:"OtherNumberType2,omitempty"`
	OtherNumberType3         *string         `json:"OtherNumberType3,omitempty"`
	Pager                    *string         `json:"Pager,omitempty"`
	PayrollId                *string         `json:"PayrollId,omitempty"`
	RegHours                 *float64        `json:"RegHours,omitempty"`
	StaffId                  *string         `json:"StaffId,omitempty"`
	StaffKey                 *string         `json:"StaffKey,omitempty" primarykey:"true"`
	StartDate                *Date           `json:"StartDate,omitempty"`
	TextColor                *string         `json:"TextColor,omitempty"`
	Addr1                    *string         `json:"Addr1,omitempty"`
	Addr2                    *string         `json:"Addr2,omitempty"`
	City                     *string         `json:"City,omitempty"`
	State                    *string         `json:"State,omitempty"`
	Zip                      *string         `json:"Zip,omitempty"`
	IsActive                 *bool           `json:"IsActive,omitempty"`
	StaffTypeKey             *string         `json:"StaffTypeKey,omitempty"`
	BillingTypeKey           *string         `json:"BillingTypeKey,omitempty"`
	UserProfileKey           *string         `json:"UserProfileKey,omitempty"`
	UserProfile              *string         `json:"UserProfile,omitempty"`
	PayPeriodGroupName       *string         `json:"PayPeriodGroupName,omitempty"`
	PayrollStartDate         *Date           `json:"PayrollStartDate,omitempty"`
	PayrollEndDate           *Date           `json:"PayrollEndDate,omitempty"`
	TimeClockStartDate       *Date           `json:"TimeClockStartDate,omitempty"`
	TimeClockEndDate         *Date           `json:"TimeClockEndDate,omitempty"`
	TimeClockKioskPIN        *string         `json:"TimeClockKioskPIN,omitempty"`
	IsAutoApproveSwap        *bool           `json:"IsAutoApproveSwap,omitempty"`
	DailyUnitAverage         *float64        `json:"DailyUnitAverage,omitempty"`
	StaffInternalId          *string         `json:"StaffInternalId,omitempty"`
	UserLastLoginDateTimeUTC *Time           `json:"UserLastLoginDateTimeUTC,omitempty"`
	SourceOfLogin            *string         `json:"SourceOfLogin,omitempty"`
	CalSyncKey               *string         `json:"CalSyncKey,omitempty"`
	Tags                     []TagCategory   `json:"Tags,omitempty"`
	TTCMTags                 []TagCategory   `json:"TTCMTags,omitempty"`
	Skillset                 []StaffSkillset `json:"Skillset,omitempty"`
	Profiles                 []Profile       `json:"Profiles,omitempty"`
}

type StaffMembers []StaffMember

type StaffSkillset struct {
	ExtractDateTime   *Time   `json:"-"`
	StaffKey          *string `json:"-"`
	StaffFirstName    *string `json:"StaffFirstName,omitempty"`
	StaffLastName     *string `json:"StaffLastName,omitempty"`
	StaffAbbreviation *string `json:"StaffAbbrev,omitempty"`
	StaffId           *string `json:"StaffId,omitempty"`
	TaskName          *string `json:"TaskName,omitempty"`
	TaskAbbreviation  *string `json:"TaskAbbrev,omitempty"`
	TaskId            *string `json:"TaskId,omitempty"`
	IsSkilledMon      *bool   `json:"IsSkilledMon,omitempty"`
	MonOccurrence     *string `json:"MonOccurrence,omitempty"`
	IsSkilledTue      *bool   `json:"IsSkilledTue,omitempty"`
	TueOccurrence     *string `json:"TueOccurrence,omitempty"`
	IsSkilledWed      *bool   `json:"IsSkilledWed,omitempty"`
	WedOccurrence     *string `json:"WedOccurrence,omitempty"`
	IsSkilledThu      *bool   `json:"IsSkilledThu,omitempty"`
	ThuOccurrence     *string `json:"ThuOccurrence,omitempty"`
	IsSkilledFri      *bool   `json:"IsSkilledFri,omitempty"`
	FriOccurrence     *string `json:"FriOccurrence,omitempty"`
	IsSkilledSat      *bool   `json:"IsSkilledSat,omitempty"`
	SatOccurrence     *string `json:"SatOccurrence,omitempty"`
	IsSkilledSun      *bool   `json:"IsSkilledSun,omitempty"`
	SunOccurrence     *string `json:"SunOccurrence,omitempty"`
}

type StaffTag struct {
}

func DefaultStaffMemberRequestQueryFields(rqf *RequestQueryFields) *RequestQueryFields {
	if rqf == nil {
		rqf = &RequestQueryFields{}
	}
	if rqf.Includes == nil {
		rqf.SetIncludes("StaffTags,TaskTags,LocationTags")
	}
	return rqf
}

func NewStaffMemberRequest(rqf *RequestQueryFields) *Request {
	requestPath := "staffmember"
	queryFields := []string{
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	if rqf != nil {
		if rqf.Includes == nil {
			rqf.SetIncludes("Skillset,Tags,TTCMTags,Profiles")
		}
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}

func (p *StaffMember) Process() error {
	ProcessStruct(p)
	for i, _ := range p.Tags {
		(&p.Tags[i]).Process()
	}
	for i, _ := range p.TTCMTags {
		(&p.TTCMTags[i]).Process()
	}

	return nil
}

func (s *StaffMember) DoStuff(db *sqlx.DB, ctx context.Context) (sql.Result, error) {
	txOptions := sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	}
	tx, err := db.BeginTxx(ctx, &txOptions)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// create table
	result, err := tx.Exec(``)
	if err != nil {
		tx.Rollback()
		return result, err
	}
	// create temp table
	// insert rows

	// create child tables
	// create child 


	if err := tx.Commit(); err != nil {
		tx.Rollback()
	}

	return nil, nil

	// there is no combination of fields for a good primary key
	// need to only insert new rows if there is a change in at least one (non-metadata) field
	// if _extract_date_time is only different field, don't insert

	// insert into staffmember ...
	// select

	// start transaction
	// create temp table
	// insert into temp table
	// use cte to join with current data and insert changes to target table
	// if err, rollback, else commit
	// not sure how to make sure no dangling temp tables

}

// import strategies:
// unique rows - audit logs, immutable rows
// full replace - snapshots where previous versions don't matter
// append changes - just need to know changes and when they take effect
// full append - full snapshots of every export

func (s *StaffMembers) XCreateTable(db *sqlx.DB, ctx context.Context) (sql.Result, error) {

	sqlStatement := `
CREATE TABLE IF NOT EXISTS staffmember (
	_raw_message text ,
	_extract_date_time timestamp with time zone ,
	abbrev text ,
	bgcolor text ,
	billsysid text ,
	compkey text ,
	contactinstructions text ,
	email text ,
	ssoid text ,
	emrid text ,
	erpid text ,
	enddate date ,
	extcallsysid text ,
	firstname text ,
	lastname text ,
	homephone text ,
	mobilephone text ,
	npi text ,
	othernumber1 text ,
	othernumber2 text ,
	othernumber3 text ,
	othernumbertype1 text ,
	othernumbertype2 text ,
	othernumbertype3 text ,
	pager text ,
	payrollid text ,
	reghours double precision ,
	staffid text ,
	staffkey text ,
	startdate date ,
	textcolor text ,
	addr1 text ,
	addr2 text ,
	city text ,
	state text ,
	zip text ,
	isactive boolean ,
	stafftypekey text ,
	billingtypekey text ,
	userprofilekey text ,
	userprofile text ,
	payperiodgroupname text ,
	payrollstartdate date ,
	payrollenddate date ,
	timeclockstartdate date ,
	timeclockenddate date ,
	timeclockkioskpin text ,
	isautoapproveswap boolean ,
	dailyunitaverage double precision ,
	staffinternalid text ,
	userlastlogindatetimeutc timestamp with time zone ,
	sourceoflogin text ,
	calsynckey text , 
PRIMARY KEY ( _raw_message, staffkey ) 

)
`

	return db.ExecContext(ctx, sqlStatement)

}

func LetsMakeHistory(db *sqlx.DB, ctx context.Context, alt bool) (sql.Result, error) {
	type TestType struct {
		Now time.Time
		ID  int
		Tag string
	}
	sqlStatement := `create table if not exists testtype (
		now timestamp with time zone,
		id bigint,
		tag text
		)`
	result, err := db.ExecContext(ctx, sqlStatement)
	if err != nil {
		return result, err
	}

	data := []TestType{}
	for i := 0; i < 20; i++ {
		var tag string
		if alt {
			tag = fmt.Sprint(i)
		} else {
			tag = fmt.Sprint(20 - i)
		}

		tt := TestType{
			Now: time.Now().UTC(),
			ID:  i,
			Tag: tag,
		}
		data = append(data, tt)
	}
	sqlStatement = `insert into testtype (now, id, tag) values ( :now , :id , :tag ) on conflict ( now, id, tag ) do nothing `
	result, err = db.NamedExecContext(ctx, sqlStatement, data)
	if err != nil {
		return result, err
	}
	return result, nil
}

var TemplateForChangesOnly string = `
--- as a DB transaction

--- create temporary table _tmp_table_name

--- insert into _tmp_table_name

--- insert into table_name with:
--- _tmp_table_name and table_name as CTE's
--- resolve conflict in a CTE
--- conflict resolution strategies:
--- changes only: new rows can't = most recent rows (by id groups)
--- 


--- insert into table_name

`

var TemplateForFullAppend string
var TemplateForNewAppend string
var TemplateForFullReplace string
var TemplateForReplaceOnConflict string
