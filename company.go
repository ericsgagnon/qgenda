package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// CompanyQuery is
type CompanyQuery struct {
	Route    string `path:"-"`
	Includes string `query:"includes"`
	Select   string `query:"$select"`
	Filter   string `query:"$filter"`
	OrderBy  string `query:"$orderby"`
	Expand   string `query:"$expand"`
}

type Query struct {
	Path  string
	Query string
	Body  string
}

// NewCompanyQuery returns a point to a CompanyQuery with default values
func NewCompanyQuery() *CompanyQuery {
	cq := &CompanyQuery{
		Route:    "/company",
		Includes: "Profiles,Organizations",
		// Select:   "",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}
	return cq
}

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

// MarshallQuery takes a query struct
func MarshallQuery(ctx context.Context, qs *interface{}) error {

	// path variables can use html/template to build paths

	// query parameters just need to be in url form

	// body variables are in ?
	// url.Par
	// var companies []Company
	// companies, err := q.GetCompanies(ctx)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(companies)
	var cq *CompanyQuery
	cq = &CompanyQuery{
		Route:    "/company",
		Includes: "includes=Profiles,Organizations",
	}

	fmt.Println(reflect.ValueOf(cq).Kind().String())
	fmt.Println(reflect.ValueOf(cq).Elem().Kind().String())
	if reflect.ValueOf(cq).Elem().Kind() == reflect.Struct {
		fmt.Println("Yeah Baby")
	}
	fmt.Println(reflect.ValueOf(cq).Elem().Type())
	fmt.Println(reflect.ValueOf(cq).Elem().NumField())
	fmt.Println(reflect.ValueOf(cq).Elem().Type().Field(0))
	fmt.Println(reflect.ValueOf(cq).Elem().Type().Field(0).Tag.Get("path"))

	// v := reflect.ValueOf(q)
	// fmt.Println(v.Kind().String())
	// vv := v.Elem()
	// fmt.Println(vv)
	// fmt.Println(vv.Kind().String())
	// var t *map[string]string
	// tt := reflect.ValueOf(t)
	// fmt.Println(tt.Kind().String())
	return nil
}

// GetCompanies uses the company endpoint to get all companies for a user
func (q *QgendaClient) GetCompanies(ctx context.Context, cq *CompanyQuery, c *[]Company) error {

	if cq == nil {
		cq = &CompanyQuery{
			Route:    "/company",
			Includes: "includes=Profiles,Organizations",
		}
	}

	// 	"includes", "Profiles,Organizations"

	// 	if qp == nil {

	// 	}

	// var c *[]Company
	// err := q.Get(ctx)

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

	return nil
}
