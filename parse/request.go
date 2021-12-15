package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
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
	// fmt.Println("test")
	x := NewRequest()
	// fmt.Println(x)
	x.RangeEndDate = timePointer(time.Now().UTC())
	x.StartDate = timePointer(time.Now().UTC().AddDate(0, 0, -5))
	v, _ := query.Values(x.RequestQueryFields)
	fmt.Println(v.Encode())

	// structFields(x)
}

// func structFields(value interface{}) {
// 	v := reflect.ValueOf(value)
// 	// values := make([]interface{}, v.NumField())
// 	fmt.Printf("Type:\t%#v\n", v.Type())
// 	fmt.Printf("Kind:\t%#v\n", v.Kind())
// 	fmt.Printf("Elem:\t%#v\n", v.Elem())
// 	fmt.Printf("Indirect(Elem):\t%#v\n", reflect.Indirect(v.Elem()))
// 	fmt.Printf("NumField:\t%#v\n", v.Elem().NumField())
// 	fmt.Printf("Elem.Kind:\t%#v\n", reflect.Indirect(v.Elem()).Kind())
// 	for i := 0; i < v.Elem().NumField(); i++ {
// 		fmt.Printf("%#v\n", v.Elem().FieldByIndex([]int{i}))
// 	}
// 	// fmt.Printf("NumMethod:\t%#v\n", v.NumMethod())
// 	// fmt.Printf("Kind:\t%#v\n", v.Kind())
// 	// fmt.Printf("Kind:\t%#v\n", v.Kind())
// 	// for i := 0; i < v.NumField(); i++ {
// 	// 	values[i] = v.Field(i).Interface()
// 	// 	fmt.Printf("%#v\n", values[i])
// 	// }

// }

func Schedule() http.Request {
	url, err := url.Parse("https://api.qgenda.com/v2/schedule")
	if err != nil {
		panic(err)
	}

	return http.Request{
		Method: http.MethodGet,
		URL:    url,
	}
}

// func ScheduleRequest(rqf *RequestQueryFields) (*Request, error) {
// 	r := NewRequest()
// 	r.Path.String = "v2/schedule"
// 	if rqf != nil {
// 		r.Query = *rqf
// 	} else {
// 		r.Query = RequestQueryFields{
// 			StartDate: NullTime{time.Now().UTC().Add(time.Hour * 168 * 2 * -1), true},
// 			EndDate:   NullTime{time.Now().UTC(), true},
// 		}
// 	}

// 	return r, nil
// }

type Request2 struct {
	Method string
	Header http.Header
	Host   string
	Path   string
	RequestQueryFields
	CompanyKey              *string    `query:"companyKey,omitempty" url:"companyKey,omitempty"`
	OrganizationKey         *int       `query:"organizationKey,omitempty" url:"organizationKey,omitempty"`
	Expand                  *string    `query:"$expand,omitempty" url:"$expand,omitempty"`
	Filter                  *string    `query:"$filter,omitempty" url:"$filter,omitempty"`
	Orderby                 *string    `query:"$orderby,omitempty" url:"$orderby,omitempty"`
	Select                  *string    `query:"$select,omitempty" url:"$select,omitempty"`
	DailyConfigurationKey   *string    `query:"dailyConfigurationKey,omitempty" url:"dailyConfigurationKey,omitempty"`
	DateFormat              *string    `query:"dateFormat,omitempty" url:"dateFormat,omitempty"`
	EndDate                 *time.Time `query:"endDate,omitempty" url:"endDate,omitempty"`
	IgnoreHoliday           *bool      `query:"ignoreHoliday,omitempty" url:"ignoreHoliday,omitempty"`
	IgnoreWeekend           *bool      `query:"ignoreWeekend,omitempty" url:"ignoreWeekend,omitempty"`
	IncludeDeletes          *bool      `query:"includeDeletes,omitempty" url:"includeDeletes,omitempty"`
	IncludeRemoved          *bool      `query:"includeRemoved,omitempty" url:"includeRemoved,omitempty"`
	Includes                *string    `query:"includes,omitempty" url:"includes,omitempty"`
	IsUniversallyLocalDates *bool      `query:"IsUniversallyLocalDates,omitempty" url:"IsUniversallyLocalDates,omitempty"`
	MaxResults              *int       `query:"maxResults,omitempty" url:"maxResults,omitempty"`
	PageToken               *string    `query:"pageToken,omitempty" url:"pageToken,omitempty"`
	RangeEndDate            *time.Time `query:"rangeEndDate,omitempty" url:"rangeEndDate,omitempty"`
	RangeStartDate          *time.Time `query:"rangeStartDate,omitempty" url:"rangeStartDate,omitempty"`
	ScheduleEndDate         *time.Time `query:"scheduleEndDate,omitempty" url:"scheduleEndDate,omitempty"`
	ScheduleStartDate       *time.Time `query:"scheduleStartDate,omitempty" url:"scheduleStartDate,omitempty"`
	SinceModifiedTimestamp  *time.Time `query:"sinceModifiedTimestamp,omitempty" url:"sinceModifiedTimestamp,omitempty"`
	StartDate               *time.Time `query:"startDate,omitempty" url:"startDate,omitempty"`
	SyncToken               *string    `query:"syncToken,omitempty" url:"syncToken,omitempty"`
}

type Request struct {
	Method string
	Header http.Header
	Host   string
	Path   string
	RequestQueryFields
}

// NewRequest returns a Request with only common members
// it is expected to a base for other 'request' functions
func NewRequest() *Request {
	r := Request{}
	r.Host = "api.qgenda.com"
	r.DateFormat = stringPointer("yyyy-MM-ddTHH:mm:ssZ")
	// Query:      RequestQueryFields{},

	return &r
}

// func (r *RequestQueryFields) Merge(rqf *RequestQueryFields) {
// 	if rqf != nil {
// 		if rqf.CompanyKey.Valid {
// 			r.CompanyKey = rqf.CompanyKey
// 		}
// 		if rqf.OrganizationKey.Valid {
// 			r.OrganizationKey = rqf.OrganizationKey
// 		}
// 		if rqf.Expand.Valid {
// 			r.Expand = rqf.Expand
// 		}
// 		if rqf.Filter.Valid {
// 			r.Filter = rqf.Filter
// 		}
// 		if rqf.Orderby.Valid {
// 			r.Orderby = rqf.Orderby
// 		}
// 		if rqf.Select.Valid {
// 			r.Select = rqf.Select
// 		}
// 		if rqf.DailyConfigurationKey.Valid {
// 			r.DailyConfigurationKey = rqf.DailyConfigurationKey
// 		}
// 		if rqf.DateFormat.Valid {
// 			r.DateFormat = rqf.DateFormat
// 		}
// 		if rqf.EndDate.Valid {
// 			r.EndDate = rqf.EndDate
// 		}
// 		if rqf.IgnoreHoliday.Valid {
// 			r.IgnoreHoliday = rqf.IgnoreHoliday
// 		}
// 		if rqf.IgnoreWeekend.Valid {
// 			r.IgnoreWeekend = rqf.IgnoreWeekend
// 		}
// 		if rqf.IncludeDeletes.Valid {
// 			r.IncludeDeletes = rqf.IncludeDeletes
// 		}
// 		if rqf.IncludeRemoved.Valid {
// 			r.IncludeRemoved = rqf.IncludeRemoved
// 		}
// 		if rqf.Includes.Valid {
// 			r.Includes = rqf.Includes
// 		}
// 		if rqf.IsUniversallyLocalDates.Valid {
// 			r.IsUniversallyLocalDates = rqf.IsUniversallyLocalDates
// 		}
// 		if rqf.MaxResults.Valid {
// 			r.MaxResults = rqf.MaxResults
// 		}
// 		if rqf.PageToken.Valid {
// 			r.PageToken = rqf.PageToken
// 		}
// 		if rqf.RangeEndDate.Valid {
// 			r.RangeEndDate = rqf.RangeEndDate
// 		}
// 		if rqf.RangeStartDate.Valid {
// 			r.RangeStartDate = rqf.RangeStartDate
// 		}
// 		if rqf.ScheduleEndDate.Valid {
// 			r.ScheduleEndDate = rqf.ScheduleEndDate
// 		}
// 		if rqf.ScheduleStartDate.Valid {
// 			r.ScheduleStartDate = rqf.ScheduleStartDate
// 		}
// 		if rqf.SinceModifiedTimestamp.Valid {
// 			r.SinceModifiedTimestamp = rqf.SinceModifiedTimestamp
// 		}
// 		if rqf.StartDate.Valid {
// 			r.StartDate = rqf.StartDate
// 		}
// 		if rqf.SyncToken.Valid {
// 			r.SyncToken = rqf.SyncToken
// 		}

// 	}
// }

// type RequestQueryFields struct {
// 	CompanyKey              NullString `query:"companyKey,omitempty" url:"companyKey,omitempty"`
// 	OrganizationKey         NullInt64  `query:"organizationKey,omitempty" url:"organizationKey,omitempty"`
// 	Expand                  NullString `query:"$expand,omitempty" url:"$expand,omitempty"`
// 	Filter                  NullString `query:"$filter,omitempty" url:"$filter,omitempty"`
// 	Orderby                 NullString `query:"$orderby,omitempty" url:"$orderby,omitempty"`
// 	Select                  NullString `query:"$select,omitempty" url:"$select,omitempty"`
// 	DailyConfigurationKey   NullString `query:"dailyConfigurationKey,omitempty" url:"dailyConfigurationKey,omitempty"`
// 	DateFormat              NullString `query:"dateFormat,omitempty" url:"dateFormat,omitempty"`
// 	EndDate                 NullTime   `query:"endDate,omitempty" url:"endDate,omitempty"`
// 	IgnoreHoliday           NullBool   `query:"ignoreHoliday,omitempty" url:"ignoreHoliday,omitempty"`
// 	IgnoreWeekend           NullBool   `query:"ignoreWeekend,omitempty" url:"ignoreWeekend,omitempty"`
// 	IncludeDeletes          NullBool   `query:"includeDeletes,omitempty" url:"includeDeletes,omitempty"`
// 	IncludeRemoved          NullBool   `query:"includeRemoved,omitempty" url:"includeRemoved,omitempty"`
// 	Includes                NullString `query:"includes,omitempty" url:"includes,omitempty"`
// 	IsUniversallyLocalDates NullBool   `query:"IsUniversallyLocalDates,omitempty" url:"IsUniversallyLocalDates,omitempty"`
// 	MaxResults              NullInt64  `query:"maxResults,omitempty" url:"maxResults,omitempty"`
// 	PageToken               NullString `query:"pageToken,omitempty" url:"pageToken,omitempty"`
// 	RangeEndDate            NullTime   `query:"rangeEndDate,omitempty" url:"rangeEndDate,omitempty"`
// 	RangeStartDate          NullTime   `query:"rangeStartDate,omitempty" url:"rangeStartDate,omitempty"`
// 	ScheduleEndDate         NullTime   `query:"scheduleEndDate,omitempty" url:"scheduleEndDate,omitempty"`
// 	ScheduleStartDate       NullTime   `query:"scheduleStartDate,omitempty" url:"scheduleStartDate,omitempty"`
// 	SinceModifiedTimestamp  NullTime   `query:"sinceModifiedTimestamp,omitempty" url:"sinceModifiedTimestamp,omitempty"`
// 	StartDate               NullTime   `query:"startDate,omitempty" url:"startDate,omitempty"`
// 	SyncToken               NullString `query:"syncToken,omitempty" url:"syncToken,omitempty"`
// }

// RequestQueryFields is an experiment in using pointer members
// for optional fields
type RequestQueryFields struct {
	CompanyKey              *string    `query:"companyKey,omitempty" url:"companyKey,omitempty"`
	OrganizationKey         *int       `query:"organizationKey,omitempty" url:"organizationKey,omitempty"`
	Expand                  *string    `query:"$expand,omitempty" url:"$expand,omitempty"`
	Filter                  *string    `query:"$filter,omitempty" url:"$filter,omitempty"`
	Orderby                 *string    `query:"$orderby,omitempty" url:"$orderby,omitempty"`
	Select                  *string    `query:"$select,omitempty" url:"$select,omitempty"`
	DailyConfigurationKey   *string    `query:"dailyConfigurationKey,omitempty" url:"dailyConfigurationKey,omitempty"`
	DateFormat              *string    `query:"dateFormat,omitempty" url:"dateFormat,omitempty"`
	EndDate                 *time.Time `query:"endDate,omitempty" url:"endDate,omitempty"`
	IgnoreHoliday           *bool      `query:"ignoreHoliday,omitempty" url:"ignoreHoliday,omitempty"`
	IgnoreWeekend           *bool      `query:"ignoreWeekend,omitempty" url:"ignoreWeekend,omitempty"`
	IncludeDeletes          *bool      `query:"includeDeletes,omitempty" url:"includeDeletes,omitempty"`
	IncludeRemoved          *bool      `query:"includeRemoved,omitempty" url:"includeRemoved,omitempty"`
	Includes                *string    `query:"includes,omitempty" url:"includes,omitempty"`
	IsUniversallyLocalDates *bool      `query:"IsUniversallyLocalDates,omitempty" url:"IsUniversallyLocalDates,omitempty"`
	MaxResults              *int       `query:"maxResults,omitempty" url:"maxResults,omitempty"`
	PageToken               *string    `query:"pageToken,omitempty" url:"pageToken,omitempty"`
	RangeEndDate            *time.Time `query:"rangeEndDate,omitempty" url:"rangeEndDate,omitempty"`
	RangeStartDate          *time.Time `query:"rangeStartDate,omitempty" url:"rangeStartDate,omitempty"`
	ScheduleEndDate         *time.Time `query:"scheduleEndDate,omitempty" url:"scheduleEndDate,omitempty"`
	ScheduleStartDate       *time.Time `query:"scheduleStartDate,omitempty" url:"scheduleStartDate,omitempty"`
	SinceModifiedTimestamp  *time.Time `query:"sinceModifiedTimestamp,omitempty" url:"sinceModifiedTimestamp,omitempty"`
	StartDate               *time.Time `query:"startDate,omitempty" url:"startDate,omitempty"`
	SyncToken               *string    `query:"syncToken,omitempty" url:"syncToken,omitempty"`
}

func (rqf RequestQueryFields) Parse() url.Values {
	u := url.Values{}
	if rqf.CompanyKey != nil {
		u.Set("companyKey", *rqf.CompanyKey)
	}
	if rqf.OrganizationKey != nil {
		u.Set("organizationKey", string(*rqf.OrganizationKey))
	}
	if rqf.Expand != nil {
		u.Set("expand", *rqf.Expand)
	}
	if rqf.Filter != nil {
		u.Set("filter", *rqf.Filter)
	}
	if rqf.Orderby != nil {
		u.Set("orderby", *rqf.Orderby)
	}
	if rqf.Select != nil {
		u.Set("select", *rqf.Select)
	}
	if rqf.DailyConfigurationKey != nil {
		u.Set("dailyConfigurationKey", *rqf.DailyConfigurationKey)
	}
	if rqf.DateFormat != nil {
		u.Set("dateFormat", *rqf.DateFormat)
	}
	if rqf.EndDate != nil {
		u.Set("endDate", (*rqf.EndDate).Format(time.RFC3339))
	}
	if rqf.IgnoreHoliday != nil {
		u.Set("ignoreHoliday", strconv.FormatBool(*rqf.IgnoreHoliday))
	}
	if rqf.IgnoreWeekend != nil {
		u.Set("ignoreWeekend", strconv.FormatBool(*rqf.IgnoreWeekend))
	}
	if rqf.IncludeDeletes != nil {
		u.Set("includeDeletes", strconv.FormatBool(*rqf.IncludeDeletes))
	}
	if rqf.IncludeRemoved != nil {
		u.Set("includeRemoved", strconv.FormatBool(*rqf.IncludeRemoved))
	}
	if rqf.Includes != nil {
		u.Set("includes", *rqf.Includes)
	}
	if rqf.IsUniversallyLocalDates != nil {
		u.Set("IsUniversallyLocalDates", strconv.FormatBool(*rqf.IsUniversallyLocalDates))
	}
	if rqf.MaxResults != nil {
		u.Set("maxResults", fmt.Sprint(*rqf.MaxResults))
	}
	if rqf.PageToken != nil {
		u.Set("pageToken", *rqf.PageToken)
	}
	if rqf.RangeEndDate != nil {
		u.Set("rangeEndDate", (*rqf.RangeEndDate).Format(time.RFC3339))
	}
	if rqf.RangeStartDate != nil {
		u.Set("rangeStartDate", (*rqf.RangeStartDate).Format(time.RFC3339))
	}
	if rqf.ScheduleEndDate != nil {
		u.Set("scheduleEndDate", (*rqf.ScheduleEndDate).Format(time.RFC3339))
	}
	if rqf.ScheduleStartDate != nil {
		u.Set("scheduleStartDate", (*rqf.ScheduleStartDate).Format(time.RFC3339))
	}
	if rqf.SinceModifiedTimestamp != nil {
		u.Set("sinceModifiedTimestamp", (*rqf.SinceModifiedTimestamp).Format(time.RFC3339))
	}
	if rqf.StartDate != nil {
		u.Set("startDate", (*rqf.StartDate).Format(time.RFC3339))
	}
	if rqf.SyncToken != nil {
		u.Set("syncToken", (*rqf.SyncToken))
	}
	return u
}

// since we're trying pointer-members, it's better to pass
// literal's as needed rather than create temporary variables
// that we might inadvertently tamper with
func stringPointer(s string) *string     { return &s }
func intPointer(i int) *int              { return &i }
func boolPointer(b bool) *bool           { return &b }
func timePointer(t time.Time) *time.Time { return &t }
func floatPointer(f float64) *float64    { return &f }

func stringFromPointer(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
func intFromPointer(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}
func boolFromPointer(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}
func timeFromPointer(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}
func floatFromPointer(f *float64) float64 {
	if f != nil {
		return *f
	}
	return float64(0)
}

func (rqf *RequestQueryFields) GetCompanyKey() string {
	return stringFromPointer(rqf.CompanyKey)
}

func (rqf *RequestQueryFields) SetCompanyKey(s string) {
	rqf.CompanyKey = stringPointer(s)
}

func (rqf *RequestQueryFields) SetExpand(s string)  { rqf.Expand = stringPointer(s) }
func (rqf *RequestQueryFields) GetExpand() string   { return stringFromPointer(rqf.Expand) }
func (rqf *RequestQueryFields) SetFilter(s string)  { rqf.Filter = stringPointer(s) }
func (rqf *RequestQueryFields) GetFilter() string   { return stringFromPointer(rqf.Filter) }
func (rqf *RequestQueryFields) SetOrderby(s string) { rqf.Orderby = stringPointer(s) }
func (rqf *RequestQueryFields) GetOrderby() string  { return stringFromPointer(rqf.Orderby) }
func (rqf *RequestQueryFields) SetSelect(s string)  { rqf.Select = stringPointer(s) }
func (rqf *RequestQueryFields) GetSelect() string   { return stringFromPointer(rqf.Select) }
func (rqf *RequestQueryFields) SetDailyConfigurationKey(s string) {
	rqf.DailyConfigurationKey = stringPointer(s)
}
func (rqf *RequestQueryFields) GetDailyConfigurationKey() string {
	return stringFromPointer(rqf.DailyConfigurationKey)
}
func (rqf *RequestQueryFields) SetDateFormat(s string) { rqf.DateFormat = stringPointer(s) }
func (rqf *RequestQueryFields) GetDateFormat() string  { return stringFromPointer(rqf.DateFormat) }
func (rqf *RequestQueryFields) SetIncludes(s string)   { rqf.Includes = stringPointer(s) }
func (rqf *RequestQueryFields) GetIncludes() string    { return stringFromPointer(rqf.Includes) }
func (rqf *RequestQueryFields) SetPageToken(s string)  { rqf.PageToken = stringPointer(s) }
func (rqf *RequestQueryFields) GetPageToken() string   { return stringFromPointer(rqf.PageToken) }
func (rqf *RequestQueryFields) SetSyncToken(s string)  { rqf.SyncToken = stringPointer(s) }
func (rqf *RequestQueryFields) GetSyncToken() string   { return stringFromPointer(rqf.SyncToken) }

// bool
func (rqf *RequestQueryFields) SetIgnoreHoliday(b bool)  { rqf.IgnoreHoliday = boolPointer(b) }
func (rqf *RequestQueryFields) GetIgnoreHoliday() bool   { return boolFromPointer(rqf.IgnoreHoliday) }
func (rqf *RequestQueryFields) SetIgnoreWeekend(b bool)  { rqf.IgnoreWeekend = boolPointer(b) }
func (rqf *RequestQueryFields) GetIgnoreWeekend() bool   { return boolFromPointer(rqf.IgnoreWeekend) }
func (rqf *RequestQueryFields) SetIncludeDeletes(b bool) { rqf.IncludeDeletes = boolPointer(b) }
func (rqf *RequestQueryFields) GetIncludeDeletes() bool  { return boolFromPointer(rqf.IncludeDeletes) }
func (rqf *RequestQueryFields) SetIncludeRemoved(b bool) { rqf.IncludeRemoved = boolPointer(b) }
func (rqf *RequestQueryFields) GetIncludeRemoved() bool  { return boolFromPointer(rqf.IncludeRemoved) }
func (rqf *RequestQueryFields) SetIsUniversallyLocalDates(b bool) {
	rqf.IsUniversallyLocalDates = boolPointer(b)
}
func (rqf *RequestQueryFields) GetIsUniversallyLocalDates() bool {
	return boolFromPointer(rqf.IsUniversallyLocalDates)
}

// int
func (rqf *RequestQueryFields) SetMaxResults(i int)      { rqf.MaxResults = intPointer(i) }
func (rqf *RequestQueryFields) GetMaxResults() int       { return intFromPointer(rqf.MaxResults) }
func (rqf *RequestQueryFields) SetOrganizationKey(i int) { rqf.OrganizationKey = intPointer(i) }
func (rqf *RequestQueryFields) GetOrganizationKey() int  { return intFromPointer(rqf.OrganizationKey) }

// time.Time
func (rqf *RequestQueryFields) SetEndDate(t time.Time)        { rqf.EndDate = timePointer(t) }
func (rqf *RequestQueryFields) GetEndDate() time.Time         { return timeFromPointer(rqf.EndDate) }
func (rqf *RequestQueryFields) SetRangeEndDate(t time.Time)   { rqf.RangeEndDate = timePointer(t) }
func (rqf *RequestQueryFields) GetRangeEndDate() time.Time    { return timeFromPointer(rqf.RangeEndDate) }
func (rqf *RequestQueryFields) SetRangeStartDate(t time.Time) { rqf.RangeStartDate = timePointer(t) }
func (rqf *RequestQueryFields) GetRangeStartDate() time.Time {
	return timeFromPointer(rqf.RangeStartDate)
}
func (rqf *RequestQueryFields) SetScheduleEndDate(t time.Time) { rqf.ScheduleEndDate = timePointer(t) }
func (rqf *RequestQueryFields) GetScheduleEndDate() time.Time {
	return timeFromPointer(rqf.ScheduleEndDate)
}
func (rqf *RequestQueryFields) SetScheduleStartDate(t time.Time) {
	rqf.ScheduleStartDate = timePointer(t)
}
func (rqf *RequestQueryFields) GetScheduleStartDate() time.Time {
	return timeFromPointer(rqf.ScheduleStartDate)
}
func (rqf *RequestQueryFields) SetSinceModifiedTimestamp(t time.Time) {
	rqf.SinceModifiedTimestamp = timePointer(t)
}
func (rqf *RequestQueryFields) GetSinceModifiedTimestamp() time.Time {
	return timeFromPointer(rqf.SinceModifiedTimestamp)
}
func (rqf *RequestQueryFields) SetStartDate(t time.Time) { rqf.StartDate = timePointer(t) }
func (rqf *RequestQueryFields) GetStartDate() time.Time  { return timeFromPointer(rqf.StartDate) }

// func QueryInt(i int) string        { return url.QueryEscape(fmt.Sprint(i)) }
// func QueryBool(b bool) string      { return url.QueryEscape(fmt.Sprint(b)) }
// func QueryTime(t time.Time) string { return url.QueryEscape(t.Format(time.RFC3339)) }
// func QueryFloat(f float64) string  { return url.QueryEscape(fmt.Sprint(f)) }
// func QueryString(s string) string  { return url.QueryEscape(s) }

// type QueryParameter string

// const (
// 	// string
// 	CompanyKey            QueryParameter = "companyKey"
// 	Expand                QueryParameter = "$expand"
// 	Filter                QueryParameter = "$filter"
// 	Orderby               QueryParameter = "$orderby"
// 	Select                QueryParameter = "$select"
// 	DailyConfigurationKey QueryParameter = "dailyConfigurationKey"
// 	DateFormat            QueryParameter = "dateFormat"
// 	Includes              QueryParameter = "includes"
// 	PageToken             QueryParameter = "pageToken"
// 	SyncToken             QueryParameter = "syncToken"
// 	// bool
// 	IgnoreHoliday           QueryParameter = "ignoreHoliday"
// 	IgnoreWeekend           QueryParameter = "ignoreWeekend"
// 	IncludeDeletes          QueryParameter = "includeDeletes"
// 	IncludeRemoved          QueryParameter = "includeRemoved"
// 	IsUniversallyLocalDates QueryParameter = "IsUniversallyLocalDates"
// 	// int
// 	MaxResults      QueryParameter = "maxResults"
// 	OrganizationKey QueryParameter = "organizationKey"
// 	// time.Time
// 	EndDate                QueryParameter = "endDate"
// 	RangeEndDate           QueryParameter = "rangeEndDate"
// 	RangeStartDate         QueryParameter = "rangeStartDate"
// 	ScheduleEndDate        QueryParameter = "scheduleEndDate"
// 	ScheduleStartDate      QueryParameter = "scheduleStartDate"
// 	SinceModifiedTimestamp QueryParameter = "sinceModifiedTimestamp"
// 	StartDate              QueryParameter = "startDate"
// )

// func ParseQueryParameter(q QueryParameter, v interface{}) (string, string) {

// 	qm := map[QueryParameter]string{}
// 	fmt.Println(qm)
// 	return "", ""
// }
