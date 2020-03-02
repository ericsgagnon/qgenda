package main

import (
	"time"
)

// "github.com/google/uuid"

// ScheduleRequestResponse has a complete request and response for schedule resource
type ScheduleRequestResponse struct {
}

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

// https://api.qgenda.com/v2/schedule

// companyKey=00000000-0000-0000-0000-000000000000
// startDate=1/1/2014
// endDate=1/31/2014&
// $select=
// $filter=IsPublished
// $orderby=Date,TaskAbbrev,StaffAbbrev
// includes=Task"

// Initialize a *RequestResponse for schedule
// scheduleRC := NewScheduleRequestConfig(&ScheduleRequestConfig{
// 	EndDate:   time.Now().UTC(),
// 	StartDate: time.Now().UTC().AddDate(0, -2, 0),
// })

// test := ScheduleRequestConfig{}
// fmt.Println(test)
// /*--------------------------------------------------------------------------*/
// d := reflect.ValueOf(scheduleRC)

// templateText := reflect.Indirect(d).FieldByName("Route").Interface().(string)
// // fmt.Println(templateText)
// // return templateText, nil
// t, err := template.New("path").Parse(templateText)
// if err != nil {
// 	log.Printf("Error Parsing Template: %v", err)
// }
// var bb bytes.Buffer
// err = t.Execute(&bb, d)
// if err != nil {
// 	log.Printf("Error Executing Template: %v", err)
// }
// p := bb.String()
// p = path.Join(p)
// p = template.HTMLEscapeString(p)
// /*--------------------------------------------------------------------------*/

// // var srrs map[time.Time]*RequestResponse
// // var b []byte
// for i := scheduleRC.StartDate; i.Before(scheduleRC.EndDate); i = i.AddDate(0, 0, 7) {
// 	wg.Add(1)
// 	go func() {

// 		fmt.Println(i)
// 		scheduleRCi := *scheduleRC
// 		scheduleRCi.StartDate = i
// 		scheduleRCi.EndDate = scheduleRCi.StartDate.AddDate(0, 0, 6)
// 		srr := NewScheduleRequestResponse(&ScheduleRequestConfig{})
// 		srr.RequestConfig = scheduleRCi

// 		srr.Requests, err = ParseRequestConfig(scheduleRCi)
// 		if err != nil {
// 			log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
// 		}
// 		if err := q.Get(ctx, srr); err != nil {
// 			log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
// 		}
// 		// fmt.Println(srr)
// 		// dataJSON := string(*srr.Response.Data)
// 		// fmt.Printf("\n%v\n", dataJSON)

// 		filename := scheduleRCi.Resource + scheduleRCi.StartDate.Format("20060102") + ".json"
// 		if err := srr.Response.ToJSONFile(filename); err != nil {
// 			log.Fatalln(err)
// 		}
// 		wg.Done()
// 	}()
// 	// srrs[i] = srr
// }

// wg.Wait()
