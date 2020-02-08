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

	//fmt.Println(q.Authorization.Get(http.CanonicalHeaderKey("Authorization")))

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
	fmt.Println(respTime.Format(time.RFC3339))
	fmt.Println(validDuration)
	fmt.Println(respTime.Add(validDuration).Format(time.RFC850))
	fmt.Println(expireTime.Format(time.RFC850))
	fmt.Println("----------------------------------------------------")

	// Set the Authorization header in the QgendaClient
	q.Authorization.Token.Set(
		http.CanonicalHeaderKey("Authorization"),
		fmt.Sprintf("bearer %v", resData["access_token"]),
	)

	q.Authorization.Expires = expireTime
	// Set the expire timestamp for the Authorization Header

	// q.Authorization.Cookies = []*http.Cookie{
	// q.Authorization = &http.Cookie{
	// 	Name:     http.CanonicalHeaderKey("Authorization"),
	// 	Value:    fmt.Sprintf("bearer %v", resData["access_token"]),
	// 	Domain:   q.BaseURL.Hostname(),
	// 	Path:     "/",
	// 	SameSite: http.SameSiteNoneMode,
	// 	// Expires: time.Now().Add(resData["expires_in"] * time.Second),
	// 	Expires: time.Time(expireTime),
	// }

	// for _, v := range ac {
	// 	fmt.Printf("%v: %v\n", v.Name, v.Value)
	// }
	// set Authorization cookie for all endpoints
	u := *q.BaseURL
	u.Path = "/"
	// q.Authorization.SetCookies(&u, ac)
	// fmt.Println(q.Authorization["Authorization"])
	// for _, v := range q.Authorization {
	// 	fmt.Printf("%v: %#v\n%v\n", v.Name, v.Expires.Format(time.RFC850), v.Value)
	// }

	fmt.Printf("Authorization: %#v\n%v\n",
		q.Authorization.Expires.Format(time.RFC3339),
		q.Authorization.Token, //[http.CanonicalHeaderKey("Authorization")],
	)
	// fmt.Printf("%v: %#v\n%v\n", q.Authorization.Name, q.Authorization.Expires.Format(time.RFC3339), q.Authorization.Value)

	// for _, v := range q.Authorization.Cookies(q.BaseURL) {
	// 	fmt.Printf("%v: %v\n%v\n", v.Name, v.Expires.Format(time.RFC3339), v.Value)
	// }

	return ctx
}
