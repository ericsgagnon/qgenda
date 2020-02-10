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

const credentialsFile = "./.tmp-creds.json"

// Auth posts credentials in a request body and creates an authorization
// bearer header in the client with the returned access token
func (q *QgendaClient) Auth(ctx context.Context) context.Context {

	// credentialsFile := "./.tmp-creds.json"
	jsonFile, err := os.Open(credentialsFile)

	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var auth AuthToken
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(b, &auth)
	// fmt.Printf("%+v", auth)

	// fmt.Println(q.Authorization.Token)
	// fmt.Println(q.Authorization.Expires)

	fmt.Printf("Authorization: %+v\n%+v\n",
		auth.Expires.Format(time.RFC3339),
		auth.Token.Get(http.CanonicalHeaderKey("Authorization")),
	)

	fmt.Printf("Credentials Expired: %t\n\n", auth.Expires.UTC().Before(time.Now().UTC()))

	if auth.Expires.UTC().Before(time.Now().UTC()) {
		fmt.Println("ttttttt")
	}
	fmt.Printf("Authorization: %#v\n%v\n",
		q.Authorization.Expires.Format(time.RFC3339),
		q.Authorization.Token.Get(http.CanonicalHeaderKey("Authorization")),
	)

	//TODO: check for Auth cookie or header, get another if missing or expired
	// request URL
	url := *q.BaseURL
	url.Path = path.Join(url.Path, "/login")

	// request
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

	j, err := json.MarshalIndent(q.Authorization, "", "  ")

	f, err := os.Create(credentialsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Write(j)

	return ctx
}

// CheckAuthToken checks the QgendaClient Auth Token for existence and expiration
func (q *QgendaClient) CheckAuthToken(ctx context.Context) context.Context {
	valid := q.Authorization.Expires.UTC().After(time.Now().UTC())
	if !valid {
		ctx = q.CheckAuthToken(ctx)
	}
	// Check for a valid authorization bearer token in QgendaClient
	// Set the Authorization header in the QgendaClient
	q.Authorization.Token.Get(
		http.CanonicalHeaderKey("Authorization"),
		fmt.Sprintf("bearer %v", resData["access_token"]),
	)

	return ctx
}

// CheckAuthFile checks the QgendaClient Auth Token for existence and expiration
func (q *QgendaClient) CheckAuthFile(ctx context.Context) context.Context {

	return ctx
}

// Login submits credentials for authorization bearer token
func (q *QgendaClient) Login(ctx context.Context) context.Context {

	return ctx
}
