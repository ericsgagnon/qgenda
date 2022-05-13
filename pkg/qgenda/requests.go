package qgenda

// calling this file requests since we need to request.go for the app itself

func NewRequestsRequest(rqf *RequestQueryFields) *Request {
	requestPath := "request"
	queryFields := []string{
		"CompanyKey",
		"StartDate",
		"EndDate",
		"DateFormat",
		"IncludeRemoved",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}

	// populate required fields

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewRequestsApprovedRequest(rqf *RequestQueryFields) *Request {
	requestPath := "request/approved"
	queryFields := []string{
		"CompanyKey",
		"StartDate",
		"EndDate",
		"MaxResults",
		"PageToken",
		"SyncToken",
		"DateFormat",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}

// Requests is plural due to name conflict with another type in the app
type Requests struct {
	RequestID       *string    `json:"RequestId,omitempty"`
	StaffNpi        *string    `json:"StaffNpi,omitempty"`
	TaskId          *string    `json:"TaskId,omitempty"`
	CallRole        *string    `json:"CallRole,omitempty"`
	CompKey         *string    `json:"CompKey,omitempty"`
	Credit          *float64   `json:"Credit,omitempty"`
	Date            *Date      `json:"Date,omitempty"`
	EndTime         *TimeOfDay `json:"EndTime,omitempty"`
	Notes           *string    `json:"Notes,omitempty"`
	RequestKey      *string    `json:"RequestKey,omitempty"`
	RequestType     *string    `json:"RequestType,omitempty"`
	RequestStatus   *string    `json:"RequestStatus,omitempty"`
	ResolutionDate  *Time      `json:"ResolutionDate,omitempty"`
	StaffAbbrev     *string    `json:"StaffAbbrev,omitempty"`
	StaffId         *string    `json:"StaffId,omitempty"`
	StaffFName      *string    `json:"StaffFName,omitempty"`
	StaffLName      *string    `json:"StaffLName,omitempty"`
	StaffKey        *string    `json:"StaffKey,omitempty"`
	StaffPayrollId  *string    `json:"StaffPayrollId,omitempty"`
	StartTime       *TimeOfDay `json:"StartTime,omitempty"`
	SubmittedByUser *string    `json:"SubmittedByUser,omitempty"`
	SubmittedDate   *Time      `json:"SubmittedDate,omitempty"`
	TaskAbbrev      *string    `json:"TaskAbbrev,omitempty"`
	TaskKey         *string    `json:"TaskKey,omitempty"`
	TaskName        *string    `json:"TaskName,omitempty"`
	StaffInternalId *string    `json:"StaffInternalId,omitempty"`
	Removed         *bool      `json:"Removed,omitempty"`
}

type RequestsApproved struct {
	NextPageToken *string    `json:"NextPageToken,omitempty"`
	NextSyncToken *string    `json:"NextSyncToken,omitempty"`
	Items         []Requests `json:"Items,omitempty"`
}

type RequestLimit struct {
	CompanyKey              *string  `json:"CompanyKey,omitempty"`
	Key                     *string  `json:"Key,omitempty"`
	StartDate               *Date    `json:"StartDate,omitempty"`
	EndDate                 *Date    `json:"EndDate,omitempty"`
	Type                    *string  `json:"Type,omitempty"`
	RecurringLength         *string  `json:"RecurringLength,omitempty"`
	CreditSource            *string  `json:"CreditSource,omitempty"`
	ErrorMessage            *string  `json:"ErrorMessage,omitempty"`
	IsActive                *bool    `json:"IsActive,omitempty"`
	DailyTotalMaxAllowedMon *float64 `json:"DailyTotalMaxAllowedMon,omitempty"`
	DailyTotalMaxAllowedTue *float64 `json:"DailyTotalMaxAllowedTue,omitempty"`
	DailyTotalMaxAllowedWed *float64 `json:"DailyTotalMaxAllowedWed,omitempty"`
	DailyTotalMaxAllowedThu *float64 `json:"DailyTotalMaxAllowedThu,omitempty"`
	DailyTotalMaxAllowedFri *float64 `json:"DailyTotalMaxAllowedFri,omitempty"`
	DailyTotalMaxAllowedSat *float64 `json:"DailyTotalMaxAllowedSat,omitempty"`
	DailyTotalMaxAllowedSun *float64 `json:"DailyTotalMaxAllowedSun,omitempty"`
	ShiftsCredit            []any    `json:"ShiftsCredit,omitempty"`
	TaskId                  *string  `json:"TaskId,omitempty"`
	TaskKey                 *string  `json:"TaskKey,omitempty"`
	TaskAbbreviation        *string  `json:"TaskAbbreviation,omitempty"`
	TaskshiftKey            *string  `json:"TaskshiftKey,omitempty"`
	IsIncluded              *bool    `json:"IsIncluded,omitempty"`
	DayOfTheWeek            *string  `json:"DayOfTheWeek,omitempty"`
	Credit                  *float64 `json:"Credit,omitempty"`
	StaffLimits             []any    `json:"StaffLimits,omitempty"`
	StaffId                 *string  `json:"StaffId,omitempty"`
	StaffInternalId         *string  `json:"StaffInternalId,omitempty"`
	StaffKey                *string  `json:"StaffKey,omitempty"`
	StaffAbbreviation       *string  `json:"StaffAbbreviation,omitempty"`
	StaffRequestLimits      *float64 `json:"StaffRequestLimits,omitempty"`
	StaffTotalLimit         *float64 `json:"StaffTotalLimit,omitempty"`
}
