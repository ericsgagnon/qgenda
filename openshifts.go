package main

import (
	"net/url"
	"time"
)

// "github.com/google/uuid"

// Request3 is light years ahead of Request2 and the entire length of the dying universe ahead of Request
type Request3 struct {
	Config interface{}
	Parse  string
	Method string
	Path   string
	Query  url.Values
	Body   url.Values
}

// NewOpenShiftsRequestResponse returns a pointer to a OpenShiftsRequestConfig with default values
func NewOpenShiftsRequestResponse() *RequestResponse {
	rr := NewRequestResponse()
	rr.RequestConfig = NewOpenShiftsRequestConfig()
	return rr
}

// NewOpenShiftsRequestConfig returns a point to a OpenShiftsRequestConfig with default values
func NewOpenShiftsRequestConfig() *OpenShiftsRequestConfig {
	r := &OpenShiftsRequestConfig{
		Resource:          "OpenShifts",
		Route:             "/OpenShifts",
		Includes:          "StaffTags,TaskTags,LocationTags",
		StartDate:         time.Now().Add(time.Hour * 168 * 2 * -1),
		EndDate:           time.Now(),
		Interval:          time.Hour * 168,
		IntervalPrecision: time.Hour * 24,
		IncludeDeletes:    true,
		// Select:         "Date,TaskAbbrev,StaffAbbrev",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}
	return r
}

// OpenShiftsRequestConfig struct captures all available request arguments for
// qgenda OpenShiftss endpoint
type OpenShiftsRequestConfig struct {
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
func (osrc OpenShiftsRequestConfig) Parse() ([]Request, error) {
	var req []Request
	for i := osrc.StartDate; i.Before(osrc.EndDate); i = i.Add(osrc.Interval) {
		srci := osrc
		srci.StartDate = i
		srci.EndDate = srci.StartDate.Add(osrc.Interval - osrc.IntervalPrecision)

		reqi, err := parseRequestConfig(srci)
		if err != nil {
			return []Request{}, err
		}
		req = append(req, reqi)
	}
	return req, nil
}

// https://api.qgenda.com/v2/OpenShifts

// companyKey=00000000-0000-0000-0000-000000000000
// startDate=1/1/2014
// endDate=1/31/2014&
// $select=
// $filter=IsPublished
// $orderby=Date,TaskAbbrev,StaffAbbrev
// includes=Task"
