package main

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"
)

// QgendaClientConfig is used to pass all necessary
// arguments to NewQgendaClient
type QgendaClientConfig struct {
	BaseURL        string
	ClientTimeout  time.Duration
	RequestTimeout time.Duration
	Email          string
	CompanyKey     string
	Password       string
}

//QgendaClient is the primary struct for handling client
// interactions with the qgenda api
type QgendaClient struct {
	BaseURL       *url.URL
	Client        *http.Client
	Credentials   *url.Values
	Values        *url.Values
	Authorization *AuthToken
	Config        QgendaClientConfig
}

// NewQgendaClient creates a QgendaClient from config values
func NewQgendaClient(qcc QgendaClientConfig) (*QgendaClient, error) {

	// part base url from a string
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
		return nil, errors.New("Error: QgendaClientConfig.Email cannot be empty")
	}
	for password = qcc.Password; password == ""; {
		return nil, errors.New("Error: QgendaClientConfig.Password cannot be empty")
	}
	for companyKey = qcc.CompanyKey; companyKey == ""; {
		return nil, errors.New("Error: QgendaClientConfig.CompanyKey cannot be empty")
	}

	authToken := &AuthToken{
		Token:   &http.Header{},
		Expires: time.Time{},
	}

	q := &QgendaClient{
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

// GetAll handles a *RequestResponse.Request and returns the data and metadata in
// *RequestResponse.Response
func (q *QgendaClient) GetAll(ctx context.Context, rr *RequestResponse) error {
	if len(rr.Requests) != len(rr.Responses) {
		rr.Responses = make([]Response, len(rr.Requests))
	}
	for i := range rr.Requests {
		if err := q.Get(ctx, rr.Requests[i], &rr.Responses[i]); err != nil {
			return err
		}
	}
	return nil
}

// Get handles a request and returns a *Response
func (q *QgendaClient) Get(ctx context.Context, r Request, response *Response) error {

	u := *q.BaseURL
	// handle authorization
	if err := q.Auth(ctx); err != nil {
		log.Printf("Error authorizing get request to %v: %v", r.Path, err)
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

	// metadata to capture data heritage
	meta := Metadata{
		APIVersion: "v2",
		Kind:       "qgenda",
		URL:        u.String(),
		Name:       r.Resource,
		Timestamp:  resTime,
	}
	response.Metadata = meta
	response.Data = b

	return nil
}
