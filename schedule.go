package main

import (
	"github.com/google/uuid"
)

// NewScheduleRequestResponse returns a pointer to a ScheduleRequestConfig with default values
func NewScheduleRequestResponse() *RequestResponse {
	rr := NewRequestResponse()
	rr.Request.Config = NewScheduleRequestConfig()
	return rr
}

// NewScheduleRequestConfig returns a point to a ScheduleRequestConfig with default values
func NewScheduleRequestConfig() *ScheduleRequestConfig {
	r := &ScheduleRequestConfig{
		Resource: "Schedule",
		Route:    "/Schedule",
		Includes: "Skillset,Tags,Profiles,TTCMTags",
		// Select:   "",
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
	Route                  string       `path:"-"`
	Includes               string       `query:"includes"`
	StartDate              DateMMDDYYYY `query:"startDate"`
	EndDate                DateMMDDYYYY `query:"endDate"`
	IncludeDeletes         bool         `query:"includeDeletes"`
	SinceModifiedTimestamp Time8601     `query:"sinceModifiedTimestamp"`
	Select                 string       `query:"$select"`
	Filter                 string       `query:"$filter"`
	OrderBy                string       `query:"$orderby"`
	Expand                 string       `query:"$expand"`
}
