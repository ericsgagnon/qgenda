package main

import (
	"time"

	"github.com/google/uuid"
)

// // StaffType represents preconfigured values
// type StaffType int

// // preconfigured values for stafftype
// const (
// 	nil StaffType = iota
// 	Physician
// 	CRNA
// 	Technologist
// 	Locum
// 	Office
// 	Resident
// 	Specialty
// 	MAPA
// 	Nurse
// 	Other
// )

// Time embeds time.Time
// qgenda doesn't comply with RFC3339
// chose this over wrapping to slightly improve
// convenience of calling time.Time's methods
// type Time struct {
// 	time.Time
// }

// // UnmarshalJSON satisfies the json.Unmarshaler interface
// func (t *Time) UnmarshalJSON(data []byte) error {

// 	tag := reflect.ValueOf(data).Type().Field(0).Tag.Get("json")
// 	fmt.Println(tag)

// 	location, err := time.LoadLocation("America/New_York")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// Ignore null, like in the main JSON package.
// 	if string(data) == "null" {
// 		return nil
// 	}
// 	// *t, err = Parse(`"`+RFC3339+`"`, string(data))
// 	t.Time, err = time.ParseInLocation("2006-01-02T15:04:05", string(data), location)
// 	return err
// }

// StaffMember represents staff, and possibly some other entities as well?
type StaffMember struct {
	Abbreviation         string    `json:"Abbrev"`
	BackgroundColor      string    `json:"BgColor"`
	BillingSystemID      uuid.UUID `json:"BillSysId"`
	CalendarSyncID       uuid.UUID `json:"CalSyncKey"`
	CompanyID            uuid.UUID `json:"CompKey"`
	Email                string    `json:"Email"`
	EMRID                string    `json:"EmrId"`
	ERPID                string    `json:"ErpId"`
	EndDate              Time      `json:"EndDate"`
	ExternalCallSystemID uuid.UUID `json:"ExtCallSysId"`
	FirstName            string    `json:"FirstName"`
	HomePhone            string    `json:"HomePhone"`
	LastName             string    `json:"LastName"`
	MobilePhone          string    `json:"MobilePhone"`
	NPI                  string    `json:"Npi"`
	PagerNumber          string    `json:"Pager"`
	PayrollID            string    `json:"PayrollId"`
	RegularHours         string    `json:"RegHours"`
	Alias                string    `json:"StaffId"`
	ID                   uuid.UUID `json:"StaffKey"`
	StartDate            Time      `json:"StartDate"`
	TextColor            string    `json:"TextColor"`
	Active               bool      `json:"IsActive"`
	StaffType            string    `json:"StaffTypeKey"`
	BillingType          string    `json:"BillingTypeKey"`
	ProfileID            uuid.UUID `json:"UserProfileKey"`
	Profile              string    `json:"UserProfile"`
	PayrollStartDate     time.Time `json:"PayrollStartDate"`
	PayrollEndDate       time.Time `json:"PayrollEndDate"`
	TimeClockStartDate   time.Time `json:"TimeClockStartDate"`
	TimeClockEndDate     time.Time `json:"TimeClockEndDate"`
	TimeClockKioskPIN    string    `json:"TimeClockKioskPIN"`
	AutoApproveSwap      bool      `json:"IsAutoApproveSwap"`
	Viewable             bool      `json:"IsViewable"`
	Schedulable          bool      `json:"IsSchedulable"`
	Address              struct {
		Line1 string `json:"Addr1"`
		Line2 string `json:"Addr2"`
		City  string `json:"City"`
		State string `json:"State"`
		Zip   string `json:"Zip"`
	}
	LastLogin struct {
		Time   time.Time `json:"UserLastLoginDateTimeUtc"`
		Source string    `json:"SourceOfLogin"`
	}
	SkillSet []SkillSet `json:"Skillset"`

	// The don't appear to be used at our insitution, not sure if they are elsewhere
	// omitting due to inability to validate
	// `json:"Tags"`
	// TTCMTags `json:"TTCMTags"`
	// `json:"CategoryKey"`
	// `json:"CategoryName"`
	// Tags `json:"Tags"`
	// `json:"Key"`
	// `json:"Name"`
	// `json:"EffectiveFromDate"`
	// `json:"EffectiveToDate"`
	// `json:"Profiles"`
	// `json:"Name"`
	// `json:"ProfileKey"`
	// `json:"DailyUnitAverage"` // don't waste your time on this
}

// Skillset captures the staff to task relationship
type SkillSet struct {
	Staff struct {
		FirstName    string `json:"StaffFirstName"`
		LastName     string `json:"StaffLastName"`
		Abbreviation string `json:"StaffAbbreviation"`
		Alias        string `json:"StaffId"`
	}
	Task struct {
		Name         string `json:"TaskName"`
		Abbreviation string `json:"TaskAbbrev"`
		Alias        string `json:"TaskId"`
	}
	Monday struct {
		Valid     bool   `json:"IsSkilledMon"`
		Frequency string `json:"MonOccurrence"`
	}
	Tuesday struct {
		Valid     bool   `json:"IsSkilledTue"`
		Frequency string `json:"TueOccurrence"`
	}
	Wednesday struct {
		Valid     bool   `json:"IsSkilledWed"`
		Frequency string `json:"WedOccurrence"`
	}
	Thursday struct {
		Valid     bool   `json:"IsSkilledThu"`
		Frequency string `json:"ThuOccurrence"`
	}
	Friday struct {
		Valid     bool   `json:"IsSkilledFri"`
		Frequency string `json:"FriOccurrence"`
	}
	Saturday struct {
		Valid     bool   `json:"IsSkilledSat"`
		Frequency string `json:"SatOccurrence"`
	}
	Sunday struct {
		Valid     bool   `json:"IsSkilledSun"`
		Frequency string `json:"SunOccurrence"`
	}
}

// func (q *QgendaClient) Get(ctx context.Context, url string)

// // GetStaffMember returns all staff members
// func (q *QgendaClient) GetStaffMember(ctx context.Context) context.Context {
// 	//TODO: check for Auth cookie or header, get another if missing or expired
// 	// request URL
// 	url := *q.BaseURL
// 	uri := "/staffmember?companyKey=" + q.Credentials.Get("companyKey") + "&includes=Skillset,Tags,Profiles,TTCMTags"
// 	url.Path = path.Join(url.Path, uri)

// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	// // req.Header.Add("Content-Type", "application/json")
// 	// req.Header.Add(
// 	// 	http.CanonicalHeaderKey("Authorization"),
// 	// 	q.Authorization.Token.Get(http.CanonicalHeaderKey("Authorization")),
// 	// )
// 	// // req.Header.Add(
// 	// // 	http.CanonicalHeaderKey("Accept-Encoding"),
// 	// // 	"*",
// 	// // )
// 	// //req.Header[http.CanonicalHeaderKey("Authorization")] = q.Auth.Token
// 	// res, err := client.Do(req)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// defer res.Body.Close()
// 	// body, err := ioutil.ReadAll(res.Body)

// 	// fmt.Println(string(body))

// 	// request
// 	// res, err := q.Client.PostForm(reqURL.String(), *q.Credentials)
// 	res, err := ctxhttp.PostForm(ctx, q.Client, url.String(), *q.Credentials)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer res.Body.Close()

// 	//response body
// 	resBody, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// authorization token is returned in the response body
// 	var resData map[string]string
// 	json.Unmarshal(resBody, &resData)

// 	// use response timestamp + valid duration to set expire time
// 	respTime, err := time.Parse(time.RFC1123, res.Header[http.CanonicalHeaderKey("date")][0])
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	validDuration, err := time.ParseDuration(resData["expires_in"] + "s")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	expireTime := respTime.Add(validDuration)

// 	// Set the Authorization header in the QgendaClient
// 	q.Authorization.Token.Set(
// 		http.CanonicalHeaderKey("Authorization"),
// 		fmt.Sprintf("bearer %v", resData["access_token"]),
// 	)

// 	q.Authorization.Expires = expireTime
// 	// set Authorization cookie for all endpoints
// 	u := *q.BaseURL
// 	u.Path = "/"

// 	fmt.Printf("Authorization: %#v\n%v\n",
// 		q.Authorization.Expires.Format(time.RFC3339),
// 		q.Authorization.Token.Get(http.CanonicalHeaderKey("Authorization")),
// 	)

// 	return ctx
// }

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
