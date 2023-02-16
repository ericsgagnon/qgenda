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
	"sort"
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
	Profiles                 []StaffProfile  `json:"Profiles,omitempty"`
}

type StaffMembers []StaffMember

func (s StaffMembers) Extract(ctx context.Context, c *Client, rqf *RequestConfig) (StaffMembers, error) {
	req := NewStaffMemberRequest(rqf)
	sms := []StaffMember{}

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
	ExtractDateTime *Time   `json:"-" db:"_extract_date_time" nullable:"false"`
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
	ExtractDateTime   *Time   `json:"-" db:"_extract_date_time" nullable:"false"`
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

type StaffProfile struct {
	// RawMessage      *string
	ExtractDateTime *Time   `json:"-" db:"_extract_date_time" nullable:"false"`
	StaffKey        *string `json:"-" nullable:"false"`
	Name            *string `json:"Name,omitempty"`
	ProfileKey      *string `json:"ProfileKey,omitempty"`
	IsViewable      *bool   `json:"IsViewable,omitempty"`
	IsSchedulable   *bool   `json:"IsSchedulable,omitempty"`
}

func DefaultStaffMemberRequestConfig(rqf *RequestConfig) *RequestConfig {
	if rqf == nil {
		rqf = &RequestConfig{}
	}
	if rqf.Includes == nil {
		rqf.SetIncludes("Tags,TTCMTags,Skillset,Profiles")
	}
	return rqf
}

func NewStaffMemberRequest(rqf *RequestConfig) *Request {
	requestPath := "staffmember"
	queryFields := []string{
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	rqf = DefaultStaffMemberRequestConfig(rqf)

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
	// s = sm
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

	if p.ExtractDateTime == nil {
		et := NewTime(nil)
		p.ExtractDateTime = &et
	}

	if p.Tags != nil {
		for i, v := range p.Tags {
			if err := Process(&v); err != nil {
				return err
			}
			p.Tags[i].ExtractDateTime = p.ExtractDateTime
			p.Tags[i].StaffKey = p.StaffKey
			sort.Slice(p.Tags[i].Tags, func(j, k int) bool { return *p.Tags[i].Tags[j].Key < *p.Tags[i].Tags[k].Key })
		}
		sort.Slice(p.Tags, func(i, j int) bool { return *p.Tags[i].CategoryKey < *p.Tags[j].CategoryKey })
	}
	if p.TTCMTags != nil {
		for i, v := range p.TTCMTags {
			if err := Process(&v); err != nil {
				return err
			}
			p.TTCMTags[i].ExtractDateTime = p.ExtractDateTime
			p.TTCMTags[i].StaffKey = p.StaffKey
			sort.Slice(p.TTCMTags[i].Tags, func(j, k int) bool { return *p.TTCMTags[i].Tags[j].Key < *p.TTCMTags[i].Tags[k].Key })
		}
		sort.Slice(p.TTCMTags, func(i, j int) bool { return *p.TTCMTags[i].CategoryKey < *p.TTCMTags[j].CategoryKey })

	}

	if p.Skillset != nil {
		for i, v := range p.Skillset {
			if err := Process(&v); err != nil {
				return err
			}
			p.Skillset[i].ExtractDateTime = p.ExtractDateTime
			p.Skillset[i].StaffKey = p.StaffKey
			// sort.Slice(p.TTCMTags[i].Tags, func(j, k int) bool { return *p.TTCMTags[i].Tags[j].Key < *p.TTCMTags[i].Tags[k].Key })
		}
		sort.Slice(p.Skillset, func(i, j int) bool { return *p.Skillset[i].TaskName < *p.Skillset[j].TaskName })
	}

	if p.Profiles != nil {
		for i, v := range p.Profiles {
			if err := Process(&v); err != nil {
				return err
			}
			p.Profiles[i].ExtractDateTime = p.ExtractDateTime
			p.Profiles[i].StaffKey = p.StaffKey
		}
		sort.Slice(p.Profiles, func(i, j int) bool { return *p.Profiles[i].ProfileKey < *p.Profiles[j].ProfileKey })
	}
	if err := p.SetRawMessage(); err != nil {
		return err
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

	table := StructToTable(StaffMember{}, basetable, schema, temp, id, nil, nil, nil)
	// sqlStatement := PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement := PGTableStatement(table, template, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprintf("%stag", basetable)
	table = StructToTable(FlatStaffTag{}, tablename, schema, temp, id, nil, nil, nil)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, template, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprintf("%sttcmtag", basetable)
	table = StructToTable(FlatStaffTag{}, tablename, schema, temp, id, nil, nil, nil)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, template, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprintf("%sskillset", basetable)
	table = StructToTable(StaffSkillset{}, tablename, schema, temp, id, nil, nil, nil)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, template, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprintf("%sprofile", basetable)
	table = StructToTable(StaffProfile{}, tablename, schema, temp, id, nil, nil, nil)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, template, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s StaffMembers) CreatePGTable(ctx context.Context, tx *sqlx.Tx, schema, tablename string, temporary bool, id string) (sql.Result, error) {
	sm := *new(StaffMember)
	return sm.CreatePGTable(ctx, tx, schema, tablename, temporary, id)
}

// InsertToPG is a StaffMembers method (as opposed to a StaffMember method) to enable bulk loads and avoid
// large numbers of single row transactions
func (s StaffMembers) InsertToPG(ctx context.Context, db *sqlx.DB, schema, tablename string) (sql.Result, error) {

	if len(s) < 1 {
		return nil, fmt.Errorf("StaffMembers.InsertToPG: length of %T < 1, nothing to do", s)
	}

	// postgres has a 65535 'parameter' limit, there is an 'unnest' work around, but for now we're just going to chunk it
	// chunkSize := 65535 / reflect.ValueOf(StaffMember{}).NumField()
	// this shouldn't be an issue for StaffMember, only implement if you need to

	tx := db.MustBegin()
	var res Result
	// make sure the tables exist
	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	sqlResult, err := s.CreatePGTable(ctx, tx, schema, tablename, false, id)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}

	// need temp tables
	// using _tmp_uuidstring_tablename to avoid any contamination when using multiple temp tables
	sqlResult, err = s.CreatePGTable(ctx, tx, schema, tablename, true, id)
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
	table := StructToTable(s[0], basetable, schema, true, id, nil, nil, nil)
	// table.UUID = id
	// sqlStatement := PGTableStatement(table, PGInsertRowsDevTpl, nil)
	insertTpl := `
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

	sqlStatement := PGTableStatement(table, insertTpl, nil)
	// fmt.Sprintln(sqlStatement)
	sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, s)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	// update table
	table.Temporary = false
	updateTpl := `
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
	sqlStatement = PGTableStatement(table, updateTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	childUpdateTpl := `
	with cte_parent_ref as (
		select distinct 
		p.staffkey,
		p._extract_date_time
		from {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Parent }} p
	), cte_current as (
		select distinct
		c.*
		from {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }} c
		order by c.staffkey, c._extract_date_time desc nulls last
	), cte_new as (
		select distinct * from _tmp_{{- .UUID -}}_{{- .Name }}
	), cte_updates as (
		select
		cn.*
		from cte_new cn
		-- inner join on parent refs
		inner join cte_parent_ref cpr on (
			    cpr.staffkey = cn.staffkey
			and cpr._extract_date_time = cn._extract_date_time
		)
		-- exclude on duplicates
		where not exists (
			select 1
			from cte_current cc where
			{{ $fields := pgincludefields .Fields -}}
			{{- range  $index, $field := $fields -}}
			{{ if ne $index 0 -}} and {{- end }} cc.{{ pgname $field }} = cn.{{ pgname $field }}
			{{ end -}}
		)
	) insert into {{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} (
		select cu.* from cte_updates cu
		where not exists (
			select 1 from 
			{{ if ne .Schema "" -}}{{ .Schema -}}.{{- end -}}{{- .Name }} dst where
			{{ $fields := pgincludefields .Fields -}}
			{{- range  $index, $field := $fields -}}
			{{ if ne $index 0 }}and {{ end -}} dst.{{ pgname $field }} = cu.{{ pgname $field }}
			{{ end -}}
		)
	)
	`

	// tag
	var tags FlatStaffTags
	for _, v := range s {
		fts := v.FlatTags()
		tags = append(tags, fts...)
	}
	tablename = fmt.Sprintf("%stag", basetable)
	table = StructToTable(FlatStaffTag{}, tablename, schema, true, id, nil, nil, nil)
	sqlStatement = PGTableStatement(table, insertTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, tags)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}
	sqlStatement = PGTableStatement(table, childUpdateTpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	// ttcmtag
	var ttcmtags FlatStaffTags
	for _, v := range s {
		fts := v.FlatTTCMTags()
		ttcmtags = append(ttcmtags, fts...)
	}
	tablename = fmt.Sprintf("%sttcmtag", basetable)
	table = StructToTable(FlatStaffTag{}, tablename, schema, true, id, nil, nil, nil)
	sqlStatement = PGTableStatement(table, insertTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, ttcmtags)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}
	sqlStatement = PGTableStatement(table, childUpdateTpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	// skillset
	var skillsets StaffSkillsets
	for _, v := range s {
		ss := v.Skillsets()
		skillsets = append(skillsets, ss...)
	}
	tablename = fmt.Sprintf("%sskillset", basetable)
	table = StructToTable(StaffSkillset{}, tablename, schema, true, id, nil, nil, nil)
	sqlStatement = PGTableStatement(table, insertTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, skillsets)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}
	sqlStatement = PGTableStatement(table, childUpdateTpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	// profile
	var profiles []StaffProfile
	for _, v := range s {
		p := v.Profiles
		profiles = append(profiles, p...)
	}
	tablename = fmt.Sprintf("%sprofile", basetable)
	table = StructToTable(StaffProfile{}, tablename, schema, true, id, nil, nil, nil)
	sqlStatement = PGTableStatement(table, insertTpl, nil)
	// fmt.Println(sqlStatement)
	sqlResult, err = tx.NamedExecContext(ctx, sqlStatement, profiles)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}
	sqlStatement = PGTableStatement(table, childUpdateTpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	err = tx.Commit()
	// fmt.Println("I'm Committed!")
	return res, err
}

func (s StaffMembers) EPL(ctx context.Context, c *Client, rqf *RequestConfig,
	db *sqlx.DB, schema, table string, newRowsOnly bool) (sql.Result, error) {
	var res Result

	sqlResult, err := PGCreateSchema(ctx, db, s, schema, table)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}
	if err != nil {
		return res, err
	}

	// sqlResult, err = s.CreatePGTable(ctx, db, schema, table)
	// if sqlResult != nil {
	// 	res = SQLResult(res, sqlResult)
	// }
	// if err != nil {
	// 	return res, err
	// }
	s, err = s.Extract(ctx, c, rqf)
	if err != nil {
		return res, err
	}
	s, err = s.Process()
	if err != nil {
		return res, err
	}
	sqlResult, err = s.InsertToPG(ctx, db, schema, table)
	if sqlResult != nil {
		res = SQLResult(res, sqlResult)
	}

	return res, err
}
