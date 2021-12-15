package qgenda

import (
    "time"
)

type Assignment struct {
	EncounterRoleKey	string	`json:"EncounterRoleKey"`
	StaffMemberKey	string	`json:"StaffMemberKey"`
	StaffMemberAbbreviation	string	`json:"StaffMemberAbbreviation"`
	RoleName	string	`json:"RoleName"`
	StaffMemberID	string	`json:"StaffMemberID"`
	StaffMemberNPI	string	`json:"StaffMemberNPI"`
	StaffMemberEMRID	string	`json:"StaffMemberEMRID"`
	IsNonStaff	bool	`json:"IsNonStaff"`
	IsPublished	bool	`json:"IsPublished"`
	IsDefaultRoomAssignment	bool	`json:"IsDefaultRoomAssignment"`
}


type Company struct {
	CompanyKey	string	`json:"CompanyKey"`
	CompanyName	string	`json:"CompanyName"`
	CompanyAbbreviation	string	`json:"CompanyAbbreviation"`
	CompanyLocation	string	`json:"CompanyLocation"`
	CompanyPhoneNumber	string	`json:"CompanyPhoneNumber"`
	Profiles	[]Profile	`json:"Profiles"`
	ProfileName	string	`json:"ProfileName"`
	ProfileKey	string	`json:"ProfileKey"`
	IsAdmin	bool	`json:"IsAdmin"`
	Organizations	[]Organization	`json:"Organizations"`
	OrgName	string	`json:"OrgName"`
	OrgKey	int	`json:"OrgKey"`
}


type EncounterRole struct {
	EncounterRoleKey	string	`json:"EncounterRoleKey"`
	Name	string	`json:"Name"`
}


type Location struct {
	CompanyKey	string	`json:"CompanyKey"`
	LocationKey	int	`json:"LocationKey"`
	Id	string	`json:"Id"`
	Name	string	`json:"Name"`
	Abbrev	string	`json:"Abbrev"`
	Address	string	`json:"Address"`
	Tags	[]Tag	`json:"Tags"`
	Notes	string	`json:"Notes"`
	TimeZone	string	`json:"TimeZone"`
}


type OpenShift struct {
	CompanyKey	string	`json:"CompanyKey"`
	ScheduleKey	string	`json:"ScheduleKey"`
	OpenShiftCount	int	`json:"OpenShiftCount"`
	CallRole	string	`json:"CallRole"`
	Credit	float64	`json:"Credit"`
	Date	time.Time	`json:"Date"`
	StartDate	time.Time	`json:"StartDate"`
	StartDateUTC	time.Time	`json:"StartDateUTC"`
	StartTime	time.Time	`json:"StartTime"`
	EndDate	time.Time	`json:"EndDate"`
	EndDateUTC	time.Time	`json:"EndDateUTC"`
	EndTime	time.Time	`json:"EndTime"`
	IsCred	bool	`json:"IsCred"`
	IsSaved	bool	`json:"IsSaved"`
	IsPublished	bool	`json:"IsPublished"`
	IsLocked	bool	`json:"IsLocked"`
	IsStruck	bool	`json:"IsStruck"`
	Notes	string	`json:"Notes"`
	IsNotePrivate	bool	`json:"IsNotePrivate"`
	TaskAbbrev	string	`json:"TaskAbbrev"`
	TaskId	string	`json:"TaskId"`
	TaskEmrId	string	`json:"TaskEmrId"`
	TaskIsPrintStart	bool	`json:"TaskIsPrintStart"`
	TaskIsPrintEnd	bool	`json:"TaskIsPrintEnd"`
	TaskExtCallSysId	string	`json:"TaskExtCallSysId"`
	TaskKey	string	`json:"TaskKey"`
	TaskName	string	`json:"TaskName"`
	TaskBillSysId	string	`json:"TaskBillSysId"`
	TaskPayrollId	string	`json:"TaskPayrollId"`
	TaskShiftKey	string	`json:"TaskShiftKey"`
	TaskType	string	`json:"TaskType"`
	TaskContactInformation	string	`json:"TaskContactInformation"`
	TaskTags	[]Tag	`json:"TaskTags"`
	LocationKey	string	`json:"LocationKey"`
	LocationName	string	`json:"LocationName"`
	LocationAbbrev	string	`json:"LocationAbbrev"`
	LocationAddress	string	`json:"LocationAddress"`
	LocationTags	[]Tag	`json:"LocationTags"`
	TimeZone	string	`json:"TimeZone"`
}


type PatientEncounter struct {
	CompanyKey	string	`json:"CompanyKey"`
	DailyConfigurationKey	string	`json:"DailyConfigurationKey"`
	EncounterKey	string	`json:"EncounterKey"`
	DailyCaseID	int	`json:"DailyCaseID"`
	Date	time.Time	`json:"Date"`
	StartTime	time.Time	`json:"StartTime"`
	EndTime	time.Time	`json:"EndTime"`
	Duration	float64	`json:"Duration"`
	IsCancelled	bool	`json:"IsCancelled"`
	Room	Room	`json:"Room"`
	RoomKey	string	`json:"RoomKey"`
	RoomName	string	`json:"RoomName"`
	LocationName	string	`json:"LocationName"`
	LocationKey	string	`json:"LocationKey"`
	Assignments	[]Assignment	`json:"Assignments"`
	ExternalID1	string	`json:"ExternalID1"`
	ExternalID2	string	`json:"ExternalID2"`
	CreatedBySource	string	`json:"CreatedBySource"`
	ModifiedBySource	string	`json:"ModifiedBySource"`
}


type PatientInformation struct {
	PatientFirstName	string	`json:"PatientFirstName"`
	PatientLastName	string	`json:"PatientLastName"`
	PatientGender	string	`json:"PatientGender"`
	PatientEmail	string	`json:"PatientEmail"`
	PatientHomePhone	string	`json:"PatientHomePhone"`
	PatientCellPhone	string	`json:"PatientCellPhone"`
	PatientAlternatePhone	string	`json:"PatientAlternatePhone"`
	PatientGuardianPhone	string	`json:"PatientGuardianPhone"`
	PatientGuardianName	string	`json:"PatientGuardianName"`
	PatientDateOfBirth	time.Time	`json:"PatientDateOfBirth"`
	PatientMedicalRecordNumber	string	`json:"PatientMedicalRecordNumber"`
	PatientClinicalNotes	string	`json:"PatientClinicalNotes"`
	PatientExtraNotes	string	`json:"PatientExtraNotes"`
	PatientAddress1	string	`json:"PatientAddress1"`
	PatientAddress2	string	`json:"PatientAddress2"`
	PatientCity	string	`json:"PatientCity"`
	PatientState	string	`json:"PatientState"`
	PatientPostalCode	string	`json:"PatientPostalCode"`
	PatientPrimaryInsuranceId	string	`json:"PatientPrimaryInsuranceId"`
	PatientSecondaryInsuranceId	string	`json:"PatientSecondaryInsuranceId"`
	PatientSocialSecurityNumber	string	`json:"PatientSocialSecurityNumber"`
	PatientMaritalStatus	string	`json:"PatientMaritalStatus"`
}


type Request struct {
	RequestId	string	`json:"RequestId"`
	StaffNpi	string	`json:"StaffNpi"`
	TaskId	string	`json:"TaskId"`
	CallRole	string	`json:"CallRole"`
	CompKey	string	`json:"CompKey"`
	Credit	float64	`json:"Credit"`
	Date	time.Time	`json:"Date"`
	EndTime	time.Time	`json:"EndTime"`
	Notes	string	`json:"Notes"`
	RequestKey	string	`json:"RequestKey"`
	RequestType	string	`json:"RequestType"`
	RequestStatus	string	`json:"RequestStatus"`
	ResolutionDate	time.Time	`json:"ResolutionDate"`
	StaffAbbrev	string	`json:"StaffAbbrev"`
	StaffFName	string	`json:"StaffFName"`
	StaffKey	string	`json:"StaffKey"`
	StaffLName	string	`json:"StaffLName"`
	StaffPayrollId	string	`json:"StaffPayrollId"`
	StartTime	time.Time	`json:"StartTime"`
	SubmittedByUser	string	`json:"SubmittedByUser"`
	SubmittedDate	time.Time	`json:"SubmittedDate"`
	TaskAbbrev	string	`json:"TaskAbbrev"`
	TaskKey	string	`json:"TaskKey"`
	TaskName	string	`json:"TaskName"`
	StaffInternalId	string	`json:"StaffInternalId"`
	Removed	bool	`json:"Removed"`
}


type ScheduleAuditLogEntry struct {
	StaffFirstName	string	`json:"StaffFirstName"`
	StaffLastName	string	`json:"StaffLastName"`
	StaffAbbreviation	string	`json:"StaffAbbreviation"`
	StaffKey	string	`json:"StaffKey"`
	TaskName	string	`json:"TaskName"`
	TaskAbbreviation	string	`json:"TaskAbbreviation"`
	TaskKey	string	`json:"TaskKey"`
	ScheduleEntryDate	time.Time	`json:"ScheduleEntryDate"`
	ScheduleEntryStartTimeUTC	time.Time	`json:"ScheduleEntryStartTimeUTC"`
	ScheduleEntryStartTimeLocal	time.Time	`json:"ScheduleEntryStartTimeLocal"`
	ScheduleEntryEndTimeUTC	time.Time	`json:"ScheduleEntryEndTimeUTC"`
	ScheduleEntryEndTimeLocal	time.Time	`json:"ScheduleEntryEndTimeLocal"`
	ScheduleEntryKey	string	`json:"ScheduleEntryKey"`
	ActivityType	string	`json:"ActivityType"`
	SourceType	string	`json:"SourceType"`
	UserFirstName	string	`json:"UserFirstName"`
	UserLastName	string	`json:"UserLastName"`
	UserKey	string	`json:"UserKey"`
	TimestampUTC	time.Time	`json:"TimestampUTC"`
	Timestamp	time.Time	`json:"Timestamp"`
	AdditionalInformation	string	`json:"AdditionalInformation"`
	Location	[]Location	`json:"Location"`
	LocationName	string	`json:"LocationName"`
	LocationAbbreviation	string	`json:"LocationAbbreviation"`
	LocationID	string	`json:"LocationID"`
	LocationKey	string	`json:"LocationKey"`
	IPAddress	string	`json:"IPAddress"`
}


type ScheduleEntry struct {
	ScheduleKey	string	`json:"ScheduleKey"`
	CallRole	string	`json:"CallRole"`
	CompKey	string	`json:"CompKey"`
	Credit	float64	`json:"Credit"`
	Date	time.Time	`json:"Date"`
	StartDateUTC	time.Time	`json:"StartDateUTC"`
	EndDateUTC	time.Time	`json:"EndDateUTC"`
	EndDate	time.Time	`json:"EndDate"`
	EndTime	time.Time	`json:"EndTime"`
	IsCred	bool	`json:"IsCred"`
	IsSaved	bool	`json:"IsSaved"`
	IsPublished	bool	`json:"IsPublished"`
	IsLocked	bool	`json:"IsLocked"`
	IsStruck	bool	`json:"IsStruck"`
	Notes	string	`json:"Notes"`
	IsNotePrivate	bool	`json:"IsNotePrivate"`
	StaffAbbrev	string	`json:"StaffAbbrev"`
	StaffBillSysId	string	`json:"StaffBillSysId"`
	StaffEmail	string	`json:"StaffEmail"`
	StaffEmrId	string	`json:"StaffEmrId"`
	StaffErpId	string	`json:"StaffErpId"`
	StaffInternalId	string	`json:"StaffInternalId"`
	StaffExtCallSysId	string	`json:"StaffExtCallSysId"`
	StaffFName	string	`json:"StaffFName"`
	StaffId	string	`json:"StaffId"`
	StaffKey	string	`json:"StaffKey"`
	StaffLName	string	`json:"StaffLName"`
	StaffMobilePhone	string	`json:"StaffMobilePhone"`
	StaffNpi	string	`json:"StaffNpi"`
	StaffPager	string	`json:"StaffPager"`
	StaffPayrollId	string	`json:"StaffPayrollId"`
	StaffTags	[]Tag	`json:"StaffTags"`
	StartDate	time.Time	`json:"StartDate"`
	StartTime	time.Time	`json:"StartTime"`
	TaskAbbrev	string	`json:"TaskAbbrev"`
	TaskBillSysId	string	`json:"TaskBillSysId"`
	TaskContactInformation	string	`json:"TaskContactInformation"`
	TaskExtCallSysId	string	`json:"TaskExtCallSysId"`
	TaskId	string	`json:"TaskId"`
	TaskEmrId	string	`json:"TaskEmrId"`
	TaskCallPriority	string	`json:"TaskCallPriority"`
	TaskDepartmentId	string	`json:"TaskDepartmentId"`
	TaskIsPrintEnd	bool	`json:"TaskIsPrintEnd"`
	TaskIsPrintStart	bool	`json:"TaskIsPrintStart"`
	TaskKey	string	`json:"TaskKey"`
	TaskName	string	`json:"TaskName"`
	TaskPayrollId	string	`json:"TaskPayrollId"`
	TaskShiftKey	string	`json:"TaskShiftKey"`
	TaskType	string	`json:"TaskType"`
	TaskTags	[]Tag	`json:"TaskTags"`
	LocationName	string	`json:"LocationName"`
	LocationAbbrev	string	`json:"LocationAbbrev"`
	LocationID	string	`json:"LocationID"`
	LocationAddress	string	`json:"LocationAddress"`
	LocationTags	[]Tag	`json:"LocationTags"`
	TimeZone	string	`json:"TimeZone"`
	LastModifiedDateUTC	time.Time	`json:"LastModifiedDateUTC"`
	IsRotationTask	bool	`json:"IsRotationTask"`
}


type ScheduleRotation struct {
	StartDate	time.Time	`json:"StartDate"`
	EndDate	time.Time	`json:"EndDate"`
	IsPublished	bool	`json:"IsPublished"`
	IsLocked	bool	`json:"IsLocked"`
	StaffKey	string	`json:"StaffKey"`
	TaskKey	string	`json:"TaskKey"`
	TimeZone	string	`json:"TimeZone"`
}


type Skillset struct {
	StaffFirstName	string	`json:"StaffFirstName"`
	StaffLastName	string	`json:"StaffLastName"`
	StaffAbbreviation	string	`json:"StaffAbbreviation"`
	StaffId	string	`json:"StaffId"`
	TaskName	string	`json:"TaskName"`
	TaskAbbrev	string	`json:"TaskAbbrev"`
	TaskId	string	`json:"TaskId"`
	IsSkilledMon	bool	`json:"IsSkilledMon"`
	MonOccurrence	string	`json:"MonOccurrence"`
	IsSkilledTue	bool	`json:"IsSkilledTue"`
	TueOccurrence	string	`json:"TueOccurrence"`
	IsSkilledWed	bool	`json:"IsSkilledWed"`
	WedOccurrence	string	`json:"WedOccurrence"`
	IsSkilledThu	bool	`json:"IsSkilledThu"`
	ThuOccurrence	string	`json:"ThuOccurrence"`
	IsSkilledFri	bool	`json:"IsSkilledFri"`
	FriOccurrence	string	`json:"FriOccurrence"`
	IsSkilledSat	bool	`json:"IsSkilledSat"`
	SatOccurrence	string	`json:"SatOccurrence"`
	IsSkilledSun	bool	`json:"IsSkilledSun"`
	SunOccurrence	string	`json:"SunOccurrence"`
}


type StandardFields struct {
	EncounterFieldKey	string	`json:"EncounterFieldKey"`
	Name	string	`json:"Name"`
	Value	string	`json:"Value"`
	Type	string	`json:"Type"`
}


type Tag struct {
	CategoryKey	int	`json:"CategoryKey"`
	CategoryName	string	`json:"CategoryName"`
	Tags	[]Tag	`json:"Tags"`
	Key	int	`json:"Key"`
	Name	string	`json:"Name"`
}


type Task struct {
	Abbrev	string	`json:"Abbrev"`
	BgColor	string	`json:"BgColor"`
	BillSysId	string	`json:"BillSysId"`
	CallRole	string	`json:"CallRole"`
	CompKey	string	`json:"CompKey"`
	ContactInformation	string	`json:"ContactInformation"`
	DepartmentId	string	`json:"DepartmentId"`
	DisplayAsAvailableOnCalendarConnection	bool	`json:"DisplayAsAvailableOnCalendarConnection"`
	EmrId	string	`json:"EmrId"`
	EnableWaitlist	bool	`json:"EnableWaitlist"`
	EndDate	time.Time	`json:"EndDate"`
	EnforceTimeCompatibility	bool	`json:"EnforceTimeCompatibility"`
	ExtCallPriority	string	`json:"ExtCallPriority"`
	ExtCallSysId	string	`json:"ExtCallSysId"`
	Extension	string	`json:"Extension"`
	FullyCompatible	bool	`json:"FullyCompatible"`
	HideOnCalendarSync	bool	`json:"HideOnCalendarSync"`
	Label	bool	`json:"Label"`
	Manual	bool	`json:"Manual"`
	Name	string	`json:"Name"`
	Notes	string	`json:"Notes"`
	NotificationList	string	`json:"NotificationList"`
	OverrideSortOrder	string	`json:"OverrideSortOrder"`
	PayrollId	string	`json:"PayrollId"`
	PayrollUnits	string	`json:"PayrollUnits"`
	ShowEndTime	bool	`json:"ShowEndTime"`
	ShowOnAdminOnly	bool	`json:"ShowOnAdminOnly"`
	TaskVisibility	string	`json:"TaskVisibility"`
	ShowOpensOnCalendar	bool	`json:"ShowOpensOnCalendar"`
	ShowStartTime	bool	`json:"ShowStartTime"`
	ShowTimesWhenEdited	bool	`json:"ShowTimesWhenEdited"`
	StartDate	time.Time	`json:"StartDate"`
	SuppressRequestEmailsAndNotifications	bool	`json:"SuppressRequestEmailsAndNotifications"`
	SyncAsAllDayEvent	bool	`json:"SyncAsAllDayEvent"`
	TaskId	string	`json:"TaskId"`
	TaskKey	string	`json:"TaskKey"`
	RequireTimePunch	bool	`json:"RequireTimePunch"`
	TextColor	string	`json:"TextColor"`
	TimeOffAfterShift	int	`json:"TimeOffAfterShift"`
	Type	string	`json:"Type"`
	Tags	[]Tag	`json:"Tags"`
	TaskShifts	[]TaskShift	`json:"TaskShifts"`
	Profiles	[]Profile	`json:"Profiles"`
	ProfileName	string	`json:"ProfileName"`
	ProfileKey	string	`json:"ProfileKey"`
	IsViewable	bool	`json:"IsViewable"`
	IsSchedulable	bool	`json:"IsSchedulable"`
	IsRotationTask	bool	`json:"IsRotationTask"`
	IsIgnoreHoursOffAfterTask	bool	`json:"IsIgnoreHoursOffAfterTask"`
}


type TaskShift struct {
	DayOfWeek	string	`json:"DayOfWeek"`
	StartTime	time.Time	`json:"StartTime"`
	EndTime	time.Time	`json:"EndTime"`
	MaxStaff	int	`json:"MaxStaff"`
	MinStaff	int	`json:"MinStaff"`
	OffAfter	int	`json:"OffAfter"`
	StatCredit	float64	`json:"StatCredit"`
	StaffCount	int	`json:"StaffCount"`
	EffectiveDate	[]EffectiveDate	`json:"EffectiveDate"`
	EffectiveFromDate	time.Time	`json:"EffectiveFromDate"`
	EffectiveToDate	time.Time	`json:"EffectiveToDate"`
	StartTime	time.Time	`json:"StartTime"`
	EndTime	time.Time	`json:"EndTime"`
	MaxStaff	int	`json:"MaxStaff"`
	MinStaff	int	`json:"MinStaff"`
	OffAfter	int	`json:"OffAfter"`
	StatCredit	float64	`json:"StatCredit"`
	IsActive	bool	`json:"IsActive"`
}
type Organization struct {
    Name string  `json:"OrgName"`
    Key  int     `json:"OrgKey"`
}

type ContactEmail struct {
	ContactName		string  `json:"ContactName"`
	ContactEmail	string  `json:"ContactEmail"`
}

type DailyCase struct {
	CompanyKey	string	`json:"CompanyKey"`
	EmrID	string	`json:"EmrId"`
	DailyCaseID	int	`json:"DailyCaseID"`
	TaskKey	string	`json:"TaskKey"`
	Task	string	`json:"Task"`
	LocationName	string	`json:"LocationName"`
	Date	time.Time	`json:"Date"`
	StartTime	time.Time	`json:"StartTime"`
	EndTime	time.Time	`json:"EndTime"`
	Surgeon	string	`json:"Surgeon"`
	Procedure	string	`json:"Procedure"`
	SupStaffKey	string	`json:"SupStaffKey"`
	SupStaffNpi	string	`json:"SupStaffNpi"`
	SupStaff	string	`json:"SupStaff"`
	DPStaffKey	string	`json:"DPStaffKey"`
	DPStaff	string	`json:"DPStaff"`
	DPStaffNpi	string	`json:"DPStaffNpi"`
	CustomTextFields	string	`json:"CustomTextFields"`
	CustomCheckboxFields	string	`json:"CustomCheckboxFields"`
	Cancelled	bool	`json:"IsCancelled"`
	Duration	float64	`json:"Duration"`
	BaseUnits	float64	`json:"BaseUnits"`
	ModifierUnits	float64	`json:"ModifierUnits"`
	TimeUnits	float64	`json:"TimeUnits"`
	PatientFirstName	string	`json:"PatientFirstName"`
	PatientLastName	string	`json:"PatientLastName"`
	PatientAge	string	`json:"PatientAge"`
	PatientGender	string	`json:"PatientGender"`
	PatientEmail	string	`json:"PatientEmail"`
	PatientHomePhone	string	`json:"PatientHomePhone"`
	PatientCellPhone	string	`json:"PatientCellPhone"`
	PatientAlternatePhone	string	`json:"PatientAlternatePhone"`
	PatientGuardianPhone	string	`json:"PatientGuardianPhone"`
	PatientGuardianName	string	`json:"PatientGuardianName"`
	PatientDOB	time.Time	`json:"PatientDOB"`
	PatientMRN	string	`json:"PatientMRN"`
	PatientClinicalNotes	string	`json:"PatientClinicalNotes"`
	PatientExtraNotes	string	`json:"PatientExtraNotes"`
	PatientAddress1	string	`json:"PatientAddress1"`
	PatientAddress2	string	`json:"PatientAddress2"`
	PatientCity	string	`json:"PatientCity"`
	PatientState	string	`json:"PatientState"`
	PatientPostalCode	string	`json:"PatientPostalCode"`
	PatientPrimaryInsuranceID	string	`json:"PatientPrimaryInsuranceID"`
	PatientSecondaryInsuranceID	string	`json:"PatientSecondaryInsuranceID"`
	PatientSocialSecurityNumber	string	`json:"PatientSocialSecurityNumber"`
	PatientMaritalStatus	string	`json:"PatientMaritalStatus"`
}

type Profile struct {
	CompanyKey	    string	            `json:"CompanyKey"`
	Name	        string	            `json:"ProfileName"`
	Key	            string	            `json:"ProfileKey"`
	StaffMembers	[]StaffMemberDetail	`json:"StaffMembers"`
	Tasks	        []Task	            `json:"Tasks"`
}

// Room is a room, hopefully
type Room struct {
	RoomKey	string	`json:"RoomKey"`
	RoomName	string	`json:"RoomName"`
	RoomAbbreviation	string	`json:"RoomAbbreviation"`
	RoomNumber	string	`json:"RoomNumber"`
	RoomFloor	string	`json:"RoomFloor"`
	LocationName	string	`json:"LocationName"`
	LocationAbbreviation	string	`json:"LocationAbbreviation"`
	LocationKey	float64	`json:"LocationKey"`
	ExternalId1	string	`json:"ExternalId1"`
	ExternalId2	string	`json:"ExternalId2"`
	ExternalId3	string	`json:"ExternalId3"`
	ExternalId4	string	`json:"ExternalId4"`
	ExternalId5	string	`json:"ExternalId5"`
	StartDate	time.Time	`json:"StartDate"`
	EndDate	time.Time		`json:"EndDate"`
	RoomShifts	[]interface{}	`json:"RoomShifts"`
	RoomTags	[]Tag	`json:"RoomTags"`
}


// StaffMemberDetail represents staff, and possibly some other entities as well?
type StaffMemberDetail struct {
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
	Tags     []Tag      `json:"Tags"`
	TTCMTags []Tag      `json:"TTCMTags"`
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
