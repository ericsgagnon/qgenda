package qgenda

type ScheduleAuditLog struct {
	StaffFirstName            *string    `json:"StaffFirstName,omitempty"`
	StaffLastName             *string    `json:"StaffLastName,omitempty"`
	StaffAbbreviation         *string    `json:"StaffAbbreviation,omitempty"`
	StaffKey                  *string    `json:"StaffKey,omitempty"`
	TaskName                  *string    `json:"TaskName,omitempty"`
	TaskAbbreviation          *string    `json:"TaskAbbreviation,omitempty"`
	TaskKey                   *string    `json:"TaskKey,omitempty"`
	ScheduleEntryDate         *Date      `json:"ScheduleEntryDate,omitempty"`
	ScheduleEntryStartTimeUTC *Time      `json:"ScheduleEntryStartTimeUTC,omitempty"`
	ScheduleEntryStartTime    *TimeOfDay `json:"ScheduleEntryStartTime,omitempty"`
	ScheduleEntryEndTimeUTC   *Time      `json:"ScheduleEntryEndTimeUTC,omitempty"`
	ScheduleEntryEndTime      *TimeOfDay `json:"ScheduleEntryEndTime,omitempty"`
	ScheduleEntryKey          *string    `json:"ScheduleEntryKey,omitempty"`
	ActivityType              *string    `json:"ActivityType,omitempty"`
	SourceType                *string    `json:"SourceType,omitempty"`
	UserFirstName             *string    `json:"UserFirstName,omitempty"`
	UserLastName              *string    `json:"UserLastName,omitempty"`
	UserKey                   *string    `json:"UserKey,omitempty"`
	TimestampUTC              *Time      `json:"TimestampUTC,omitempty"`
	Timestamp                 *Time      `json:"Timestamp,omitempty"`
	AdditionalInformation     *string    `json:"AdditionalInformation,omitempty"`
	Locations                 []Location `json:"Locations,omitempty"`
	IPAddress                 *string    `json:"IPAddress,omitempty"`
}

func NewScheduleAuditLogRequest(rqf *RequestConfig) *Request {
	requestPath := "schedule/auditLog"
	queryFields := []string{
		"CompanyKey",
		"ScheduleStartDate",
		"ScheduleEndDate",
		"DateFormat",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
