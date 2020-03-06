package main

// StaffTargetConfig is intended to be used as inputs to
// api requests to the company endpoints
type StaffTargetConfig struct {
	Resource string `resource:"-"`
	Route    string `path:"-"`
	Includes string `query:"includes"`
	Select   string `query:"$select"`
	Filter   string `query:"$filter"`
	OrderBy  string `query:"$orderby"`
	Expand   string `query:"$expand"`
}

// NewStaffTargetRequestConfig returns a pointer to a StaffTargetConfig with default values
func NewStaffTargetRequestConfig(rc *StaffTargetConfig) *StaffTargetConfig {
	if rc == nil {
		rc = &StaffTargetConfig{}
	}

	r := &StaffTargetConfig{
		Resource: "StaffTarget",
		Route:    "/stafftarget",
		Includes: "Staff,Profiles,Locations",
		// Select:   "",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}

	fillDefaults(rc, r)
	return rc
}

// NewStaffTargetRequestResponse returns a pointer to a ScheduleRequestConfig with default values
func NewStaffTargetRequestResponse(rc *StaffTargetConfig) *RequestResponse {
	rr := NewRequestResponse()
	rr.RequestConfig = NewStaffTargetRequestConfig(rc)
	return rr
}

// Parse parses the RequestConfig into one or more Requests
func (rc StaffTargetConfig) Parse() ([]Request, error) {
	var req []Request
	reqi, err := parseRequestConfig(rc)
	if err != nil {
		return []Request{}, err
	}
	req = append(req, reqi)
	return req, nil
}
