package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

//https://restapi.qgenda.com/?version=latest

type qgendaAPI struct {
	BaseURL url.URL
	Client  *http.Client
}

func main() {

	var qa qgendaAPI
	u, err := url.Parse("https://api.qgenda.com")
	if err != nil {
		log.Fatal(err)
	}
	qa.BaseURL = *u
	qa.Client = &http.Client{}

	// qa.BaseURL, err :=
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	fmt.Println(qa)

	// client := &http.Client{}

	// baseURL :=
	// method := "POST"

	// qgendaEmail := os.Getenv("QGENDA_EMAIL")
	// fmt.Println(qgendaEmail)
	// qgendaPassword := os.Getenv("QGENDA_PASSWORD")

	// formData := url.Values{
	// 	"name": {"test"},
	// }

	// fmt.Println(formData)

	// payload := fmt.Sprintf("email=%s&password=%s", qgendaEmail, qgendaPassword)

	// payloadReader := strings.NewReader(payload)

	// req, err := http.NewRequest(method, baseURL, payloadReader)

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// res, err := client.Do(req)
	// defer res.Body.Close()
	// body, err := ioutil.ReadAll(res.Body)

	// fmt.Println(string(body))

	// u := &url.URL{
	// 	Scheme: "https",
	// 	Host:   "api.qgenda.com",
	// 	User: url.UserPassword(
	// 		os.Getenv("QGENDA_EMAIL"),
	// 		os.Getenv("QGENDA_PASSWORD"),
	// 	),
	// }
	// fmt.Println(u)
	// qgendaCompanyKey := os.Getenv("QGENDA_COMPANY_KEY")

	// payload := strings.NewReader("email=test@test.com&password=test123")
	/*-----------------------------------------------------------------------*/

	// url = "https://api.qgenda.com/v2/schedule?companyKey=00000000-0000-0000-0000-000000000000&startDate=1/1/2014&endDate=1/31/2014&$select=Date,TaskAbbrev,StaffAbbrev&$filter=IsPublished&$orderby=Date,TaskAbbrev,StaffAbbrev&includes=Task"
	// method = "GET"

	// client = &http.Client{}
	// req, err = http.NewRequest(method, url, nil)

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// // req.Header.Add("Authorization", "bearer eyJhbGciOiJBMjU2S1ciLCJlbmMiO...")
	// req.Header.Add("Authorization", "bearer eyJhbGciOiJBMjU2S1ciLCJlbmMiO...")

	// res, err = client.Do(req)
	// defer res.Body.Close()
	// body, err = ioutil.ReadAll(res.Body)

	// fmt.Println(string(body))
}
