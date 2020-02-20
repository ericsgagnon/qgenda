package main

import (
	"bytes"
	"fmt"
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
	Config interface{}
	Method string
	Path   string
	Query  url.Values
	Body   url.Values
}

// NewRequest initializes a Request and returns a pointer
func NewRequest() *Request {
	r := &Request{
		Config: &struct{}{},
		Method: "",
		Path:   "",
		Query:  url.Values{},
		Body:   url.Values{},
	}
	return r
}

// ParseRequest takes a *QueryConfig and builds the path, query, and body of the request
func (r *Request) ParseRequest() error {
	var err error
	r.Method = http.MethodGet
	// Encode path
	r.Path, err = EncodePath(r.Config)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	// Encode query
	r.Query, err = EncodeURLValues(r.Config, "query")
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	// Encode body
	r.Body, err = EncodeURLValues(r.Config, "body")
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	return nil
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
