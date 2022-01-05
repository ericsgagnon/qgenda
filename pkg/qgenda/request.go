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
	r.SetDateFormat("yyyy-MM-ddTHH:mm:ssZ")
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
