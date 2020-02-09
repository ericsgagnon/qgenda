package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

// Auth posts credentials in a request body and creates an authorization
// bearer header in the client with the returned access token
func (q *QgendaClient) Auth(ctx context.Context) context.Context {

	//TODO: check for Auth cookie or header, get another if missing or expired
	// request URL
	reqURL := *q.BaseURL
	reqURL.Path = path.Join(reqURL.Path, "/login")

	// request
	// res, err := q.Client.PostForm(reqURL.String(), *q.Credentials)
	res, err := ctxhttp.PostForm(ctx, q.Client, reqURL.String(), *q.Credentials)
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

	return ctx
}
