package main

import (
	"time"
)

// "github.com/google/uuid"

// NewScheduleRequestResponse returns a pointer to a ScheduleRequestConfig with default values
func NewScheduleRequestResponse(src *ScheduleRequestConfig) *RequestResponse {
	rr := NewRequestResponse()
	rr.RequestConfig = NewScheduleRequestConfig(src)
	return rr
}

// NewScheduleRequestConfig returns a point to a ScheduleRequestConfig with default values
func NewScheduleRequestConfig(src *ScheduleRequestConfig) *ScheduleRequestConfig {
	if src == nil {
		src = &ScheduleRequestConfig{}
	}
	r := &ScheduleRequestConfig{
		Resource:          "Schedule",
		Route:             "/schedule",
		Includes:          "StaffTags,TaskTags,LocationTags",
		StartDate:         time.Now().Add(time.Hour * 168 * 2 * -1),
		EndDate:           time.Now(),
		Interval:          time.Hour * 168,
		IntervalPrecision: time.Hour * 24,
		IncludeDeletes:    true,
		// Select:         ,
		// Filter:   "",
		OrderBy: "Date,TaskAbbrev,StaffAbbrev",
		// Expand:   "",
	}
	fillDefaults(src, r)
	return src
}

// ScheduleRequestConfig struct captures all available request arguments for
// qgenda Schedules endpoint
type ScheduleRequestConfig struct {
	Resource               string
	Route                  string        `path:"-"`
	Includes               string        `query:"includes"`
	StartDate              time.Time     `query:"startDate" format:"01/02/2006" iteration:"start"`
	EndDate                time.Time     `query:"endDate" format:"01/02/2006" iteration:"end"`
	Interval               time.Duration `iteration:"interval"`
	IntervalPrecision      time.Duration `iteration:"precision"`
	IncludeDeletes         bool          `query:"includeDeletes"`
	SinceModifiedTimestamp time.Time     `query:"sinceModifiedTimestamp" format:"2006-01-02T15:04:05Z"`
	Select                 string        `query:"$select"`
	Filter                 string        `query:"$filter"`
	OrderBy                string        `query:"$orderby"`
	Expand                 string        `query:"$expand"`
}

// Parse parses the RequestConfig into one or more Requests
func (src ScheduleRequestConfig) Parse() ([]Request, error) {
	var req []Request
	for i := src.StartDate; i.Before(src.EndDate); i = i.Add(src.Interval) {
		srci := src
		srci.StartDate = i
		srci.EndDate = srci.StartDate.Add(src.Interval - src.IntervalPrecision)

		reqi, err := parseRequestConfig(srci)
		if err != nil {
			return []Request{}, err
		}
		req = append(req, reqi)
	}
	return req, nil
}
