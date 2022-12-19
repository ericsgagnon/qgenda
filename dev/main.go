package main

import (
	"bytes"
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
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
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
	// pgx way
	db, err := pgx.Connect(ctx, os.Getenv("PG_CONNECTION_STRING"))
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping(ctx)
	defer db.Close(ctx)

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

	// configure request
	srrqf := &qgenda.RequestConfig{}
	srrqf.SetStartDate(time.Now().UTC().Add(-1 * 14 * 24 * time.Hour))
	srrqf.SetEndDate(time.Now().UTC())
	srrqf.SetSinceModifiedTimestamp(time.Now().UTC().Add(-1 * 14 * 24 * time.Hour))
	sr := qgenda.NewScheduleRequest(srrqf)

	// get data
	resp, err := c.Do(ctx, sr)
	for k, v := range resp.Header {
		fmt.Printf("%20s %-80s\n", k, "-")
		for vi, vv := range v {
			fmt.Printf("\t%3d: %40s\n", vi, vv)
		}
	}
	resp.Header.Get(http.CanonicalHeaderKey("Date"))

	if err != nil {
		log.Println(err)
	}
	data, err := io.ReadAll(resp.Body)
	// data2 := *&data
	var sch []qgenda.Schedule
	if err := json.Unmarshal(data, &sch); err != nil {
		log.Println(err)
	}

	// process data
	qgenda.Process(sch)

	// load data
	jsonOut, err := json.MarshalIndent(sch, "", "\t")
	if err != nil {
		log.Println(err)
	}
	os.WriteFile("schedule.json", jsonOut, 0644)

	// var rawscheduledata []byte
	// if err := json.Unmarshal(data, &sch); err != nil {
	// 	log.Println(err)
	// }
	// load data
	// jsonOut2, err := json.MarshalIndent(data2, "", "\t")
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(string(data))
	os.WriteFile("rawschedule.json", data, 0644)

	// var locationdata []byte
	// lrqf := &qgenda.RequestConfig{}
	// lr := qgenda.NewLocationRequest(lrqf)
	// resp, err = c.Do(ctx, lr)
	// if err != nil {
	// 	log.Println(err)
	// }
	// data, err = io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// }
	// var out bytes.Buffer
	// if err := json.Indent(&out, data, "", "\t"); err != nil {
	// 	log.Println(err)
	// }
	// if err := os.WriteFile("locations.json", out.Bytes(), 0644); err != nil {
	// 	log.Println(err)
	// }

	if err := MakeItHappen(ctx, c, qgenda.NewCompanyRequest(&qgenda.RequestConfig{}), "company.json"); err != nil {
		log.Println(err)
	}
	// if err := MakeItHappen(ctx, c, qgenda.NewDailyCaseRequest(&qgenda.RequestConfig{}), "dailycase.json"); err != nil {
	// 	log.Println(err)
	// }
	// if err := MakeItHappen(ctx, c, qgenda.NewDailyDailyConfigurationRequest(&qgenda.RequestConfig{}), "dailydailyconfiguration.json"); err != nil {
	// 	log.Println(err)
	// }
	// if err := MakeItHappen(ctx, c, qgenda.NewDailyPatientEncounterRequest(&qgenda.RequestConfig{}), "dailypatientencounter.json"); err != nil {
	// 	log.Println(err)
	// }
	// if err := MakeItHappen(ctx, c, qgenda.NewDailyRoomRequest(&qgenda.RequestConfig{}), "dailyroom.json"); err != nil {
	// 	log.Println(err)
	// }
	// if err := MakeItHappen(ctx, c, qgenda.NewLocationRequest(&qgenda.RequestConfig{}), "locations.json"); err != nil {
	// 	log.Println(err)
	// }
	// if err := MakeItHappen(ctx, c, qgenda.NewLocationStaffRequest(&qgenda.RequestConfig{}), "locationstaff.json"); err != nil {
	// 	log.Println(err)
	// }
	// if err := MakeItHappen(ctx, c, qgenda.NewLocationTasksRequest(&qgenda.RequestConfig{}), "locationtasks.json"); err != nil {
	// 	log.Println(err)
	// }
	if err := MakeItHappen(ctx, c, qgenda.NewOpenShiftsRequest(&qgenda.RequestConfig{}), "openshifts.json"); err != nil {
		log.Println(err)
	}
	// if err := MakeItHappen(ctx, c, qgenda.NewOrganizationRequest(&qgenda.RequestConfig{}), "organization.json"); err != nil {
	// 	log.Println(err)
	// }
	// if err := MakeItHappen(ctx, c, qgenda.NewPayRateRequest(&qgenda.RequestConfig{}), "payrate.json"); err != nil {
	// 	log.Println(err)
	// }
	if err := MakeItHappen(ctx, c, qgenda.NewProfileRequest(&qgenda.RequestConfig{}), "profile.json"); err != nil {
		log.Println(err)
	}
	if err := MakeItHappen(ctx, c, qgenda.NewRequestsApprovedRequest(
		&qgenda.RequestConfig{
			StartDate: qgenda.Pointer(time.Now().UTC().Add(-1 * 14 * 24 * time.Hour)),
			EndDate:   qgenda.Pointer(time.Now().UTC()),
		}), "requestapproved.json"); err != nil {
		log.Println(err)
	}
	if err := MakeItHappen(ctx, c, qgenda.NewRequestsRequest(&qgenda.RequestConfig{}), "request.json"); err != nil {
		log.Println(err)
	}
	// configure request
	// rr := qgenda.NewRequestsRequest(
	// 	&qgenda.RequestConfig{
	// 		StartDate: qgenda.Pointer(time.Now().UTC().Add(-1 * 14 * 24 * time.Hour)),
	// 		EndDate:   qgenda.Pointer(time.Now().UTC()),
	// 	})
	// rr := qgenda.NewRequestsRequest(nil)
	// // get data
	// resp, err = c.Do(ctx, rr)
	// if err != nil {
	// 	log.Println(err)
	// }
	// data, err = io.ReadAll(resp.Body)
	// // data2 := *&data
	// var requests []qgenda.Requests
	// if err := json.Unmarshal(data, &requests); err != nil {
	// 	log.Println(err)
	// }

	// // process data
	// qgenda.Process(requests)

	// // load data
	// jsonOut, err = json.MarshalIndent(requests, "", "\t")
	// if err != nil {
	// 	log.Println(err)
	// }
	// os.WriteFile("requests.json", jsonOut, 0644)

	if err := HandleStructuredData[qgenda.StaffMember](ctx, c, qgenda.NewStaffMemberRequest(&qgenda.RequestConfig{}), "out/staffmember.json"); err != nil {
		log.Println(err)
	}
	if err := HandleStructuredData[qgenda.Schedule](ctx, c, qgenda.NewScheduleRequest(&qgenda.RequestConfig{}), "out/schedule.json"); err != nil {
		log.Println(err)
	}
	if err := HandleStructuredData[qgenda.Requests](ctx, c, qgenda.NewRequestsRequest(&qgenda.RequestConfig{}), "out/requests.json"); err != nil {
		log.Println(err)
	}
	// if err := MakeItHappen(ctx, c, qgenda.NewStaffMemberLocationRequest(&qgenda.RequestConfig{}), "staffmemberlocation.json"); err != nil {
	// 	log.Println(err)
	// }
	if err := MakeItHappen(ctx, c, qgenda.NewScheduleAuditLogRequest(&qgenda.RequestConfig{}), "scheduleauditlog.json"); err != nil {
		log.Println(err)
	}
	if err := MakeItHappen(ctx, c, qgenda.NewScheduleRequest(&qgenda.RequestConfig{
		StartDate:              qgenda.Pointer(time.Now().UTC().Add(-1 * 14 * 24 * time.Hour)),
		EndDate:                qgenda.Pointer(time.Now().UTC()),
		SinceModifiedTimestamp: qgenda.Pointer(time.Now().UTC().Add(-1 * 14 * 24 * time.Hour)),
	}), "rawschedule.json"); err != nil {
		log.Println(err)
	}

	if err := MakeItHappen(ctx, c, qgenda.NewStaffMemberRequest(&qgenda.RequestConfig{}), "out/rawstaffmember.json"); err != nil {
		log.Println(err)
	}
	// if err := MakeItHappen(ctx, c, qgenda.NewStaffMemberRequestLimitRequest(&qgenda.RequestConfig{}), "staffmemberrequestlimit.json"); err != nil {
	// 	log.Println(err)
	// }
	// if err := MakeItHappen(ctx, c, qgenda.NewStaffMemberStaffIdRequest(&qgenda.RequestConfig{}), "staffmemberstaffid.json"); err != nil {
	// 	log.Println(err)
	// }
	// if err := MakeItHappen(ctx, c, qgenda.NewStaffTargetRequest(&qgenda.RequestConfig{}), "stafftarget.json"); err != nil {
	// 	log.Println(err)
	// }
	if err := MakeItHappen(ctx, c, qgenda.NewTagRequest(&qgenda.RequestConfig{}), "tags.json"); err != nil {
		log.Println(err)
	}

	// if err := MakeItHappen(ctx, c, qgenda.NewTaskLocationRequest(&qgenda.RequestConfig{}), "tasklocation.json"); err != nil {
	// 	log.Println(err)
	// }
	if err := MakeItHappen(ctx, c, qgenda.NewTaskRequest(&qgenda.RequestConfig{}), "task.json"); err != nil {
		log.Println(err)
	}
	// if err := MakeItHappen(ctx, c, qgenda.NewTimeEventRequest(&qgenda.RequestConfig{}), "timeevent.json"); err != nil {
	// 	log.Println(err)
	// }
	if err := MakeItHappen(ctx, c, qgenda.NewUserRequest(&qgenda.RequestConfig{}), "user.json"); err != nil {
		log.Println(err)
	}
	// if err := MakeItHappen(ctx, c, qgenda.NewUserRequest(&qgenda.RequestConfig{Expand: qgenda.Pointer("Companies")}), "userExpandedCompanies.json"); err != nil {
	// 	log.Println(err)
	// }
	// if err := MakeItHappen(ctx, c, qgenda.NewUserRequest(&qgenda.RequestConfig{Expand: qgenda.Pointer("Companies/Locations")}), "userExpandedCompaniesLocations.json"); err != nil {
	// 	log.Println(err)
	// }
	// // ScheduleAuditLog
	// salrqf := &qgenda.RequestConfig{}
	// sal := qgenda.NewScheduleAuditLogRequest(salrqf)
	// sal.SetScheduleStartDate(time.Now().UTC().Add(-1 * 14 * 24 * time.Hour))
	// sal.SetScheduleEndDate(time.Now().UTC())

	// resp, err = c.Do(ctx, sal)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// data, err = io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// // fmt.Println("scheduleAuditLog success???")
	// jsonOut, err = json.MarshalIndent(data, "", "\t")
	// if err != nil {
	// 	log.Println(err)
	// }
	// // fmt.Println(string(jsonOut))
	// os.WriteFile("scheduleAuditLog.json", jsonOut, 0644)

	// // Tag
	// talrqf := &qgenda.RequestConfig{}
	// tr := qgenda.NewTagRequest(talrqf)
	// resp, err = c.Do(ctx, tr)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// data, err = io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println("tags success???")
	// os.WriteFile("tags.json", data, 0644)

	// sch, err := qgenda.ScheduleFromHTTPResponse(resp)

}

func MakeItHappen(ctx context.Context, c *qgenda.Client, r *qgenda.Request, file string) error {
	resp, err := c.Do(ctx, r)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "\t"); err != nil {
		return err
	}
	if err := os.WriteFile(file, out.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func HandleStructuredData[T any](ctx context.Context, c *qgenda.Client, r *qgenda.Request, file string) error {
	resp, err := c.Do(ctx, r)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	da := []T{}
	log.Printf("%#v\n", da)
	if err := json.Unmarshal(data, &da); err != nil {
		return err
	}

	// process data
	// if err := qgenda.Process(da); err != nil {
	// for i, _ := range da {
	// 	if err := da[i].Process(); err != nil {
	// 		log.Printf("HandleStructuredData %T %s\n", da, err)
	// 	}
	// }
	if err := qgenda.Process(da); err != nil {
		log.Printf("HandleStructuredData %T %s\n", da, err)
	}

	// load data
	jsonOut, err := json.MarshalIndent(da, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(file, jsonOut, 0644)

}

// Parameters is a key-value map to represent arguments
// it is generally used to pass arguments for getting or sending
// data in data models
type Parameters map[any]any

// ///////////////////////////////////////////
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

// type Schedule struct {
// 	QgendaScheduleEndpoint   Endpoint
// 	QgendaScheduleRequest    struct{}
// 	QgendaScheduleResponse   struct{}
// 	Data                     interface{}
// 	PostgresScheduleEndpoint Endpoint
// 	PostgresScheduleRequest  struct{}
// 	PostgresScheduleResponse struct{}
// 	OracleScheduleEnpoint    Endpoint
// 	OracleScheduleRequest    struct{}
// 	OracleScheduleResponse   struct{}
// 	ProtobufScheduleEndpoint Endpoint
// 	ProtobufScheduleRequest  struct{}
// 	ProtobufScheduleResponse struct{}
// }

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
// v, _ := query.Values(x.RequestConfig)
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
