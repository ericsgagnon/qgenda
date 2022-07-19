package qgenda

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type StaffMember struct {
	RawMessage               *string         `json:"-" db:"_raw_message"`
	ExtractDateTime          *Time           `json:"_extract_date_time" db:"_extract_date_time" primarykey:"table" idtype:"order"`
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
	StaffKey                 *string         `json:"StaffKey,omitempty" primarykey:"table" idtype:"group"`
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
	Tags                     []StaffTag      `json:"Tags,omitempty" dbTable:"stafftag" includeFields:"primarykey"`
	TTCMTags                 []StaffTag      `json:"TTCMTags,omitempty" dbTable:"staffttcmtag" includeFields:"primarykey"`
	Skillset                 []StaffSkillset `json:"Skillset,omitempty" dbTable:"staffskillset" includeFields:"primarykey"`
	Profiles                 []Profile       `json:"Profiles,omitempty"`
}

type StaffMembers []StaffMember

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

// StaffTagDB is basically TagCategory with the StaffMember.StaffKey and ExtractDateTime added for dependent tables in DB's
type FlatStaffTag struct {
	ExtractDateTime *Time   `json:"-" nullable:"false"`
	StaffKey        *string `json:"-" nullable:"false"`
	CategoryKey     *int64  `json:"CategoryKey" nullable:"false"`
	CategoryName    *string `json:"CategoryName" nullable:"false"`
	TagKey          *int64  `json:"Key" db:"tagkey" nullable:"false"`
	TagName         *string `json:"Name" db:"tagname" nullable:"false"`
	// LastModifiedDateUTC *Time   `json:"LastModifiedDateUTC,omitempty" nullable:"false"`
}

type FlatStaffTags []FlatStaffTag

func (s StaffMember) FlatTags() FlatStaffTags {
	fsts := []FlatStaffTag{}
	if len(s.Tags) == 0 {
		return fsts
	}
	for _, t := range s.Tags {
		for _, tt := range t.Tags {
			fst := FlatStaffTag{
				ExtractDateTime: s.ExtractDateTime,
				StaffKey:        s.StaffKey,
				CategoryKey:     t.CategoryKey,
				CategoryName:    t.CategoryName,
				TagKey:          tt.Key,
				TagName:         tt.Name,
			}
			fsts = append(fsts, fst)

		}
	}
	return fsts
}

func (s StaffMember) FlatTTCMTags() FlatStaffTags {
	fsts := []FlatStaffTag{}
	if len(s.TTCMTags) == 0 {
		return fsts
	}
	for _, t := range s.TTCMTags {
		for _, tt := range t.Tags {
			fst := FlatStaffTag{
				ExtractDateTime: s.ExtractDateTime,
				StaffKey:        s.StaffKey,
				CategoryKey:     t.CategoryKey,
				CategoryName:    t.CategoryName,
				TagKey:          tt.Key,
				TagName:         tt.Name,
			}
			fsts = append(fsts, fst)

		}
	}
	return fsts
}

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

type StaffSkillsets []StaffSkillset

func (s StaffMember) Skillsets() StaffSkillsets {
	sks := []StaffSkillset{}
	if len(s.Skillset) == 0 {
		return sks
	}
	for _, sk := range s.Skillset {
		ski := StaffSkillset{
			ExtractDateTime:   sk.ExtractDateTime,
			StaffKey:          sk.StaffKey,
			StaffFirstName:    sk.StaffFirstName,
			StaffLastName:     sk.StaffLastName,
			StaffAbbreviation: sk.StaffAbbreviation,
			StaffId:           sk.StaffId,
			TaskName:          sk.TaskName,
			TaskAbbreviation:  sk.TaskAbbreviation,
			TaskId:            sk.TaskId,
			IsSkilledMon:      sk.IsSkilledMon,
			MonOccurrence:     sk.MonOccurrence,
			IsSkilledTue:      sk.IsSkilledTue,
			TueOccurrence:     sk.TueOccurrence,
			IsSkilledWed:      sk.IsSkilledWed,
			WedOccurrence:     sk.WedOccurrence,
			IsSkilledThu:      sk.IsSkilledThu,
			ThuOccurrence:     sk.ThuOccurrence,
			IsSkilledFri:      sk.IsSkilledFri,
			FriOccurrence:     sk.FriOccurrence,
			IsSkilledSat:      sk.IsSkilledSat,
			SatOccurrence:     sk.SatOccurrence,
			IsSkilledSun:      sk.IsSkilledSun,
			SunOccurrence:     sk.SunOccurrence,
		}
		sks = append(sks, ski)
	}
	return sks
}

// func (s FlatStaffTagRelation) ToTable(name string, schema string, temp bool, constraints map[string]string, tags map[string][]string) Table {
// 	ftr := FlatTagRelation{}

// 	table := StructToTable(ftr, name, schema, temp, constraints, tags)
// 	fields := StructToFields()
// }

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

	if p.Skillset != nil {
		for _, v := range p.Skillset {
			// (&p.TTCMTags[i]).Process()
			if err := Process(&v); err != nil {
				return err
			}
		}
	}

	if p.Profiles != nil {
		for _, v := range p.Profiles {
			// (&p.TTCMTags[i]).Process()
			if err := Process(&v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s StaffMember) CreatePGTable(ctx context.Context, tx *sqlx.Tx, schema, tablename string, temp bool) (sql.Result, error) {

	basetable := "staffmember"
	if tablename != "" {
		basetable = tablename
	}
	var res Result

	sqlResult, err := tx.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema))
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	table := StructToTable(Schedule{}, basetable, schema, temp, nil, nil)
	sqlStatement := PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprintf("%stag", basetable)
	table = StructToTable(FlatStaffTag{}, tablename, schema, temp, nil, nil)
	sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprintf("%sttcmtag", basetable)
	table = StructToTable(FlatStaffTag{}, tablename, schema, temp, nil, nil)
	sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprintf("%sskillset", basetable)
	table = StructToTable(StaffSkillset{}, tablename, schema, temp, nil, nil)
	sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	// staffprofile doesn't appear to be implemented
	// tablename = fmt.Sprintf("%sprofile", basetable)
	// table = StructToTable(Profile{}, tablename, schema, temp, nil, nil)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	// result, err = tx.ExecContext(ctx, sqlStatement)
	// if err != nil {
	// 	return result, err
	// }
	// sqlResult = SQLResult(result)

	return res, nil
}

var smInsertRowsPGSQL = `

`

// InsertToPG is a StaffMembers method (as opposed to a StaffMember method) to enable bulk loads and avoid
// large numbers of single row transactions
func (s StaffMembers) InsertToPG(ctx context.Context, db *sqlx.DB, schema, tablename string) (sql.Result, error) {

	if len(s) < 1 {
		return nil, fmt.Errorf("StaffMembers.InsertToPG: length of %T < 1, nothing to do", s)
	}

	tx := db.MustBegin()
	var res Result

	// make sure the tables exist
	sqlResult, err := s[0].CreatePGTable(ctx, tx, schema, tablename, false)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}

	// need temp tables
	sqlResult, err = s[0].CreatePGTable(ctx, tx, schema, tablename, true)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}

	// primary table
	basetable := "staffmember"
	if tablename != "" {
		basetable = tablename
	}
	table := StructToTable(s[0], basetable, schema, true, nil, nil)
	sqlStatement := PGTableStatement(table, PGInsertRowsDevTpl, nil)
	fmt.Println(sqlStatement)
	sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, s)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	err = tx.Commit()
	return res, err
}

// func (s StaffMember) CreatePGTable(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {

// 	result, err := pgCreateTable(ctx, db, s, schema, table)
// 	if err != nil {
// 		return result, err
// 	}
// 	sqlResult := SQLResult(result)

// 	tablename := fmt.Sprintf("%stag", table)
// 	value := []StaffTag{}
// 	result, err = PGCreateTable(ctx, db, value, schema, tablename)
// 	if err != nil {
// 		return result, err
// 	}
// 	sqlResult = SQLResult(sqlResult, result)

// 	tablename = fmt.Sprintf("%sttcmtag", table)
// 	result, err = PGCreateTable(ctx, db, []StaffTag{}, schema, tablename)
// 	if err != nil {
// 		return result, err
// 	}
// 	sqlResult = SQLResult(sqlResult, result)

// 	tablename = fmt.Sprintf("%sskillset", table)
// 	result, err = PGCreateTable(ctx, db, []StaffSkillset{}, schema, tablename)
// 	if err != nil {
// 		return result, err
// 	}
// 	sqlResult = SQLResult(sqlResult, result)

// 	tablename = fmt.Sprintf("%sprofile", table)

// 	result, err = PGCreateTable(ctx, db, []Profile{}, schema, tablename)
// 	// result, err := PGCreateTable(ctx, db, value, schema, tablename)
// 	if err != nil {
// 		return result, err
// 	}
// 	sqlResult = SQLResult(sqlResult, result)
// 	return sqlResult, nil
// 	// Profiles                 []Profile       `json:"Profiles,omitempty"`

// 	// create temporary table if not exists _tmp_testtype (
// 	// 	now timestamp with time zone,
// 	// 	id bigint,
// 	// 	tag text
// 	// 	)
// 	// ;
// 	// insert into _tmp_testtype values ( now(), (random() * 100 )::bigint , ( random() * 100 )::text);
// 	// insert into _tmp_testtype select * from testtype;
// 	// -- select * from pg_temp._tmp_testtype;

// 	// -- select count(*) from testtype;
// 	// with cte_row_numbers as (

// 	// select
// 	// row_number() over (partition by tt.id order by tt.now desc) rn,
// 	// tt.now,
// 	// tt.id,
// 	// tt.tag
// 	// -- from testtype tt
// 	// from _tmp_testtype tt
// 	// ),
// 	// cte_most_recent as (
// 	// select * from cte_row_numbers where rn = 1
// 	// ),
// 	// cte_anti_joined as (
// 	// select
// 	// r.now,
// 	// r.id,
// 	// r.tag
// 	// from cte_most_recent r
// 	// where not exists (
// 	// select 1
// 	// from testtype tt where
// 	// 	r.id = tt.id
// 	// and r.tag = tt.tag
// 	// )
// 	// )
// 	// insert into testtype
// 	// select * from cte_anti_joined

// 	// ;

// 	// select * from testtype tt order by tt.now desc;

// }

// func (s StaffMembers) InsertToPG(pgc *PGClient, ctx context.Context) (sql.Result, error) {
// 	txOptions := sql.TxOptions{
// 		Isolation: 0,
// 		ReadOnly:  false,
// 	}
// 	fmt.Print(txOptions)
// 	table := StructToTable(s[0], "staffmember", pgc.Config.Schema, false, nil, nil)
// 	var sqlResult sql.Result
// 	chunkSize := 65535 / reflect.ValueOf(s[0]).NumField()
// 	var res Result
// 	for i := 0; i < len(s); i = i + chunkSize {
// 		j := i + chunkSize
// 		if j > len(s) {
// 			j = len(s)
// 		}
// 		values := s[i:j]
// 		pgc.Tx = pgc.MustBegin()

// 		// need a temp table for this
// 		result, err := s[i].CreatePGTable(pgc, ctx, true)
// 		sqlResult = SQLResult(sqlResult, result)
// 		if err != nil {
// 			return sqlResult, err
// 		}
// 		tmp := table
// 		tmp.Temporary = true
// 		// pgc.InsertRows(ctx, tmp, values, true) // this doesn't work! silly 'generics'
// 		// result, err = PGInsertTx(ctx, pgc.Tx, tmp, "", values)
// 		sqlStatement := PGTableStatement(tmp, PGInsertRowsDevTpl, nil)
// 		result, err = pgc.Tx.NamedExecContext(ctx, sqlStatement, values)
// 		sqlResult = SQLResult(sqlResult, result)
// 		if err != nil {
// 			return sqlResult, err
// 		}

// 		// insert all to temp table
// 		sqlResult, err := pgc.Tx.NamedExecContext(
// 			ctx,
// 			// PGInsertStatement(value[0], schema, table),
// 			PGStatement(values[0], table.Schema, table.Name, pgInsertTpl),
// 			s[i:j],
// 		)
// 		res = SQLResult(res, sqlResult)
// 		if err != nil {
// 			return res, err
// 		}

// 		// for _, field := range table.Fields {
// 		// 	switch field.Kind {
// 		// 	case "slice", "map":
// 		// 		continue
// 		// 		// v := reflect.ValueOf(value[i]).FieldByName(field.Name)
// 		// 		// PGInsertRowsDev(ctx, tx, table, v.InterfaceData())
// 		// 		// name := strings.ToLower(fmt.Sprintf("%s%s", table.Name, reflect.ValueOf()))
// 		// 		// table := StructToTable()
// 		// 		// PGInsertRowsDev(ctx, tx, table, )
// 		// 	default:
// 		// 		continue
// 		// 	}
// 		// }

// 	}

// 	// tx, err := db.BeginTxx(ctx, &txOptions)
// 	// if err != nil {
// 	// 	tx.Rollback()
// 	// 	return nil, err
// 	// }
// 	// // create table
// 	// result, err := tx.Exec(``)
// 	// if err != nil {
// 	// 	tx.Rollback()
// 	// 	return result, err
// 	// }
// 	// // create temp table
// 	// // insert rows

// 	// // create child tables
// 	// // create child

// 	// if err := tx.Commit(); err != nil {
// 	// 	tx.Rollback()
// 	// }

// 	return nil, nil

// 	// there is no combination of fields for a good primary key
// 	// need to only insert new rows if there is a change in at least one (non-metadata) field
// 	// if _extract_date_time is only different field, don't insert

// 	// insert into staffmember ...
// 	// select

// 	// start transaction
// 	// create temp table
// 	// insert into temp table
// 	// use cte to join with current data and insert changes to target table
// 	// if err, rollback, else commit
// 	// not sure how to make sure no dangling temp tables

// }

// func (s StaffMembers) InsertToPGDev(pgc *PGClient, ctx context.Context) (sql.Result, error) {
// 	txOptions := sql.TxOptions{
// 		Isolation: 0,
// 		ReadOnly:  false,
// 	}
// 	fmt.Print(txOptions)
// 	table := StructToTable(s[0], "staffmember", pgc.Config.Schema, false, nil, nil)
// 	var sqlResult sql.Result
// 	chunkSize := 65535 / reflect.ValueOf(s[0]).NumField()
// 	var res Result
// 	for i := 0; i < len(s); i = i + chunkSize {
// 		j := i + chunkSize
// 		if j > len(s) {
// 			j = len(s)
// 		}
// 		values := s[i:j]
// 		pgc.Tx = pgc.MustBegin()

// 		// need a temp table for this
// 		result, err := s[i].CreatePGTable(pgc, ctx, true)
// 		sqlResult = SQLResult(sqlResult, result)
// 		if err != nil {
// 			return sqlResult, err
// 		}
// 		tmp := table
// 		tmp.Temporary = true
// 		// pgc.InsertRows(ctx, tmp, values, true) // this doesn't work! silly 'generics'
// 		// result, err = PGInsertTx(ctx, pgc.Tx, tmp, "", values)
// 		sqlStatement := PGTableStatement(tmp, PGInsertRowsDevTpl, nil)
// 		result, err = pgc.Tx.NamedExecContext(ctx, sqlStatement, values)
// 		sqlResult = SQLResult(sqlResult, result)
// 		if err != nil {
// 			return sqlResult, err
// 		}

// 		// insert all to temp table
// 		sqlResult, err := pgc.Tx.NamedExecContext(
// 			ctx,
// 			// PGInsertStatement(value[0], schema, table),
// 			PGStatement(values[0], table.Schema, table.Name, pgInsertTpl),
// 			s[i:j],
// 		)
// 		res = SQLResult(res, sqlResult)
// 		if err != nil {
// 			return res, err
// 		}

// 		// for _, field := range table.Fields {
// 		// 	switch field.Kind {
// 		// 	case "slice", "map":
// 		// 		continue
// 		// 		// v := reflect.ValueOf(value[i]).FieldByName(field.Name)
// 		// 		// PGInsertRowsDev(ctx, tx, table, v.InterfaceData())
// 		// 		// name := strings.ToLower(fmt.Sprintf("%s%s", table.Name, reflect.ValueOf()))
// 		// 		// table := StructToTable()
// 		// 		// PGInsertRowsDev(ctx, tx, table, )
// 		// 	default:
// 		// 		continue
// 		// 	}
// 		// }

// 	}

// 	// tx, err := db.BeginTxx(ctx, &txOptions)
// 	// if err != nil {
// 	// 	tx.Rollback()
// 	// 	return nil, err
// 	// }
// 	// // create table
// 	// result, err := tx.Exec(``)
// 	// if err != nil {
// 	// 	tx.Rollback()
// 	// 	return result, err
// 	// }
// 	// // create temp table
// 	// // insert rows

// 	// // create child tables
// 	// // create child

// 	// if err := tx.Commit(); err != nil {
// 	// 	tx.Rollback()
// 	// }

// 	return nil, nil

// 	// there is no combination of fields for a good primary key
// 	// need to only insert new rows if there is a change in at least one (non-metadata) field
// 	// if _extract_date_time is only different field, don't insert

// 	// insert into staffmember ...
// 	// select

// 	// start transaction
// 	// create temp table
// 	// insert into temp table
// 	// use cte to join with current data and insert changes to target table
// 	// if err, rollback, else commit
// 	// not sure how to make sure no dangling temp tables

// }

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
// primary key: _extract_date_time, _raw_message
// on conflict, do nothing
// no need to query, all data is full export, full import

var TemplateForNewAppend string // insert all non-duplicates, either by primary keys or entire rows
// use temp table to insert all data
// use cte's to anti-join last inserted rows based on primary key identifier, sorted by primary key sorter
// on conflict, do nothing

var TemplateForFullReplace string // completely replace the table, basically, truncate then insert - snapshots
// truncate table
// insert new rows
// never need to query db for qf's

// replace rows on conflict - snapshots
// insert
var TemplateForReplaceOnConflict string = `
INSERT INTO {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Table }} (
	{{ joinss  }}
) VALUES (

)
`

// full replace

// var pgInsertTpl = `
// INSERT INTO {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Table }} (
// {{- range  $index, $field := .Fields -}}
// {{- if ne $index 0 -}},{{- end }}
// 	{{ pgname $field }}
// {{- end }}
// ) VALUES (
// {{- range  $index, $field := .Fields -}}
// {{- if ne $index 0 -}},{{- end }}
// 	:{{ pgname $field }}
// {{- end }}
// ) ON CONFLICT (
// {{- $primarykey := join .PrimaryKey  ", " -}}
// {{ if ne $primarykey "" }}
// {{ $primarykey }}
// {{ else }}
// {{- range  $index, $field := .Fields -}}
// {{- if ne $index 0 -}},{{- end }}
// 	{{ pgname $field }}
// {{- end -}}
// {{- end }}
// ) DO NOTHING
// `
