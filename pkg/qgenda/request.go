package qgenda

import (
	"bytes"
	"net/http"
	"net/url"
	"path"
	"time"
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
	RequestConfig
}

// NewRequest returns a Request with only common members
// it is expected to a base for other 'request' functions
func NewRequest() *Request {
	r := Request{
		Scheme:        "https",
		Header:        http.Header{},
		Host:          "api.qgenda.com",
		Path:          "v2",
		RequestConfig: RequestConfig{},
	}
	// r.SetDateFormat("yyyy-MM-ddTHH:mm:ssZ")
	return &r
}

func (r *Request) ToURL() *url.URL {
	return &url.URL{
		Scheme:   r.Scheme,
		Host:     r.Host,
		Path:     r.Path,
		RawQuery: r.RequestConfig.Parse().Encode(),
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

func NewRequestWithQueryField(requestPath string, allowableQueryFields []string, rqf *RequestConfig) *Request {
	aqfMap := map[string]interface{}{}
	for _, v := range allowableQueryFields {
		aqfMap[v] = struct{}{}
	}

	r := NewRequest()
	r.Path = path.Join(r.Path, requestPath)
	// if rqf.Select != nil {
	// 	r.SetSelect(rqf.GetSelect())
	// }
	qf := RequestConfig{}
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
	r.RequestConfig = qf
	return r
}

// these return requests that are specific to each dataset we're interested in
// they can be moved to their appropriate files as each are created

func NewRotationsRequest(rqf *RequestConfig) *Request {
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

func NewRequestLimitRequest(rqf *RequestConfig) *Request {
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
func NewTaskRequest(rqf *RequestConfig) *Request {
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
func NewTaskLocationRequest(rqf *RequestConfig) *Request {
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
func NewDailyPatientEncounterRequest(rqf *RequestConfig) *Request {
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
func NewDailyDailyConfigurationRequest(rqf *RequestConfig) *Request {
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
func NewDailyDailyConfigurationDailyConfigurationKeyRequest(rqf *RequestConfig) *Request {
	requestPath := "daily/dailyconfiguration/:dailyConfigurationKey"
	queryFields := []string{
		"CompanyKey",
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}
func NewDailyRoomRequest(rqf *RequestConfig) *Request {
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
func NewDailyCaseRequest(rqf *RequestConfig) *Request {
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
func NewPayRateRequest(rqf *RequestConfig) *Request {
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
func NewTimeEventRequest(rqf *RequestConfig) *Request {
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

func NewOrganizationRequest(rqf *RequestConfig) *Request {
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

func NewCompanyRequest(rqf *RequestConfig) *Request {
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
func NewStaffTargetRequest(rqf *RequestConfig) *Request {
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
func NewProfileRequest(rqf *RequestConfig) *Request {
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
func NewUserRequest(rqf *RequestConfig) *Request {
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
