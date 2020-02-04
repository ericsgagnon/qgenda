package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//https://restapi.qgenda.com/?version=latest

//QgendaClient is the primary struct for handling client
// interactions with the qgenda api
type QgendaClient struct {
	BaseURL *url.URL
	// LoginURL   *url.URL
	Client     *http.Client
	Values     *url.Values
	Email      string
	CompanyKey string
	Password   string
}

func main() {

	var q QgendaClient
	q.SetBaseURL("https://api.qgenda.com/v2")

}

// SetBaseURL sets the base url for the qgenda rest api client
func (q *QgendaClient) SetBaseURL(s string) {
	baseURL, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	q.BaseURL = baseURL

}

// Login uses environment variables to exchange client credentials for
// a bearer token
func (q *QgendaClient) Login() {
	q.Password = os.Getenv("QGENDA_PASSWORD")
	q.Email = os.Getenv("QGENDA_EMAIL")
	q.CompanyKey = os.Getenv("QGENDA_COMPANY_KEY")
	requestString := fmt.Sprintf("email=%v&password=%v", q.Email, q.Password)
	requestBody := strings.NewReader(requestString)
	loginURL := fmt.Sprintf("%v/login", q.BaseURL.String())
	fmt.Println(loginURL)

	method := "POST"

	req, err := http.NewRequest(method, loginURL, requestBody)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := q.Client.Do(req)
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(string(resBody))

}
