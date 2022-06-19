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
	RawMessage               *string         `json:"-" db:"_raw_message"`
	ExtractDateTime          *Time           `json:"-" db:"_extract_date_time" dependentID:"true"`
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
	StaffKey                 *string         `json:"StaffKey,omitempty" dependentID:"true"`
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
	Tags                     []StaffTag      `json:"Tags,omitempty"`
	TTCMTags                 []StaffTag      `json:"TTCMTags,omitempty"`
	Skillset                 []StaffSkillset `json:"Skillset,omitempty"`
	Profiles                 []Profile       `json:"Profiles,omitempty"`
}

type StaffMembers []StaffMember

type StaffSkillset struct {
	ExtractDateTime   *Time   `json:"-" nullable:"false"`
	StaffKey          *string `json:"-" nullable:"false"`
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

// StaffTag is basically TagCategory with the StaffMember.StaffKey and ExtractDateTime added for dependent tables in DB's
type StaffTag struct {
	ExtractDateTime     *Time   `json:"-" nullable:"false"`
	StaffKey            *string `json:"-" nullable:"false"`
	LastModifiedDateUTC *Time   `json:"LastModifiedDateUTC,omitempty" nullable:"false"`
	CategoryKey         *int64  `json:"CategoryKey" nullable:"false"`
	CategoryName        *string `json:"CategoryName" nullable:"false"`
	Tags                []struct {
		Key  *int64  `json:"Key" db:"tagkey" nullable:"false"`
		Name *string `json:"Name" db:"tagname" nullable:"false"`
	}
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
	if p.Tags != nil {
		for _, v := range p.Tags {
			// (&p.Tags[i]).Process()
			if err := Process(&v); err != nil {
				return err
			}
		}
	}
	if p.TTCMTags != nil {
		for _, v := range p.TTCMTags {
			// (&p.TTCMTags[i]).Process()
			if err := Process(&v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s StaffMember) CreatePGTable(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {

	result, err := pgCreateTable(ctx, db, s, schema, table)
	if err != nil {
		return result, err
	}
	sqlResult := SQLResult(result)

	tablename := fmt.Sprintf("%stag", table)
	value := []StaffTag{}
	result, err = PGCreateTable(ctx, db, value, schema, tablename)
	if err != nil {
		return result, err
	}
	sqlResult = SQLResult(sqlResult, result)

	tablename = fmt.Sprintf("%sttcmtag", table)
	result, err = PGCreateTable(ctx, db, []StaffTag{}, schema, tablename)
	if err != nil {
		return result, err
	}
	sqlResult = SQLResult(sqlResult, result)

	tablename = fmt.Sprintf("%sskillset", table)
	result, err = PGCreateTable(ctx, db, []StaffSkillset{}, schema, tablename)
	if err != nil {
		return result, err
	}
	sqlResult = SQLResult(sqlResult, result)

	tablename = fmt.Sprintf("%sprofile", table)

	result, err = PGCreateTable(ctx, db, []Profile{}, schema, tablename)
	// result, err := PGCreateTable(ctx, db, value, schema, tablename)
	if err != nil {
		return result, err
	}
	sqlResult = SQLResult(sqlResult, result)
	return sqlResult, nil
	// Profiles                 []Profile       `json:"Profiles,omitempty"`

	// create temporary table if not exists _tmp_testtype (
	// 	now timestamp with time zone,
	// 	id bigint,
	// 	tag text
	// 	)
	// ;
	// insert into _tmp_testtype values ( now(), (random() * 100 )::bigint , ( random() * 100 )::text);
	// insert into _tmp_testtype select * from testtype;
	// -- select * from pg_temp._tmp_testtype;

	// -- select count(*) from testtype;
	// with cte_row_numbers as (

	// select
	// row_number() over (partition by tt.id order by tt.now desc) rn,
	// tt.now,
	// tt.id,
	// tt.tag
	// -- from testtype tt
	// from _tmp_testtype tt
	// ),
	// cte_most_recent as (
	// select * from cte_row_numbers where rn = 1
	// ),
	// cte_anti_joined as (
	// select
	// r.now,
	// r.id,
	// r.tag
	// from cte_most_recent r
	// where not exists (
	// select 1
	// from testtype tt where
	// 	r.id = tt.id
	// and r.tag = tt.tag
	// )
	// )
	// insert into testtype
	// select * from cte_anti_joined

	// ;

	// select * from testtype tt order by tt.now desc;

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
	txOptions := sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	}
	tx, err := db.BeginTxx(ctx, &txOptions)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	var result sql.Result
	res, err := tx.ExecContext(ctx, sqlStatement)
	result = SQLResult(result, res)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	// need a temp table
	fmt.Println("creating temp table")
	sqlStatement = `create temporary table if not exists _tmp_testtype (
		now timestamp with time zone,
		id bigint,
		tag text
		) -- ON COMMIT PRESERVE ROWS`
	res, err = tx.ExecContext(ctx, sqlStatement)
	result = SQLResult(result, res)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
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
	fmt.Println("inserting into temp table")
	sqlStatement = `insert into _tmp_testtype (now, id, tag) values ( :now , :id , :tag ) `
	res, err = tx.NamedExecContext(ctx, sqlStatement, data)
	result = SQLResult(result, res)
	if err != nil {
		return result, err
	}
	fmt.Println("about to commit")
	// if err := tx.Commit(); err != nil {
	// 	if err := tx.Rollback(); err != nil {
	// 		return nil, err
	// 	}
	// 	return nil, err
	// }
	// sqlStatement = `
	// with cte_row_numbers as (
	// 	select
	// 	rn rownumber()
	// )

	// `
	// res, err = tx.NamedExecContext(ctx, sqlStatement, data)
	// result = SQLResult(result, res)
	// if err != nil {
	// 	return result, err
	// }
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}

	return result, nil
}

// insert non-duplicate rows, based on entire row, ordered by extract date time
// suitable for source that returns only the most recent version of some piece of data:
// eg staffmember doesn't include a log, just the most recent version
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

var TemplateForFullAppend string // insert all data, regardless of duplicates
var TemplateForNewAppend string  // insert all non-duplicates, either by primary keys or entire rows

var TemplateForFullReplace string       // completely replace the table, basically, truncate then insert - snapshots
var TemplateForReplaceOnConflict string // replace rows on conflict - snapshots
