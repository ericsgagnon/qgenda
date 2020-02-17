package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
	"text/template"

	// "io/ioutil"
	"log"
	// "net/http"
	"os"
	// "strings"
	"time"
)

//https://restapi.qgenda.com/?version=latest

var err error

// "https://api.qgenda.com/v2/schedule/openshifts?companyKey=00000000-0000-0000-0000-000000000000&startDate=1/1/2012&endDate=1/31/2012&includes=LocationTags"
// "https://api.qgenda.com/v2/schedule/openshifts?companyKey=00000000-0000-0000-0000-000000000000&startDate=1/1/2014&endDate=1/31/2014&$select=Date,TaskAbbrev,OpenShiftCount&$filter=IsPublished&$orderby=Date,TaskAbbrev,OpenShiftCount&includes=Task"
func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)
	ctx := context.Background()
	// Set a duration.
	// duration := 150 * time.Millisecond

	// // Create a context that is both manually cancellable and will signal
	// // a cancel at the specified duration.
	// ctx, cancel := context.WithTimeout(context.Background(), duration)
	// defer cancel()

	// use environment variables to provide credentials
	q, err := NewQgendaClient(
		QgendaClientConfig{
			BaseURL:       "https://api.qgenda.com/v2",
			ClientTimeout: time.Second * 10,
			// grab credentials from environment variables
			Email:      os.Getenv("QGENDA_EMAIL"),
			CompanyKey: os.Getenv("QGENDA_COMPANY_KEY"),
			Password:   os.Getenv("QGENDA_PASSWORD"),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	err = q.Auth(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// var cq CompanyQuery
	// cq = CompanyQuery{
	// 	Route:    "/company",
	// 	Includes: "Profiles,Organizations",
	// 	Select:   "example select statement",
	// 	Filter:   "justafilter",
	// 	OrderBy:  "orderingordering",
	// 	Expand:   "expanding",
	// }
	cq := NewCompanyQuery()
	StructToQuery(cq)
}

// NewCompanyQuery returns a point to a CompanyQuery with default values
func NewCompanyQuery() *CompanyQuery {
	cq := &CompanyQuery{
		Route:    "/company",
		Includes: "Profiles,Organizations",
		// Select:   "example select statement",
		// Filter:   "justafilter",
		// OrderBy:  "orderingordering",
		// Expand:   "expanding",
	}
	return cq
}

// EncodePath uses html template to interpolate path values for an endpoint
func EncodePath(data interface{}) (string, error) {
	d := reflect.ValueOf(data)
	templateText := reflect.Indirect(d).FieldByName("Route").Interface().(string)
	// fmt.Println(templateText)
	// return templateText, nil
	t, err := template.New("path").Parse(templateText)
	if err != nil {
		log.Printf("Error Parsing Template: %v", err)
		return "", err
	}
	var bb bytes.Buffer
	err = t.Execute(&bb, d)
	if err != nil {
		log.Printf("Error Executing Template: %v", err)
		return "", err
	}
	p := bb.String()
	p = path.Join(p)
	return p, nil
}

// EncodeURLValues extracts struct values that match the provided tag and encodes them into a
// an escaped string of the url.Values
func EncodeURLValues(data interface{}, tag string) (string, error) {
	d := reflect.ValueOf(data)
	dv := reflect.Indirect(d)
	uv := url.Values{}
	for i := 0; i < dv.NumField(); i++ {
		query, ok := dv.Type().Field(i).Tag.Lookup(tag)
		if ok {
			val := fmt.Sprintf("%v", dv.Field(i).Interface())
			if val != "" {
				uv.Add(query, val)
			}
			//qv[query] = []string{val}

		}
	}

	u := uv.Encode()
	return u, nil

}

// StructToQuery awesomeness
func StructToQuery(qs interface{}) {

	// v := reflect.ValueOf(qs).Elem()
	// fmt.Printf("%#v\n", v)

	p, err := EncodePath(qs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(p)

	// Build path
	// fmt.Println("------------------------------------------------------------------------")
	// pathTemplate := v.FieldByName("Route").Interface().(string)
	// t, err := template.New("path").Parse(pathTemplate)
	// if err != nil {
	// 	log.Printf("Terrible Error: %v", err)
	// }

	// var bb bytes.Buffer
	// err = t.Execute(&bb, v)
	// if err != nil {
	// 	log.Printf("Terrible Error: %v", err)
	// }
	// urlPath := bb.String()
	// urlPath = path.Join(urlPath)
	// fmt.Println(urlPath)

	// Encode query
	fmt.Println("------------------------------------------------------------------------")
	u, err := EncodeURLValues(qs, "query")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u)
	// qv := url.Values{}
	// for i := 0; i < v.NumField(); i++ {
	// 	query, ok := v.Type().Field(i).Tag.Lookup("query")
	// 	if ok {
	// 		val := fmt.Sprintf("%v", v.Field(i).Interface())
	// 		//qv[query] = []string{val}
	// 		qv.Add(query, val)

	// 	}
	// }

	// urlQuery := qv.Encode()
	// fmt.Sprintf(urlQuery)

	// Encode body
	fmt.Println("------------------------------------------------------------------------")
	b, err := EncodeURLValues(qs, "body")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b)

	// for i := 0; i < v.NumField(); i++ {
	// 	ni := v.Type().Field(i).Name
	// 	ti := v.Type().Field(i).Type
	// 	vi := v.Field(i).Interface()
	// 	tgi := v.Type().Field(i).Tag.Get("query")
	// 	fmt.Printf("Name: %v\tType: %v\tValue: %v\tTag: %v\n", ni, ti, vi, tgi)
	// }
	// p, err := EncodePath(v)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(p)

	// var qv map[string][]string
}

// Get handles all aspects of the http get request and handling the response
func (q *QgendaClient) Get(ctx context.Context, url string, qp *url.Values, s *[]interface{}) error {

	if err := q.Auth(ctx); err != nil {
		log.Printf("Error authorizing get request to %v: %v", url, err)
		return err
	}
	qp.Add("companyKey", q.Credentials.Get("companyKey"))

	endpoint := path.Join(url, qp.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, strings.NewReader(qp.Encode()))
	if err != nil {
		log.Printf("Error in request to %v: %v", url, err)
		return err
	}
	req.Header = q.Authorization.Token.Clone()
	res, err := q.Client.Do(req)
	if err != nil {
		log.Printf("Error retrieving response from %v: %v", endpoint, err)
		return err
	}
	// TODO: improve reading response for larger requests
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response from %v: %v", endpoint, err)
		return err
	}
	defer res.Body.Close()
	if err := json.Unmarshal(body, s); err != nil {
		log.Printf("Error unmarshalling response from %v", err)
		return err
	}

	return nil
}

// fmt.Printf("\n----------------------------------------------------------------------------------\n")
// fmt.Printf("\n%v\n", res.Request)
// fmt.Printf("\n----------------------------------------------------------------------------------\n")
// fmt.Printf("\n%v\n", q.Client.Transport)
// req.Header.Add(http.CanonicalHeaderKey("Accept-Encoding"), "*gzip")
// req.Header.Add(http.CanonicalHeaderKey("Accept-Encoding"), "*")
// req.Header.Add(http.CanonicalHeaderKey("Content-Type"), "application/json")

// g, err := gzip.NewReader(res.Body)
// if err != nil {
// 	log.Fatalln(err)
// }
// if _, err := io.Copy(os.Stdout, g); err != nil {
// 	log.Fatal(err)
// }

// if err := g.Close(); err != nil {
// 	log.Fatal(err)
// }
// print token
// fmt.Printf("Authorization: %#v\n%v\n",
// 	q.Authorization.Expires.Format(time.RFC3339),
// 	q.Authorization.Token.Get(http.CanonicalHeaderKey("Authorization")),
// )
// for _, v := range q.Authorization.Cookies(q.BaseURL) {
// 	fmt.Printf("%v: %v\n%v\n", v.Name, v.Expires, v.Value)
// }

// fmt.Println("---------------------------------------------------------")

// fmt.Println("---------------------------------------------------------")

// }

// body, err := ioutil.ReadAll(res.Body)
// if err != nil {
// 	log.Fatal(err)

// }

// fmt.Println(string(body))

// url := "https://api.qgenda.com/v2/company?includes=Profiles,Organizations"
// url := "https://api.qgenda.com/v2/staffmember?companyKey=" + q.Credentials.Get("companyKey") + "&includes=Skillset,Tags,Profiles,TTCMTags"
// fmt.Println(url)
// t := map[string][]string.(q.Credentials)["companyKey"]
// companyKey = "8c44c075-d894-4b00-9ae7-3b3842226626"
// profileKey = "7f4d8aa0-292d-43b9-bec9-d253624c7de0"

//url := "https://api.qgenda.com/v2/facility?companyKey=" + q.Credentials.Get("companyKey") + "&includes=TaskShift"
// url := "https://api.qgenda.com/v2/location?companyKey=" + q.Credentials.Get("companyKey")
// method := "GET"

// payload := strings.NewReader("")

// client := &http.Client{}
// req, err := http.NewRequest(method, url, payload)

// if err != nil {
// 	fmt.Println(err)
// }
// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// // req.Header.Add("Content-Type", "application/json")
// req.Header.Add(
// 	http.CanonicalHeaderKey("Authorization"),
// 	q.Authorization.Token.Get(http.CanonicalHeaderKey("Authorization")),
// )
// // req.Header.Add(
// // 	http.CanonicalHeaderKey("Accept-Encoding"),
// // 	"*",
// // )
// //req.Header[http.CanonicalHeaderKey("Authorization")] = q.Auth.Token
// res, err := client.Do(req)
// if err != nil {
// 	log.Fatal(err)
// }
// defer res.Body.Close()
// body, err := ioutil.ReadAll(res.Body)

// fmt.Println(string(body))
// ioutil.WriteFile("samples/staffmembers.json", body, 0777)

// date := "2100-01-01T00:00:00"
// dateTime, err := time.ParseInLocation(time.RFC3339, date, time.Local)
// if err != nil {
// 	log.Fatal(err)

// }
// fmt.Println(dateTime)
