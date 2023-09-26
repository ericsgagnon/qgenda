package qgenda

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/exiledavatar/gotoolkit/meta"
	"github.com/jmoiron/sqlx"
)

type Staff struct {
	// ------- metadata ------------------- //
	RawMessage       *string `json:"-" db:"_raw_message" pgtype:"text" idhash:"-"`
	ProcessedMessage *string `json:"-" db:"_processed_message" pgtype:"text" idhash:"-"` // RawMessage processed, omits 'message' metadata and 'noisy' fields (eg lastlogin)
	SourceQuery      *string `json:"_source_query,omitempty" db:"_source_query" pgtype:"text" idhash:"-"`
	ExtractDateTime  *Time   `json:"_extract_date_time,omitempty" db:"_extract_date_time" pgtype:"timestamp with time zone" idhash:"-"`
	IDHash           *string `json:"_id_hash,omitempty" db:"_id_hash" pgtype:"text" primarykey:"true" idhash:"-" ` // hash of identifying fields
	// ------------------------------------ //
	Abbrev                   *string       `json:"Abbrev,omitempty" db:"abbrev" pgtype:"text" idhash:"true"`
	BgColor                  *string       `json:"BgColor,omitempty" db:"bgcolor" pgtype:"text" idhash:"true"`
	BillSysID                *string       `json:"BillSysId,omitempty" db:"billsysid" pgtype:"text" idhash:"true"`
	CompKey                  *string       `json:"CompKey,omitempty" db:"compkey" pgtype:"text" idhash:"true"`
	ContactInstructions      *string       `json:"Contact Instructions,omitempty" db:"contact" pgtype:"text" idhash:"true"`
	Email                    *string       `json:"Email,omitempty" db:"email" pgtype:"text" idhash:"true"`
	SsoID                    *string       `json:"SsoId,omitempty" db:"ssoid" pgtype:"text" idhash:"true"`
	EmrID                    *string       `json:"EmrId,omitempty" db:"emrid" pgtype:"text" idhash:"true"`
	ErpID                    *string       `json:"ErpId,omitempty" db:"erpid" pgtype:"text" idhash:"true"`
	EndDate                  *Date         `json:"EndDate,omitempty" db:"enddate" pgtype:"date" idhash:"true" `
	ExtCallSysId             *string       `json:"ExtCallSysId,omitempty" db:"extcallsysid" pgtype:"text" idhash:"true" `
	FirstName                *string       `json:"FirstName,omitempty" db:"firstname" pgtype:"text" idhash:"true" `
	LastName                 *string       `json:"LastName,omitempty" db:"lastname" pgtype:"text" idhash:"true" `
	HomePhone                *string       `json:"HomePhone,omitempty" db:"homephone" pgtype:"text" idhash:"true" `
	MobilePhone              *string       `json:"MobilePhone,omitempty" db:"mobilephone" pgtype:"text" idhash:"true" `
	Npi                      *string       `json:"Npi,omitempty" db:"npi" pgtype:"text" idhash:"true" `
	OtherNumber1             *string       `json:"OtherNumber1,omitempty" db:"othernumber1" pgtype:"text" idhash:"true" `
	OtherNumber2             *string       `json:"OtherNumber2,omitempty" db:"othernumber2" pgtype:"text" idhash:"true" `
	OtherNumber3             *string       `json:"OtherNumber3,omitempty" db:"othernumber3" pgtype:"text" idhash:"true" `
	OtherNumberType1         *string       `json:"OtherNumberType1,omitempty" db:"othernumbertype1" pgtype:"text" idhash:"true" `
	OtherNumberType2         *string       `json:"OtherNumberType2,omitempty" db:"othernumbertype2" pgtype:"text" idhash:"true" `
	OtherNumberType3         *string       `json:"OtherNumberType3,omitempty" db:"othernumbertype3" pgtype:"text" idhash:"true" `
	Pager                    *string       `json:"Pager,omitempty" db:"pager" pgtype:"text" idhash:"true" `
	PayrollId                *string       `json:"PayrollId,omitempty" db:"payrollid" pgtype:"text" idhash:"true" `
	RegHours                 *float64      `json:"RegHours,omitempty" db:"reghours" pgtype:"numeric" idhash:"-" `
	StaffId                  *string       `json:"StaffId,omitempty" db:"staffid" pgtype:"text" idhash:"true" `
	StaffKey                 *string       `json:"StaffKey,omitempty" db:"staffkey" pgtype:"text" idhash:"true" `
	StartDate                *Date         `json:"StartDate,omitempty" db:"startdate" pgtype:"date" idhash:"true" `
	TextColor                *string       `json:"TextColor,omitempty" db:"textcolor" pgtype:"text" idhash:"true" `
	Addr1                    *string       `json:"Addr1,omitempty" db:"addr1" pgtype:"text" idhash:"true" `
	Addr2                    *string       `json:"Addr2,omitempty" db:"addr2" pgtype:"text" idhash:"true" `
	City                     *string       `json:"City,omitempty" db:"city" pgtype:"text" idhash:"true" `
	State                    *string       `json:"State,omitempty" db:"state" pgtype:"text" idhash:"true" `
	Zip                      *string       `json:"Zip,omitempty" db:"zip" pgtype:"text" idhash:"true" `
	IsActive                 *bool         `json:"IsActive,omitempty" db:"isactive" pgtype:"boolean" idhash:"true" `
	StaffTypeKey             *string       `json:"StaffTypeKey,omitempty" db:"stafftypekey" pgtype:"text" idhash:"true" `
	BillingTypeKey           *string       `json:"BillingTypeKey,omitempty" db:"billingtypekey" pgtype:"text" idhash:"true" `
	UserProfileKey           *string       `json:"UserProfileKey,omitempty" db:"userprofilekey" pgtype:"text" idhash:"true" `
	UserProfile              *string       `json:"UserProfile,omitempty" db:"userprofile" pgtype:"text" idhash:"true" `
	PayPeriodGroupName       *string       `json:"PayPeriodGroupName,omitempty" db:"payperiodgroupname" pgtype:"text" idhash:"true" `
	PayrollStartDate         *Date         `json:"PayrollStartDate,omitempty" db:"payrollstartdate" pgtype:"date" idhash:"true" `
	PayrollEndDate           *Date         `json:"PayrollEndDate,omitempty" db:"payrollenddate" pgtype:"date" idhash:"true" `
	TimeClockStartDate       *Date         `json:"TimeClockStartDate,omitempty" db:"timeclockstartdate" pgtype:"date" idhash:"true" `
	TimeClockEndDate         *Date         `json:"TimeClockEndDate,omitempty" db:"timeclockenddate" pgtype:"date" idhash:"true" `
	TimeClockKioskPIN        *string       `json:"TimeClockKioskPIN,omitempty" db:"timeclockkioskpin" pgtype:"text" idhash:"true" `
	IsAutoApproveSwap        *bool         `json:"IsAutoApproveSwap,omitempty" db:"isautoapproveswap" pgtype:"boolean" idhash:"true" `
	DailyUnitAverage         *float64      `json:"DailyUnitAverage,omitempty" db:"dailyunitaverage" pgtype:"numeric" idhash:"-"`
	StaffInternalId          *string       `json:"StaffInternalId,omitempty" db:"staffinternalid" pgtype:"text" idhash:"true" `
	UserLastLoginDateTimeUTC *Time         `json:"UserLastLoginDateTimeUTC,omitempty" db:"userlastlogindatetimeutc" pgtype:"timestamp with time zone" idhash:"-"`
	SourceOfLogin            *string       `json:"SourceOfLogin,omitempty" db:"sourceoflogin" pgtype:"text" idhash:"-" `
	CalSyncKey               *string       `json:"CalSyncKey,omitempty" db:"calsynckey" pgtype:"text" idhash:"true" `
	Tags                     StaffTags     `json:"Tags,omitempty" table:"stafftag" db:"-" pgtype:"text" idhash:"true"`
	TTCMTags                 StaffTags     `json:"TTCMTags,omitempty" table:"staffttcmtag" db:"-" idhash:"true"`
	Skillset                 StaffSkills   `json:"Skillset,omitempty" table:"staffskillset" db:"-" idhash:"true"`
	Profiles                 StaffProfiles `json:"Profiles,omitempty" table:"staffprofile" db:"-" idhash:"true"`
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
	// Propagate fields and sort children
	// Tags
	for i, _ := range s.Tags {
		s.Tags[i].ExtractDateTime = s.ExtractDateTime
		s.Tags[i].StaffIDHash = s.IDHash
		s.Tags[i].StaffKey = s.StaffKey
	}
	s.Tags.Sort()

	// TTCMTags
	for i, _ := range s.TTCMTags {
		s.TTCMTags[i].ExtractDateTime = s.ExtractDateTime
		s.TTCMTags[i].StaffIDHash = s.IDHash
		s.TTCMTags[i].StaffKey = s.StaffKey
	}
	s.TTCMTags.Sort()

	// Skillset
	for i, _ := range s.Skillset {
		s.Skillset[i].ExtractDateTime = s.ExtractDateTime
		s.Skillset[i].StaffIDHash = s.IDHash
		s.Skillset[i].StaffKey = s.StaffKey
	}
	s.Skillset.Sort()

	// Profiles
	for i, _ := range s.Profiles {
		s.Profiles[i].ExtractDateTime = s.ExtractDateTime
		s.Profiles[i].StaffIDHash = s.IDHash
		s.Profiles[i].StaffKey = s.StaffKey
	}
	s.Profiles.Sort()

	if err := s.SetMessage(); err != nil {
		return err
	}
	if err := s.SetIDHash(); err != nil {
		return err
	}

	for i, _ := range s.Tags {
		s.Tags[i].StaffIDHash = s.IDHash
	}
	for i, _ := range s.TTCMTags {
		s.TTCMTags[i].StaffIDHash = s.IDHash
	}
	for i, _ := range s.Skillset {
		s.Skillset[i].StaffIDHash = s.IDHash
	}
	for i, _ := range s.Profiles {
		s.Profiles[i].StaffIDHash = s.IDHash
	}

	if err := meta.ProcessStruct(s); err != nil {
		return fmt.Errorf("error processing %T:\t%q", s, err)
	}

	return nil
}

func (s *Staff) SetMessage() error {
	// take a copy and strip message fields, for good measure
	ss := *s
	ss.RawMessage = nil
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

func (s *Staff) SetIDHash() error {
	idh := meta.ToValueMap(*s, "idhash").Hash()
	s.IDHash = &idh
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

func (s Staff) CreatePGTable(ctx context.Context, db *sqlx.DB, schema, table string) (sql.Result, error) {
	return CreatePGTable(ctx, db, s, schema, table)
}

// Staff is a snapshot dataset that does full import for all changes, it always requests all records by default
func (s Staff) GetPGStatus(ctx context.Context, db *sqlx.DB, schema, table string) (*RequestConfig, error) {
	return DefaultStaffRequestConfig(), nil
}

// func (s Staff) RequestConfig() *RequestConfig {
// 	requestPath := "staffmember"
// 	allowedFields := []string{
// 		"Includes",
// 		"Select",
// 		"Filter",
// 		"OrderBy",
// 		"Expand",
// 	}
// 	rc := NewRequestConfig(requestPath, allowedFields)
// 	rc.SetIncludes("Tags,TTCMTags,Skillset,Profiles")
// 	return rc

// }

// func (s Staff) NewRequestConfig(rc *RequestConfig) *RequestConfig {
// 	return s.RequestConfig().Merge(rc)
// }

// func (s Staff) NewRequest(rc *RequestConfig) *Request {
// 	rc = s.NewRequestConfig(rc)
// 	return NewRequest(rc)
// }
