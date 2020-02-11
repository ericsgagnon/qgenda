package main

import (
	"context"
	"fmt"
	"net/http"

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

	err = q.Auth(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// print token
	fmt.Printf("Authorization: %#v\n%v\n",
		q.Authorization.Expires.Format(time.RFC3339),
		q.Authorization.Token.Get(http.CanonicalHeaderKey("Authorization")),
	)
	// for _, v := range q.Authorization.Cookies(q.BaseURL) {
	// 	fmt.Printf("%v: %v\n%v\n", v.Name, v.Expires, v.Value)
	// }

	fmt.Println("---------------------------------------------------------")

	fmt.Println("---------------------------------------------------------")

	// res, err := q.Client.Get("https://api.qgenda.com/v2/company")
	// if err != nil {
	// 	log.Fatal(err)

	// }

	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)

	// }

	// fmt.Println(string(body))

	// url := "https://api.qgenda.com/v2/company?includes=Profiles,Organizations"
	// url := "https://api.qgenda.com/v2/staffmember?companyKey=" + q.Credentials.Get("companyKey") + "&includes=Skillset,Tags,Profiles,TTCMTags"
	// fmt.Println(url)
	// t := map[string][]string.(q.Credentials)["companyKey"]
	// companyKey = "8c44c075-d894-4b00-9ae7-3b3842226626"
	// profileKey = "7f4d8aa0-292d-43b9-bec9-d253624c7de0"

	//url := "https://api.qgenda.com/v2/facility?companyKey=" + q.Credentials.Get("companyKey") + "&includes=TaskShift"
	// url := "https://api.qgenda.com/v2/location?companyKey=" + q.Credentials.Get("companyKey")
	// method := "GET"

	// payload := strings.NewReader("")

	// client := &http.Client{}
	// req, err := http.NewRequest(method, url, payload)

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// // req.Header.Add("Content-Type", "application/json")
	// req.Header.Add(
	// 	http.CanonicalHeaderKey("Authorization"),
	// 	q.Authorization.Token.Get(http.CanonicalHeaderKey("Authorization")),
	// )
	// // req.Header.Add(
	// // 	http.CanonicalHeaderKey("Accept-Encoding"),
	// // 	"*",
	// // )
	// //req.Header[http.CanonicalHeaderKey("Authorization")] = q.Auth.Token
	// res, err := client.Do(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer res.Body.Close()
	// body, err := ioutil.ReadAll(res.Body)

	// fmt.Println(string(body))
	// ioutil.WriteFile("samples/staffmembers.json", body, 0777)

	// date := "2100-01-01T00:00:00"
	// dateTime, err := time.ParseInLocation(time.RFC3339, date, time.Local)
	// if err != nil {
	// 	log.Fatal(err)

	// }
	// fmt.Println(dateTime)

}
