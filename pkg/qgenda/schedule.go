package qgenda

import (
	"time"
)

type Schedule struct {
	ScheduleKey  *string    `json:"ScheduleKey,omitempty"`
	CallRole     *string    `json:"CallRole,omitempty"`
	CompKey      *string    `json:"CompKey,omitempty"`
	Credit       *float64   `json:"Credit,omitempty"`
	Date         *Date      `json:"Date,omitempty"`
	StartDateUTC *time.Time `json:"StartDateUTC,omitempty"`
	EndDateUTC   *time.Time `json:"EndDateUTC,omitempty"`
	EndDate      *Date      `json:"EndDate,omitempty"`
	EndTime      *TimeOfDay `json:"EndTime,omitempty"`
	// Date         *string    `json:"Date,omitempty"`
	// StartDateUTC *string `json:"StartDateUTC,omitempty"`
	// EndDateUTC   *string `json:"EndDateUTC,omitempty"`
	// EndDate      *string    `json:"EndDate,omitempty"`
	// EndTime           *string `json:"EndTime,omitempty"`
	IsCred            *bool      `json:"IsCred,omitempty"`
	IsPublished       *bool      `json:"IsPublished,omitempty"`
	IsLocked          *bool      `json:"IsLocked,omitempty"`
	IsStruck          *bool      `json:"IsStruck,omitempty"`
	Notes             *string    `json:"Notes,omitempty"`
	IsNotePrivate     *bool      `json:"IsNotePrivate,omitempty"`
	StaffAbbrev       *string    `json:"StaffAbbrev,omitempty"`
	StaffBillSysId    *string    `json:"StaffBillSysId,omitempty"`
	StaffEmail        *string    `json:"StaffEmail,omitempty"`
	StaffEmrId        *string    `json:"StaffEmrId,omitempty"`
	StaffErpId        *string    `json:"StaffErpId,omitempty"`
	StaffInternalId   *string    `json:"StaffInternalId,omitempty"`
	StaffExtCallSysId *string    `json:"StaffExtCallSysId,omitempty"`
	StaffFName        *string    `json:"StaffFName,omitempty"`
	StaffId           *string    `json:"StaffId,omitempty"`
	StaffKey          *string    `json:"StaffKey,omitempty"`
	StaffLName        *string    `json:"StaffLName,omitempty"`
	StaffMobilePhone  *string    `json:"StaffMobilePhone,omitempty"`
	StaffNpi          *string    `json:"StaffNpi,omitempty"`
	StaffPager        *string    `json:"StaffPager,omitempty"`
	StaffPayrollId    *string    `json:"StaffPayrollId,omitempty"`
	StaffTags         []any      `json:"StaffTags,omitempty"`
	StartDate         *Date      `json:"StartDate,omitempty"`
	StartTime         *TimeOfDay `json:"StartTime,omitempty"`
	// StartDate *string    `json:"StartDate,omitempty"`
	// StartTime              *string `json:"StartTime,omitempty"`
	TaskAbbrev             *string `json:"TaskAbbrev,omitempty"`
	TaskBillSysId          *string `json:"TaskBillSysId,omitempty"`
	TaskContactInformation *string `json:"TaskContactInformation,omitempty"`
	TaskExtCallSysId       *string `json:"TaskExtCallSysId,omitempty"`
	TaskId                 *string `json:"TaskId,omitempty"`
	TaskKey                *string `json:"TaskKey,omitempty"`
	TaskName               *string `json:"TaskName,omitempty"`
	TaskPayrollId          *string `json:"TaskPayrollId,omitempty"`
	TaskEmrId              *string `json:"TaskEmrId,omitempty"`
	TaskCallPriority       *string `json:"TaskCallPriority,omitempty"`
	TaskDepartmentId       *string `json:"TaskDepartmentId,omitempty"`
	TaskIsPrintEnd         *bool   `json:"TaskIsPrintEnd,omitempty"`
	TaskIsPrintStart       *bool   `json:"TaskIsPrintStart,omitempty"`
	TaskShiftKey           *string `json:"TaskShiftKey,omitempty"`
	TaskType               *string `json:"TaskType,omitempty"`
	TaskTags               []any   `json:"TaskTags,omitempty"`
	LocationName           *string `json:"LocationName,omitempty"`
	LocationAbbrev         *string `json:"LocationAbbrev,omitempty"`
	LocationID             *string `json:"LocationID,omitempty"`
	LocationAddress        *string `json:"LocationAddress,omitempty"`
	TimeZone               *string `json:"TimeZone,omitempty"`
	// LastModifiedDateUTC    *time.Time `json:"LastModifiedDateUTC,omitempty"`
	LastModifiedDateUTC *string `json:"LastModifiedDateUTC,omitempty"`
	LocationTags        []any   `json:"LocationTags,omitempty"`
	IsRotationTask      *bool   `json:"IsRotationTask"`
}

type ScheduleAuditLog struct {
	StaffFirstName            *string       `json:"StaffFirstName,omitempty"`
	StaffLastName             *string       `json:"StaffLastName,omitempty"`
	StaffAbbreviation         *string       `json:"StaffAbbreviation,omitempty"`
	StaffKey                  *string       `json:"StaffKey,omitempty"`
	TaskName                  *string       `json:"TaskName,omitempty"`
	TaskAbbreviation          *string       `json:"TaskAbbreviation,omitempty"`
	TaskKey                   *string       `json:"TaskKey,omitempty"`
	ScheduleEntryDate         *string       `json:"ScheduleEntryDate,omitempty"`
	ScheduleEntryStartTimeUTC *time.Time    `json:"ScheduleEntryStartTimeUTC,omitempty"`
	ScheduleEntryStartTime    *string       `json:"ScheduleEntryStartTime,omitempty"`
	ScheduleEntryEndTimeUTC   *time.Time    `json:"ScheduleEntryEndTimeUTC,omitempty"`
	ScheduleEntryEndTime      *string       `json:"ScheduleEntryEndTime,omitempty"`
	ScheduleEntryKey          *string       `json:"ScheduleEntryKey,omitempty"`
	ActivityType              *string       `json:"ActivityType,omitempty"`
	SourceType                *string       `json:"SourceType,omitempty"`
	UserFirstName             *string       `json:"UserFirstName,omitempty"`
	UserLastName              *string       `json:"UserLastName,omitempty"`
	UserKey                   *string       `json:"UserKey,omitempty"`
	TimestampUTC              *string       `json:"TimestampUTC,omitempty"`
	Timestamp                 *string       `json:"Timestamp,omitempty"`
	AdditionalInformation     *string       `json:"AdditionalInformation,omitempty"`
	Locations                 []interface{} `json:"Locations,omitempty"`
	IPAddress                 *string       `json:"IPAddress,omitempty"`
}

type OpenShift struct {
	CompanyKey             *string    `json:"CompanyKey,omitempty"`
	ScheduleKey            *string    `json:"ScheduleKey,omitempty"`
	OpenShiftCount         *int64     `json:"OpenShiftCount,omitempty"`
	CallRole               *string    `json:"CallRole,omitempty"`
	Credit                 *float64   `json:"Credit,omitempty"`
	Date                   *time.Time `json:"Date,omitempty"`
	StartDate              *time.Time `json:"StartDate,omitempty"`
	StartDateUTC           *time.Time `json:"StartDateUTC,omitempty"`
	StartTime              *time.Time `json:"StartTime,omitempty"`
	EndDate                *time.Time `json:"EndDate,omitempty"`
	EndDateUTC             *time.Time `json:"EndDateUTC,omitempty"`
	EndTime                *time.Time `json:"EndTime,omitempty"`
	IsCred                 *bool      `json:"IsCred,omitempty"`
	IsSaved                *bool      `json:"IsSaved,omitempty"`
	IsPublished            *bool      `json:"IsPublished,omitempty"`
	IsLocked               *bool      `json:"IsLocked,omitempty"`
	IsStruck               *bool      `json:"IsStruck,omitempty"`
	Notes                  *string    `json:"Notes,omitempty"`
	IsNotePrivate          *bool      `json:"IsNotePrivate,omitempty"`
	TaskAbbrev             *string    `json:"TaskAbbrev,omitempty"`
	TaskId                 *string    `json:"TaskId,omitempty"`
	TaskEmrId              *string    `json:"TaskEmrId,omitempty"`
	TaskIsPrintStart       *bool      `json:"TaskIsPrintStart,omitempty"`
	TaskIsPrintEnd         *bool      `json:"TaskIsPrintEnd,omitempty"`
	TaskExtCallSysId       *string    `json:"TaskExtCallSysId,omitempty"`
	TaskKey                *string    `json:"TaskKey,omitempty"`
	TaskName               *string    `json:"TaskName,omitempty"`
	TaskBillSysId          *string    `json:"TaskBillSysId,omitempty"`
	TaskPayrollId          *string    `json:"TaskPayrollId,omitempty"`
	TaskShiftKey           *string    `json:"TaskShiftKey,omitempty"`
	TaskType               *string    `json:"TaskType,omitempty"`
	TaskContactInformation *string    `json:"TaskContactInformation,omitempty"`
	TaskTags               []any      `json:"TaskTags,omitempty"`
	LocationKey            *string    `json:"LocationKey,omitempty"`
	LocationName           *string    `json:"LocationName,omitempty"`
	LocationAbbrev         *string    `json:"LocationAbbrev,omitempty"`
	LocationAddress        *string    `json:"LocationAddress,omitempty"`
	LocationTags           []any      `json:"LocationTags,omitempty"`
	TimeZone               *string    `json:"TimeZone,omitempty"`
}

type Location struct {
	CompanyKey  *string     `json:"CompanyKey,omitempty"`
	LocationKey *int64      `json:"LocationKey,omitempty"`
	ID          *string     `json:"Id,omitempty"`
	Name        *string     `json:"Name,omitempty"`
	Address     *string     `json:"Address,omitempty"`
	Abbrev      *string     `json:"Abbrev,omitempty"`
	Notes       *string     `json:"Notes,omitempty"`
	TimeZone    *string     `json:"TimeZone,omitempty"`
	Tags        interface{} `json:"Tags,omitempty"`
}

type TagCompany struct {
	CompanyName   string        `json:"CompanyName"`
	CompanyKey    string        `json:"CompanyKey"`
	TagCategories []TagCategory `json:"Tags"`
}

type TagCategory struct {
	CategoryKey                    int    `json:"CategoryKey"`
	CategoryName                   string `json:"CategoryName"`
	CategoryDateCreated            string `json:"CategoryDateCreated"`
	CategoryDateLastModified       string `json:"CategoryDateLastModified"`
	IsAvailableForCreditAllocation bool   `json:"IsAvailableForCreditAllocation"`
	IsAvailableForLDailySum        bool   `json:"IsAvailableForLDailySum"`
	IsAvailableForHoliday          bool   `json:"IsAvailableForHoliday"`
	IsAvailableForLocation         bool   `json:"IsAvailableForLocation"`
	IsAvailableForProfile          bool   `json:"IsAvailableForProfile"`
	IsAvailableForSeries           bool   `json:"IsAvailableForSeries"`
	IsAvailableForStaff            bool   `json:"IsAvailableForStaff"`
	IsAvailableForStaffLocation    bool   `json:"IsAvailableForStaffLocation"`
	IsAvailableForStaffTarget      bool   `json:"IsAvailableForStaffTarget"`
	IsAvailableForRequestLimit     bool   `json:"IsAvailableForRequestLimit"`
	IsAvailableForTask             bool   `json:"IsAvailableForTask"`
	IsTTCMCategory                 bool   `json:"IsTTCMCategory"`
	IsSingleTaggingOnly            bool   `json:"IsSingleTaggingOnly"`
	IsPermissionCategory           bool   `json:"IsPermissionCategory"`
	IsUsedForStats                 bool   `json:"IsUsedForStats"`
	CategoryBackgroundColor        string `json:"CategoryBackgroundColor"`
	CategoryTextColor              string `json:"CategoryTextColor"`
	Tags                           []Tag  `json:"Tags"`
}

type Tag struct {
	Key              int    `json:"Key"`
	Name             string `json:"Name"`
	DateCreated      string `json:"DateCreated"`
	DateLastModified string `json:"DateLastModified"`
	BackgroundColor  string `json:"BackgroundColor"`
	TextColor        string `json:"TextColor"`
}
