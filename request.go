package main

import (
	"net/url"
)

// Request holds the processed (escaped) values for each element
// of the api requests
type Request struct {
	// Config   interface{} //RequestConfigurator
	Resource string
	Method   string
	Path     string
	Query    url.Values
	Body     url.Values
}

// NewRequest initializes a Request and returns a pointer
func NewRequest() Request {
	// var rc *RequestConfigurator
	r := Request{
		// Config: &struct{}{},
		// Config: rc,
		Resource: "",
		Method:   "",
		Path:     "",
		Query:    url.Values{},
		Body:     url.Values{},
	}
	return r
}

// NewRequests generalizes NewResuts to a slice of pointers to requests
func NewRequests() []Request {
	return []Request{NewRequest()}
}
