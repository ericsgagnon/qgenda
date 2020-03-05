package main

// LocationRequestConfig is intended to be used as inputs to
// api requests to the company endpoints
type LocationRequestConfig struct {
	Resource string `resource:"-"`
	Route    string `path:"-"`
	Includes string `query:"includes"`
	Select   string `query:"$select"`
	Filter   string `query:"$filter"`
	OrderBy  string `query:"$orderby"`
	Expand   string `query:"$expand"`
}

// NewLocationRequestConfig returns a pointer to a LocationRequestConfig with default values
func NewLocationRequestConfig(rc *LocationRequestConfig) *LocationRequestConfig {
	if rc == nil {
		rc = &LocationRequestConfig{}
	}

	r := &LocationRequestConfig{
		Resource: "Location",
		Route:    "/location",
		Includes: "Tags",
		// Select:   "",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}

	fillDefaults(rc, r)
	return rc
}

// NewLocationRequestResponse returns a pointer to a ScheduleRequestConfig with default values
func NewLocationRequestResponse(rc *LocationRequestConfig) *RequestResponse {
	rr := NewRequestResponse()
	rr.RequestConfig = NewLocationRequestConfig(rc)
	return rr
}

// Parse parses the RequestConfig into one or more Requests
func (rc LocationRequestConfig) Parse() ([]Request, error) {
	var req []Request
	reqi, err := parseRequestConfig(rc)
	if err != nil {
		return []Request{}, err
	}
	req = append(req, reqi)
	return req, nil
}

/*---------------------------------------------------------------------------------------*/

// // Company contains basic company info
// type Company struct {
// 	ID            uuid.UUID      `json:"CompanyKey"`
// 	Name          string         `json:"CompanyName"`
// 	Abbreviation  string         `json:"CompanyAbbr"`
// 	CreatedTime   TimeUTC        `json:"DateCreatedUtc,omitempty"`
// 	Location      string         `json:"CompanyLocation,omitempty"`
// 	PhoneNumber   string         `json:"CompanyPhoneNumber,omitempty"`
// 	Profiles      []Profile      `json:"Profiles,omitempty"`
// 	Organizations []Organization `json:"Organizations,omitempty"`
// }

// // Profile appears to link a user role to a company...
// type Profile struct {
// 	Name  string    `json:"ProfileName"`
// 	Key   uuid.UUID `json:"ProfileKey"`
// 	Admin bool      `json:"IsAdmin"`
// }

// // Organization appears to linke multiple companies and users
// type Organization struct {
// 	Name string `json:"OrgName"`
// 	Key  int    `json:"OrgKey"`
// }
