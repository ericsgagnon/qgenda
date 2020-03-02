package main

// RequestResponse captures all aspects of a request-response cycle
// a RequestResponse should have a single RequestConfig and one or more
// Requests and Responses. Once Requests are fulfilled with Responses,
// the number and order of Responses should mimic Requests.
type RequestResponse struct {
	RequestConfig RequestConfigurator
	Requests      []Request
	Responses     []Response
}

// NewRequestResponse initializes a RequestResponse and returns a pointer
// it is intended to be used in other constructors to avoid null pointer
func NewRequestResponse() *RequestResponse {
	rr := &RequestResponse{
		//RequestConfig: ,
		Requests:  NewRequests(),
		Responses: NewResponses(),
	}
	return rr
}
