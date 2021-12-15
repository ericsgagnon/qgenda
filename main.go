package main

import (
	"context"
	"fmt"
	"sync"

	// "io/ioutil"
	"log"
	// "net/http"
	"os"
	// "strings"
	"time"
)

//https://restapi.qgenda.com/?version=latest

var err error
var wg sync.WaitGroup

// https://identity.getpostman.com/login?continue=https%3A%2F%2Fgo.postman.co%2Fnetwork%2Fimport%3Fcollection%3D1543481-097862de-81df-4145-a69e-2a328f02e487-S1TVYJ7R%26referrer%3Dhttps%253A%252F%252Frestapi.qgenda.com%252F%253Fversion%253Dlatest%2523e174c1d9-3793-4273-8c89-6194e363050a%26versionTag%3Dlatest%26environment%3D1543481-1a427bac-5902-4ebd-b709-06d4cef1d6ed-S1TVYJ7R%26traceId%3Dundefined
// "https://api.qgenda.com/v2/schedule/openshifts?companyKey=00000000-0000-0000-0000-000000000000&startDate=1/1/2012&endDate=1/31/2012&includes=LocationTags"
// "https://api.qgenda.com/v2/schedule/openshifts?companyKey=00000000-0000-0000-0000-000000000000&startDate=1/1/2014&endDate=1/31/2014&$select=Date,TaskAbbrev,OpenShiftCount&$filter=IsPublished&$orderby=Date,TaskAbbrev,OpenShiftCount&includes=Task"
func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)

	ctx := context.Background()
	// Set a duration.
	// duration := 150 * time.Millisecond

	// // Create a context that is both manually cancellable and will signal
	// // a cancel at the specified duration.
	// ctx, cancel := context.WithTimeout(context.Background(), duration)
	// defer cancel()

	// use environment variables to provide credentials
	q, err := NewQgendaClient(
		QgendaClientConfig{
			BaseURL:       "https://api.qgenda.com/v2",
			ClientTimeout: time.Second * 10,
			// grab credentials from environment variables
			Email:      os.Getenv("QGENDA_EMAIL"),
			CompanyKey: os.Getenv("QGENDA_COMPANY_KEY"),
			Password:   os.Getenv("QGENDA_PASSWORD"),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	// initial login
	err = q.Auth(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize a *RequestResponse for company
	// crr := NewCompanyRequestResponse()
	// // parse the *RequestResponse.Request.Config
	// if err := crr.Request.ParseRequest(); err != nil {
	// 	log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
	// }
	// if err := q.Get(ctx, crr); err != nil {
	// 	log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
	// }
	// if err := crr.Response.ToJSONFile(""); err != nil {
	// 	log.Fatalln(err)
	// }

	// // Initialize a *RequestResponse for company
	// smrr := NewStaffMemberRequestResponse()
	// // parse the *RequestResponse.Request.Config
	// if err := smrr.Request.ParseRequest(); err != nil {
	// 	log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
	// }
	// if err := q.Get(ctx, smrr); err != nil {
	// 	log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
	// }
	// if err := smrr.Response.ToJSONFile(""); err != nil {
	// 	log.Fatalln(err)
	// }

	// // Initialize a *RequestResponse for schedule
	// src := NewScheduleRequestConfig()
	// // grab a year of schedule
	// src.EndDate = time.Now().UTC()
	// src.StartDate = src.EndDate.AddDate(0, -2, 0)
	// srcOut, err := yaml.Marshal(src)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(string(srcOut))
	// // var srrs map[time.Time]*RequestResponse
	// // var b []byte
	// for i := src.StartDate; i.Before(src.EndDate); i = i.AddDate(0, 0, 7) {
	// 	wg.Add(1)
	// 	go func() {

	// 		fmt.Println(i)
	// 		srci := *src
	// 		srci.StartDate = i
	// 		srci.EndDate = srci.StartDate.AddDate(0, 0, 6)
	// 		srr := NewScheduleRequestResponse()
	// 		srr.Request.Config = srci

	// 		if err := srr.Request.ParseRequest(); err != nil {
	// 			log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
	// 		}
	// 		if err := q.Get(ctx, srr); err != nil {
	// 			log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
	// 		}
	// 		// fmt.Println(srr)
	// 		// dataJSON := string(*srr.Response.Data)
	// 		// fmt.Printf("\n%v\n", dataJSON)

	// 		filename := srci.Resource + srci.StartDate.Format("20060102") + ".json"
	// 		if err := srr.Response.ToJSONFile(filename); err != nil {
	// 			log.Fatalln(err)
	// 		}
	// 		wg.Done()
	// 	}()
	// 	// srrs[i] = srr
	// 	wg.Wait()
	// }

	rr := NewRequestResponse()
	rr.Request.Config = &struct {
		Resource       string
		Route          string    `path:"-"`
		Includes       string    `query:"includes"`
		Company        string    `query:"companyKey"`
		IncludeDeletes bool      `query:"includeDeletes"`
		Select         string    `query:"$select"`
		Filter         string    `query:"$filter"`
		OrderBy        string    `query:"$orderby"`
		Expand         string    `query:"$expand"`
		StartDate      time.Time `query:"startDate" format:"01/02/2006"`
		EndDate        time.Time `query:"endDate" format:"01/02/2006"`
		// SinceModifiedTimestamp time.Time `query:"sinceModifiedTimestamp" format:"2006-01-02T15:04:05Z"`
	}{
		Resource: "Test",
		Route:    "/tags",

		// Includes: "Skillset,Tags,Profiles,TTCMTags",
		// Includes: "TaskShifts,Skillset,Tags,Profiles,TTCMTags,Locations,Staff",
		// Company: q.Config.CompanyKey,
		// IncludeDeletes: true,
		// StartDate: time.Now().Add(time.Hour * 168 * 2 * -1),
		// EndDate:   time.Now(),
		// Select:         "Date,TaskAbbrev,StaffAbbrev",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}

	rr.Request.ParseRequest()
	fmt.Println(rr.Request.Body)
	if err := q.Get(ctx, rr); err != nil {
		log.Fatalf("Error getting Generic RequestResponse: %v", err)
	}
	fmt.Println(string(*rr.Response.Data))
	if err := rr.Response.ToJSONFile("test.json"); err != nil {
		log.Fatalln(err)
	}

}

// fmt.Println(i)
// go func(x time.Time) {
// 	fmt.Println(x)
// 	srrs[i].Request.Config = src
// 	wg.Done()
// }(i)
// parse the *RequestResponse.Request.Config

// fmt.Println(srr.Request.Config.(ScheduleRequestConfig).EndDate)

// startDate := time.Now().UTC().Add(time.Hour * 24 * 7 * 7 * -1)
// endDate := time.Now().UTC().AddDate(0, 0, 7*7)
// fmt.Println(startDate)
// fmt.Println(endDate)

// for i := startDate; i.Before(endDate); i = i.AddDate(0, 1, 7) {
// 	wg.Add(1)
// 	fmt.Println(i)
// 	go func(x time.Time) {
// 		fmt.Println(x)
// 		wg.Done()
// 	}(i)
// }
// func printTime(x time.Time) {
// 	fmt.Println(x)
// 	wg.Done()
// }

// fmt.Println(crr.Request.String())
// fmt.Println(crr.Request)

// crrJSON, err := json.Marshal(crr)
// if err != nil {
// 	log.Fatalln(err)
// }
// fmt.Println(string(crrJSON))

// crrYAML, err := yaml.Marshal(crr)
// if err != nil {
// 	log.Fatalln(err)
// }
// fmt.Println(string(crrYAML))

// fmt.Println(sprintRequestConfigurator(crr.Request.Config))

// // Initialize a *RequestResponse for staffmembers
// srr := NewStaffMemberRequestResponse()
// // parse the *RequestResponse.Request.Config
// if err := srr.Request.ParseRequest(); err != nil {
// 	log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
// }
// if err := q.Get(ctx, srr); err != nil {
// 	log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
// }
// if err := srr.Response.ToJSONFile(""); err != nil {
// 	log.Fatalln(err)
// }

// src := NewScheduleRequestConfig()
// fmt.Println(sprintRequestConfigurator2(src))
// 	u, err := EncodeURLValues(src, "query")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	for k, v := range u {
// 		fmt.Printf("|%-25v|%-60v|\n", k, v)
// 	}
// x := &time.Time{}
// fmt.Printf("\n%v\n", x.)
// tagValue := dv.Type().Field(i).Tag.Get(tag)
// fieldType := strings.Split(tagValue, ",")
// fmt.Println(fieldType)
// fmt.Println(dv.Type().Field(i).Tag.Get("format"))
// fieldType := dv.Type().Field(i).Tag.Get("format")
// fieldFormat := dv.Type().Field(i).Tag.Get("format")

// fmt.Println(dv.Type().Field(i).Tag)
