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

	/* Company -------------------------------------------------------------------------*/
	if err := q.DownloadAndSaveCompany(ctx, nil, ""); err != nil {
		log.Fatalln(err)
	}
	/* StaffMembers --------------------------------------------------------------------*/
	// if err := q.DownloadAndSaveStaffMember(ctx, nil, ""); err != nil {
	// 	log.Fatalln(err)
	// }
	if err := q.DownloadAndSaveStaffMember(ctx, nil, fmt.Sprintf("staffmember-%s.json", time.Now().Format("20060102150405Z0700"))); err != nil {
		log.Fatalln(err)
	}

	/* Schedule ------------------------------------------------------------------------*/
	scheduleRC := &ScheduleRequestConfig{
		StartDate: time.Now().UTC().AddDate(-1, -3, 0),
		EndDate:   time.Now().UTC().AddDate(0, 6, 0),
	}
	if err := q.DownloadAndSaveSchedule(ctx, scheduleRC, ""); err != nil {
		log.Fatalln(err)
	}
	/* OpenShifts ----------------------------------------------------------------------*/
	openShiftsRC := &OpenShiftsRequestConfig{
		StartDate: time.Now().UTC().AddDate(-1, -3, 0),
		EndDate:   time.Now().UTC().AddDate(0, 6, 0),
	}
	if err := q.DownloadAndSaveOpenShifts(ctx, openShiftsRC, ""); err != nil {
		log.Fatalln(err)
	}
	/* Rotations -----------------------------------------------------------------------*/
	rotationsRC := &RotationsRequestConfig{
		RangeStartDate:    time.Now().UTC().AddDate(-1, -3, 0),
		RangeEndDate:      time.Now().UTC().AddDate(0, 6, 0),
		Interval:          time.Hour * 24 * 180,
		IntervalPrecision: time.Hour * 24,
	}
	if err := q.DownloadAndSaveRotations(ctx, rotationsRC, ""); err != nil {
		log.Fatalln(err)
	}
	/* Requests ------------------------------------------------------------------------*/
	requestRC := &RequestRequestConfig{
		StartDate:         time.Now().UTC().AddDate(-1, -3, 0),
		EndDate:           time.Now().UTC().AddDate(0, 6, 0),
		Interval:          time.Hour * 24 * 180,
		IntervalPrecision: time.Hour * 24,
	}
	if err := q.DownloadAndSaveRequests(ctx, requestRC, ""); err != nil {
		log.Fatalln(err)
	}
	/* Tasks ----------------------------------------------------------------------------*/
	// if err := q.DownloadAndSaveTasks(ctx, nil, ""); err != nil {
	// 	log.Fatalln(err)
	// }
	if err := q.DownloadAndSaveTasks(ctx, nil, fmt.Sprintf("tasks-%s.json", time.Now().Format("20060102150405Z0700"))); err != nil {
		log.Fatalln(err)
	}

	// Reshape the data

}

// DownloadAndSaveCompany is only intended to be used interactively in main
func (q *QgendaClient) DownloadAndSaveCompany(ctx context.Context, rc *CompanyRequestConfig, filename string) error {
	rr := NewCompanyRequestResponse(rc)
	if err := rr.Parse(); err != nil {
		log.Printf("Error parsing RequestConfig to Request: %v", err)
		return err
	}
	if err := q.GetAll(ctx, rr); err != nil {
		log.Printf("Error getting requests: %v", err)
		return err
	}
	if err := rr.ResponsesToJSONFile(filename); err != nil {
		log.Printf("Error writing to JSONFile: %v", err)
		return err
	}
	return nil
}

// DownloadAndSaveStaffMember is only intended to be used interactively in main
func (q *QgendaClient) DownloadAndSaveStaffMember(ctx context.Context, rc *StaffMemberRequestConfig, filename string) error {
	rr := NewStaffMemberRequestResponse(rc)
	if err := rr.Parse(); err != nil {
		log.Printf("Error parsing RequestConfig to Request: %v", err)
		return err
	}
	if err := q.GetAll(ctx, rr); err != nil {
		log.Printf("Error getting requests: %v", err)
		return err
	}
	if err := rr.ResponsesToJSONFile(filename); err != nil {
		log.Printf("Error writing to JSONFile: %v", err)
		return err
	}
	return nil
}

// DownloadAndSaveSchedule is only intended to be used interactively in main
func (q *QgendaClient) DownloadAndSaveSchedule(ctx context.Context, rc *ScheduleRequestConfig, filename string) error {
	rr := NewScheduleRequestResponse(rc)
	if err := rr.Parse(); err != nil {
		log.Printf("Error parsing RequestConfig to Request: %v", err)
		return err
	}
	if err := q.GetAll(ctx, rr); err != nil {
		log.Printf("Error getting requests: %v", err)
		return err
	}
	if err := rr.ResponsesToJSONFile(filename); err != nil {
		log.Printf("Error writing to JSONFile: %v", err)
		return err
	}
	return nil
}

// DownloadAndSaveOpenShifts is only intended to be used interactively in main
func (q *QgendaClient) DownloadAndSaveOpenShifts(ctx context.Context, rc *OpenShiftsRequestConfig, filename string) error {
	rr := NewOpenShiftsRequestResponse(rc)
	if err := rr.Parse(); err != nil {
		log.Printf("Error parsing RequestConfig to Request: %v", err)
		return err
	}
	if err := q.GetAll(ctx, rr); err != nil {
		log.Printf("Error getting requests: %v", err)
		return err
	}
	if err := rr.ResponsesToJSONFile(filename); err != nil {
		log.Printf("Error writing to JSONFile: %v", err)
		return err
	}
	return nil
}

// DownloadAndSaveRotations is only intended to be used interactively in main
func (q *QgendaClient) DownloadAndSaveRotations(ctx context.Context, rc *RotationsRequestConfig, filename string) error {
	rr := NewRotationsRequestResponse(rc)
	if err := rr.Parse(); err != nil {
		log.Printf("Error parsing RequestConfig to Request: %v", err)
		return err
	}
	if err := q.GetAll(ctx, rr); err != nil {
		log.Printf("Error getting requests: %v", err)
		return err
	}
	if err := rr.ResponsesToJSONFile(filename); err != nil {
		log.Printf("Error writing to JSONFile: %v", err)
		return err
	}
	return nil
}

// DownloadAndSaveRequests is only intended to be used interactively in main
func (q *QgendaClient) DownloadAndSaveRequests(ctx context.Context, rc *RequestRequestConfig, filename string) error {
	rr := NewRequestRequestResponse(rc)
	if err := rr.Parse(); err != nil {
		log.Printf("Error parsing RequestConfig to Request: %v", err)
		return err
	}
	if err := q.GetAll(ctx, rr); err != nil {
		log.Printf("Error getting requests: %v", err)
		return err
	}
	if err := rr.ResponsesToJSONFile(filename); err != nil {
		log.Printf("Error writing to JSONFile: %v", err)
		return err
	}
	return nil
}

// DownloadAndSaveTasks is only intended to be used interactively in main
func (q *QgendaClient) DownloadAndSaveTasks(ctx context.Context, rc *TaskRequestConfig, filename string) error {
	rr := NewTaskRequestResponse(rc)
	if err := rr.Parse(); err != nil {
		log.Printf("Error parsing RequestConfig to Request: %v", err)
		return err
	}
	if err := q.GetAll(ctx, rr); err != nil {
		log.Printf("Error getting requests: %v", err)
		return err
	}
	if err := rr.ResponsesToJSONFile(filename); err != nil {
		log.Printf("Error writing to JSONFile: %v", err)
		return err
	}
	return nil
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

// Don't seem to have access to these resources...
// /* DailyCase -----------------------------------------------------------------------*/
// dailyCaseRR := NewDailyCaseRequestResponse(nil)
// if err := dailyCaseRR.Parse(); err != nil {
// 	log.Fatalln(err)
// }
// if err := q.GetAll(ctx, dailyCaseRR); err != nil {
// 	log.Fatalln(err)
// }
// if err := dailyCaseRR.ResponsesToJSONFile(""); err != nil {
// 	log.Fatalln(err)
// }

// /* Facility ------------------------------------------------------------------------*/
// facilityRR := NewFacilityRequestResponse(nil)
// if err := facilityRR.Parse(); err != nil {
// 	log.Fatalln(err)
// }
// if err := q.GetAll(ctx, facilityRR); err != nil {
// 	log.Fatalln(err)
// }
// if err := facilityRR.ResponsesToJSONFile(""); err != nil {
// 	log.Fatalln(err)
// }
// /* StaffTarget ---------------------------------------------------------------------*/
// staffTargetRR := NewStaffTargetRequestResponse(nil)
// if err := staffTargetRR.Parse(); err != nil {
// 	log.Fatalln(err)
// }
// if err := q.GetAll(ctx, staffTargetRR); err != nil {
// 	log.Fatalln(err)
// }
// if err := staffTargetRR.ResponsesToJSONFile(""); err != nil {
// 	log.Fatalln(err)
// }
// /* TimeEvent ------------------------------------------------------------------------*/
// timeEventRR := NewTimeEventRequestResponse(&TimeEventRequestConfig{
// 	StartDate: time.Now().UTC().AddDate(0, -1, 0),
// })
// if err := timeEventRR.Parse(); err != nil {
// 	log.Fatalln(err)
// }
// if err := q.GetAll(ctx, timeEventRR); err != nil {
// 	log.Fatalln(err)
// }
// for _, v := range timeEventRR.Responses {
// 	fmt.Println(string(v.Data))
// }
// if err := timeEventRR.ResponsesToJSONFile(""); err != nil {
// 	log.Fatalln(err)
// }
// /* Location --------------------------------------------------------------------------*/
// locationRR := NewLocationRequestResponse(nil)
// if err := locationRR.Parse(); err != nil {
// 	log.Fatalln(err)
// }
// if err := q.GetAll(ctx, locationRR); err != nil {
// 	log.Fatalln(err)
// }
// for _, v := range locationRR.Responses {
// 	fmt.Println(string(v.Data))
// }
// if err := locationRR.ResponsesToJSONFile(""); err != nil {
// 	log.Fatalln(err)
// }
// /* Tags ------------------------------------------------------------------------------*/
// tagsRR := NewTagsRequestResponse(nil)
// if err := tagsRR.Parse(); err != nil {
// 	log.Fatalln(err)
// }
// if err := q.GetAll(ctx, tagsRR); err != nil {
// 	log.Fatalln(err)
// }
// if err := tagsRR.ResponsesToJSONFile(""); err != nil {
// 	log.Fatalln(err)
// }
