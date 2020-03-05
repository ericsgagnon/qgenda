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

// Parse creates a RequestResponses Requests from its RequestConfig
// It uses the Parse method of the RequestConfigurator
func (rr *RequestResponse) Parse() error {
	var err error
	if rr.Requests, err = rr.RequestConfig.Parse(); err != nil {
		return err
	}
	return nil
}

// ResponsesToJSONFile writes all of a RequestResponses' Responses to a file
// if an empty string is passed as filename, it defaults to the 'Resource' from
// the request, also defaults to the user's cache directory
func (rr *RequestResponse) ResponsesToJSONFile(filename string) error {
	for _, v := range rr.Responses {
		if err := v.ToJSONFile(filename); err != nil {
			return err
		}
	}
	return nil
}
