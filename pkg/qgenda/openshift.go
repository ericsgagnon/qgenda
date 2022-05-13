package qgenda

type OpenShift struct {
	CompanyKey             *string       `json:"CompanyKey,omitempty"`
	ScheduleKey            *string       `json:"ScheduleKey,omitempty"`
	OpenShiftCount         *int64        `json:"OpenShiftCount,omitempty"`
	CallRole               *string       `json:"CallRole,omitempty"`
	Credit                 *float64      `json:"Credit,omitempty"`
	Date                   *Date         `json:"Date,omitempty"`
	StartDate              *Date         `json:"StartDate,omitempty"`
	StartDateUTC           *Time         `json:"StartDateUTC,omitempty"`
	StartTime              *TimeOfDay    `json:"StartTime,omitempty"`
	EndDate                *Date         `json:"EndDate,omitempty"`
	EndDateUTC             *Time         `json:"EndDateUTC,omitempty"`
	EndTime                *TimeOfDay    `json:"EndTime,omitempty"`
	IsCred                 *bool         `json:"IsCred,omitempty"`
	IsSaved                *bool         `json:"IsSaved,omitempty"`
	IsPublished            *bool         `json:"IsPublished,omitempty"`
	IsLocked               *bool         `json:"IsLocked,omitempty"`
	IsStruck               *bool         `json:"IsStruck,omitempty"`
	Notes                  *string       `json:"Notes,omitempty"`
	IsNotePrivate          *bool         `json:"IsNotePrivate,omitempty"`
	TaskAbbrev             *string       `json:"TaskAbbrev,omitempty"`
	TaskId                 *string       `json:"TaskId,omitempty"`
	TaskEmrId              *string       `json:"TaskEmrId,omitempty"`
	TaskIsPrintStart       *bool         `json:"TaskIsPrintStart,omitempty"`
	TaskIsPrintEnd         *bool         `json:"TaskIsPrintEnd,omitempty"`
	TaskExtCallSysId       *string       `json:"TaskExtCallSysId,omitempty"`
	TaskKey                *string       `json:"TaskKey,omitempty"`
	TaskName               *string       `json:"TaskName,omitempty"`
	TaskBillSysId          *string       `json:"TaskBillSysId,omitempty"`
	TaskPayrollId          *string       `json:"TaskPayrollId,omitempty"`
	TaskShiftKey           *string       `json:"TaskShiftKey,omitempty"`
	TaskType               *string       `json:"TaskType,omitempty"`
	TaskContactInformation *string       `json:"TaskContactInformation,omitempty"`
	TaskTags               []TagCategory `json:"TaskTags,omitempty"`
	LocationKey            *string       `json:"LocationKey,omitempty"`
	LocationName           *string       `json:"LocationName,omitempty"`
	LocationAbbrev         *string       `json:"LocationAbbrev,omitempty"`
	LocationAddress        *string       `json:"LocationAddress,omitempty"`
	LocationTags           []Location    `json:"LocationTags,omitempty"`
	TimeZone               *string       `json:"TimeZone,omitempty"`
}

func NewOpenShiftsRequest(rqf *RequestQueryFields) *Request {
	requestPath := "schedule/openshifts"
	queryFields := []string{
		"StartDate",
		"EndDate",
		"CompanyKey",
		"DateFormat",
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	if rqf != nil {
		if rqf.Includes == nil {
			rqf.SetIncludes("TaskTags,LocationTags")
		}
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
