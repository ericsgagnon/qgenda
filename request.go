package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
	"text/template"
	"time"
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

// RequestConfigurator implements common methods for ResourceRequestConfigs
type RequestConfigurator interface{}

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

// // EncodeURLValues extracts struct values that match tag and returns them in a url.Values
// func EncodeURLValues(data interface{}, tag string) (url.Values, error) {
// 	d := reflect.ValueOf(data)
// 	dv := reflect.Indirect(d)
// 	uv := url.Values{}
// 	for i := 0; i < dv.NumField(); i++ {
// 		query, ok := dv.Type().Field(i).Tag.Lookup(tag)
// 		tagValue := dv.Type().Field(i).Tag.Get(tag)
// 		fieldType := strings.Split(tagValue, ",")
// 		fmt.Println(fieldType)
// 		// fmt.Println(dv.Type().Field(i).Tag.Get("format"))
// 		// fieldType := dv.Type().Field(i).Tag.Get("format")
// 		// fieldFormat := dv.Type().Field(i).Tag.Get("format")

// 		fmt.Println(dv.Type().Field(i).Tag)
// 		if ok {
// 			val := fmt.Sprintf("%v", dv.Field(i).Interface())
// 			if val != "" {
// 				uv.Add(query, val)
// 			}
// 		}
// 	}
// 	// u := uv.Encode()

// 	return uv, nil
// }

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

// String returns a prettier version of Request
func (r Request) String() string {

	out :=
		fmt.Sprintf(
			"Request:\n  Config: %v\n  Method: %v\n  Path:%v\n",
			r.Config, r.Method, r.Path) +
			sprintURLValues(r.Query) +
			sprintURLValues(r.Body)
	return out
}

// sprintURLValues is a helper function to make printing structs that
// contain them easier
func sprintURLValues(u url.Values) string {

	var maxKeyLength int
	var maxValLength int
	for k, v := range u {
		if len(k) > maxKeyLength {
			maxKeyLength = len(k)
		}
		val := strings.Join(v, ", ")
		if len(val) > maxValLength {
			maxValLength = len(val)
		}
	}
	format := `` + "  %-" + fmt.Sprintf("%d", maxKeyLength) + "v: %v\n" + ``

	var out string
	for k, v := range u {
		val := strings.Join(v, ", ")
		if len(val) > 0 {
			out += fmt.Sprintf(format, k, val)
		}
	}
	return out
}

// sprintRequestConfigurator returns a prettier version of RequestConfigurator
func sprintRequestConfigurator(rc interface{}) string {
	// still struggling through pkg reflect
	var out string
	d := reflect.ValueOf(rc)
	dv := reflect.Indirect(d)
	for i := 0; i < dv.NumField(); i++ {
		structField := dv.Type().Field(i)
		field := reflect.Indirect(dv.Field(i))

		if field.String() != "" {
			out += fmt.Sprintf("  %v: %v  %v\n", structField.Name, field, structField.Tag)
		}
		// field := reflect.Indirect(dv.Field(i))
		// var val string
		// if query, ok := structField.Tag.Lookup(tag); ok {
		// }
	}
	return out
}

// fieldType := field.Type().String()
// fieldFormat := structField.Tag.Get("format")
// fieldValue := field.Interface()
