package main

import (
	"context"

	// "io/ioutil"
	"log"
	// "net/http"
	"os"
	// "strings"
	"time"
)

//https://restapi.qgenda.com/?version=latest

var err error

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

	// Initialize a *RequestResponse for staffmembers
	srr := NewStaffMemberRequestResponse()
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
}
