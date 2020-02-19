package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/google/uuid"
)

// ItemList is a universal container to data and metadata
type ItemList struct {
	MetaData *Metadata   `json:"Metadata"`
	Items    interface{} `json:"Items"`
}

// CompanyRequest is intended to be used as inputs to
// api requests to the company endpoints
type CompanyRequest struct {
	Route    string `path:"-"`
	Includes string `query:"includes"`
	Select   string `query:"$select"`
	Filter   string `query:"$filter"`
	OrderBy  string `query:"$orderby"`
	Expand   string `query:"$expand"`
}

// NewCompanyRequest returns a point to a CompanyRequest with default values
func NewCompanyRequest() *CompanyRequest {
	cr := &CompanyRequest{
		Route:    "/company",
		Includes: "Profiles,Organizations",
		// Select:   "",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}
	return cr
}

// Company contains basic company info
type Company struct {
	ID            uuid.UUID      `json:"CompanyKey"`
	Name          string         `json:"CompanyName"`
	Abbreviation  string         `json:"CompanyAbbr"`
	CreatedTime   TimeUTC        `json:"DateCreatedUtc,omitempty"`
	Location      string         `json:"CompanyLocation,omitempty"`
	PhoneNumber   string         `json:"CompanyPhoneNumber,omitempty"`
	Profiles      []Profile      `json:"Profiles,omitempty"`
	Organizations []Organization `json:"Organizations,omitempty"`
}

// Profile appears to link a user role to a company...
type Profile struct {
	Name  string    `json:"ProfileName"`
	Key   uuid.UUID `json:"ProfileKey"`
	Admin bool      `json:"IsAdmin"`
}

// Organization appears to linke multiple companies and users
type Organization struct {
	Name string `json:"OrgName"`
	Key  int    `json:"OrgKey"`
}

// GetCompanies uses the company endpoint to get all companies for a user
// func (q *QgendaClient) GetCompanies(ctx context.Context, cr *CompanyRequest, c *[]Company) error {
func (q *QgendaClient) GetCompanies(ctx context.Context, cr *CompanyRequest, il *ItemList) error {

	if cr == nil {
		cr = NewCompanyRequest()
	}
	r, err := ParseRequest(cr)
	fmt.Println(r)
	if err != nil {
		log.Fatal(err)
	}
	bb, meta, err := q.Get(ctx, r)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(meta)
	meta.Name = "CompanyList"

	// fmt.Println(string(bb))
	var c *[]Company
	if err := json.Unmarshal(bb, &c); err != nil {
		log.Printf("Error unmarshalling response from %v", err)
		return err
	}

	il.MetaData = meta
	il.Items = c

	return nil
}

// ToJSONFile writes an itemlist to a file in json format, using metadata to form
// the filename
func (il *ItemList) ToJSONFile(p string, f string) error {
	// create directory, or use "data" as default
	if p == "" {
		p = "data"
	}
	if err := os.MkdirAll(p, 0777); err != nil {
		log.Printf("Error making directory %v: %#v", p, err)
		return err
	}

	// build filename if not provided
	if f == "" {
		// f = il.MetaData.Name + "-" + il.MetaData.Timestamp.UTC().Format("20060102T150405Z07:00") + ".json"
		f = il.MetaData.Name + ".json"
	}
	f = path.Join(p, f)

	mm, err := json.MarshalIndent(il, "", "  ")
	if err != nil {
		log.Printf("Error marshalling to json: %+v", err)
		return err
	}

	if err := ioutil.WriteFile(f, mm, 0755); err != nil {
		log.Printf("Error writing file %v to disk: %v", f, err)
		return err
	}

	return nil
}

// FromJSONFile reads an itemlist from a jsonfile
func FromJSONFile(f string, il *ItemList) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.Printf("Error Reading file %v: %v", f, err)
		return err
	}

	if err := json.Unmarshal(b, il); err != nil {
		log.Printf("Error Unmarshaling file %v: %v", f, err)
		return err
	}
	return nil
}

// CheckJSONFileAge
