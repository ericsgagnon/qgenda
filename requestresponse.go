package main

// RequestResponse captures all aspects of a request-response cycle
type RequestResponse struct {
	Request  *Request
	Response *Response
}

// NewRequestResponse initializes a RequestResponse and returns a pointer
// it is intended to be used in other constructors to avoid null pointer
func NewRequestResponse() *RequestResponse {
	rr := &RequestResponse{
		Request:  NewRequest(),
		Response: NewResponse(),
	}
	return rr
}
