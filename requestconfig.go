package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"path"
	"reflect"
	"text/template"
	"time"
)

// RequestConfigurator defines expected behavior of the various [Resource]RequestConfigs
type RequestConfigurator interface {
	Parse() ([]Request, error)
}

// ParseRequestConfig takes an interface and builds a *Request
// while this accepts an interface it is designed to use the
// [Resource]RequestConfig's defined in this package.
// func ParseRequestConfig(rc interface{}) (Request, error) {
// 	// Encode path
// 	p, err := EncodePath(rc)
// 	if err != nil {
// 		log.Printf("%v\n", err)
// 		return Request{}, err
// 	}
// 	// fmt.Println(p)

// 	// Encode query
// 	q, err := EncodeURLValues(rc, "query")
// 	if err != nil {
// 		log.Printf("%v\n", err)
// 		return Request{}, err
// 	}
// 	// fmt.Println(q)
// 	// Encode body
// 	b, err := EncodeURLValues(rc, "body")
// 	if err != nil {
// 		log.Printf("%v\n", err)
// 		return Request{}, err
// 	}
// 	data := reflect.ValueOf(rc)
// 	resource := reflect.Indirect(data).FieldByName("Route").Interface().(string)
// 	// fmt.Println(b)
// 	r := Request{
// 		Resource: resource,
// 		Method:   "GET",
// 		Path:     p,
// 		Query:    q,
// 		Body:     b,
// 	}
// 	return r, err
// }

// parseRequestConfig takes an interface and builds a *Request
// while this accepts an interface it is designed to use the
// [Resource]RequestConfig's defined in this package.
func parseRequestConfig(rc interface{}) (Request, error) {
	// Encode path
	p, err := EncodePath(rc)
	if err != nil {
		log.Printf("%v\n", err)
		return Request{}, err
	}
	// fmt.Println(p)

	// Encode query
	q, err := EncodeURLValues(rc, "query")
	if err != nil {
		log.Printf("%v\n", err)
		return Request{}, err
	}
	// fmt.Println(q)
	// Encode body
	b, err := EncodeURLValues(rc, "body")
	if err != nil {
		log.Printf("%v\n", err)
		return Request{}, err
	}
	// fmt.Println(b)
	r := Request{
		Resource: "", //rc.Resource,
		Method:   "GET",
		Path:     p,
		Query:    q,
		Body:     b,
	}
	return r, err
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

// Resource returns the interfaces
// func Resource(data interface{}, tag string) (string, error) {

// 	// still struggling through pkg reflect
// 	d := reflect.ValueOf(data)
// 	dv := reflect.Indirect(d)
// 	var resource string
// 	for i := 0; i < dv.NumField(); i++ {
// 		structField := dv.Type().Field(i)
// 		field := reflect.Indirect(dv.Field(i))
// 		var val string
// 		if query, ok := structField.Tag.Lookup(tag); ok {
// 			fieldType := field.Type().String()
// 			fieldFormat := structField.Tag.Get("format")
// 			fieldValue := field.Interface()
// 			switch {
// 			case fieldType == "time.Time" && !fieldValue.(time.Time).IsZero():
// 				if fieldFormat != "" {
// 					val = fieldValue.(time.Time).Format(fieldFormat)
// 				} else {
// 					val = fieldValue.(time.Time).Format(time.RFC3339)
// 				}
// 			default:
// 				val = fmt.Sprint(fieldValue)
// 			}
// 			if val != "" {
// 				uv.Add(query, val)
// 			}
// 		}
// 	}
// 	return uv, nil
// }

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

// fillDefaults is a helper meant to fill in a structs zero-valued fields with default values
func fillDefaults(data interface{}, defaults interface{}) {

	dst := reflect.Indirect(reflect.ValueOf(data))
	def := reflect.ValueOf(defaults).Elem()

	for i := 0; i < def.NumField(); i++ {
		// use field name since there's no type checking on the arguments
		fieldName := def.Type().Field(i).Name
		defField := reflect.Indirect(def.FieldByName(fieldName))
		dstField := reflect.Indirect(dst.FieldByName(fieldName))
		// only copy non-zero defaults to destination values that are zero
		if !defField.IsZero() && dstField.IsZero() {
			dst.FieldByName(fieldName).Set(def.FieldByName(fieldName))
		}

	}
}

// // RequestConfig is a placeholder for specific resource RequestConfigs
// type RequestConfig struct {
// 	RequestConfigurator
// 	// interface{}
// }

// // RequestConfigurator is really just a placeholder
// // TODO: implement common methods for RequestConfigs
// type RequestConfigurator interface {
// 	// Get(s string) interface{}

// 	// String() string
// 	//Parse() (*Request, error)
// 	// ParseQuery() *url.Values
// 	// ParseBody() *url.Values
// }

// // ParseRequestConfigurator takes a *RequestConfigurator returns a *Request
// func ParseRequestConfigurator(rc RequestConfigurator) (*Request, error) {

// 	// Encode path
// 	p, err := EncodePath(rc)
// 	if err != nil {
// 		log.Printf("%v\n", err)
// 		return nil, err
// 	}
// 	// fmt.Println(p)

// 	// Encode query
// 	q, err := EncodeURLValues(rc, "query")
// 	if err != nil {
// 		log.Printf("%v\n", err)
// 		return nil, err
// 	}
// 	// fmt.Println(q)
// 	// Encode body
// 	b, err := EncodeURLValues(rc, "body")
// 	if err != nil {
// 		log.Printf("%v\n", err)
// 		return nil, err
// 	}
// 	// fmt.Println(b)
// 	r := &Request{
// 		Method: "GET",
// 		Path:   p,
// 		Query:  q,
// 		Body:   b,
// 	}

// 	return r, err

// }
