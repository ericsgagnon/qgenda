package qgenda

// type DailyCase struct {
// 	 CompanyKey *string `json:"CompanyKey,omitempty"`
// 	 EmrId *string `json:"EmrId,omitempty"`
// 	 DailyCaseID *int64 `json:"DailyCaseID,omitempty"`
// 	 TaskKey GUID `json:"TaskKey,omitempty"`
// 	 Task *string `json:"Task,omitempty"`
// 	 Date Date `json:"Date,omitempty"`
// 	 StartTime TimeOfDay `json:"StartTime,omitempty"`
// 	 EndTime TimeOfDay `json:"EndTime,omitempty"`
// 	 Surgeon *string `json:"Surgeon,omitempty"`
// 	 Procedure *string `json:"Procedure,omitempty"`
// 	 SupStaffKey []any `json:"SupStaffKey,omitempty"`
// 	 SupStaffNpi []any `json:"SupStaffNpi,omitempty"`
// 	 SupStaff []any `json:"SupStaff,omitempty"`
// 	 DPStaffKey []any `json:"DPStaffKey,omitempty"`
// 	 DPStaff []any `json:"DPStaff,omitempty"`
// 	 DPStaffNpi []any `json:"DPStaffNpi,omitempty"`
// 	 CustomTextFields []any `json:"CustomTextFields,omitempty"`
// 	 CustomCheckboxFields []any `json:"CustomCheckboxFields,omitempty"`
// 	 IsCancelled *bool `json:"IsCancelled,omitempty"`
// 	 Duration Number `json:"Duration,omitempty"`
// 	 BaseUnits Number `json:"BaseUnits,omitempty"`
// 	 ModifierUnits Number `json:"ModifierUnits,omitempty"`
// 	 TimeUnits Number `json:"TimeUnits,omitempty"`
// 	 PatientAge *string `json:"PatientAge,omitempty"`
// 	 PatientGender *string `json:"PatientGender,omitempty"`
// 	 PatientEmail *string `json:"PatientEmail,omitempty"`
// 	 PatientHomePhone *string `json:"PatientHomePhone,omitempty"`
// 	 PatientCellPhone *string `json:"PatientCellPhone,omitempty"`
// 	 PatientAlternatePhone *string `json:"PatientAlternatePhone,omitempty"`
// 	 PatientGuardianPhone *string `json:"PatientGuardianPhone,omitempty"`
// 	 PatientDOB Date `json:"PatientDOB,omitempty"`
// 	 PatientMRN *string `json:"PatientMRN,omitempty"`
// 	 PatientClinicalNotes *string `json:"PatientClinicalNotes,omitempty"`
// 	 PatientExtraNotes *string `json:"PatientExtraNotes,omitempty"`
// 	 PatientAddress1 *string `json:"PatientAddress1,omitempty"`
// 	 PatientAddress2 *string `json:"PatientAddress2,omitempty"`
// 	 PatientCity *string `json:"PatientCity,omitempty"`
// 	 PatientState *string `json:"PatientState,omitempty"`
// 	 PatientPostalCode *string `json:"PatientPostalCode,omitempty"`
// 	 PatientPrimaryInsuranceID *string `json:"PatientPrimaryInsuranceID,omitempty"`
// 	 PatientSecondaryInsuranceID *string `json:"PatientSecondaryInsuranceID,omitempty"`
// 	 PatientSocialSecurityNumber *string `json:"PatientSocialSecurityNumber,omitempty"`
// 	 PatientMaritalStatus *string `json:"PatientMaritalStatus,omitempty"`
// 	 DailyCaseID *int64 `json:"DailyCaseID,omitempty"`
// }

// type Location struct {
// 	CompanyKey  *string       `json:"CompanyKey,omitempty"`
// 	LocationKey Int           `json:"LocationKey,omitempty"`
// 	Id          *string       `json:"Id,omitempty"`
// 	Abbrev      *string       `json:"Abbrev,omitempty"`
// 	Address     *string       `json:"Address,omitempty"`
// 	Tags        []TagCategory `json:"Tags,omitempty"`
// 	Notes       *string       `json:"Notes,omitempty"`
// 	TimeZone    *string       `json:"TimeZone,omitempty"`
// }

// type Company struct {
// 	CompanyKey     *string  `json:"CompanyKey,omitempty"`
// 	CompanyAbbr    *string  `json:"CompanyAbbr,omitempty"`
// 	DateCreatedUtc *Time `json:"DateCreatedUtc,omitempty"`
// }

type StaffLocation struct {
	Location Location `json:"Location:,omitempty"`
	//  CompanyKey GUID `json:"CompanyKey,omitempty"`
	//  LocationKey Int `json:"LocationKey,omitempty"`
	//  Id *string `json:"Id,omitempty"`
	//  Address *string `json:"Address,omitempty"`
	//  Abbrev *string `json:"Abbrev,omitempty"`
	//  Notes *string `json:"Notes,omitempty"`
	//  TimeZone *string (see Appendix) `json:"TimeZone,omitempty"`
	//  Tags []any `json:"Tags,omitempty"`
	Staff any `json:"Staff:,omitempty"`
	//  StaffKey GUID `json:"StaffKey,omitempty"`
	//  Id *string `json:"Id,omitempty"`
	//  Email *string `json:"Email,omitempty"`
	//  Tags []any `json:"Tags,omitempty"`
	//  IsCredentialed *bool `json:"IsCredentialed,omitempty"`
	//  InactiveDate Date `json:"InactiveDate,omitempty"`
	//  'Credentials: []any `json:"Credentials:,omitempty"`'
	//  IsPending *bool `json:"IsPending,omitempty"`
	//  StartDate Date `json:"StartDate,omitempty"`
	//  EndDate Date `json:"EndDate,omitempty"`
	//  Tags []any `json:"Tags,omitempty"`
	//  CategoryKey int `json:"CategoryKey,omitempty"`
	//  'Tags: []any `json:"Tags:,omitempty"`'
	//  Key int `json:"Key,omitempty"`
	//  EffectiveFromDate Datetime `json:"EffectiveFromDate,omitempty"`
	//  EffectiveToDate Datetime `json:"EffectiveToDate,omitempty"`
}

// type TagCategory struct {
// 	 CategoryKey int `json:"CategoryKey,omitempty"`
// 	 Tags []Tag `json:"Tags:,omitempty"`
// 	 Key int `json:"Key,omitempty"`
// }

// type TagDetailsByCompany struct {
// 	 CompanyKey *string `json:"CompanyKey,omitempty"`
// 	 'Tags: []any `json:"Tags:,omitempty"`'
// 	 CategoryKey int `json:"CategoryKey,omitempty"`
// 	 CategoryDateCreated DateTime `json:"CategoryDateCreated,omitempty"`
// 	 CategoryDateModified DateTime `json:"CategoryDateModified,omitempty"`
// 	 IsAvailableForCreditAllocation *bool `json:"IsAvailableForCreditAllocation,omitempty"`
// 	 IsAvailableForDailySum *bool `json:"IsAvailableForDailySum,omitempty"`
// 	 IsAvailableForHoliday *bool `json:"IsAvailableForHoliday,omitempty"`
// 	 IsAvailableForLocation *bool `json:"IsAvailableForLocation,omitempty"`
// 	 IsAvailableForProfile *bool `json:"IsAvailableForProfile,omitempty"`
// 	 IsAvailableForSeries *bool `json:"IsAvailableForSeries,omitempty"`
// 	 IsAvailableForStaff *bool `json:"IsAvailableForStaff,omitempty"`
// 	 IsAvailableForStaffLocation *bool `json:"IsAvailableForStaffLocation,omitempty"`
// 	 IsAvailableForStaffTarget *bool `json:"IsAvailableForStaffTarget,omitempty"`
// 	 IsAvailableForRequestLimit *bool `json:"IsAvailableForRequestLimit,omitempty"`
// 	 IsAvailableForTask *bool `json:"IsAvailableForTask,omitempty"`
// 	 IsTTCMCategory *bool `json:"IsTTCMCategory,omitempty"`
// 	 IsSingleTaggingOnly *bool `json:"IsSingleTaggingOnly,omitempty"`
// 	 IsPermissionCategory *bool `json:"IsPermissionCategory,omitempty"`
// 	 IsUsedForStats *bool `json:"IsUsedForStats,omitempty"`
// 	 IsUsedForFiltering *bool `json:"IsUsedForFiltering,omitempty"`
// 	 CategoryBackgroundColor *string `json:"CategoryBackgroundColor,omitempty"`
// 	 CategoryTextColor *string `json:"CategoryTextColor,omitempty"`
// 	 'Tags: []any `json:"Tags:,omitempty"`'
// 	 Key int `json:"Key,omitempty"`
// 	 DateCreated DateTime `json:"DateCreated,omitempty"`
// 	 DateLastModified DateTime `json:"DateLastModified,omitempty"`
// 	 BackgroundColor *string `json:"BackgroundColor,omitempty"`
// 	 TextColor *string `json:"TextColor,omitempty"`
// }

type Skillset struct {
	StaffAbbreviation *string `json:"StaffAbbreviation,omitempty"`
	StaffId           *string `json:"StaffId,omitempty"`
	TaskAbbrev        *string `json:"TaskAbbrev,omitempty"`
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

type ScheduleRotation struct {
	StartDate   *Date   `json:"StartDate,omitempty"`
	EndDate     *Date   `json:"EndDate,omitempty"`
	IsPublished *bool   `json:"IsPublished,omitempty"`
	IsLocked    *bool   `json:"IsLocked,omitempty"`
	StaffKey    *string `json:"StaffKey,omitempty"`
	TaskKey     *string `json:"TaskKey,omitempty"`
	TimeZone    *string `json:"TimeZone,omitempty"`
}

type Company struct {
	CompanyKey          *string   `json:"CompanyKey,omitempty"`
	CompanyAbbreviation *string   `json:"CompanyAbbreviation,omitempty"`
	CompanyLocation     *string   `json:"CompanyLocation,omitempty"`
	CompanyPhoneNumber  *string   `json:"CompanyPhoneNumber,omitempty"`
	Profiles            []Profile `json:"Profiles:,omitempty"`
	Organizations       []any     `json:"Organizations:,omitempty"`
	//  ProfileKey GUID `json:"ProfileKey,omitempty"`
	//  IsAdmin *bool `json:"IsAdmin,omitempty"`
	//  OrgKey *int64 `json:"OrgKey,omitempty"`
}

type Profile struct {
	Name         *string `json:"Name,omitempty"`
	CompanyKey   *string `json:"CompanyKey,omitempty"`
	ProfileKey   *string `json:"ProfileKey,omitempty"`
	StaffMembers []Staff `json:"StaffMembers,omitempty"`
	// StaffAbbreviation *string `json:"StaffAbbreviation,omitempty"`
	// StaffId           *string `json:"StaffId,omitempty"`
	// StaffKey          *string `json:"StaffKey,omitempty"`
	IsViewable    *bool `json:"IsViewable,omitempty"`
	IsSchedulable *bool `json:"IsSchedulable,omitempty"`
	Tasks         []any `json:"Tasks,omitempty"`
	// TaskAbbreviation  *string `json:"TaskAbbreviation,omitempty"`
	// TaskId            *string `json:"TaskId,omitempty"`
	// TaskKey           *string `json:"TaskKey,omitempty"`
	// IsViewable        *bool   `json:"IsViewable,omitempty"`
	// IsSchedulable     *bool   `json:"IsSchedulable,omitempty"`
}

type User struct {
	Email                   *string    `json:"Email,omitempty"`
	LastLoginDate           *Time      `json:"LastLoginDate,omitempty"`
	Companies               []Company  `json:"Companies,omitempty"`
	CompanyKey              *string    `json:"CompanyKey,omitempty"`
	IsRegistered            *bool      `json:"IsRegistered,omitempty"`
	ProfileLastModifiedDate *Time      `json:"ProfileLastModifiedDate,omitempty"`
	ProfileLastModifiedBy   *string    `json:"ProfileLastModifiedBy,omitempty"`
	Type                    *string    `json:"Type,omitempty"`
	StaffAbbreviation       *string    `json:"StaffAbbreviation,omitempty"`
	Locations               []Location `json:"Locations,omitempty"`
	LocationID              *string    `json:"LocationID,omitempty"`
}

type NotificationList struct {
	CompanyKey         *string `json:"CompanyKey,omitempty"`
	NotificationListID *int64  `json:"NotificationListID,omitempty"`
	ContactEmails      []any   `json:"ContactEmails,omitempty"`
	ContactEmail       *string `json:"ContactEmail,omitempty"`
	Tags               []any   `json:"Tags,omitempty"`
	TagCategory        *string `json:"TagCategory,omitempty"`
}

type TimeEvent struct {
	ActualClockIn          *Time   `json:"ActualClockIn,omitempty"`
	ActualClockOut         *Time   `json:"ActualClockOut,omitempty"`
	CompanyKey             *string `json:"CompanyKey,omitempty"`
	Date                   *Date   `json:"Date,omitempty"`
	DayOfWeek              *string `json:"DayOfWeek,omitempty"`
	Duration               *int64  `json:"Duration,omitempty"`
	EffectiveClockIn       *Time   `json:"EffectiveClockIn,omitempty"`
	EffectiveClockOut      *Time   `json:"EffectiveClockOut,omitempty"`
	IsClockInGeoVerified   *bool   `json:"IsClockInGeoVerified,omitempty"`
	IsClockOutGeoVerified  *bool   `json:"IsClockOutGeoVerified,omitempty"`
	IsEarly                *bool   `json:"IsEarly,omitempty"`
	IsLate                 *bool   `json:"IsLate,omitempty"`
	IsExcessiveDuration    *bool   `json:"IsExcessiveDuration,omitempty"`
	IsExtended             *bool   `json:"IsExtended,omitempty"`
	IsStruck               *bool   `json:"IsStruck,omitempty"`
	IsUnplanned            *bool   `json:"IsUnplanned,omitempty"`
	FlagsResolved          *bool   `json:"FlagsResolved,omitempty"`
	Notes                  *string `json:"Notes,omitempty"`
	ReasonCode             *string `json:"ReasonCode,omitempty"`
	ReasonCodeId           *string `json:"ReasonCodeId,omitempty"`
	ScheduleEntry          *string `json:"ScheduleEntry,omitempty"`
	ScheduleEntryKey       *string `json:"ScheduleEntryKey,omitempty"`
	StaffKey               *string `json:"StaffKey,omitempty"`
	TaskKey                *string `json:"TaskKey,omitempty"`
	TaskShiftKey           *string `json:"TaskShiftKey,omitempty"`
	TimePunchEventKey      *int64  `json:"TimePunchEventKey,omitempty"`
	TimeZone               *string `json:"TimeZone,omitempty"`
	ActualClockInLocal     *Time   `json:"ActualClockInLocal,omitempty"`
	ActualClockOutLocal    *Time   `json:"ActualClockOutLocal,omitempty"`
	EffectiveClockInLocal  *Time   `json:"EffectiveClockInLocal,omitempty"`
	EffectiveClockOutLocal *Time   `json:"EffectiveClockOutLocal,omitempty"`
	LastModifiedDate       *Time   `json:"LastModifiedDate,omitempty"`
}
