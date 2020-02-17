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

// EncodeURLValues extracts struct values that match the provided tag and encodes them into a
// an escaped string
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

// Get handles all aspects of the http get request and handling the response
func (q *QgendaClient) Get(ctx context.Context, rs *Request) ([]byte, error) {

	url := *q.BaseURL
	r := *rs
	if err := q.Auth(ctx); err != nil {
		log.Printf("Error authorizing get request to %v: %v", r.Path, err)
		return nil, err
	}
	r.Query.Add("companyKey", q.Credentials.Get("companyKey"))
	// authTokenHeader := url.Values{}
	// authTokenHeader.Add("companyKey", q.Credentials.Get("companyKey"))
	url.RawQuery = r.Query.Encode()
	url.Path = path.Join(url.Path, r.Path)
	fmt.Println(url.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		log.Printf("Error in request to %v: %v", url, err)
		return nil, err
	}
	req.Header = q.Authorization.Token.Clone()
	res, err := q.Client.Do(req)
	if err != nil {
		log.Printf("Error retrieving response from %v: %v", url, err)
		return nil, err
	}
	// TODO: improve reading response for larger requests
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response from %v: %v", url, err)
		return nil, err
	}
	defer res.Body.Close()

	return b, nil
}

// func (q *QgendaClient) Get(ctx context.Context, , s *[]interface{}) error {

// // Get handles all aspects of the http get request and handling the response
// func (q *QgendaClient) Get(ctx context.Context, url string, qp *url.Values, s *[]interface{}) error {

// 	if err := q.Auth(ctx); err != nil {
// 		log.Printf("Error authorizing get request to %v: %v", url, err)
// 		return err
// 	}
// 	qp.Add("companyKey", q.Credentials.Get("companyKey"))

// 	endpoint := path.Join(url, qp.Encode())
// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, strings.NewReader(qp.Encode()))
// 	if err != nil {
// 		log.Printf("Error in request to %v: %v", url, err)
// 		return err
// 	}
// 	req.Header = q.Authorization.Token.Clone()
// 	res, err := q.Client.Do(req)
// 	if err != nil {
// 		log.Printf("Error retrieving response from %v: %v", endpoint, err)
// 		return err
// 	}
// 	// TODO: improve reading response for larger requests
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Printf("Error reading response from %v: %v", endpoint, err)
// 		return err
// 	}
// 	defer res.Body.Close()
// 	if err := json.Unmarshal(body, s); err != nil {
// 		log.Printf("Error unmarshalling response from %v", err)
// 		return err
// 	}

// 	return nil
// }
