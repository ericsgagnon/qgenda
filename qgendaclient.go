package main

import (
	"errors"
	"net/http"
	"net/url"
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

// AuthToken is the header that holds the authorization: bearer token and expire timestamp
type AuthToken struct {
	Token   *http.Header
	Expires time.Time
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
	// Authorization *http.Header
	// Authorization *cookiejar.Jar
	// Authorization *http.Cookie
}

// NewQgendaClient creates a QgendaClient from config values
func NewQgendaClient(qcc QgendaClientConfig) (*QgendaClient, error) {

	// part base url from a string
	bu, err := url.Parse(qcc.BaseURL)
	if err != nil {
		return nil, err
	}

	// create a somewhat safe cookie jar
	// jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	// if err != nil {
	// 	return nil, err
	// }

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

	// acd := &http.Cookie{}
	// acd := &http.Header{}

	q := &QgendaClient{
		BaseURL: bu,
		Client: &http.Client{
			Timeout: cto,
			// Jar:     jar,
		},
		Credentials: &url.Values{
			"email":      {email},
			"companyKey": {companyKey},
			"password":   {password},
		},
		Values: &url.Values{},
		// Authorization: &[]http.Header{},
		// Authorization: jar,
		Authorization: &AuthToken{},
		Config:        qcc,
	}
	return q, nil
}
