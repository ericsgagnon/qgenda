package main

import (
	"github.com/google/uuid"
)

// NewStaffMemberRequestResponse returns a pointer to a StaffMemberRequestConfig with default values
func NewStaffMemberRequestResponse() *RequestResponse {
	rr := NewRequestResponse()
	rr.Request.Config = NewStaffMemberRequestConfig()
	return rr
}

// NewStaffMemberRequestConfig returns a point to a StaffMemberRequestConfig with default values
func NewStaffMemberRequestConfig() *StaffMemberRequestConfig {
	r := &StaffMemberRequestConfig{
		Resource: "StaffMember",
		Route:    "/staffmember",
		Includes: "Skillset,Tags,Profiles,TTCMTags",
		// Select:   "",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}
	return r
}

// StaffMemberRequestConfig struct captures all available request arguments for
// qgenda StaffMembers endpoint
type StaffMemberRequestConfig struct {
	Resource string
	Route    string `path:"-"`
	Includes string `query:"includes"`
	Select   string `query:"$select"`
	Filter   string `query:"$filter"`
	OrderBy  string `query:"$orderby"`
	Expand   string `query:"$expand"`
}

// StaffMember represents staff, and possibly some other entities as well?
type StaffMember struct {
	Abbreviation         string    `json:"Abbrev"`
	BackgroundColor      string    `json:"BgColor"`
	BillingSystemID      string    `json:"BillSysId,omitempty"`
	CalendarSyncID       string    `json:"CalSyncKey,omitempty"`
	CompanyID            uuid.UUID `json:"CompKey,omitempty"`
	Email                string    `json:"Email"`
	EMRID                string    `json:"EmrId"`
	ERPID                string    `json:"ErpId"`
	EndDate              Time      `json:"EndDate"`
	ExternalCallSystemID string    `json:"ExtCallSysId,omitempty"`
	FirstName            string    `json:"FirstName"`
	HomePhone            string    `json:"HomePhone"`
	LastName             string    `json:"LastName"`
	MobilePhone          string    `json:"MobilePhone"`
	NPI                  string    `json:"Npi"`
	PagerNumber          string    `json:"Pager"`
	PayrollID            string    `json:"PayrollId"`
	RegularHours         string    `json:"RegHours"`
	Alias                string    `json:"StaffId"`
	ID                   uuid.UUID `json:"StaffKey,omitempty"`
	StartDate            Time      `json:"StartDate"`
	TextColor            string    `json:"TextColor"`
	Active               bool      `json:"IsActive"`
	StaffType            string    `json:"StaffTypeKey"`
	BillingType          string    `json:"BillingTypeKey"`
	ProfileID            uuid.UUID `json:"UserProfileKey,omitempty"`
	Profile              string    `json:"UserProfile"`
	PayrollStartDate     Time      `json:"PayrollStartDate"`
	PayrollEndDate       Time      `json:"PayrollEndDate"`
	TimeClockStartDate   Time      `json:"TimeClockStartDate"`
	TimeClockEndDate     Time      `json:"TimeClockEndDate"`
	TimeClockKioskPIN    string    `json:"TimeClockKioskPIN"`
	AutoApproveSwap      bool      `json:"IsAutoApproveSwap"`
	DailyUnitAverage     float64   `json:"DailyUnitAverage"`
	Viewable             bool      `json:"IsViewable"`
	Schedulable          bool      `json:"IsSchedulable"`
	Test                 Time      `json:"UserLastLoginDateTimeUtc"`
	Address              struct {
		Line1 string `json:"Addr1"`
		Line2 string `json:"Addr2"`
		City  string `json:"City"`
		State string `json:"State"`
		Zip   string `json:"Zip"`
	}
	LastLogin struct {
		Time   TimeUTC `json:"UserLastLoginDateTimeUtc"`
		Source string  `json:"SourceOfLogin"`
	}
	SkillSet []SkillSet `json:"Skillset"`
	Profiles []Profile  `json:"Profiles"`
	// Tags     []Tag      `json:"Tags"`
	// TTCMTags []Tag      `json:"TTCMTags"`
}

// SkillSet captures the staff to task relationship
type SkillSet struct {
	Staff struct {
		FirstName    string `json:"StaffFirstName"`
		LastName     string `json:"StaffLastName"`
		Abbreviation string `json:"StaffAbbreviation"`
		Alias        string `json:"StaffId"`
	}
	Task struct {
		Name         string `json:"TaskName"`
		Abbreviation string `json:"TaskAbbrev"`
		Alias        string `json:"TaskId"`
	}
	Monday struct {
		Valid     bool   `json:"IsSkilledMon"`
		Frequency string `json:"MonOccurrence"`
	}
	Tuesday struct {
		Valid     bool   `json:"IsSkilledTue"`
		Frequency string `json:"TueOccurrence"`
	}
	Wednesday struct {
		Valid     bool   `json:"IsSkilledWed"`
		Frequency string `json:"WedOccurrence"`
	}
	Thursday struct {
		Valid     bool   `json:"IsSkilledThu"`
		Frequency string `json:"ThuOccurrence"`
	}
	Friday struct {
		Valid     bool   `json:"IsSkilledFri"`
		Frequency string `json:"FriOccurrence"`
	}
	Saturday struct {
		Valid     bool   `json:"IsSkilledSat"`
		Frequency string `json:"SatOccurrence"`
	}
	Sunday struct {
		Valid     bool   `json:"IsSkilledSun"`
		Frequency string `json:"SunOccurrence"`
	}
}

// GetStaffMembers returns all staff members
// func (q *QgendaClient) GetStaffMembers(ctx context.Context, sr *StaffMemberRequest, il *ItemList) error {
// 	if sr == nil {
// 		sr = NewStaffMemberRequest()
// 	}
// 	r, err := ParseRequest(sr)
// 	// fmt.Println(r)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	bb, meta, err := q.Get(ctx, r)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// fmt.Println(meta)
// 	meta.Name = "StaffMemberList"
// 	// if err := ioutil.WriteFile("data/sml.json", bb, 0777); err != nil {
// 	// 	log.Printf("Couldn't response to disk: %v", err)
// 	// 	return err
// 	// }
// 	// fmt.Println(string(bb))
// 	var sm []StaffMember
// 	// c := []StaffMember{}
// 	if err := json.Unmarshal(bb, &sm); err != nil {
// 		log.Printf("Error unmarshalling response from %v", err)
// 		return err
// 	}
// 	// fmt.Printf("\n\n%+v\n\n", sm)
// 	il.MetaData = meta
// 	il.Items = sm

// 	return nil
// }

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
