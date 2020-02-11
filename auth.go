package main

import (
	"context"
	"encoding/json"
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

// Auth posts credentials in a request body and creates an authorization
// bearer header in the client with the returned access token
func (q *QgendaClient) Auth(ctx context.Context) error {
	fmt.Println("Step 1")
	valid := q.Authorization.Valid(ctx)

	// check token validity
	if !valid {
		fmt.Println("Step 2")
		// check cached file for valid token
		err := q.Authorization.ReadFile(ctx)
		fmt.Printf("%v\n", q.Authorization)
		if err == nil && q.Authorization.Valid(ctx) {
			fmt.Println("Step 3")
			return nil
		}
	} else {
		// login
		fmt.Println("Step 4")
		err = q.Login(ctx)
		if err != nil {
			fmt.Println("Step 5")
			return err
		}
	}
	fmt.Println("Step 6")
	return nil
}

// Valid checks if the AuthToken is valid
func (t *AuthToken) Valid(ctx context.Context) bool {
	// checks if AuthToken exists and will expire more than a minute from now
	if t == nil || t.Token.Get(http.CanonicalHeaderKey("Authorization")) == "" || t.Expires.UTC().After(time.Now().Add(time.Minute).UTC()) {
		//TODO: Add jwt format validation
		return false
	}
	return true

}

// WriteFile writes the AuthToken to a file cache
func (t *AuthToken) WriteFile(ctx context.Context) {

	j, err := json.MarshalIndent(t, "", "  ")

	f, err := os.Create(credentialsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Write(j)

}

// ReadFile reads the AuthToken from a file cache
func (t *AuthToken) ReadFile(ctx context.Context) error {
	// credentialsFile := "./.tmp-creds.json"
	// stat, err := os.Stat(credentialsFile)
	f, err := os.Open(credentialsFile)
	if os.IsNotExist(err) {
		// log.Printf("AuthToken cache file %v not found: %v\n", credentialsFile, err)
		fmt.Println("File no existing")
		return err
	}
	defer f.Close()

	// read our opened jsonFile as a byte array.
	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Trouble reading da file")
		return err
	}
	json.Unmarshal(b, t)
	fmt.Printf("Cached Token:\n%v\n", t)
	return nil
}

// Login submits credentials for authorization bearer token
func (q *QgendaClient) Login(ctx context.Context) error {

	// request URL
	url := *q.BaseURL
	url.Path = path.Join(url.Path, "/login")

	// res, err := q.Client.PostForm(url.String(), *q.Credentials)
	res, err := ctxhttp.PostForm(ctx, q.Client, url.String(), *q.Credentials)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//response body
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// authorization token is returned in the response body
	var resData map[string]string
	json.Unmarshal(resBody, &resData)

	// use response timestamp + valid duration to set expire time
	respTime, err := time.Parse(time.RFC1123, res.Header[http.CanonicalHeaderKey("date")][0])
	if err != nil {
		return err
	}
	validDuration, err := time.ParseDuration(resData["expires_in"] + "s")
	if err != nil {
		return err
	}
	expireTime := respTime.Add(validDuration)

	// Set the Authorization header in the QgendaClient
	q.Authorization.Token.Set(
		http.CanonicalHeaderKey("Authorization"),
		fmt.Sprintf("bearer %v", resData["access_token"]),
	)

	q.Authorization.Expires = expireTime
	// set Authorization cookie for all endpoints
	u := *q.BaseURL
	u.Path = "/"

	// troubleshooting
	fmt.Printf("\nLogged In - Token:\n%v", q.Authorization)

	return nil
}

// Check for a valid authorization bearer token in QgendaClient
// Set the Authorization header in the QgendaClient
// t.Get(http.CanonicalHeaderKey("Authorization"))
// fmt.Printf("%+v", auth)

// fmt.Println(q.Authorization.Token)
// fmt.Println(q.Authorization.Expires)

// fmt.Printf("Authorization: %+v\n%+v\n",
// 	auth.Expires.Format(time.RFC3339),
// 	auth.Token.Get(http.CanonicalHeaderKey("Authorization")),
// )

// fmt.Printf("Credentials Expired: %t\n\n", auth.Expires.UTC().Before(time.Now().UTC()))

// fmt.Printf("Authorization: %#v\n%v\n",
// 	q.Authorization.Expires.Format(time.RFC3339),
// 	q.Authorization.Token.Get(http.CanonicalHeaderKey("Authorization")),
// )
