package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"time"

	"golang.org/x/net/publicsuffix"
)

//https://restapi.qgenda.com/?version=latest

//QgendaClient is the primary struct for handling client
// interactions with the qgenda api
type QgendaClient struct {
	BaseURL       *url.URL
	Client        *http.Client
	Credentials   *url.Values
	Values        *url.Values
	Email         string
	CompanyKey    string
	Password      string
	Authorization *http.Header
}

var err error

func main() {

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	// grab credentials from environment variables
	q := &QgendaClient{
		BaseURL: &url.URL{
			Scheme: "https",
			Host:   "api.qgenda.com",
			Path:   "v2",
		},
		Client: &http.Client{
			Timeout: time.Second * 10,
			Jar:     jar,
		},
		Credentials: &url.Values{
			"email":      {os.Getenv("QGENDA_EMAIL")},
			"companyKey": {os.Getenv("QGENDA_COMPANY_KEY")},
			"password":   {os.Getenv("QGENDA_PASSWORD")},
		},
		Values:        &url.Values{},
		Email:         os.Getenv("QGENDA_EMAIL"),
		CompanyKey:    os.Getenv("QGENDA_COMPANY_KEY"),
		Password:      os.Getenv("QGENDA_PASSWORD"),
		Authorization: &http.Header{},
	}
	// fmt.Println(q)
	q.Login()
	// fmt.Println(q.Authorization.Get(http.CanonicalHeaderKey("Authorization")))
	// ck := q.Client.Jar.Cookies(&url.URL{
	// 	Scheme: "https",
	// 	Host:   "api.qgenda.com",
	// })

	// tt := *q.BaseURL
	// tt.Path = ""
	// ck := q.Client.Jar.Cookies(&tt)

	// for _, v := range q.Client.Jar.Cookies(q.BaseURL) {
	for _, v := range q.Client.Jar.Cookies(&url.URL{}) {
		fmt.Printf("%v: %v\n", v.Name, v.Value)
	}

}

// Login posts credentials in a request body and creates an authorization
// bearer header in the client with the returned access token
func (q *QgendaClient) Login() {

	//fmt.Println(q.Authorization.Get(http.CanonicalHeaderKey("Authorization")))
	//TODO: check for Auth cookie or header, get another if missing or expired
	// request URL
	reqURL := *q.BaseURL
	reqURL.Path = path.Join(reqURL.Path, "/login")

	// request
	res, err := q.Client.PostForm(reqURL.String(), *q.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var resData map[string]string
	json.Unmarshal(resBody, &resData)

	// bearerToken := resData["access_token"]
	validDuration, err := time.ParseDuration(resData["expires_in"] + "s")
	if err != nil {
		log.Fatal(err)
	}

	// Set the Authorization header in the QgendaClient
	q.Authorization.Set(
		http.CanonicalHeaderKey("Authorization"),
		fmt.Sprintf("bearer %v", resData["access_token"]),
	)

	ac := []*http.Cookie{
		&http.Cookie{
			Name:     http.CanonicalHeaderKey("Authorization"),
			Value:    fmt.Sprintf("bearer %v", resData["access_token"]),
			Domain:   q.BaseURL.Hostname(),
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			// Expires: time.Now().Add(resData["expires_in"] * time.Second),
			Expires: time.Now().Add(validDuration),
		},
	}

	// for _, v := range ac {
	// 	fmt.Printf("%v: %v\n", v.Name, v.Value)
	// }
	q.Client.Jar.SetCookies(q.BaseURL, ac)
}

// fmt.Sprintln(ac)
// q.Client.Jar.SetCookies(q.BaseURL)
// fmt.Println(q.Authorization)
// fmt.Println(resData["access_token"])
// for k, v := range resData {
// 	fmt.Println("-----------------------")
// 	fmt.Printf("%#v: %#v\n", k, v)
// }
// request body
// v := &url.Values{}
// v.Add("email", q.Email)
// v.Add("password", q.Password)
// v.Add("companyKey", q.CompanyKey)
// reqBody := strings.NewReader(v.Encode())
// fmt.Println(reqBody)
// res, err := q.Client.PostForm(reqURL.String(), *v)
