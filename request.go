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
	"time"
)

// Request holds the processed (escaped) values for each element
// of the api requests
type Request struct {
	// Config RequestConfigurator
	Config interface{}
	Method string
	Path   string
	Query  url.Values
	Body   url.Values
}

// Request2 holds the processed (escaped) values for each element
// of the api requests
type Request2 struct {
	Config RequestConfigurator
	// Config interface{}
	Method string
	Path   string
	Query  url.Values
	Body   url.Values
}

// NewRequest initializes a Request and returns a pointer
func NewRequest() *Request {
	// var rc *RequestConfigurator
	r := &Request{
		// Config: &struct{}{},
		// Config: rc,
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

	// still struggling through pkg reflect
	d := reflect.ValueOf(data)
	dv := reflect.Indirect(d)
	uv := url.Values{}
	for i := 0; i < dv.NumField(); i++ {
		structField := dv.Type().Field(i)
		field := reflect.Indirect(dv.Field(i))
		var val string
		if query, ok := structField.Tag.Lookup(tag); ok {
			fieldType := field.Type().String()
			fieldFormat := structField.Tag.Get("format")
			fieldValue := field.Interface()
			switch {
			case fieldType == "time.Time" && !fieldValue.(time.Time).IsZero():
				if fieldFormat != "" {
					val = fieldValue.(time.Time).Format(fieldFormat)
				} else {
					val = fieldValue.(time.Time).Format(time.RFC3339)
				}
			default:
				val = fmt.Sprint(fieldValue)
			}
			if val != "" {
				uv.Add(query, val)
			}
		}
	}
	return uv, nil
}

// ParseRequest takes a *RequestConfigurator returns a *Request
func ParseRequest(rc RequestConfigurator) (*Request, error) {
	// Encode path
	p, err := EncodePath(rc)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}
	// fmt.Println(p)

	// Encode query
	q, err := EncodeURLValues(rc, "query")
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}
	// fmt.Println(q)
	// Encode body
	b, err := EncodeURLValues(rc, "body")
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}
	// fmt.Println(b)
	r := &Request{
		Config: rc,
		Method: "GET",
		Path:   p,
		Query:  q,
		Body:   b,
	}
	return r, err
}
