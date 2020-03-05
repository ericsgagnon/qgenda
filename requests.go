package main

import (
	"fmt"
	"time"
)

// NewRequestRequestResponse returns a pointer to a RequestRequestConfig with default values
func NewRequestRequestResponse(rc *RequestRequestConfig) *RequestResponse {
	rr := NewRequestResponse()
	rr.RequestConfig = NewRequestRequestConfig(rc)
	return rr
}

// NewRequestRequestConfig returns a point to a RequestRequestConfig with default values
func NewRequestRequestConfig(rc *RequestRequestConfig) *RequestRequestConfig {
	if rc == nil {
		rc = &RequestRequestConfig{}
	}
	r := &RequestRequestConfig{
		Resource:          "Request",
		Route:             "/request",
		StartDate:         time.Now().Add(time.Hour * 168 * 2 * -1),
		EndDate:           time.Now(),
		Interval:          time.Hour * 168,
		IntervalPrecision: time.Hour * 24,
		// Select:         "Date,TaskAbbrev,StaffAbbrev",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}
	fillDefaults(rc, r)
	return rc
}

// RequestRequestConfig struct captures all available request arguments for
// qgenda Requests endpoint
type RequestRequestConfig struct {
	Resource          string
	Route             string        `path:"-"`
	StartDate         time.Time     `query:"startDate" format:"01/02/2006" iteration:"start"`
	EndDate           time.Time     `query:"endDate" format:"01/02/2006" iteration:"end"`
	Interval          time.Duration `iteration:"interval"`
	IntervalPrecision time.Duration `iteration:"precision"`
	Select            string        `query:"$select"`
	Filter            string        `query:"$filter"`
	OrderBy           string        `query:"$orderby"`
	Expand            string        `query:"$expand"`
}

// Parse parses the RequestConfig into one or more Requests
func (rc RequestRequestConfig) Parse() ([]Request, error) {
	var req []Request
	for i := rc.StartDate; i.Before(rc.EndDate); i = i.Add(rc.Interval) {
		srci := rc
		srci.StartDate = i
		srci.EndDate = srci.StartDate.Add(rc.Interval - rc.IntervalPrecision)
		srci.Resource = fmt.Sprintf("%v-%v-%v",
			srci.Resource,
			srci.StartDate.Format("20060102"),
			srci.EndDate.Format("20060102"),
		)
		reqi, err := parseRequestConfig(srci)
		if err != nil {
			return []Request{}, err
		}
		req = append(req, reqi)
	}
	return req, nil
}

// https://api.qgenda.com/v2/Request

// companyKey=00000000-0000-0000-0000-000000000000
// startDate=1/1/2014
// endDate=1/31/2014&
// $select=
// $filter=IsPublished
// $orderby=Date,TaskAbbrev,StaffAbbrev
// includes=Task"
