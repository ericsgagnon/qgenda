package main

import (
	"time"

	"github.com/google/uuid"
)

// CompanyQuery is
type CompanyQuery struct {
	Route    string `schema:"-"`
	Includes string `schema:"includes"`
	Select   string `schema:"$select"`
	Filter   string `schema:"$filter"`
	OrderBy  string `schema:"$orderby"`
	Expand   string `schema:"$expand"`
}

// /company

// Company contains basic company info
type Company struct {
	ID            uuid.UUID      `json:"CompanyKey"`
	Name          string         `json:"CompanyName"`
	Abbreviation  string         `json:"CompanyAbbr"`
	CreatedTime   time.Time      `json:"DateCreatedUtc"`
	Location      string         `json:"CompanyLocation"`
	PhoneNumber   string         `json:"CompanyPhoneNumber"`
	Profiles      []Profile      `json:"Profiles"`
	Organizations []Organization `json:"Organizations"`
}

// Profile appears to link a user role to a company...
type Profile struct {
	Name  string    `json:"ProfileName"`
	Key   uuid.UUID `json:"ProfileKey"`
	Admin bool      `json:"IsAdmin"`
}

// Organization appears to linke multiple companies and users
type Organization struct {
	Name string    `json:"OrgName"`
	ID   uuid.UUID `json:"OrgKey"`
}
