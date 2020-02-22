package main

import "time"

// "github.com/google/uuid"

// NewScheduleRequestResponse returns a pointer to a ScheduleRequestConfig with default values
func NewScheduleRequestResponse() *RequestResponse {
	rr := NewRequestResponse()
	rr.Request.Config = NewScheduleRequestConfig()
	return rr
}

// NewScheduleRequestConfig returns a point to a ScheduleRequestConfig with default values
func NewScheduleRequestConfig() *ScheduleRequestConfig {
	r := &ScheduleRequestConfig{
		Resource:  "Schedule",
		Route:     "/schedule",
		Includes:  "StaffTags,TaskTags,LocationTags",
		StartDate: time.Now().Add(time.Hour * 168 * 2 * -1),
		EndDate:   time.Now(),
		Select:    "Date,TaskAbbrev,StaffAbbrev",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}
	return r
}

// ScheduleRequestConfig struct captures all available request arguments for
// qgenda Schedules endpoint
type ScheduleRequestConfig struct {
	Resource               string
	Route                  string    `path:"-"`
	Includes               string    `query:"includes"`
	StartDate              time.Time `query:"startDate" format:"01/02/2006"`
	EndDate                time.Time `query:"endDate"`
	IncludeDeletes         bool      `query:"includeDeletes"`
	SinceModifiedTimestamp time.Time `query:"sinceModifiedTimestamp"`
	Select                 string    `query:"$select"`
	Filter                 string    `query:"$filter"`
	OrderBy                string    `query:"$orderby"`
	Expand                 string    `query:"$expand"`
}

// https://api.qgenda.com/v2/schedule

// companyKey=00000000-0000-0000-0000-000000000000
// startDate=1/1/2014
// endDate=1/31/2014&
// $select=
// $filter=IsPublished
// $orderby=Date,TaskAbbrev,StaffAbbrev
// includes=Task"
