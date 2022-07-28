package qgenda

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StaffMember struct {
	RawMessage               *string         `json:"-" db:"_raw_message" primarykey:"table" idtype:"group"`
	ExtractDateTime          *Time           `json:"_extract_date_time,omitempty" db:"_extract_date_time" primarykey:"table" idtype:"order"`
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

func (s StaffMembers) Extract(ctx context.Context, c *Client, rqf *RequestQueryFields) (StaffMembers, error) {
	req := NewStaffMemberRequest(rqf)
	sms := StaffMembers{}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &sms); err != nil {
		return nil, err
	}
	// user response header for extract time
	// fmt.Println(resp.Header.Get("date"))
	et, err := ParseTime(resp.Header.Get("date"))
	if err != nil {
		return nil, err
	}
	for i, _ := range sms {
		sms[i].ExtractDateTime = &et
	}
	log.Printf("staffmembers: extractdatetime: %s\ttotal: %d", et, len(sms))
	return sms, nil

}

func (s StaffMembers) Process() (StaffMembers, error) {
	for i, _ := range s {
		if err := s[i].Process(); err != nil {
			return nil, err
		}
	}
	return s, nil
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

func (s *StaffMember) UnmarshalJSON(b []byte) error {
	// custom unmarshaling is purely to assign a copy of the json representation in .RawMessage
	// idea from:
	// https://biscuit.ninja/posts/go-avoid-an-infitine-loop-with-custom-json-unmarshallers/
	type SM StaffMember
	var sm SM

	if err := json.Unmarshal(b, &sm); err != nil {
		return err
	}
	smm := StaffMember(sm)
	// we need to omit ExtractDateTime from StaffMember.RawMessage
	// the easiest way to do that is to set it to nil and remarshal
	if err := (&smm).SetRawMessage(); err != nil {
		return err
	}
	// var smrm SM
	// smrm = sm
	// // if err := json.Unmarshal(b, &smrm); err != nil {
	// // 	return err
	// // }
	// smrm.RawMessage = nil
	// smrm.ExtractDateTime = nil
	// smrmBytes, err := json.Marshal(smrm)
	// if err != nil {
	// 	return err
	// }

	// var bb bytes.Buffer
	// if err := json.Compact(&bb, smrmBytes); err != nil {
	// 	return err
	// }
	// rawMessage := bb.String()
	// sm.RawMessage = &rawMessage

	*s = smm

	return nil
}

// LoadFile is used to import any cached files. Because ExtractDateTime isn't actually part of
// StaffMember, it will also use the file's modTime as a proxy for ExtractDateTime.
func (s *StaffMembers) LoadFile(filename string) error {
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}
	modTime := fi.ModTime()

	b, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	sm := []StaffMember{}
	if err := json.Unmarshal(b, &sm); err != nil {
		log.Println(err)
	}
	for i, v := range sm {
		if v.ExtractDateTime == nil {
			proxyExtractDateTime := NewTime(modTime)
			v.ExtractDateTime = &proxyExtractDateTime
		}
		sm[i] = v
	}
	*s = sm
	return nil
}

func (s StaffMembers) WriteFile(filename string) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	var bb bytes.Buffer
	if err := json.Compact(&bb, b); err != nil {
		return err
	}
	if err := os.WriteFile(filename, bb.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

// SetRawMessage creates a value copy of the receiver, sets RawMessage and ExtractDateTime to nil,
// does a compact marshaling, and assigns the result to RawMessage
func (s *StaffMember) SetRawMessage() error {
	sm := *s
	sm.RawMessage = nil
	sm.ExtractDateTime = nil
	b, err := json.Marshal(sm)
	if err != nil {
		return err
	}
	var bb bytes.Buffer
	if err := json.Compact(&bb, b); err != nil {
		return err
	}
	rawMessage := bb.String()
	s.RawMessage = &rawMessage
	return nil
}

func (p *StaffMember) Process() error {
	ProcessStruct(p)
	// manage empty metadata fields
	if err := p.SetRawMessage(); err != nil {
		return err
	}
	// if p.RawMessage == nil {
	// 	rm, err := json.Marshal(*p)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	var bb bytes.Buffer
	// 	if err := json.Compact(&bb, rm); err != nil {
	// 		return err
	// 	}
	// 	rawMessage := bb.String()
	// 	p.RawMessage = &rawMessage
	// }

	if p.ExtractDateTime == nil {
		et := NewTime(nil)
		p.ExtractDateTime = &et
	}

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

func (s StaffMember) CreatePGTable(ctx context.Context, tx *sqlx.Tx, schema, tablename string, temp bool, id string) (sql.Result, error) {

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

	template := `
	CREATE {{- if .Temporary }} TEMPORARY TABLE IF NOT EXISTS _tmp_{{- .UUID -}}_{{- .Name -}}
	{{ else }} TABLE IF NOT EXISTS {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }}
	{{- end }} (
	{{- $fields := pgincludefields .Fields -}}
	{{- range  $index, $field := $fields -}}
	{{- if ne $index 0 -}},{{- end }}
		{{ pgname $field }} {{ pgtype $field.Type }} {{ if $field.Unique }} unique {{ end -}} {{- if not $field.Nullable -}} not null {{- end }}
	{{- end -}}
	)
	`
	// {{- if not .Temporary }}
	// {{- $pk := .Constraints.primarykey -}}
	// {{ if $pk }},
	// PRIMARY KEY ( {{ $pk }} )
	// {{- end -}}
	// {{- end }}

	table := StructToTable(StaffMember{}, basetable, schema, temp, id, nil, nil)
	// sqlStatement := PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement := PGTableStatement(table, template, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprintf("%stag", basetable)
	table = StructToTable(FlatStaffTag{}, tablename, schema, temp, id, nil, nil)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, template, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprintf("%sttcmtag", basetable)
	table = StructToTable(FlatStaffTag{}, tablename, schema, temp, id, nil, nil)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, template, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprintf("%sskillset", basetable)
	table = StructToTable(StaffSkillset{}, tablename, schema, temp, id, nil, nil)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, template, nil)
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

func (s StaffMember) SQLStatementInsertPGRowsTemp() string {

	return ""
}

var pgInsertSMRowsTpl = `
INSERT INTO 
{{- if .Temporary }}  _tmp_{{- .Name -}} 
{{- else }} {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name -}}
{{- end }} (
	{{- $fields := pgincludefields .Fields -}}
	{{- range  $index, $field := $fields -}}
	{{- if ne $index 0 -}},{{- end }}
		{{ pgname $field }}
	{{- end }} 
	) VALUES (
	{{- range  $index, $field := .Fields -}}
	{{- if ne $index 0 -}},{{- end }}
		:{{ pgname $field }}
	{{- end }}	
)
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
	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	sqlResult, err := s[0].CreatePGTable(ctx, tx, schema, tablename, false, id)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}

	// need temp tables
	// using _tmp_uuidstring_tablename to avoid any contamination when using multiple temp tables
	sqlResult, err = s[0].CreatePGTable(ctx, tx, schema, tablename, true, id)
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
	// temp table
	// table := StructToTable(s[0], basetable, schema, true, nil, nil)
	table := StructToTable(s[0], basetable, schema, true, id, nil, nil)
	// table.UUID = id
	// sqlStatement := PGTableStatement(table, PGInsertRowsDevTpl, nil)
	template := `
	INSERT INTO 
	{{- if .Temporary }}  _tmp_{{- .UUID -}}_{{- .Name -}} 
	{{- else }} {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name -}}
	{{- end }} (
		{{- $fields := pgincludefields .Fields -}}
		{{- range  $index, $field := $fields -}}
		{{- if ne $index 0 -}},{{- end }}
			{{ pgname $field }}
		{{- end }} 
		) VALUES (
		{{- range  $index, $field := pgincludefields .Fields -}}
		{{- if ne $index 0 -}},{{- end }}
			:{{ pgname $field }}
		{{- end }}	
	)
	`
	// {{- if not .Temporary }}
	// ON CONFLICT (
	// 	{{- $primarykey := join .PrimaryKey  ", " -}}
	// 	{{ if ne $primarykey "" }}
	// 	{{ $primarykey }}
	// 	{{ else }}
	// 	{{- range  $index, $field := .Fields -}}
	// 	{{- if ne $index 0 -}},{{- end }}
	// 		{{ pgname $field }}
	// 	{{- end -}}
	// 	{{- end }}
	// 	) DO NOTHING
	// {{ end }}

	sqlStatement := PGTableStatement(table, template, nil)
	fmt.Sprintln(sqlStatement)
	sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, s)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	// update table
	table.Temporary = false
	template = `
	with cte_most_recent as (
		select distinct on (s.staffkey)
		s._raw_message,
		s.staffkey
		from {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }} s
		order by s.staffkey, s._extract_date_time desc nulls last
	), cte_new as (
		select distinct * from _tmp_{{- .UUID -}}_{{- .Name }}
	), cte_updates as (
		select
		cn.*
		from cte_new cn
		where not exists (
			select 1
			from cte_most_recent cmr 
			where cmr._raw_message = cn._raw_message
			and   cmr.staffkey     = cn.staffkey
		)
	) insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (
		select cu.* from cte_updates cu
	)
	`

	// with cte_partitioned_row_numbers as (
	// 	select * ,
	// 	row_number() over ( partition by {{ fieldswithtagvalue .Fields "idtype" "group" | pgnames | joinss " , " }} order by {{ fieldswithtagvalue .Fields "idtype" "order" | pgnames | joinss " desc , " }} desc ) rn
	// 	FROM {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }}
	// ), cte_last_inserted_rows as (
	// 	select * from cte_partitioned_row_numbers cprn where cprn.rn = 1
	// ), cte_new_rows as (
	// 	select distinct * from _tmp_{{- .Name }}
	// ), cte_anti_joined as (
	// 	select
	// 	cnr.*
	// 	from cte_new_rows cnr
	// 	where not exists (
	// 		select 1
	// 		from cte_last_inserted_rows clir where
	// 		{{- $pk := .Fields.WithTagValue "primarykey" "table" -}}
	// 		{{- $joinfields := $pk.WithoutTagValue "idtype" "order" | pgincludefields | pgnames -}}
	// 		{{- range  $index, $field := $joinfields  }}
	// 		{{ if ne $index 0 }} and {{ end -}} clir.{{ $field }} = cnr.{{ $field }}
	// 		{{- end }}
	// 	)
	// ) insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (
	// 	select caj.* from cte_anti_joined caj
	// )
	// sqlStatement = PGTableStatement(table, PGInsertChangesOnlyDevTpl, nil)
	sqlStatement = PGTableStatement(table, template, nil)
	fmt.Println(sqlStatement)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	err = tx.Commit()
	fmt.Println("I'm Committed!")
	return res, err
}

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
