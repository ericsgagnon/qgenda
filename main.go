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

	"gopkg.in/yaml.v2"
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

	// scheduleRC := NewScheduleRequestConfig(&ScheduleRequestConfig{
	// 	StartDate: time.Now().UTC().AddDate(-1, 0, 0),
	// })
	// scheduleR, err := scheduleRC.Parse()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	scheduleRR := NewScheduleRequestResponse(&ScheduleRequestConfig{
		StartDate: time.Now().UTC().AddDate(0, -1, 0),
	})
	scheduleRR.Requests, err = scheduleRR.RequestConfig.Parse()
	if err != nil {
		log.Fatalln(err)
	}
	PrintYAML(scheduleRR)
	// preconfigure response slice capacity so we can index values and use goroutines
	scheduleRR.Responses = make([]Response, len(scheduleRR.Requests))
	for i := range scheduleRR.Requests {
		if q.GetRequest(ctx, scheduleRR.Requests[i], &scheduleRR.Responses[i]) != nil {
			log.Fatalln(err)
		}
	}
	PrintYAML(scheduleRR.Responses)
	// request, err := ParseRequestConfig(scheduleRC)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// PrintYAML(request)
	// companyRC := NewCompanyRequestConfig()

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

	// scheduleRC = &ScheduleRequestConfig{
	// 	EndDate:   time.Now().UTC(),
	// 	StartDate: time.Now().UTC().AddDate(0, -2, 0),
	// }

	// fillDefaults(scheduleRC, NewScheduleRequestConfig(nil))
	// PrintYAML(scheduleRC)
}

// PrintYAML is a convenience function to print yaml-ized versions of
// variables to stdout (console). It is not meant for 'real' use.
func PrintYAML(in interface{}) {
	inYAML, err := yaml.Marshal(in)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(inYAML))
}

// defSF := def.Type().Field(i)
// defFN := defSF.Name
// defFT := defSF.Type.String()
// defField := reflect.Indirect(def.Field(i))
// defFV := defField.Interface()

// isZero := defField.IsZero()
// fmt.Printf("%#v\n", isZero)
// switch {
// case defFT == "time.Time" && !defFV.(time.Time).IsZero():
// 	fmt.Printf("%v:\t%v\n", defFN, defFV)

// default:
// 	fmt.Sprint(defFV)
// }
// fmt.Println(defFN)
// fmt.Printf("\n%#v\n", data)
// iterate through names/values and compare with argument struct

// replace argument fields with contstructor defaults if arguments are blank
// or ??

// blankSRC := ScheduleRequestConfig{}
// e := reflect.ValueOf(blankSRC)
// fmt.Printf("\n%#v\n", e)
// d := reflect.ValueOf(data)
// fmt.Printf("\n%#v\n", d)

// // ev := reflect.Indirect(e)
// // fmt.Printf("\n%#v\n", ev)
// dv := reflect.Indirect(d)
// // fmt.Printf("\n%#v\n", dv)

// // fmt.Printf("\nblank:\t%#v\targument:%#v\n", ev.NumField(), dv.NumField())
// // fmt.Printf("\nblank:\t%#v\targument:%#v\n", e.NumField(), d.NumField())
// uv := url.Values{}
// for i := 0; i < dv.NumField(); i++ {
// 	structField := dv.Type().Field(i)
// 	// fmt.Printf("\n%#v\n", structField)

// 	// fmt.Println(structField.Name)

// 	field := reflect.Indirect(dv.Field(i))
// 	// fmt.Printf("\n%#v\n", field)
// 	var val string
// 	if query, ok := structField.Tag.Lookup(""); ok {
// 		fieldType := field.Type().String()
// 		fieldFormat := structField.Tag.Get("format")
// 		fieldValue := field.Interface()
// 		switch {
// 		case fieldType == "time.Time" && !fieldValue.(time.Time).IsZero():
// 			if fieldFormat != "" {
// 				val = fieldValue.(time.Time).Format(fieldFormat)
// 			} else {
// 				val = fieldValue.(time.Time).Format(time.RFC3339)
// 			}
// 		default:
// 			val = fmt.Sprint(fieldValue)
// 		}
// 		if val != "" {
// 			uv.Add(query, val)
// 		}
// 	}
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
