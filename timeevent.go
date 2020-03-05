package main

import (
	"fmt"
	"time"
)

// NewTimeEventRequestResponse returns a pointer to a TimeEventRequestConfig with default values
func NewTimeEventRequestResponse(rc *TimeEventRequestConfig) *RequestResponse {
	rr := NewRequestResponse()
	rr.RequestConfig = NewTimeEventRequestConfig(rc)
	return rr
}

// NewTimeEventRequestConfig returns a point to a TimeEventRequestConfig with default values
func NewTimeEventRequestConfig(rc *TimeEventRequestConfig) *TimeEventRequestConfig {
	if rc == nil {
		rc = &TimeEventRequestConfig{}
	}
	r := &TimeEventRequestConfig{
		Resource:          "TimeEvent",
		Route:             "/timeevent",
		Includes:          "ScheduleEntry,Task,StaffMember",
		StartDate:         time.Now().Add(time.Hour * 168 * 2 * -1),
		EndDate:           time.Now(),
		Interval:          time.Hour * 168,
		IntervalPrecision: time.Hour * 24,
		// Select:         ,
		// Filter:   "",
		// OrderBy: "Date,TaskAbbrev,StaffAbbrev",
		// Expand:   "",
	}
	fillDefaults(rc, r)
	return rc
}

// TimeEventRequestConfig struct captures all available request arguments for
// qgenda Schedules endpoint
type TimeEventRequestConfig struct {
	Resource          string
	Route             string        `path:"-"`
	Includes          string        `query:"includes"`
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
func (rc TimeEventRequestConfig) Parse() ([]Request, error) {
	var req []Request
	for i := rc.StartDate; i.Before(rc.EndDate); i = i.Add(rc.Interval) {
		rci := rc
		rci.StartDate = i
		rci.EndDate = rci.StartDate.Add(rc.Interval - rc.IntervalPrecision)
		rci.Resource = fmt.Sprintf("%v-%v-%v",
			rci.Resource,
			rci.StartDate.Format("20060102"),
			rci.EndDate.Format("20060102"),
		)
		reqi, err := parseRequestConfig(rci)
		if err != nil {
			return []Request{}, err
		}
		req = append(req, reqi)
	}
	return req, nil
}

// Stop overoptimizing!!
// func testIterationParse(rc RequestConfigurator) ([]Request, error) {
// 	tag := "iteration"
// 	var req []Request
// 	d := reflect.ValueOf(rc)
// 	dv := reflect.Indirect(d)

// 	for i := 0; i < dv.NumField(); i++ {
// 		structField := dv.Type().Field(i)
// 		field := reflect.Indirect(dv.Field(i))
// 		var iteratorConfig map[string]string
// 		if value, ok := structField.Tag.Lookup(tag); ok {
// 			switch value {
// 			case "start":
// 				iteratorConfig["start"] = field.Type().Name()
// 			case "end":
// 			case "interval":
// 			case "precision":
// 			default:

// 			}
// 		}
// 	}

// }
