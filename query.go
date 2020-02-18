package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"text/template"
	"time"
)

// Request holds the processed (escaped) values for each element
// of the api requests
type Request struct {
	Method string
	Path   string
	Query  url.Values
	Body   url.Values
}

// EncodePath uses html template to interpolate path values for an endpoint
func EncodePath(data interface{}) (string, error) {
	d := reflect.ValueOf(data)
	templateText := reflect.Indirect(d).FieldByName("Route").Interface().(string)
	// fmt.Println(templateText)
	// return templateText, nil
	t, err := template.New("path").Parse(templateText)
	if err != nil {
		log.Printf("Error Parsing Template: %v", err)
		return "", err
	}
	var bb bytes.Buffer
	err = t.Execute(&bb, d)
	if err != nil {
		log.Printf("Error Executing Template: %v", err)
		return "", err
	}
	p := bb.String()
	p = path.Join(p)
	p = template.HTMLEscapeString(p)
	return p, nil
}

// EncodeURLValues extracts struct values that match tag and returns them in a url.Values
func EncodeURLValues(data interface{}, tag string) (url.Values, error) {
	d := reflect.ValueOf(data)
	dv := reflect.Indirect(d)
	uv := url.Values{}
	for i := 0; i < dv.NumField(); i++ {
		query, ok := dv.Type().Field(i).Tag.Lookup(tag)
		if ok {
			val := fmt.Sprintf("%v", dv.Field(i).Interface())
			if val != "" {
				uv.Add(query, val)
			}
		}
	}
	// u := uv.Encode()
	return uv, nil
}

// ParseRequest takes a *QueryConfig and builds the path, query, and body of the request
func ParseRequest(qs interface{}) (*Request, error) {
	// Encode path
	p, err := EncodePath(qs)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}
	// fmt.Println(p)

	// Encode query
	q, err := EncodeURLValues(qs, "query")
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}
	// fmt.Println(q)
	// Encode body
	b, err := EncodeURLValues(qs, "body")
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}
	// fmt.Println(b)
	r := &Request{
		Method: "GET",
		Path:   p,
		Query:  q,
		Body:   b,
	}
	return r, err
}

// Metadata captures relevant metadata from each response
type Metadata struct {
	APIVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	Timestamp  time.Time `json:"time"`
}

// Get handles all aspects of the http get request and handling the response
func (q *QgendaClient) Get(ctx context.Context, rs *Request) ([]byte, *Metadata, error) {

	u := *q.BaseURL
	r := *rs
	if err := q.Auth(ctx); err != nil {
		log.Printf("Error authorizing get request to %v: %v", r.Path, err)
		return nil, nil, err
	}
	r.Query.Add("companyKey", q.Credentials.Get("companyKey"))
	// authTokenHeader := url.Values{}
	// authTokenHeader.Add("companyKey", q.Credentials.Get("companyKey"))
	u.RawQuery = r.Query.Encode()
	u.Path = path.Join(u.Path, r.Path)
	fmt.Println(u.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		log.Printf("Error in request to %v: %v", u, err)
		return nil, nil, err
	}
	req.Header = q.Authorization.Token.Clone()
	res, err := q.Client.Do(req)
	if err != nil {
		log.Printf("Error retrieving response from %v: %v", u, err)
		return nil, nil, err
	}
	// TODO: improve reading response for larger requests
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response from %v: %v", u, err)
		return nil, nil, err
	}
	defer res.Body.Close()
	resTime, err := http.ParseTime(res.Header.Get("Date"))
	if err != nil {
		log.Printf("Error parsing date header in response: %v", err)
		// accept time.Now as a 'rough' estimate of now
		resTime = time.Now()
	}

	// uu, err := url.QueryUnescape(u.String())
	// if err != nil {
	// 	log.Printf("Error unescaping url %v: %v", u.String(), err)
	// }

	// fmt.Printf("\n\n%v\n\n", uu)
	meta := &Metadata{
		APIVersion: "v2",
		Kind:       "qgenda",
		URL:        u.String(),
		Name:       "",
		Timestamp:  resTime,
	}

	return b, meta, nil
}
