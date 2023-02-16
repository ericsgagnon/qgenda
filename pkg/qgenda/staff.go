package qgenda

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/jmoiron/sqlx"
)

type Staff struct {
	// ------- metadata ------------------- //
	RawMessage       *string `json:"-" db:"_raw_message"`
	ProcessedMessage *string `json:"-" db:"_processed_message"` // RawMessage processed, omits 'message' metadata and 'noisy' fields (eg lastlogin)
	SourceQuery      *string `json:"_source_query,omitempty" db:"_source_query"`
	ExtractDateTime  *Time   `json:"_extract_date_time,omitempty" db:"_extract_date_time"`
	IDHash           *string `json:"-" primarykey:"table"  db:"_id_hash"` // hash of identifying fields: for staffmember, this is the processed message hash
	// ------------------------------------ //
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
	Tags                     []XStaffTags    `json:"Tags,omitempty" dbTable:"stafftag" includeFields:"primarykey"`
	TTCMTags                 []XStaffTags    `json:"TTCMTags,omitempty" dbTable:"staffttcmtag" includeFields:"primarykey"`
	Skillset                 []StaffSkill    `json:"Skillset,omitempty" dbTable:"staffskillset" includeFields:"primarykey"`
	Profiles                 []XStaffProfile `json:"Profiles,omitempty"`
}

func (s *Staff) UnmarshalJSON(b []byte) error {
	// alias technique to avoid infinite recursion
	type Alias Staff
	var a Alias

	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}

	var bb bytes.Buffer
	if err := json.Compact(&bb, b); err != nil {
		return err
	}
	rawMessage := bb.String()

	dest := Staff(a)
	dest.RawMessage = &rawMessage

	*s = dest
	return nil

}

func (s Staff) MarshalJSON() ([]byte, error) {
	type Alias Staff
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

func (s *Staff) Process() error {
	if err := ProcessStruct(s); err != nil {
		return fmt.Errorf("error processing %T:\t%q", s, err)
	}
	// TODO: Need to deal with IDHash, maybe multiple times???
	// Tags
	if len(s.Tags) > 0 {
		for i, _ := range s.Tags {
			s.Tags[i].ExtractDateTime = s.ExtractDateTime
			s.Tags[i].IDHash = s.IDHash
			s.Tags[i].StaffKey = s.StaffKey
			if err := s.Tags[i].Process(); err != nil {
				return err
			}
			sort.SliceStable(s.Tags, func(i, j int) bool {
				return *(s.Tags[i].CategoryKey) < *(s.Tags[j].CategoryKey)
			})
		}
	}

	// TTCMTags
	if len(s.TTCMTags) > 0 {
		for i, _ := range s.TTCMTags {
			s.TTCMTags[i].ExtractDateTime = s.ExtractDateTime
			s.TTCMTags[i].IDHash = s.IDHash
			s.TTCMTags[i].StaffKey = s.StaffKey
			if err := s.TTCMTags[i].Process(); err != nil {
				return err
			}
			sort.SliceStable(s.TTCMTags, func(i, j int) bool {
				return *(s.TTCMTags[i].CategoryKey) < *(s.TTCMTags[j].CategoryKey)
			})
		}
	}

	// Skillset
	if len(s.Skillset) > 0 {
		for i, _ := range s.Skillset {
			s.Skillset[i].ExtractDateTime = s.ExtractDateTime
			s.Skillset[i].IDHash = s.IDHash
			s.Skillset[i].StaffKey = s.StaffKey
			if err := s.Skillset[i].Process(); err != nil {
				return err
			}
			sort.SliceStable(s.Skillset, func(i, j int) bool {
				return *(s.Skillset[i].TaskName) < *(s.Skillset[j].TaskName)
			})
		}
	}

	// Profiles
	if len(s.Profiles) > 0 {
		for i, _ := range s.Profiles {
			s.Profiles[i].ExtractDateTime = s.ExtractDateTime
			s.Profiles[i].IDHash = s.IDHash
			s.Profiles[i].StaffKey = s.StaffKey
			if err := s.Profiles[i].Process(); err != nil {
				return err
			}
			sort.SliceStable(s.Profiles, func(i, j int) bool {
				return *(s.Profiles[i].ProfileKey) < *(s.Profiles[j].ProfileKey)
			})
		}
	}
	if err := s.SetMessage(); err != nil {
		return err
	}
	if err := s.SetIDHash(); err != nil {
		return err
	}

	for i, _ := range s.Tags {
		s.Tags[i].IDHash = s.IDHash
	}
	for i, _ := range s.TTCMTags {
		s.TTCMTags[i].IDHash = s.IDHash
	}
	for i, _ := range s.Skillset {
		s.Skillset[i].IDHash = s.IDHash
	}
	for i, _ := range s.Profiles {
		s.Profiles[i].IDHash = s.IDHash
	}

	return nil
}

func (s *Staff) SetMessage() error {
	// take a copy and strip metadata, for good measure
	ss := *s
	ss.RawMessage = nil
	ss.ProcessedMessage = nil
	// ss.MessageHash = nil

	b, err := json.Marshal(ss)
	// b, err := json.Marshal(s)
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

func (s *Staff) SetIDHash() error {
	if s.ProcessedMessage == nil {
		return fmt.Errorf("ProcessedMessage is empty, cannot hash")
	}

	// TODO copy s, remove metadata and noisy fields
	ss := *s
	ss.RawMessage = nil
	ss.ProcessedMessage = nil
	ss.SourceQuery = nil
	ss.ExtractDateTime = nil
	ss.IDHash = nil
	ss.UserLastLoginDateTimeUTC = nil
	ss.DailyUnitAverage = nil
	// id := map[string]any{}

	b, err := json.Marshal(ss)
	if err != nil {
		return err
	}
	var bb bytes.Buffer
	if err := json.Compact(&bb, b); err != nil {
		return err
	}

	// id := fmt.Sprint(*s.ScheduleKey, "-", (*s.LastModifiedDateUTC).Time.Format(time.RFC3339Nano))
	h := Hash(bb.String())
	s.IDHash = &h
	return nil
}

func DefaultStaffRequestConfig() *RequestConfig {
	requestPath := "staffmember"
	allowedFields := []string{
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	rc := NewRequestConfig(requestPath, allowedFields)
	rc.SetIncludes("Tags,TTCMTags,Skillset,Profiles")
	return rc
}

func NewStaffRequestConfig(rc *RequestConfig) *RequestConfig {
	return DefaultStaffRequestConfig().Merge(rc)

}

func NewStaffRequest(rc *RequestConfig) *Request {
	rc = NewStaffRequestConfig(rc)
	return NewRequest(rc)
}

func (s Staff) PGCreateTable(ctx context.Context, tx *sqlx.Tx, schema, tablename string, temporary bool, id string) (sql.Result, error) {
	basetable := "staff"
	if tablename != "" {
		basetable = tablename
	}
	var res Result

	if !temporary && schema != "" {
		sqlResult, err := tx.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema))
		res = SQLResult(res, sqlResult)
		if err != nil {
			return res, err
		}
	}

	tablename = basetable
	tpl := `{{- $table := . -}}
	CREATE {{- if .Temporary }} TEMPORARY TABLE IF NOT EXISTS _tmp_{{- .UUID -}}_{{- .Name -}}
	{{ else }} TABLE IF NOT EXISTS {{ .Schema -}}{{- if ne .Schema "" -}}.{{- end -}}{{- .Name }}
	{{- end }} (
		{{- $fields := pgincludefields .Fields -}}
	{{- range  $index, $field := $fields -}}
	{{- if ne $index 0 -}},{{- end }}
	{{ pgname $field }} {{ pgtype $field.Type }} {{ if not $table.Temporary }}{{ if $field.Unique }} unique {{ end -}} {{- if not $field.Nullable -}} not null {{- end }}{{- end }}
	{{- end -}}
	{{- if not .Temporary }}{{- if .PrimaryKey }}, 
		PRIMARY KEY ({{ .PrimaryKey }}){{- end -}}{{- end }}
	)	
	`
	table := StructToTable(Staff{}, tablename, schema, temporary, id, nil, nil, nil)
	sqlStatement := PGTableStatement(table, tpl, nil)
	// fmt.Println(sqlStatement)

	sqlResult, err := tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprint(basetable, "tag")
	table = StructToTable(XStaffTag{}, tablename, schema, temporary, id, nil, nil, nil)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, tpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprint(basetable, "ttcmtag")
	table = StructToTable(XStaffTag{}, tablename, schema, temporary, id, nil, nil, nil)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, tpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprint(basetable, "skill")
	table = StructToTable(StaffSkill{}, tablename, schema, temporary, id, nil, nil, nil)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, tpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	tablename = fmt.Sprint(basetable, "profile")
	table = StructToTable(XStaffProfile{}, tablename, schema, temporary, id, nil, nil, nil)
	// sqlStatement = PGTableStatement(table, PGCreateTableDevTpl, nil)
	sqlStatement = PGTableStatement(table, tpl, nil)
	sqlResult, err = tx.ExecContext(ctx, sqlStatement)
	res = SQLResult(res, sqlResult)
	if err != nil {
		return res, err
	}

	return nil, nil
}

// PGGetCDC returns an instance of staff with the max extractdatetime, but we don't actually use any
// constraints for CDC for this table, so it is for info only
func (s *Staff) PGGetCDC(ctx context.Context, db *sqlx.DB, schema, table string) (*Staff, error) {
	var ss *Staff
	if schema != "" {
		table = schema + "." + table
	}
	query := fmt.Sprintf("SELECT MAX ( _extract_date_time ) FROM %s ", table)
	if err := db.GetContext(ctx, ss, query); err != nil {
		return nil, err
	}
	return ss, nil
}

// ToRequestConfig returns the default staff request config, since we don't do CDC on this table
func (s Staff) ToRequestConfig() (*RequestConfig, error) {
	rc := DefaultStaffRequestConfig()
	return rc, nil
}
