package main

import (
	"time"
)

// RotationsRequestConfig struct captures all available request arguments for
// qgenda Rotations endpoint
type RotationsRequestConfig struct {
	Resource          string
	Route             string        `path:"-"`
	Includes          string        `query:"includes"`
	RangeStartDate    time.Time     `query:"rangeStartDate" format:"01/02/2006" iteration:"start"`
	RangeEndDate      time.Time     `query:"rangeEndDate" format:"01/02/2006" iteration:"end"`
	Interval          time.Duration `iteration:"interval"`
	IntervalPrecision time.Duration `iteration:"precision"`
	IgnoreHoliday     bool          `query:"ignoreHoliday"`
	IgnoreWeekend     bool          `query:"ignoreWeekend"`
	Select            string        `query:"$select"`
	Filter            string        `query:"$filter"`
	OrderBy           string        `query:"$orderby"`
	Expand            string        `query:"$expand"`
}

// NewRotationsRequestConfig returns a pointer to a RotationsRequestConfig with default values
func NewRotationsRequestConfig(rc *RotationsRequestConfig) *RotationsRequestConfig {
	if rc == nil {
		rc = &RotationsRequestConfig{}
	}

	r := &RotationsRequestConfig{
		Resource:          "Rotations",
		Route:             "schedule/rotations",
		RangeStartDate:    time.Now().Add(time.Hour * 168 * 12 * -1),
		RangeEndDate:      time.Now(),
		Interval:          time.Hour * 168 * 4,
		IntervalPrecision: time.Hour * 24,
		IgnoreHoliday:     false,
		IgnoreWeekend:     false,
	}

	fillDefaults(rc, r)
	return rc
}

// NewRotationsRequestResponse returns a pointer to a RotationsRequestConfig with default values
func NewRotationsRequestResponse(rc *RotationsRequestConfig) *RequestResponse {
	rr := NewRequestResponse()
	rr.RequestConfig = NewRotationsRequestConfig(rc)
	return rr
}

// // Parse parses the RequestConfig into one or more Requests
// func (rc RotationsRequestConfig) Parse() ([]Request, error) {
// 	var req []Request
// 	reqi, err := parseRequestConfig(rc)
// 	if err != nil {
// 		return []Request{}, err
// 	}
// 	req = append(req, reqi)
// 	return req, nil
// }

// Parse parses the RequestConfig into one or more Requests
func (rc RotationsRequestConfig) Parse() ([]Request, error) {
	var req []Request
	for i := rc.RangeStartDate; i.Before(rc.RangeEndDate); i = i.Add(rc.Interval) {
		rci := rc
		rci.RangeStartDate = i
		rci.RangeEndDate = rci.RangeStartDate.Add(rc.Interval - rc.IntervalPrecision)

		reqi, err := parseRequestConfig(rci)
		if err != nil {
			return []Request{}, err
		}
		req = append(req, reqi)
	}
	return req, nil
}
