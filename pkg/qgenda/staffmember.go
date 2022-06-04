package qgenda

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
	StaffKey                 *string         `json:"StaffKey,omitempty"`
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

type StaffSkillset struct {
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

