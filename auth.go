package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

// AuthToken is the header that holds the authorization: bearer token and expire timestamp
type AuthToken struct {
	Token   *http.Header `json:"Token"`
	Expires time.Time    `json:"Expires"`
}

// use a preset file for caching the AuthToken to disk
const credentialsFile = "./.auth-token.json"

// Auth is the top level authorization method, it checks the validity of an existing token
// reads a cached token from disk into AuthToken, or logs in and retrieves a new token
// as needed
func (q *QgendaClient) Auth(ctx context.Context) error {

	log.Printf("----------------------------------------------------------------------------")
	log.Printf("Check AuthToken validity")
	if q.Authorization.Valid(ctx) {
		return nil
	}
	log.Printf("----------------------------------------------------------------------------")
	log.Printf("Check cached AuthToken validity")

	err := q.Authorization.ReadFile(ctx)
	// fmt.Printf(
	// 	"\n\nRead File: %v\n%v\n\n",
	// 	q.Authorization.Expires.String(),
	// 	q.Authorization.Token.Get(http.CanonicalHeaderKey("Authorization")),
	// )

	if err == nil {
		if q.Authorization.Valid(ctx) {
			return nil
		}
	}
	// else {
	// 	log.Printf("I'm a pretty princess: %v", err)
	// }

	log.Printf("----------------------------------------------------------------------------")
	log.Printf("Login")
	err = q.Login(ctx)
	return err
}

// Valid checks if the AuthToken is valid
func (t *AuthToken) Valid(ctx context.Context) bool {
	// checks if AuthToken exists and will expire more than a minute from now
	switch {
	case t == nil:
		log.Printf("AuthToken not initialized")
		return false
	case t.Token.Get(http.CanonicalHeaderKey("Authorization")) == "":
		log.Printf("AuthToken Authorization header is empty")
		fallthrough
	case t.Expires.IsZero():
		log.Printf("AuthToken expiration at zero %v", t.Expires.UTC().String())
		return false
	case t.Expires.UTC().Before(time.Now().Add(time.Minute).UTC()):
		log.Printf("AuthToken expired %v", t.Expires.UTC().String())
		return false
	case t.Expires.UTC().After(time.Now().Add(time.Minute).UTC()):
		log.Printf("AuthToken will expire %v", t.Expires.UTC().String())
		//TODO: Add jwt format validation
		return true
	}
	return false
}

// WriteFile writes the AuthToken to a file cache
func (t *AuthToken) WriteFile(ctx context.Context) error {

	j, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		log.Printf("Error marshalling AuthToken to json: %v", err)
	}

	f, err := os.Create(credentialsFile)
	if err != nil {
		log.Printf("Error opening file %v: %v", credentialsFile, err)
		return err
	}
	defer f.Close()

	if _, err := f.Write(j); err != nil {
		log.Printf("Unable to write AuthToken:\n%v\n", err)
		return err
	}
	return nil
}

// ReadFile reads the AuthToken from a file cache
func (t *AuthToken) ReadFile(ctx context.Context) error {

	f, err := os.Open(credentialsFile)
	if err != nil {
		log.Printf("Error opening AuthToken cache file %v: %v\n", credentialsFile, err)
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("Error reading AuthToken cache file %v: %v\n", credentialsFile, err)
		return err
	}

	tkn := &AuthToken{
		Token:   &http.Header{},
		Expires: time.Time{},
	}
	if err := json.Unmarshal(b, tkn); err != nil {
		log.Printf("Error unmarshalling cached AuthToken: %v", err)
	}
	// fmt.Printf("\nin ReadFile - tkn:\n%v", tkn)
	if tkn.Expires.UTC().Before(t.Expires.UTC()) {
		m := fmt.Sprintf("Cached AuthToken expired %v", tkn.Expires.UTC().String())
		log.Printf(m)
		return errors.New(m)
	}
	t.Token = tkn.Token
	t.Expires = tkn.Expires
	t = tkn
	// fmt.Printf("\nin ReadFile - AuthToken:\n%v", t)
	log.Printf("Cached AuthToken appears valid until %v", t.Expires.UTC().String())
	return nil

}

// Login submits credentials for authorization bearer token
func (q *QgendaClient) Login(ctx context.Context) error {

	// request URL
	route := "/login"
	endpoint := *q.BaseURL
	endpoint.Path = path.Join(endpoint.Path, route)

	// manually set a 'start' time for the AuthToken
	// to be conservative with expiration time, the start timestamp is
	// set just before the token is requested
	start := time.Now().UTC()

	// res, err := q.Client.PostForm(url.String(), *q.Credentials)
	res, err := ctxhttp.PostForm(ctx, q.Client, endpoint.String(), *q.Credentials)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Printf("Logged In")

	//response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading auth token from response: %v", err)
		return err
	}

	// authorization token is returned in the response body
	var tokenMap map[string]string
	if err := json.Unmarshal(body, &tokenMap); err != nil {
		log.Printf("Error unmarshalling json response from login endpoint: %v", err)
		return err
	}

	// fmt.Printf("\nin Login:\n\ntokenMap:\n%v\n\n", tokenMap)

	// Set the Authorization header in the QgendaClient
	q.Authorization.Token.Set(
		http.CanonicalHeaderKey("Authorization"),
		fmt.Sprintf("bearer %v", tokenMap["access_token"]),
	)

	fmt.Printf("\nin Login:\n\nAuthToken:\n%v\n\n", q.Authorization.Token)
	// use response timestamp + valid duration to set expire time
	// resTime, err := time.Parse(time.RFC1123, res.Header[http.CanonicalHeaderKey("date")][0])
	// if err != nil {
	// 	log.Printf("Error parsing response timestamp from response headers: %v", err)
	// 	return err
	// }
	validDuration, err := time.ParseDuration(tokenMap["expires_in"] + "s")
	if err != nil {
		log.Printf("Error parsing token valid duration: %v", err)
		return err
	}
	// q.Authorization.Expires = resTime.Add(validDuration)
	q.Authorization.Expires = start.Add(validDuration)
	log.Printf("Token Updated")

	if q.Authorization.Valid(ctx) {
		log.Printf("Token appears valid")
	}

	// write AuthToken to file
	if err := q.Authorization.WriteFile(ctx); err != nil {
		log.Printf("Unable to write AuthToken to %v:  %v", credentialsFile, err)
		return err
	}
	log.Printf("AuthToken written to %v", credentialsFile)
	return nil
}
