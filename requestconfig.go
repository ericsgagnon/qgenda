package main

import (
	"log"
)

// RequestConfigurator is really just a placeholder
// TODO: implement common methods for RequestConfigs
type RequestConfigurator interface {
	// Get(s string) interface{}

	// String() string
	// Parse() Request
	// ParseQuery() *url.Values
	// ParseBody() *url.Values
}

// ParseRequestConfig takes a *RequestConfigurator and builds the path, query, and body of the request
func ParseRequestConfig(rc RequestConfigurator) (*Request, error) {
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
		Method: "GET",
		Path:   p,
		Query:  q,
		Body:   b,
	}
	return r, err
}
