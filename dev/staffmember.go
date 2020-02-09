package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// // StaffType represents preconfigured values
// type StaffType int

// // preconfigured values for stafftype
// const (
// 	nil StaffType = iota
// 	Physician
// 	CRNA
// 	Technologist
// 	Locum
// 	Office
// 	Resident
// 	Specialty
// 	MAPA
// 	Nurse
// 	Other
// )

// Time embeds time.Time
// qgenda doesn't comply with RFC3339
// chose this over wrapping to slightly improve
// convenience of calling time.Time's methods
type Time struct {
	time.Time
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *Time) UnmarshalJSON(data []byte) error {

	tag := reflect.ValueOf(data).Type().Field(0).Tag.Get("json")
	fmt.Println(tag)

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// *t, err = Parse(`"`+RFC3339+`"`, string(data))
	t.Time, err = time.ParseInLocation("2006-01-02T15:04:05", string(data), location)
	return err
}

// StaffMember represents staff, and possibly some other entities as well?
type StaffMember struct {
	Abbreviation         string    `json:"Abbrev"`
	BackgroundColor      string    `json:"BgColor"`
	BillingSystemID      uuid.UUID `json:"BillSysId"`
	CalendarSyncID       uuid.UUID `json:"CalSyncKey"`
	CompanyID            uuid.UUID `json:"CompKey"`
	Email                string    `json:"Email"`
	EMRID                string    `json:"EmrId"`
	ERPID                string    `json:"ErpId"`
	EndDate              Time      `json:"EndDate"`
	ExternalCallSystemID uuid.UUID `json:"ExtCallSysId"`
	FirstName            string    `json:"FirstName"`
	HomePhone            string    `json:"HomePhone"`
	LastName             string    `json:"LastName"`
	MobilePhone          string    `json:"MobilePhone"`
	NPI                  string    `json:"Npi"`
	PagerNumber          string    `json:"Pager"`
	PayrollID            string    `json:"PayrollId"`
	RegularHours         string    `json:"RegHours"`
	Alias                string    `json:"StaffId"`
	ID                   uuid.UUID `json:"StaffKey"`
	StartDate            Time      `json:"StartDate"`
	TextColor            string    `json:"TextColor"`
	Active               bool      `json:"IsActive"`
	StaffType            string    `json:"StaffTypeKey"`
	BillingType          string    `json:"BillingTypeKey"`
	ProfileID            uuid.UUID `json:"UserProfileKey"`
	Profile              string    `json:"UserProfile"`
	PayrollStartDate     time.Time `json:"PayrollStartDate"`
	PayrollEndDate       time.Time `json:"PayrollEndDate"`
	TimeClockStartDate   time.Time `json:"TimeClockStartDate"`
	TimeClockEndDate     time.Time `json:"TimeClockEndDate"`
	TimeClockKioskPIN    string    `json:"TimeClockKioskPIN"`
	AutoApproveSwap      bool      `json:"IsAutoApproveSwap"`
	Viewable             bool      `json:"IsViewable"`
	Schedulable          bool      `json:"IsSchedulable"`
	Address              struct {
		Line1 string `json:"Addr1"`
		Line2 string `json:"Addr2"`
		City  string `json:"City"`
		State string `json:"State"`
		Zip   string `json:"Zip"`
	}
	LastLogin struct {
		Time   time.Time `json:"UserLastLoginDateTimeUtc"`
		Source string    `json:"SourceOfLogin"`
	}
	// The don't appear to be used at our insitution, not sure if they are elsewhere
	// omitting due to inability to validate
	// `json:"Tags"`
	// TTCMTags `json:"TTCMTags"`
	// `json:"CategoryKey"`
	// `json:"CategoryName"`
	// Tags `json:"Tags"`
	// `json:"Key"`
	// `json:"Name"`
	// `json:"EffectiveFromDate"`
	// `json:"EffectiveToDate"`
	// `json:"Skillset"`
	// `json:"Profiles"`
	// `json:"Name"`
	// `json:"ProfileKey"`
	// `json:"DailyUnitAverage"` // don't waste your time on this
}
