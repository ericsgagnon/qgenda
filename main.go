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
	crr := NewCompanyRequestResponse()
	// parse the *RequestResponse.Request.Config
	if err := crr.Request.ParseRequest(); err != nil {
		log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
	}
	if err := q.Get(ctx, crr); err != nil {
		log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
	}
	if err := crr.Response.ToJSONFile(""); err != nil {
		log.Fatalln(err)
	}

	// Initialize a *RequestResponse for company
	smrr := NewStaffMemberRequestResponse()
	// parse the *RequestResponse.Request.Config
	if err := smrr.Request.ParseRequest(); err != nil {
		log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
	}
	if err := q.Get(ctx, smrr); err != nil {
		log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
	}
	if err := smrr.Response.ToJSONFile(""); err != nil {
		log.Fatalln(err)
	}

	// Initialize a *RequestResponse for schedule
	srr := NewScheduleRequestResponse()
	// parse the *RequestResponse.Request.Config
	if err := srr.Request.ParseRequest(); err != nil {
		log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
	}
	if err := q.Get(ctx, srr); err != nil {
		log.Fatalf("Error parsing *RequestResponse.Request.Config: %v", err)
	}
	if err := srr.Response.ToJSONFile(""); err != nil {
		log.Fatalln(err)
	}

	startDate := time.Now().UTC().Add(time.Hour * 24 * 7 * 7 * -1)
	endDate := time.Now().UTC().AddDate(0, 0, 7*7)
	fmt.Println(startDate)
	fmt.Println(endDate)

	for i := startDate; i.Before(endDate); i = i.AddDate(0, 1, 7) {
		wg.Add(1)
		fmt.Println(i)
		go func(x time.Time) {
			fmt.Println(x)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func printTime(x time.Time) {
	fmt.Println(x)
	wg.Done()
}

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
