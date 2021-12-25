package qgenda

import (
	// "context"
	"context"
	"errors"
	"io/ioutil"
	"log"
	"path"
	"reflect"

	// "io/ioutil"
	// "log"
	"net/http"
	"net/url"

	// "path"
	// "reflect"
	"time"
)

// ClientConfig is used to pass all necessary
// arguments to NewClient
type ClientConfig struct {
	BaseURL        string
	ClientTimeout  time.Duration
	RequestTimeout time.Duration
	Email          string
	CompanyKey     string
	Password       string
}

// Client is the primary struct for handling client
// interactions with the qgenda api
type Client struct {
	BaseURL       *url.URL
	Client        *http.Client
	Credentials   *url.Values
	Values        *url.Values
	Authorization *AuthToken
	Config        ClientConfig
}

// NewQgendaClient creates a QgendaClient from config values
func NewClient(qcc ClientConfig) (*Client, error) {

	// parse base url from a string
	bu, err := url.Parse(qcc.BaseURL)
	if err != nil {
		return nil, err
	}

	// provide reasonable default client timeout
	var cto time.Duration
	if time.Duration(qcc.ClientTimeout) < time.Second*1 {
		cto = time.Second * 10
	} else {
		cto = qcc.ClientTimeout
	}

	// check for non nil values of credentials
	var email, password, companyKey string
	for email = qcc.Email; email == ""; {
		return nil, errors.New("Error: ClientConfig.Email cannot be empty")
	}
	for password = qcc.Password; password == ""; {
		return nil, errors.New("Error: ClientConfig.Password cannot be empty")
	}
	for companyKey = qcc.CompanyKey; companyKey == ""; {
		return nil, errors.New("Error: ClientConfig.CompanyKey cannot be empty")
	}

	authToken := &AuthToken{
		Token:   &http.Header{},
		Expires: time.Time{},
	}

	q := &Client{
		BaseURL: bu,
		Client: &http.Client{
			Timeout: cto,
		},
		Credentials: &url.Values{
			"email":      {email},
			"companyKey": {companyKey},
			"password":   {password},
		},
		Values:        &url.Values{},
		Authorization: authToken,
		Config:        qcc,
	}
	return q, nil
}

func parseRequest(c *Client, r *http.Request) (*http.Request, error) {

	return nil, nil
}

func ParseParameters(p ...Parameters) (*url.Values, error) {
	
}

func get(c *Client, ctx context.Context, r *http.Request) (*http.Response, error) {

	u := *c.BaseURL
	// handle authorization
	if err := c.Auth(ctx); err != nil {
		log.Printf("Error authorizing get request to %v: %v", r.URL, err)
		return err
	}
	r.Query.Add("companyKey", q.Credentials.Get("companyKey"))
	// build and send http request
	u.RawQuery = r.Query.Encode()
	u.Path = path.Join(u.Path, r.Path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		log.Printf("Error in request to %v: %v", u, err)
		return err
	}
	req.Header = q.Authorization.Token.Clone()
	res, err := q.Client.Do(req)
	if err != nil {
		log.Printf("Error retrieving response from %v: %v", u, err)
		return err
	}

	// handle response
	// TODO: improve reading response for larger requests
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response from %v: %v", u, err)
		return err
	}
	defer res.Body.Close()
	resTime, err := http.ParseTime(res.Header.Get("Date"))
	if err != nil {
		log.Printf("Error parsing date header in response: %v", err)
		// accept time.Now as a 'rough' estimate of now
		resTime = time.Now()
	}

	requestConfig := reflect.ValueOf(r.Config)
	resourceName := reflect.Indirect(requestConfig).FieldByName("Resource").Interface().(string)
	// metadata to capture data heritage
	meta := &Metadata{
		APIVersion: "v2",
		Kind:       "qgenda",
		URL:        u.String(),
		Name:       resourceName,
		Timestamp:  resTime,
	}
	rr.Response.Metadata = meta
	rr.Response.Data = &b

	return nil
}

// // Get handles a *RequestResponse.Request and returns the data and metadata in
// // *RequestResponse.Response
// func (q *Client) Get(ctx context.Context, rr *RequestResponse) error {

// 	u := *q.BaseURL
// 	r := *rr.Request
// 	// handle authorization
// 	if err := q.Auth(ctx); err != nil {
// 		log.Printf("Error authorizing get request to %v: %v", r.Path, err)
// 		return err
// 	}
// 	r.Query.Add("companyKey", q.Credentials.Get("companyKey"))
// 	// build and send http request
// 	u.RawQuery = r.Query.Encode()
// 	u.Path = path.Join(u.Path, r.Path)
// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
// 	if err != nil {
// 		log.Printf("Error in request to %v: %v", u, err)
// 		return err
// 	}
// 	req.Header = q.Authorization.Token.Clone()
// 	res, err := q.Client.Do(req)
// 	if err != nil {
// 		log.Printf("Error retrieving response from %v: %v", u, err)
// 		return err
// 	}

// 	// handle response
// 	// TODO: improve reading response for larger requests
// 	b, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Printf("Error reading response from %v: %v", u, err)
// 		return err
// 	}
// 	defer res.Body.Close()
// 	resTime, err := http.ParseTime(res.Header.Get("Date"))
// 	if err != nil {
// 		log.Printf("Error parsing date header in response: %v", err)
// 		// accept time.Now as a 'rough' estimate of now
// 		resTime = time.Now()
// 	}

// 	requestConfig := reflect.ValueOf(r.Config)
// 	resourceName := reflect.Indirect(requestConfig).FieldByName("Resource").Interface().(string)
// 	// metadata to capture data heritage
// 	meta := &Metadata{
// 		APIVersion: "v2",
// 		Kind:       "qgenda",
// 		URL:        u.String(),
// 		Name:       resourceName,
// 		Timestamp:  resTime,
// 	}
// 	rr.Response.Metadata = meta
// 	rr.Response.Data = &b

// 	return nil
// }
