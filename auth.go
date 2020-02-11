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
	Token   *http.Header
	Expires time.Time
}

const credentialsFile = "./.auth-token.json"

// Auth posts credentials in a request body and creates an authorization
// bearer header in the client with the returned access token
func (q *QgendaClient) Auth(ctx context.Context) {

	//
	valid := q.Authorization.Valid(ctx)
	fmt.Println(valid)
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
func (t *AuthToken) ReadFile(ctx context.Context) {
	// credentialsFile := "./.tmp-creds.json"
	f, err := os.Open(credentialsFile)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// read our opened jsonFile as a byte array.
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(b, t)
}

// Login submits credentials for authorization bearer token
func (q *QgendaClient) Login(ctx context.Context) {

	// request URL
	url := *q.BaseURL
	url.Path = path.Join(url.Path, "/login")

	// res, err := q.Client.PostForm(url.String(), *q.Credentials)
	res, err := ctxhttp.PostForm(ctx, q.Client, url.String(), *q.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	//response body
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// authorization token is returned in the response body
	var resData map[string]string
	json.Unmarshal(resBody, &resData)

	// use response timestamp + valid duration to set expire time
	respTime, err := time.Parse(time.RFC1123, res.Header[http.CanonicalHeaderKey("date")][0])
	if err != nil {
		log.Fatal(err)
	}
	validDuration, err := time.ParseDuration(resData["expires_in"] + "s")
	if err != nil {
		log.Fatal(err)
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

	fmt.Printf("Authorization: %#v\n%v\n",
		q.Authorization.Expires.Format(time.RFC3339),
		q.Authorization.Token.Get(http.CanonicalHeaderKey("Authorization")),
	)

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
