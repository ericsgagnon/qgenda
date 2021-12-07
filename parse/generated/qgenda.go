package qgenda



type Assignment struct {
	EncounterRoleKey	string
	StaffMemberKey	string
	StaffMemberAbbreviation	string
	RoleName	string
	StaffMemberID	string
	StaffMemberNPI	string
	StaffMemberEMRID	string
	IsNonStaff	boolean
	IsPublished	boolean
	IsDefaultRoomAssignment	boolean
}


type Company struct {
	CompanyKey	string
	CompanyName	string
	CompanyAbbreviation	string
	CompanyLocation	string
	CompanyPhoneNumber	string
	Profiles	[]Profiles
	ProfileName	string
	ProfileKey	string
	IsAdmin	boolean
	Organizations	[]Organizations
	OrgName	string
	OrgKey	int
}


type DailyCase struct {
	CompanyKey	string
	EmrId	string
	DailyCaseID	int
	TaskKey	string
	Task	string
	LocationName	string
	Date	time.Time
	StartTime	time.Time
	EndTime	time.Time
	Surgeon	string
	Procedure	string
	SupStaffKey	string
	SupStaffNpi	string
	SupStaff	string
	DPStaffKey	string
	DPStaff	string
	DPStaffNpi	string
	CustomTextFields	string
	CustomCheckboxFields	string
	IsCancelled	boolean
	Duration	float64
	BaseUnits	float64
	ModifierUnits	float64
	TimeUnits	float64
	PatientFirstName	string
	PatientLastName	string
	PatientAge	string
	PatientGender	string
	PatientEmail	string
	PatientHomePhone	string
	PatientCellPhone	string
	PatientAlternatePhone	string
	PatientGuardianPhone	string
	PatientGuardianName	string
	PatientDOB	time.Time
	PatientMRN	string
	PatientClinicalNotes	string
	PatientExtraNotes	string
	PatientAddress1	string
	PatientAddress2	string
	PatientCity	string
	PatientState	string
	PatientPostalCode	string
	PatientPrimaryInsuranceID	string
	PatientSecondaryInsuranceID	string
	PatientSocialSecurityNumber	string
	PatientMaritalStatus	string
	DailyCaseID	int
}


type DailyConfiguration struct {
	DailyConfigurationKey	string
	Name	string
	EncounterRoles	[]EncounterRoles
	EncounterFieldSettings	[]EncounterFieldSettings
	EncounterPhiFieldSettings	[]EncounterPhiFieldSettings
}


type EncounterRole struct {
	EncounterRoleKey	string
	Name	string
}


type Location struct {
	CompanyKey	string
	LocationKey	int
	Id	string
	Name	string
	Abbrev	string
	Address	string
	Tags	[]Tags
	Notes	string
	TimeZone	string
}


type NotificationList struct {
	Name	string
	CompanyKey	string
	NotificationListID	int
	ContactEmails	[]ContactEmails
	ContactName	string
	ContactEmail	string
	Tags	[]Tags
	TagCategory	string
	TagName	string
}


type OpenShift struct {
	CompanyKey	string
	ScheduleKey	string
	OpenShiftCount	int
	CallRole	string
	Credit	float64
	Date	time.Time
	StartDate	time.Time
	StartDateUTC	time.Time
	StartTime	time.Time
	EndDate	time.Time
	EndDateUTC	time.Time
	EndTime	time.Time
	IsCred	boolean
	IsSaved	boolean
	IsPublished	boolean
	IsLocked	boolean
	IsStruck	boolean
	Notes	string
	IsNotePrivate	boolean
	TaskAbbrev	string
	TaskId	string
	TaskEmrId	string
	TaskIsPrintStart	boolean
	TaskIsPrintEnd	boolean
	TaskExtCallSysId	string
	TaskKey	string
	TaskName	string
	TaskBillSysId	string
	TaskPayrollId	string
	TaskShiftKey	string
	TaskType	string
	TaskContactInformation	string
	TaskTags	[]TaskTags
	LocationKey	string
	LocationName	string
	LocationAbbrev	string
	LocationAddress	string
	LocationTags	[]LocationTags
	TimeZone	string
}


type PatientEncounter struct {
	CompanyKey	string
	DailyConfigurationKey	string
	EncounterKey	string
	DailyCaseID	int
	Date	time.Time
	StartTime	time.Time
	EndTime	time.Time
	Duration	float64
	IsCancelled	boolean
	Room	room
	RoomKey	string
	RoomName	string
	LocationName	string
	LocationKey	string
	Assignments	[]Assignments
	ExternalID1	string
	ExternalID2	string
	CreatedBySource	string
	ModifiedBySource	string
}


type PatientInformation struct {
	PatientFirstName	string
	PatientLastName	string
	PatientGender	string
	PatientEmail	string
	PatientHomePhone	string
	PatientCellPhone	string
	PatientAlternatePhone	string
	PatientGuardianPhone	string
	PatientGuardianName	string
	PatientDateOfBirth	time.Time
	PatientMedicalRecordNumber	string
	PatientClinicalNotes	string
	PatientExtraNotes	string
	PatientAddress1	string
	PatientAddress2	string
	PatientCity	string
	PatientState	string
	PatientPostalCode	string
	PatientPrimaryInsuranceId	string
	PatientSecondaryInsuranceId	string
	PatientSocialSecurityNumber	string
	PatientMaritalStatus	string
}


type PayRate struct {
	CompanyKey	string
	PayRateKey	int
	PayCodeLabel	string
	PayCodeID	string
	Rate	float64
	RateType	string
	LocationId	string
	LocationName	string
	StaffID	string
	StaffPayrollID	string
	StaffLastName	string
	StaffFirstName	string
	StaffAbbreviation	string
	StartDate	time.Time
	EndDate	time.Time
	StaffTagCategoryName	string
	StaffTagName	string
}


type Profile struct {
	CompanyKey	string
	ProfileName	string
	ProfileKey	string
	StaffMembers	[]StaffMembers
	StaffFirstName	string
	StaffLastName	string
	StaffAbbreviation	string
	StaffId	string
	StaffKey	string
	IsViewable	boolean
	IsSchedulable	boolean
	Tasks	[]Tasks
	TaskName	string
	TaskAbbreviation	string
	TaskId	string
	TaskKey	string
	IsViewable	boolean
	IsSchedulable	boolean
}


type Request struct {
	RequestId	string
	StaffNpi	string
	TaskId	string
	CallRole	string
	CompKey	string
	Credit	float64
	Date	time.Time
	EndTime	time.Time
	Notes	string
	RequestKey	string
	RequestType	string
	RequestStatus	string
	ResolutionDate	time.Time
	StaffAbbrev	string
	StaffFName	string
	StaffKey	string
	StaffLName	string
	StaffPayrollId	string
	StartTime	time.Time
	SubmittedByUser	string
	SubmittedDate	time.Time
	TaskAbbrev	string
	TaskKey	string
	TaskName	string
	StaffInternalId	string
	Removed	boolean
}


type RequestApproved struct {
	nextSyncToken	string
	items[]	[]items
	CompanyKey	string
	RequestId	string
	RequestKey	string
	RequestType	string
	Date	string
	StartTime	string
	EndTime	string
	Credit	int
	TaskKey	string
	TaskName	string
	TaskAbbrev	string
	TaskId	string
	CallRole	string
	StaffKey	string
	StaffFName	string
	StaffLName	string
	StaffAbbrev	string
	StaffInternalId	string
	StaffPayrollId	string
	StaffNpi	string
	Notes	string
	SubmittedByUser	string
	SubmittedDate	string
	ResolutionDate	string
	FormFields[]	[]FormFields
	Name	string
	Value	string
	RequestLimitValues[]	[]RequestLimitValues
	RequestLimitKey	string
	Value	int
	RequestNotes[]	[]RequestNotes
	Timestamp	string
	UserName	string
	Note	string
	Removed	boolean
}


type RequestLimit struct {
	CompanyKey	string
	Name	string
	Key	string
	StartDate	time.Time
	EndDate	time.Time
	Type	string
	RecurringLength	string
	CreditSource	string
	ErrorMessage	string
	IsActive	boolean
	DailyTotalMaxAllowedMon	float64
	DailyTotalMaxAllowedTue	float64
	DailyTotalMaxAllowedWed	float64
	DailyTotalMaxAllowedThu	float64
	DailyTotalMaxAllowedFri	float64
	DailyTotalMaxAllowedSat	float64
	DailyTotalMaxAllowedSun	float64
	ShiftsCredit	[]ShiftsCredit
	TaskId	string
	TaskKey	string
	TaskName	string
	TaskAbbreviation	string
	TaskshiftKey	string
	IsIncluded	boolean
	DayOfTheWeek	string
	Credit	float64
	StaffLimits	[]StaffLimits
	StaffId	string
	StaffInternalId	string
	StaffKey	string
	StaffFirstName	string
	StaffLastName	string
	StaffAbbreviation	string
	StaffRequestLimits	float64
	StaffTotalLimit	float64
}


type Room struct {
	RoomKey	string
	RoomName	string
	RoomAbbreviation	string
	RoomNumber	string
	RoomFloor	string
	LocationName	string
	LocationAbbreviation	string
	LocationKey	float64
	ExternalId1	string
	ExternalId2	string
	ExternalId3	string
	ExternalId4	string
	ExternalId5	string
	StartDate	time.Time
	EndDate	time.Time
	RoomShifts	[]RoomShifts
	RoomTags	[]RoomTags
}


type ScheduleAuditLogEntry struct {
	StaffFirstName	string
	StaffLastName	string
	StaffAbbreviation	string
	StaffKey	string
	TaskName	string
	TaskAbbreviation	string
	TaskKey	string
	ScheduleEntryDate	time.Time
	ScheduleEntryStartTimeUTC	time.Time
	ScheduleEntryStartTimeLocal	time.Time
	ScheduleEntryEndTimeUTC	time.Time
	ScheduleEntryEndTimeLocal	time.Time
	ScheduleEntryKey	string
	ActivityType	string
	SourceType	string
	UserFirstName	string
	UserLastName	string
	UserKey	string
	TimestampUTC	time.Time
	Timestamp	time.Time
	AdditionalInformation	string
	Location(Array)	[]Location
	LocationName	string
	LocationAbbreviation	string
	LocationID	string
	LocationKey	string
	IPAddress	string
}


type ScheduleEntry struct {
	ScheduleKey	string
	CallRole	string
	CompKey	string
	Credit	float64
	Date	time.Time
	StartDateUTC	time.Time
	EndDateUTC	time.Time
	EndDate	time.Time
	EndTime	time.Time
	IsCred	boolean
	IsSaved	boolean
	IsPublished	boolean
	IsLocked	boolean
	IsStruck	boolean
	Notes	string
	IsNotePrivate	boolean
	StaffAbbrev	string
	StaffBillSysId	string
	StaffEmail	string
	StaffEmrId	string
	StaffErpId	string
	StaffInternalId	string
	StaffExtCallSysId	string
	StaffFName	string
	StaffId	string
	StaffKey	string
	StaffLName	string
	StaffMobilePhone	string
	StaffNpi	string
	StaffPager	string
	StaffPayrollId	string
	StaffTags	[]StaffTags
	StartDate	time.Time
	StartTime	time.Time
	TaskAbbrev	string
	TaskBillSysId	string
	TaskContactInformation	string
	TaskExtCallSysId	string
	TaskId	string
	TaskEmrId	string
	TaskCallPriority	string
	TaskDepartmentId	string
	TaskIsPrintEnd	boolean
	TaskIsPrintStart	boolean
	TaskKey	string
	TaskName	string
	TaskPayrollId	string
	TaskShiftKey	string
	TaskType	string
	TaskTags	[]TaskTags
	LocationName	string
	LocationAbbrev	string
	LocationID	string
	LocationAddress	string
	LocationTags	[]LocationTags
	TimeZone	string
	LastModifiedDateUTC	time.Time
	IsRotationTask	boolean
}


type ScheduleRotation struct {
	StartDate	time.Time
	EndDate	time.Time
	IsPublished	boolean
	IsLocked	boolean
	StaffKey	string
	TaskKey	string
	TimeZone	string
}


type Skillset struct {
	StaffFirstName	string
	StaffLastName	string
	StaffAbbreviation	string
	StaffId	string
	TaskName	string
	TaskAbbrev	string
	TaskId	string
	IsSkilledMon	boolean
	MonOccurrence	string
	IsSkilledTue	boolean
	TueOccurrence	string
	IsSkilledWed	boolean
	WedOccurrence	string
	IsSkilledThu	boolean
	ThuOccurrence	string
	IsSkilledFri	boolean
	FriOccurrence	string
	IsSkilledSat	boolean
	SatOccurrence	string
	IsSkilledSun	boolean
	SunOccurrence	string
}


type StaffLocation struct {
	Location	object
	CompanyKey	string
	LocationKey	int
	Id	string
	Name	string
	Address	string
	Abbrev	string
	Notes	string
	TimeZone	string
	Tags	[]Tags
	Staff	
	StaffKey	string
	Id	string
	FirstName	string
	LastName	string
	Email	string
	Tags	[]Tags
	IsCredentialed	boolean
	InactiveDate	time.Time
	Credentials	[]Credentials
	IsPending	boolean
	StartDate	time.Time
	EndDate	time.Time
	Tags	[]Tags
	CategoryKey	int
	CategoryName	string
	Tags	[]Tags
	Key	int
	Name	string
	EffectiveFromDate	time.Time
	EffectiveToDate	time.Time
}


type StaffMemberDetail struct {
	Abbrev	string
	BgColor	string
	BillSysId	string
	CompKey	string
	Contact Instructions	string
	Email	string
	SsoId	string
	EmrId	string
	ErpId	string
	EndDate	time.Time
	ExtCallSysId	string
	FirstName	string
	HomePhone	string
	LastName	string
	MobilePhone	string
	Npi	string
	OtherNumber1	string
	OtherNumber2	string
	OtherNumber3	string
	OtherNumberType1	string
	OtherNumberType2	string
	OtherNumberType3	string
	Pager	string
	PayrollId	string
	RegHours	float64
	StaffId	string
	StaffKey	string
	StartDate	time.Time
	TextColor	string
	Addr1	string
	Addr2	string
	City	string
	State	string
	Zip	string
	IsActive	boolean
	StaffTypeKey	string
	BillingTypeKey	string
	UserProfileKey	string
	UserProfile	string
	PayrollStartDate	time.Time
	PayrollEndDate	time.Time
	TimeClockStartDate	time.Time
	TimeClockEndDate	time.Time
	TimeClockKioskPIN	string
	IsAutoApproveSwap	boolean
	DailyUnitAverage	float64
	StaffInternalId	string
	UserLastLoginDateTimeUTC	time.Time
	SourceOfLogin	string
	CalSyncKey	string
	Tags	[]Tags
	TTCMTags	[]TTCMTags
	CategoryKey	int
	CategoryName	string
	Tags	[]Tags
	Key	int
	Name	string
	EffectiveFromDate	time.Time
	EffectiveToDate	time.Time
	Skillset	[]Skillset
	Profiles	[]Profiles
	Name	string
	ProfileKey	string
	IsViewable	boolean
	IsSchedulable	boolean
}


type StaffTarget struct {
	CompanyKey	string
	Name	string
	Key	int
	OccurrenceNumber	int
	OccurrenceType	string
	StartDate	time.Time
	EndDate	time.Time
	IncludeInAlgorithm	boolean
	TargetType	string
	WeekDefinitionType	string
	CountType	string
	DefaultMin	float64
	DefaultIdeal	float64
	DefaultMax	float64
	DefaultMin5WeekMonth	float64
	DefaultIdeal5WeekMonth	float64
	DefaultMax5WeekMonth	float64
	ManualValidationCheck	string
	ErrorMessage	string
	Tasks	[]Tasks
	TaskName	string
	TaskAbbreviation	string
	TaskKey	string
	TaskID	string
	TaskShiftKey	string
	DayOfTheWeek	string
	Staff	[]Staff
	StaffFirstName	string
	StaffLastName	string
	StaffAbbreviation	string
	StaffId	string
	Targets	[]Targets
	EffectiveDate	time.Time
	Min	float64
	Ideal	float64
	Max	float64
	Min5WeekMonth	float64
	Ideal5WeekMonth	float64
	Max5WeekMonth	float64
	IsActive	boolean
	Profiles	[]Profiles
	ProfileName	string
	ProfileKey	string
	ShowOnBottomPanel	boolean
	ShowOnStatsTab	boolean
	ValidateOnAdminTab	boolean
	Locations	[]Locations
	LocationName	string
	LocationID	string
}


type StandardFields struct {
	EncounterFieldKey	string
	Name	string
	Value	string
	Type	string
}


type Tag struct {
	CategoryKey	int
	CategoryName	string
	Tags	[]Tags
	Key	int
	Name	string
}


type TagDetailsByCompany struct {
	CompanyName	string
	CompanyKey	string
	Tags	[]Tags
	CategoryKey	int
	CategoryName	string
	CategoryDateCreated	time.Time
	CategoryDateModified	time.Time
	IsAvailableForCreditAllocation	boolean
	IsAvailableForDailySum	boolean
	IsAvailableForHoliday	boolean
	IsAvailableForLocation	boolean
	IsAvailableForProfile	boolean
	IsAvailableForSeries	boolean
	IsAvailableForStaff	boolean
	IsAvailableForStaffLocation	boolean
	IsAvailableForStaffTarget	boolean
	IsAvailableForRequestLimit	boolean
	IsAvailableForTask	boolean
	IsTTCMCategory	boolean
	IsSingleTaggingOnly	boolean
	IsPermissionCategory	boolean
	IsUsedForStats	boolean
	IsUsedForFiltering	boolean
	CategoryBackgroundColor	string
	CategoryTextColor	string
	Tags	[]Tags
	Key	int
	Name	string
	DateCreated	time.Time
	DateLastModified	time.Time
	BackgroundColor	string
	TextColor	string
}


type Task struct {
	Abbrev	string
	BgColor	string
	BillSysId	string
	CallRole	string
	CompKey	string
	ContactInformation	string
	DepartmentId	string
	DisplayAsAvailableOnCalendarConnection	boolean
	EmrId	string
	EnableWaitlist	boolean
	EndDate	time.Time
	EnforceTimeCompatibility	boolean
	ExtCallPriority	string
	ExtCallSysId	string
	Extension	string
	FullyCompatible	boolean
	HideOnCalendarSync	boolean
	Label	boolean
	Manual	boolean
	Name	string
	Notes	string
	NotificationList	string
	OverrideSortOrder	string
	PayrollId	string
	PayrollUnits	string
	ShowEndTime	boolean
	ShowOnAdminOnly	boolean
	TaskVisibility	string
	ShowOpensOnCalendar	boolean
	ShowStartTime	boolean
	ShowTimesWhenEdited	boolean
	StartDate	time.Time
	SuppressRequestEmailsAndNotifications	boolean
	SyncAsAllDayEvent	boolean
	TaskId	string
	TaskKey	string
	RequireTimePunch	boolean
	TextColor	string
	TimeOffAfterShift	int
	Type	string
	Tags	[]Tags
	TaskShifts	[]TaskShifts
	Profiles	[]Profiles
	ProfileName	string
	ProfileKey	string
	IsViewable	boolean
	IsSchedulable	boolean
	IsRotationTask	boolean
	IsIgnoreHoursOffAfterTask	boolean
}


type TaskShift struct {
	DayOfWeek	string
	StartTime	time.Time
	EndTime	time.Time
	MaxStaff	int
	MinStaff	int
	OffAfter	int
	StatCredit	float64
	StaffCount	int
	EffectiveDate	[]EffectiveDate
	EffectiveFromDate	time.Time
	EffectiveToDate	time.Time
	StartTime	time.Time
	EndTime	time.Time
	MaxStaff	int
	MinStaff	int
	OffAfter	int
	StatCredit	float64
	IsActive	boolean
}


type TimeEvent struct {
	ActualClockIn	time.Time
	ActualClockOut	time.Time
	CompanyKey	string
	Date	time.Time
	DayOfWeek	string
	Duration	int
	EffectiveClockIn	time.Time
	EffectiveClockOut	time.Time
	IsClockInGeoVerified	boolean
	IsClockOutGeoVerified	boolean
	IsEarly	boolean
	IsLate	boolean
	IsExcessiveDuration	boolean
	IsExtended	boolean
	IsStruck	boolean
	IsUnplanned	boolean
	FlagsResolved	boolean
	Notes	string
	ReasonCode	string
	ReasonCodeId	string
	ScheduleEntry	string
	ScheduleEntryKey	string
	StaffKey	string
	TaskKey	string
	TaskShiftKey	string
	TimePunchEventKey	int
	TimeZone	string
	ActualClockInLocal	time.Time
	ActualClockOutLocal	time.Time
	EffectiveClockInLocal	time.Time
	EffectiveClockOutLocal	time.Time
	LastModifiedDate	time.Time
}


type User struct {
	Email	string
	FirstName	string
	LastName	string
	LastLoginDate	time.Time
	Companies	[]Companies
	CompanyKey	string
	IsRegistered	boolean
	ProfileName	string
	ProfilePermissionConfigurationName	string
	ProfileLastModifiedDate	time.Time
	ProfileLastModifiedBy	string
	Type	string
	StaffFirstName	string
	StaffLastName	string
	StaffAbbreviation	string
	Locations	[]Locations
	LocationID	string
	LocationName	string
}
