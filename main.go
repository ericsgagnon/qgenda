package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ericsgagnon/qgenda/pkg/qgenda"
)

// steps:
// export qgenda collection from postman to src/qgenda_restapi.postman_collection.json
// cat src/qgenda_restapi.postman_collection.json | yq eval '
//     .item.[] |
//     select( .name == "API Calls" ) |
//     .item.[].item.[] |
//     select( .request.method == "GET" ) |
//     [ select( .request.url.path.[] | contains( ":" ) | not ) ]
// ' -P - > src/qgenda-api-get.yaml
// note that either our login only has limited access or many endpoints aren't implemented for us

func main() {
	ctx := context.Background()
	qcc := &qgenda.ClientConfig{
		Email:    os.Getenv("QGENDA_EMAIL"),
		Password: os.Getenv("QGENDA_PASSWORD"),
	}
	c, err := qgenda.NewClient(qcc)
	if err != nil {
		log.Fatalln(err)
	}
	c.Auth()
	// Schedule

	// scheduleStartDate := time.Now().UTC().Add(-1 * 14 * 24 * time.Hour)
	// scheduleEndDate := time.Now().UTC()
	// srrqf := &qgenda.RequestQueryFields{
	// 	ScheduleStartDate: &scheduleStartDate,
	// 	ScheduleEndDate:   &scheduleEndDate,
	// }
	srrqf := &qgenda.RequestQueryFields{}
	srrqf.SetStartDate(time.Now().UTC().Add(-1 * 14 * 24 * time.Hour))
	srrqf.SetEndDate(time.Now().UTC())
	srrqf.SetSinceModifiedTimestamp(time.Now().UTC().Add(-1 * 14 * 24 * time.Hour))
	sr := qgenda.NewScheduleRequest(srrqf)
	// srJSON, err := json.MarshalIndent(sr, "", "\t")
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(string(srJSON))
	resp, err := c.Do(ctx, sr)
	if err != nil {
		log.Println(err)
	}
	data, err := io.ReadAll(resp.Body)
	// fmt.Println(resp.Status)
	// fmt.Println(string(data))
	var sch []qgenda.Schedule
	if err := json.Unmarshal(data, &sch); err != nil {
		log.Println(err)
	}
	// fmt.Println(sch)
	qgenda.Process(sch)
	jsonOut, err := json.MarshalIndent(sch, "", "\t")
	if err != nil {
		log.Println(err)
	}
	os.WriteFile("schedule.json", jsonOut, 0644)

	// ScheduleAuditLog
	salrqf := &qgenda.RequestQueryFields{}
	sal := qgenda.NewScheduleAuditLogRequest(salrqf)
	sal.SetScheduleStartDate(time.Now().UTC().Add(-1 * 14 * 24 * time.Hour))
	sal.SetScheduleEndDate(time.Now().UTC())

	resp, err = c.Do(ctx, sal)
	if err != nil {
		log.Fatalln(err)
	}
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("scheduleAuditLog success???")
	os.WriteFile("scheduleAuditLog.json", data, 0644)

	// Tag
	talrqf := &qgenda.RequestQueryFields{}
	tr := qgenda.NewTagRequest(talrqf)
	resp, err = c.Do(ctx, tr)
	if err != nil {
		log.Fatalln(err)
	}
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("tags success???")
	os.WriteFile("tags.json", data, 0644)

	// sch, err := qgenda.ScheduleFromHTTPResponse(resp)

}

// Parameters is a key-value map to represent arguments
// it is generally used to pass arguments for getting or sending
// data in data models
type Parameters map[any]any

type Inner struct {
	Value string
}

func (i *Inner) Print() {
	fmt.Println("inner")
}

/////////////////////////////////////////////
type App struct {
	Config      interface{}
	Clients     []*http.Client
	DataObjects []DataObject
}

type DataObject struct {
	Schema    struct{}
	Endpoints []Endpoint
}

type Endpoint struct {
	URL *url.URL
}

func (e *Endpoint) Request(u *url.Values) *http.Request {
	r := http.Request{}
	return &r
}

func (do *DataObject) Request(s string) *http.Request {
	return &http.Request{}
}

type Schedule struct {
	QgendaScheduleEndpoint   Endpoint
	QgendaScheduleRequest    struct{}
	QgendaScheduleResponse   struct{}
	Data                     interface{}
	PostgresScheduleEndpoint Endpoint
	PostgresScheduleRequest  struct{}
	PostgresScheduleResponse struct{}
	OracleScheduleEnpoint    Endpoint
	OracleScheduleRequest    struct{}
	OracleScheduleResponse   struct{}
	ProtobufScheduleEndpoint Endpoint
	ProtobufScheduleRequest  struct{}
	ProtobufScheduleResponse struct{}
}

// Model encapsulates the following elements:
// - data: go representation of the data - prefers structs
// - endpoints: translations to and from external representations or systems
// - process: sequence of zero of more operations to validate/transform data
type Model interface {
	Endpoints() []Endpoint
	Endpoint(s string) Endpoint
	Data() *any
	Process() error
}

// fmt.Println(qgenda.Config{})
// fmt.Println("test")
// x := qgenda.NewRequest()
// fmt.Println(x)
// x.SetRangeEndDate(time.Now().UTC())
// x.StartDate = timePointer(time.Now().UTC().AddDate(0, 0, -5))
// x.SetStartDate(time.Now().UTC().AddDate(0, 0, -5))
// v, _ := query.Values(x.RequestQueryFields)
// fmt.Println(v.Encode())
// fmt.Println(x.Parse().Encode())
// y := qgenda.NewScheduleRequest(nil)
// fmt.Println(y.ToHTTPRequest().URL.String())

// z := int(3)
// fmt.Println(z)

// zz := new(int)
// zz = int(3)
// fmt.Println(zz)

// zz := Parameters{}
// zz["bool"] = true
// zz["int"] = 3
// zz["string"] = "string"
// for k, v := range zz {
// 	fmt.Printf("%#v[%T]:\t%#v\n", k, v, v)
// }

// us := "https://restapi.qgenda.com/v2/schedule/?CompanyKey=12345678&startDate=2021-12-01"
// u1, err := url.Parse(us)
// if err != nil {
// 	log.Fatalln(err)

// }
// fmt.Printf("%#v\n", u1)

// u2, err := url.ParseRequestURI(us)
// if err != nil {
// 	log.Fatalln(err)
// }
// fmt.Printf("%#v\n", u2)
// fmt.Println(u2)

// u3, err := url.ParseQuery(us)
// if err != nil {
// 	log.Fatalln(err)

// }
// fmt.Println(u3)
// for k, v := range u3 {
// 	fmt.Printf("%#v:\t%#v\n", k, v)
// }
