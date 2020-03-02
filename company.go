package main

// CompanyRequestConfig is intended to be used as inputs to
// api requests to the company endpoints
type CompanyRequestConfig struct {
	Resource string `resource:"-"`
	Route    string `path:"-"`
	Includes string `query:"includes"`
	Select   string `query:"$select"`
	Filter   string `query:"$filter"`
	OrderBy  string `query:"$orderby"`
	Expand   string `query:"$expand"`
}

// NewCompanyRequestConfig returns a pointer to a CompanyRequestConfig with default values
func NewCompanyRequestConfig() *CompanyRequestConfig {
	cr := &CompanyRequestConfig{
		Resource: "Company",
		Route:    "/company",
		Includes: "Profiles,Organizations",
		// Select:   "",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}
	return cr
}

// // Parse builds a *Request from a *RequestConfigurator
// func (rc *CompanyRequestConfig) Parse() (*Request, error) {
// 	return ParseRequestConfig(rc)
// }

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
