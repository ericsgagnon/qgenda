package main

// NewFacilityRequestResponse returns a pointer to a ScheduleRequestConfig with default values
func NewFacilityRequestResponse(rc *FacilityRequestConfig) *RequestResponse {
	rr := NewRequestResponse()
	rr.RequestConfig = NewFacilityRequestConfig(rc)
	return rr
}

// NewFacilityRequestConfig returns a pointer to a FacilityRequestConfig with default values
func NewFacilityRequestConfig(rc *FacilityRequestConfig) *FacilityRequestConfig {
	if rc == nil {
		rc = &FacilityRequestConfig{}
	}

	r := &FacilityRequestConfig{
		Resource: "Facility",
		Route:    "/facility",
		Includes: "TaskShift",
		// Select:   "",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}

	fillDefaults(rc, r)
	return rc
}

// FacilityRequestConfig struct captures all available request arguments for
// qgenda StaffMembers endpoint
type FacilityRequestConfig struct {
	Resource string
	Route    string `path:"-"`
	Includes string `query:"includes"`
	Select   string `query:"$select"`
	Filter   string `query:"$filter"`
	OrderBy  string `query:"$orderby"`
	Expand   string `query:"$expand"`
}

// Parse parses the RequestConfig into one or more Requests
func (rc FacilityRequestConfig) Parse() ([]Request, error) {
	var req []Request
	reqi, err := parseRequestConfig(rc)
	if err != nil {
		return []Request{}, err
	}
	req = append(req, reqi)
	return req, nil
}
