package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

//https://restapi.qgenda.com/?version=latest

//QgendaClient is the primary struct for handling client
// interactions with the qgenda api
type QgendaClient struct {
	BaseURL    *url.URL
	Client     *http.Client
	Values     *url.Values
	Email      string
	CompanyKey string
	Password   string
}

var err error

func main() {

	// grab credentials from environment variables
	q := &QgendaClient{
		// BaseURL:    "https://api.qgenda.com/v2",
		BaseURL: &url.URL{
			Scheme: "https",
			Host:   "api.qgenda.com",
			Path:   "v2",
		},
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		Values:     &url.Values{},
		Email:      os.Getenv("QGENDA_EMAIL"),
		CompanyKey: os.Getenv("QGENDA_COMPANY_KEY"),
		Password:   os.Getenv("QGENDA_PASSWORD"),
	}
	// fmt.Println(q)
	q.Login()

}

// Login posts credentials in a request body and sets the
// authorization bearer header with the returned access token
func (q *QgendaClient) Login() {

	// request URL
	reqURL := *q.BaseURL
	reqURL.Path = path.Join(reqURL.Path, "/login")
	fmt.Println(reqURL.String())

	// request body
	v := &url.Values{}
	v.Add("email", q.Email)
	v.Add("password", q.Password)
	v.Add("companyKey", q.CompanyKey)
	reqBody := strings.NewReader(v.Encode())
	fmt.Println(reqBody)

	// request
	res, err := q.Client.PostForm(reqURL.String(), *v)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\n\n%+v\n\n", string(resBody))
}

// req, err := http.NewRequest(http.MethodPost, reqURL.String(), reqBody)
// if err != nil {
// 	log.Fatalln(err)
// }
// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// res, err := q.Client.Do(req)
// if err != nil {
// 	log.Fatalln(err)
// }

// fmt.Sprint(string(resBody))
// fmt.Println(q.Client)
