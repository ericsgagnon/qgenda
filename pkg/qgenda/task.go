package qgenda

type Task struct {
	Abbrev                                 *string      `json:"Abbrev,omitempty"`
	BgColor                                *string      `json:"BgColor,omitempty"`
	BillSysId                              *string      `json:"BillSysId,omitempty"`
	CallRole                               *string      `json:"CallRole,omitempty"`
	CompKey                                *string      `json:"CompKey,omitempty"`
	ContactInformation                     *string      `json:"ContactInformation,omitempty"`
	DepartmentId                           *string      `json:"DepartmentId,omitempty"`
	DisplayAsAvailableOnCalendarConnection *bool        `json:"DisplayAsAvailableOnCalendarConnection,omitempty"`
	EmrId                                  *string      `json:"EmrId,omitempty"`
	EnableWaitlist                         *bool        `json:"EnableWaitlist,omitempty"`
	EndDate                                Date         `json:"EndDate,omitempty"`
	EnforceTimeCompatibility               *bool        `json:"EnforceTimeCompatibility,omitempty"`
	ExtCallPriority                        *string      `json:"ExtCallPriority,omitempty"`
	ExtCallSysId                           *string      `json:"ExtCallSysId,omitempty"`
	Extension                              *string      `json:"Extension,omitempty"`
	FullyCompatible                        *bool        `json:"FullyCompatible,omitempty"`
	HideOnCalendarSync                     *bool        `json:"HideOnCalendarSync,omitempty"`
	Label                                  *bool        `json:"Label,omitempty"`
	Manual                                 *bool        `json:"Manual,omitempty"`
	Notes                                  *string      `json:"Notes,omitempty"`
	NotificationList                       *string      `json:"NotificationList,omitempty"`
	OverrideSortOrder                      *string      `json:"OverrideSortOrder,omitempty"`
	PayrollId                              *string      `json:"PayrollId,omitempty"`
	PayrollUnits                           *string      `json:"PayrollUnits,omitempty"`
	ShowEndTime                            *bool        `json:"ShowEndTime,omitempty"`
	ShowOnAdminOnly                        *bool        `json:"ShowOnAdminOnly,omitempty"`
	TaskVisibility                         *string      `json:"TaskVisibility,omitempty"`
	ShowOpensOnCalendar                    *bool        `json:"ShowOpensOnCalendar,omitempty"`
	ShowStartTime                          *bool        `json:"ShowStartTime,omitempty"`
	ShowTimesWhenEdited                    *bool        `json:"ShowTimesWhenEdited,omitempty"`
	StartDate                              Date         `json:"StartDate,omitempty"`
	SuppressRequestEmailsAndNotifications  *bool        `json:"SuppressRequestEmailsAndNotifications,omitempty"`
	SyncAsAllDayEvent                      *bool        `json:"SyncAsAllDayEvent,omitempty"`
	TaskId                                 *string      `json:"TaskId,omitempty"`
	TaskKey                                *string      `json:"TaskKey,omitempty"`
	RequireTimePunch                       *bool        `json:"RequireTimePunch,omitempty"`
	TextColor                              *string      `json:"TextColor,omitempty"`
	TimeOffAfterShift                      *int64       `json:"TimeOffAfterShift,omitempty"`
	Type                                   *string      `json:"Type,omitempty"`
	Tags                                   TaskTags     `json:"Tags,omitempty"`
	TaskShifts                             TaskShifts   `json:"TaskShifts,omitempty"`
	Profiles                               TaskProfiles `json:"Profiles,omitempty"`
	ProfileKey                             *string      `json:"ProfileKey,omitempty"`
	IsViewable                             *bool        `json:"IsViewable,omitempty"`
	IsSchedulable                          *bool        `json:"IsSchedulable,omitempty"`
	IsRotationTask                         *bool        `json:"IsRotationTask,omitempty"`
	IsIgnoreHoursOffAfterTask              *bool        `json:"IsIgnoreHoursOffAfterTask,omitempty"`
}

type TaskTags []TaskTag

type TaskTag struct {
	ExtractDateTime *Time   `json:"-" db:"_extract_date_time" pgtype:"timestamp with time zone"`
	TaskIDHash      *string `json:"-" db:"_task_id_hash" pgtype:"text" parentprimarykey:"true" idhash:"true"`
	IDHash          *string `json:"-" db:"_id_hash" pgtype:"text" primarykey:"true"` // hash of identifying fields
	// --
	TaskKey      *string `json:"-" db:"taskkey" pgtype:"text" idhash:"true"`
	CategoryKey  *int64  `json:"-" db:"categorykey" pgtype:"bigint" idhash:"true"`
	CategoryName *string `json:"-" db:"categoryname" pgtype:"text" idhash:"true"`
	TagKey       *int64  `json:"Key" db:"tagkey" pgtype:"bigint" idhash:"true"`
	TagName      *string `json:"Name" db:"tagname" pgtype:"text" idhash:"true"`
}

type TaskShifts []TaskShift

type TaskShift struct {
	DayOfWeek     *string         `json:"DayOfWeek,omitempty"`
	StartTime     *TimeOfDay      `json:"StartTime,omitempty"`
	EndTime       *TimeOfDay      `json:"EndTime,omitempty"`
	MaxStaff      *int64          `json:"MaxStaff,omitempty"`
	MinStaff      *int64          `json:"MinStaff,omitempty"`
	OffAfter      *int64          `json:"OffAfter,omitempty"`
	StatCredit    *float64        `json:"StatCredit,omitempty"`
	StaffCount    *int64          `json:"StaffCount,omitempty"`
	EffectiveDate []EffectiveDate `json:"EffectiveDate:,omitempty"`
	IsActive      *bool           `json:"IsActive,omitempty"`
}

type EffectiveDate struct {
	EffectiveFromDate *Date      `json:"EffectiveFromDate,omitempty"`
	EffectiveToDate   *Date      `json:"EffectiveToDate,omitempty"`
	StartTime         *TimeOfDay `json:"StartTime,omitempty"`
	EndTime           *TimeOfDay `json:"EndTime,omitempty"`
	MaxStaff          *int64     `json:"MaxStaff,omitempty"`
	MinStaff          *int64     `json:"MinStaff,omitempty"`
	OffAfter          *int64     `json:"OffAfter,omitempty"`
	StatCredit        *float64   `json:"StatCredit,omitempty"`
}

type TaskProfiles []TaskProfile

type TaskProfile struct {
}
