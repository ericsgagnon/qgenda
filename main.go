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

	// // get companies
	// cq := NewCompanyRequest()
	// cr, err := ParseRequest(cq)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(cr)
	// bb, err := q.Get2(ctx, cr)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(bb))
	// // unmarshall companies data
	// var c []Company
	// if err := json.Unmarshal(bb, &c); err != nil {
	// 	log.Printf("Error unmarshalling response from %v", err)
	// 	//return err
	// }
	// fmt.Printf("\n\n%v\n", c)

	// fmt.Printf("\n\n%v\n", c[0].CreatedTime.Valid)
	var c []Company
	if err := q.GetCompanies(ctx, nil, &c); err != nil {
		log.Fatal(err)
	}

}
