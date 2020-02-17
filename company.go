package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

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

// func (q *QgendaClient) GetCompanies(ctx context.Context, cq *CompanyRequest, c *[]Company) error {
// 	//func (q *QgendaClient) Get(ctx context.Context, r Request) error Request

// }

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
func (q *QgendaClient) GetCompanies(ctx context.Context, cr *CompanyRequest, c *[]Company) error {

	if cr == nil {
		cr = NewCompanyRequest()
	}
	r, err := ParseRequest(cr)
	fmt.Println(r)
	if err != nil {
		log.Fatal(err)
	}
	bb, err := q.Get(ctx, r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bb))
	if err := json.Unmarshal(bb, &c); err != nil {
		log.Printf("Error unmarshalling response from %v", err)
		return err
	}

	return nil
}

// unmarshall companies data
// var c []Company
// fmt.Printf("\n\n%v\n", c)

// fmt.Printf("\n\n%v\n", c[0].CreatedTime.Valid)

// if c == nil {
// 	return fmt.Errorf(fmt.Sprintf("Error: must provide %T", c))
// }

// req, err := http.NewRequestWithContext(ctx, "GET", "https://api.qgenda.com/v2/company", strings.NewReader("?includes=Profiles,Organizations&companyKey=8c44c075-d894-4b00-9ae7-3b3842226626"))
// if err != nil {
// 	log.Fatalln(err)

// }
// req.Header = q.Authorization.Token.Clone()
// res, err := q.Client.Do(req)
// if err != nil {
// 	log.Printf("Error getting companies %v", err)
// 	return nil, err
// }
// body, err := ioutil.ReadAll(res.Body)
// if err != nil {
// 	log.Fatal(err)
// }
// defer res.Body.Close()
// ioutil.WriteFile("samples/company.json", body, 0777)
// // fmt.Println(string(body))
// if err := json.Unmarshal(body, &c); err != nil {
// 	log.Printf("Error unmarshalling companies from response: %v", err)
// 	return nil, err
// }
