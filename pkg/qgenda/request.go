package qgenda

import (
	"bytes"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/google/go-querystring/query"
)

// Request holds all non-client info for each request type.
// It is used by every api call and is customized by
// initialization functions for resource.
type Request struct {
	Scheme string
	Method string
	Header http.Header
	Host   string
	Path   string
	Body   []byte
	RequestQueryFields
}

// NewRequest returns a Request with only common members
// it is expected to a base for other 'request' functions
func NewRequest() *Request {
	r := Request{
		Scheme:             "https",
		Header:             http.Header{},
		Host:               "api.qgenda.com",
		Path:               "v2",
		RequestQueryFields: RequestQueryFields{},
	}
	// r.SetDateFormat("yyyy-MM-ddTHH:mm:ssZ")
	return &r
}

func (r *Request) ToURL() *url.URL {
	return &url.URL{
		Scheme:   r.Scheme,
		Host:     r.Host,
		Path:     r.Path,
		RawQuery: r.RequestQueryFields.Parse().Encode(),
	}
}

func (r *Request) Encode() string {
	return r.ToURL().String()
}

func (r *Request) ToHTTPRequest() *http.Request {
	req, err := http.NewRequest(r.Method, r.Encode(), bytes.NewReader(r.Body))
	if err != nil {
		panic(err)
	}
	for k, v := range r.Header {
		req.Header[k] = v
	}
	return req
}

// AppendPath is a convenience func to append sub-paths to a base path
func (r *Request) AppendPath(p string) {
	r.Path = path.Join(r.Path, p)
}

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
	EndDate                 *time.Time `query:"endDate,omitempty" url:"endDate,omitempty" layout:"2006-01-02T15:04:05Z"` //layout:"01/02/2006"
	IgnoreHoliday           *bool      `query:"ignoreHoliday,omitempty" url:"ignoreHoliday,omitempty"`
	IgnoreWeekend           *bool      `query:"ignoreWeekend,omitempty" url:"ignoreWeekend,omitempty"`
	IncludeDeletes          *bool      `query:"includeDeletes,omitempty" url:"includeDeletes,omitempty"`
	IncludeRemoved          *bool      `query:"includeRemoved,omitempty" url:"includeRemoved,omitempty"`
	Includes                *string    `query:"includes,omitempty" url:"includes,omitempty"`
	IsUniversallyLocalDates *bool      `query:"IsUniversallyLocalDates,omitempty" url:"IsUniversallyLocalDates,omitempty"`
	MaxResults              *int       `query:"maxResults,omitempty" url:"maxResults,omitempty"`
	PageToken               *string    `query:"pageToken,omitempty" url:"pageToken,omitempty"`
	RangeEndDate            *time.Time `query:"rangeEndDate,omitempty" url:"rangeEndDate,omitempty" layout:"2006-01-02T15:04:05Z"`           //layout:"01/02/2006"
	RangeStartDate          *time.Time `query:"rangeStartDate,omitempty" url:"rangeStartDate,omitempty" layout:"2006-01-02T15:04:05Z"`       //layout:"01/02/2006"
	ScheduleEndDate         *time.Time `query:"scheduleEndDate,omitempty" url:"scheduleEndDate,omitempty" layout:"2006-01-02T15:04:05Z"`     //layout:"01/02/2006"
	ScheduleStartDate       *time.Time `query:"scheduleStartDate,omitempty" url:"scheduleStartDate,omitempty" layout:"2006-01-02T15:04:05Z"` //layout:"01/02/2006"
	SinceModifiedTimestamp  *time.Time `query:"sinceModifiedTimestamp,omitempty" url:"sinceModifiedTimestamp,omitempty" layout:"2006-01-02T15:04:05Z"`
	StartDate               *time.Time `query:"startDate,omitempty" url:"startDate,omitempty"`
	SyncToken               *string    `query:"syncToken,omitempty" url:"syncToken,omitempty"`
}

func (rqf *RequestQueryFields) Parse() url.Values {
	v, err := query.Values(rqf)
	if err != nil {
		panic(err)
	}
	return v
}

func (rqf *RequestQueryFields) ToQuery() url.Values {
	v, err := query.Values(rqf)
	if err != nil {
		panic(err)
	}
	return v
}

// Setters and Getters - Using Set* and Get* to avoid conflicts with member names
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

func (r *Request) Parse() (*http.Request, error) {
	return http.NewRequest(r.Method, r.RequestQueryFields.Parse().Encode(), nil)

}

func NewRequestWithQueryField(requestPath string, allowableQueryFields []string, rqf *RequestQueryFields) *Request {
	aqfMap := map[string]interface{}{}
	for _, v := range allowableQueryFields {
		aqfMap[v] = struct{}{}
	}

	r := NewRequest()
	r.Path = path.Join(r.Path, requestPath)
	// if rqf.Select != nil {
	// 	r.SetSelect(rqf.GetSelect())
	// }
	qf := RequestQueryFields{}
	if rqf != nil {
		if _, ok := aqfMap["CompanyKey"]; ok && rqf.CompanyKey != nil {
			qf.SetCompanyKey(rqf.GetCompanyKey())
		}
		if _, ok := aqfMap["OrganizationKey"]; ok && rqf.OrganizationKey != nil {
			qf.SetOrganizationKey(rqf.GetOrganizationKey())
		}
		if _, ok := aqfMap["Expand"]; ok && rqf.Expand != nil {
			qf.SetExpand(rqf.GetExpand())
		}
		if _, ok := aqfMap["Filter"]; ok && rqf.Filter != nil {
			qf.SetFilter(rqf.GetFilter())
		}
		if _, ok := aqfMap["Orderby"]; ok && rqf.Orderby != nil {
			qf.SetOrderby(rqf.GetOrderby())
		}
		if _, ok := aqfMap["Select"]; ok && rqf.Select != nil {
			qf.SetSelect(rqf.GetSelect())
		}
		if _, ok := aqfMap["DailyConfigurationKey"]; ok && rqf.DailyConfigurationKey != nil {
			qf.SetDailyConfigurationKey(rqf.GetDailyConfigurationKey())
		}
		if _, ok := aqfMap["DateFormat"]; ok {
			if rqf.DateFormat != nil {
				qf.SetDateFormat(rqf.GetDateFormat())
			} else {
				qf.SetDateFormat("yyyy-MM-ddTHH:mm:ssZ")
				// qf.SetDateFormat("MM/dd/yyyy")
			}
		}
		if _, ok := aqfMap["IgnoreHoliday"]; ok {
			if rqf.IgnoreHoliday != nil {
				qf.SetIgnoreHoliday(rqf.GetIgnoreHoliday())
			} else {
				qf.SetIgnoreHoliday(false)
			}
		}
		if _, ok := aqfMap["IgnoreWeekend"]; ok {
			if rqf.IgnoreWeekend != nil {
				qf.SetIgnoreWeekend(rqf.GetIgnoreWeekend())
			} else {
				qf.SetIgnoreWeekend(false)
			}
		}
		if _, ok := aqfMap["IncludeDeletes"]; ok {
			if rqf.IncludeDeletes != nil {
				qf.SetIncludeDeletes(rqf.GetIncludeDeletes())
			} else {
				qf.SetIncludeDeletes(true)
			}
		}
		if _, ok := aqfMap["IncludeRemoved"]; ok {
			if rqf.IncludeRemoved != nil {
				qf.SetIncludeRemoved(rqf.GetIncludeRemoved())
			} else {
				qf.SetIncludeRemoved(true)
			}
		}
		if _, ok := aqfMap["Includes"]; ok && rqf.Includes != nil {
			qf.SetIncludes(rqf.GetIncludes())
		}
		if _, ok := aqfMap["IsUniversallyLocalDates"]; ok && rqf.IsUniversallyLocalDates != nil {
			qf.SetIsUniversallyLocalDates(rqf.GetIsUniversallyLocalDates())
		}
		if _, ok := aqfMap["MaxResults"]; ok && rqf.MaxResults != nil {
			qf.SetMaxResults(rqf.GetMaxResults())
		}
		if _, ok := aqfMap["PageToken"]; ok && rqf.PageToken != nil {
			qf.SetPageToken(rqf.GetPageToken())
		}
		if _, ok := aqfMap["RangeStartDate"]; ok {
			if rqf.RangeStartDate != nil {
				qf.SetRangeStartDate(rqf.GetRangeStartDate())
			} else {
				qf.SetRangeStartDate(time.Now().UTC().Add(time.Hour * 24 * 14 * -1))
			}
		}
		if _, ok := aqfMap["RangeEndDate"]; ok {
			if rqf.RangeEndDate != nil {
				qf.SetRangeEndDate(rqf.GetRangeEndDate())
			} else {
				qf.SetRangeEndDate(qf.GetRangeStartDate().Add(time.Hour * 24 * 14))
			}
		}
		if _, ok := aqfMap["ScheduleStartDate"]; ok {
			if rqf.ScheduleStartDate != nil {
				qf.SetScheduleStartDate(rqf.GetScheduleStartDate())
			} else {
				qf.SetScheduleStartDate(time.Now().UTC().Add(time.Hour * 24 * 14 * -1))
			}
		}
		if _, ok := aqfMap["ScheduleEndDate"]; ok {
			if rqf.ScheduleEndDate != nil {
				qf.SetScheduleEndDate(rqf.GetScheduleEndDate())
			} else {
				qf.SetScheduleEndDate(qf.GetScheduleStartDate().Add(time.Hour * 24 * 14))
			}
		}
		if _, ok := aqfMap["SinceModifiedTimestamp"]; ok {
			if rqf.SinceModifiedTimestamp != nil {
				qf.SetSinceModifiedTimestamp(rqf.GetSinceModifiedTimestamp())
			} else {
				qf.SetSinceModifiedTimestamp(time.Now().UTC().Add(time.Hour * 24 * 14 * -1))
			}
		}
		if _, ok := aqfMap["StartDate"]; ok {
			if rqf.StartDate != nil {
				qf.SetStartDate(rqf.GetStartDate())
			} else {
				qf.SetStartDate(time.Now().UTC().Add(time.Hour * 24 * 14 * -1))
			}
		}
		if _, ok := aqfMap["EndDate"]; ok {
			if rqf.EndDate != nil {
				qf.SetEndDate(rqf.GetEndDate())
			} else {
				qf.SetEndDate(qf.GetStartDate().Add(time.Hour * 24 * 14))
			}
		}

		if _, ok := aqfMap["SyncToken"]; ok && rqf.SyncToken != nil {
			qf.SetSyncToken(rqf.GetSyncToken())
		}
	} else {
		return NewRequestWithQueryField(requestPath, allowableQueryFields, &qf)
	}
	r.RequestQueryFields = qf
	return r
}

func NewRotationsRequest(rqf *RequestQueryFields) *Request {
	requestPath := "schedule/rotations"
	queryFields := []string{
		"CompanyKey",
		"RangeStartDate",
		"RangeEndDate",
		"IgnoreHoliday",
		"IgnoreWeekend",
		"DateFormat",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}

func NewRequestLimitRequest(rqf *RequestQueryFields) *Request {
	requestPath := "requestlimit"
	queryFields := []string{
		"DateFormat",
		"CompanyKey",
		"StartDate",
		"EndDate",
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	if rqf != nil {
		if rqf.Includes == nil {
			rqf.SetIncludes("ShiftsCredit,StaffLimits")
		}
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewTaskRequest(rqf *RequestQueryFields) *Request {
	requestPath := "task"
	queryFields := []string{
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	if rqf != nil {
		if rqf.Includes == nil {
			rqf.SetIncludes("Profiles,Tags,TaskShifts")
		}
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewTaskLocationRequest(rqf *RequestQueryFields) *Request {
	requestPath := "task/:taskid/location"
	queryFields := []string{
		"CompanyKey",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewDailyPatientEncounterRequest(rqf *RequestQueryFields) *Request {
	requestPath := "daily/patientencounter"
	queryFields := []string{
		"CompanyKey",
		"DailyConfigurationKey",
		"StartDate",
		"EndDate",
		"DateFormat",
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	if rqf != nil {
		if rqf.Includes == nil {
			rqf.SetIncludes("StandardFields,PatientInformation")
		}
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewDailyDailyConfigurationRequest(rqf *RequestQueryFields) *Request {
	requestPath := "daily/dailyconfiguration"
	queryFields := []string{
		"CompanyKey",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewDailyDailyConfigurationDailyConfigurationKeyRequest(rqf *RequestQueryFields) *Request {
	requestPath := "daily/dailyconfiguration/:dailyConfigurationKey"
	queryFields := []string{
		"CompanyKey",
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewDailyRoomRequest(rqf *RequestQueryFields) *Request {
	requestPath := "daily/room"
	queryFields := []string{
		"CompanyKey",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewDailyCaseRequest(rqf *RequestQueryFields) *Request {
	requestPath := "dailycase"
	queryFields := []string{
		"CompanyKey",
		"StartDate",
		"EndDate",
		"DateFormat",
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	if rqf != nil {
		if rqf.Includes == nil {
			rqf.SetIncludes("Task,Supervisors,DirectProviders")
		}
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewPayRateRequest(rqf *RequestQueryFields) *Request {
	requestPath := "payrate"
	queryFields := []string{
		"CompanyKey",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewTimeEventRequest(rqf *RequestQueryFields) *Request {
	requestPath := "timeevent"
	queryFields := []string{
		"CompanyKey",
		"StartDate",
		"EndDate",
		"DateFormat",
		"IsUniversallyLocalDates",
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	if rqf != nil {
		if rqf.Includes == nil {
			rqf.SetIncludes("ScheduleEntry,Task,StaffMember")
		}
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}

func NewOrganizationRequest(rqf *RequestQueryFields) *Request {
	requestPath := "organization"
	queryFields := []string{
		"OrganizationKey",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}

func NewCompanyRequest(rqf *RequestQueryFields) *Request {
	requestPath := "company"
	queryFields := []string{
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	if rqf != nil {
		if rqf.Includes == nil {
			rqf.SetIncludes("Profiles,Organizations")
		}
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewStaffTargetRequest(rqf *RequestQueryFields) *Request {
	requestPath := "stafftarget"
	queryFields := []string{
		"Includes",
		"CompanyKey",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	if rqf != nil {
		if rqf.Includes == nil {
			rqf.SetIncludes("Staff,Profiles,Locations")
		}
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewProfileRequest(rqf *RequestQueryFields) *Request {
	requestPath := "profile"
	queryFields := []string{
		"CompanyKey",
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	if rqf != nil {
		if rqf.Includes == nil {
			rqf.SetIncludes("Staff,Tasks")
		}
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewUserRequest(rqf *RequestQueryFields) *Request {
	requestPath := "user/"
	queryFields := []string{
		"CompanyKey",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
