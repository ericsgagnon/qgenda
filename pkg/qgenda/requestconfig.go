package qgenda

import (
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

// RequestConfig is an experiment in using pointer members
// for optional fields
type RequestConfig struct {
	Path                    *string    `yaml:"-" query:"-" url:"-"`
	AllowedFields           []string   `yaml:"-" query:"-" url:"-"`
	CompanyKey              *string    `yaml:"companyKey,omitempty" query:"companyKey,omitempty" url:"companyKey,omitempty"`
	OrganizationKey         *int       `yaml:"organizationKey,omitempty" query:"organizationKey,omitempty" url:"organizationKey,omitempty"`
	Expand                  *string    `yaml:"$expand,omitempty" query:"$expand,omitempty" url:"$expand,omitempty"`
	Filter                  *string    `yaml:"$filter,omitempty" query:"$filter,omitempty" url:"$filter,omitempty"`
	Orderby                 *string    `yaml:"$orderby,omitempty" query:"$orderby,omitempty" url:"$orderby,omitempty"`
	Select                  *string    `yaml:"$select,omitempty" query:"$select,omitempty" url:"$select,omitempty"`
	DailyConfigurationKey   *string    `yaml:"dailyConfigurationKey,omitempty" query:"dailyConfigurationKey,omitempty" url:"dailyConfigurationKey,omitempty"`
	DateFormat              *string    `yaml:"dateFormat,omitempty" query:"dateFormat,omitempty" url:"dateFormat,omitempty"`
	EndDate                 *time.Time `yaml:"endDate,omitempty" query:"endDate,omitempty" url:"endDate,omitempty" layout:"2006-01-02T15:04:05Z"` //layout:"01/02/2006"
	IgnoreHoliday           *bool      `yaml:"ignoreHoliday,omitempty" query:"ignoreHoliday,omitempty" url:"ignoreHoliday,omitempty"`
	IgnoreWeekend           *bool      `yaml:"ignoreWeekend,omitempty" query:"ignoreWeekend,omitempty" url:"ignoreWeekend,omitempty"`
	IncludeDeletes          *bool      `yaml:"includeDeletes,omitempty" query:"includeDeletes,omitempty" url:"includeDeletes,omitempty"`
	IncludeRemoved          *bool      `yaml:"includeRemoved,omitempty" query:"includeRemoved,omitempty" url:"includeRemoved,omitempty"`
	Includes                *string    `yaml:"includes,omitempty" query:"includes,omitempty" url:"includes,omitempty"`
	IsUniversallyLocalDates *bool      `yaml:"IsUniversallyLocalDates,omitempty" query:"IsUniversallyLocalDates,omitempty" url:"IsUniversallyLocalDates,omitempty"`
	MaxResults              *int       `yaml:"maxResults,omitempty" query:"maxResults,omitempty" url:"maxResults,omitempty"`
	PageToken               *string    `yaml:"pageToken,omitempty" query:"pageToken,omitempty" url:"pageToken,omitempty"`
	RangeEndDate            *time.Time `yaml:"rangeEndDate,omitempty" query:"rangeEndDate,omitempty" url:"rangeEndDate,omitempty" layout:"2006-01-02T15:04:05Z"`                //layout:"01/02/2006"
	RangeStartDate          *time.Time `yaml:"rangeStartDate,omitempty" query:"rangeStartDate,omitempty" url:"rangeStartDate,omitempty" layout:"2006-01-02T15:04:05Z"`          //layout:"01/02/2006"
	ScheduleEndDate         *time.Time `yaml:"scheduleEndDate,omitempty" query:"scheduleEndDate,omitempty" url:"scheduleEndDate,omitempty" layout:"2006-01-02T15:04:05Z"`       //layout:"01/02/2006"
	ScheduleStartDate       *time.Time `yaml:"scheduleStartDate,omitempty" query:"scheduleStartDate,omitempty" url:"scheduleStartDate,omitempty" layout:"2006-01-02T15:04:05Z"` //layout:"01/02/2006"
	SinceModifiedTimestamp  *time.Time `yaml:"sinceModifiedTimestamp,omitempty" query:"sinceModifiedTimestamp,omitempty" url:"sinceModifiedTimestamp,omitempty" layout:"2006-01-02T15:04:05Z"`
	StartDate               *time.Time `yaml:"startDate,omitempty" query:"startDate,omitempty" url:"startDate,omitempty"`
	SyncToken               *string    `yaml:"syncToken,omitempty" query:"syncToken,omitempty" url:"syncToken,omitempty"`
	NewDataOnly             bool       `yaml:"newDataOnly,omitempty" query:"-" url:"-"`
	UseCache                bool       `yaml:"useCache,omitempty" query:"-" url:"-"`
	CacheFilenames          []string   `yaml:"cacheFilenames" query:"-" url:"-"`
	BatchSize               int        `yaml:"batchSize" query:"-" url:"-"`
}

// NewRequestConfig is intended for each data type to customize request configurations
// by standardizing how path and allowed query fields are handled
func NewRequestConfig(requestPath string, allowedFields []string) *RequestConfig {
	rc := RequestConfig{
		Path:          &requestPath,
		AllowedFields: allowedFields,
	}
	rc.FilterFields()
	return &rc
}

// Merge overwrites the receiver's fields with values from the requestconfig passed to it
// It is a noop if no requestconfig is passed
func (rc *RequestConfig) Merge(rcfg *RequestConfig) *RequestConfig {
	afm := map[string]bool{}
	if rc.AllowedFields != nil {
		for _, v := range rc.AllowedFields {
			afm[v] = true
		}
	}
	if rcfg != nil {
		if _, ok := afm["CompanyKey"]; ok && rcfg.CompanyKey != nil {
			rc.SetCompanyKey(rcfg.GetCompanyKey())
		}
		if _, ok := afm["OrganizationKey"]; ok && rcfg.OrganizationKey != nil {
			rc.SetOrganizationKey(rcfg.GetOrganizationKey())
		}
		if _, ok := afm["Expand"]; ok && rcfg.Expand != nil {
			rc.SetExpand(rcfg.GetExpand())
		}
		if _, ok := afm["Filter"]; ok && rcfg.Filter != nil {
			rc.SetFilter(rcfg.GetFilter())
		}
		if _, ok := afm["Orderby"]; ok && rcfg.Orderby != nil {
			rc.SetOrderby(rcfg.GetOrderby())
		}
		if _, ok := afm["Select"]; ok && rcfg.Select != nil {
			rc.SetSelect(rcfg.GetSelect())
		}
		if _, ok := afm["DailyConfigurationKey"]; ok && rcfg.DailyConfigurationKey != nil {
			rc.SetDailyConfigurationKey(rcfg.GetDailyConfigurationKey())
		}
		if _, ok := afm["DateFormat"]; ok {
			if rcfg.DateFormat != nil {
				rc.SetDateFormat(rcfg.GetDateFormat())
			} else {
				rc.SetDateFormat("yyyy-MM-ddTHH:mm:ssZ")
				// rc.SetDateFormat("MM/dd/yyyy")
			}
		}
		if _, ok := afm["IgnoreHoliday"]; ok {
			if rcfg.IgnoreHoliday != nil {
				rc.SetIgnoreHoliday(rcfg.GetIgnoreHoliday())
			} else {
				rc.SetIgnoreHoliday(false)
			}
		}
		if _, ok := afm["IgnoreWeekend"]; ok {
			if rcfg.IgnoreWeekend != nil {
				rc.SetIgnoreWeekend(rcfg.GetIgnoreWeekend())
			} else {
				rc.SetIgnoreWeekend(false)
			}
		}
		if _, ok := afm["IncludeDeletes"]; ok {
			if rcfg.IncludeDeletes != nil {
				rc.SetIncludeDeletes(rcfg.GetIncludeDeletes())
			} else {
				rc.SetIncludeDeletes(true)
			}
		}
		if _, ok := afm["IncludeRemoved"]; ok {
			if rcfg.IncludeRemoved != nil {
				rc.SetIncludeRemoved(rcfg.GetIncludeRemoved())
			} else {
				rc.SetIncludeRemoved(true)
			}
		}
		if _, ok := afm["Includes"]; ok && rcfg.Includes != nil {
			rc.SetIncludes(rcfg.GetIncludes())
		}
		if _, ok := afm["IsUniversallyLocalDates"]; ok && rcfg.IsUniversallyLocalDates != nil {
			rc.SetIsUniversallyLocalDates(rcfg.GetIsUniversallyLocalDates())
		}
		if _, ok := afm["MaxResults"]; ok && rcfg.MaxResults != nil {
			rc.SetMaxResults(rcfg.GetMaxResults())
		}
		if _, ok := afm["PageToken"]; ok && rcfg.PageToken != nil {
			rc.SetPageToken(rcfg.GetPageToken())
		}
		if _, ok := afm["RangeStartDate"]; ok {
			if rcfg.RangeStartDate != nil {
				rc.SetRangeStartDate(rcfg.GetRangeStartDate())
			} else {
				rc.SetRangeStartDate(time.Now().UTC().Add(time.Hour * 24 * 14 * -1))
			}
		}
		if _, ok := afm["RangeEndDate"]; ok {
			if rcfg.RangeEndDate != nil {
				rc.SetRangeEndDate(rcfg.GetRangeEndDate())
			} else {
				rc.SetRangeEndDate(rc.GetRangeStartDate().Add(time.Hour * 24 * 14))
			}
		}
		if _, ok := afm["ScheduleStartDate"]; ok {
			if rcfg.ScheduleStartDate != nil {
				rc.SetScheduleStartDate(rcfg.GetScheduleStartDate())
			} else {
				rc.SetScheduleStartDate(time.Now().UTC().Add(time.Hour * 24 * 14 * -1))
			}
		}
		if _, ok := afm["ScheduleEndDate"]; ok {
			if rcfg.ScheduleEndDate != nil {
				rc.SetScheduleEndDate(rcfg.GetScheduleEndDate())
			} else {
				rc.SetScheduleEndDate(rc.GetScheduleStartDate().Add(time.Hour * 24 * 14))
			}
		}
		if _, ok := afm["SinceModifiedTimestamp"]; ok {
			if rcfg.SinceModifiedTimestamp != nil {
				rc.SetSinceModifiedTimestamp(rcfg.GetSinceModifiedTimestamp())
			} else {
				rc.SetSinceModifiedTimestamp(time.Now().UTC().Add(time.Hour * 24 * 14 * -1))
			}
		}
		if _, ok := afm["StartDate"]; ok {
			if rcfg.StartDate != nil {
				rc.SetStartDate(rcfg.GetStartDate())
			} else {
				rc.SetStartDate(time.Now().UTC().Add(time.Hour * 24 * 14 * -1))
			}
		}
		if _, ok := afm["EndDate"]; ok {
			if rcfg.EndDate != nil {
				rc.SetEndDate(rcfg.GetEndDate())
			} else {
				rc.SetEndDate(rc.GetStartDate().Add(time.Hour * 24 * 14))
			}
		}

		if _, ok := afm["SyncToken"]; ok && rcfg.SyncToken != nil {
			rc.SetSyncToken(rcfg.GetSyncToken())
		}

	}

	return rc
}

func (rc *RequestConfig) Parse() url.Values {
	v, err := query.Values(rc)
	if err != nil {
		panic(err)
	}
	return v
}

func (rc *RequestConfig) ToQuery() url.Values {
	v, err := query.Values(rc)
	if err != nil {
		panic(err)
	}
	return v
}

// FilterFields sets all fields not in .AllowedFields to nil
// it does nothing if .AllowedFields is nil
func (rc *RequestConfig) FilterFields() *RequestConfig {
	if rc.AllowedFields == nil {
		return rc
	}
	afm := map[string]interface{}{}
	for _, v := range rc.AllowedFields {
		afm[v] = struct{}{}
	}

	return rc
}

// Setters and Getters - Using Set* and Get* to avoid conflicts with member names
func (rc *RequestConfig) GetCompanyKey() string             { return toValue(rc.CompanyKey) }
func (rc *RequestConfig) SetCompanyKey(s string)            { rc.CompanyKey = toPointer(s) }
func (rc *RequestConfig) SetExpand(s string)                { rc.Expand = toPointer(s) }
func (rc *RequestConfig) GetExpand() string                 { return toValue(rc.Expand) }
func (rc *RequestConfig) SetFilter(s string)                { rc.Filter = toPointer(s) }
func (rc *RequestConfig) GetFilter() string                 { return toValue(rc.Filter) }
func (rc *RequestConfig) SetOrderby(s string)               { rc.Orderby = toPointer(s) }
func (rc *RequestConfig) GetOrderby() string                { return toValue(rc.Orderby) }
func (rc *RequestConfig) SetSelect(s string)                { rc.Select = toPointer(s) }
func (rc *RequestConfig) GetSelect() string                 { return toValue(rc.Select) }
func (rc *RequestConfig) SetDailyConfigurationKey(s string) { rc.DailyConfigurationKey = toPointer(s) }
func (rc *RequestConfig) GetDailyConfigurationKey() string  { return toValue(rc.DailyConfigurationKey) }
func (rc *RequestConfig) SetDateFormat(s string)            { rc.DateFormat = toPointer(s) }
func (rc *RequestConfig) GetDateFormat() string             { return toValue(rc.DateFormat) }
func (rc *RequestConfig) SetIncludes(s string)              { rc.Includes = toPointer(s) }
func (rc *RequestConfig) GetIncludes() string               { return toValue(rc.Includes) }
func (rc *RequestConfig) SetPageToken(s string)             { rc.PageToken = toPointer(s) }
func (rc *RequestConfig) GetPageToken() string              { return toValue(rc.PageToken) }
func (rc *RequestConfig) SetSyncToken(s string)             { rc.SyncToken = toPointer(s) }
func (rc *RequestConfig) GetSyncToken() string              { return toValue(rc.SyncToken) }

// bool
func (rc *RequestConfig) SetIgnoreHoliday(b bool)  { rc.IgnoreHoliday = toPointer(b) }
func (rc *RequestConfig) GetIgnoreHoliday() bool   { return toValue(rc.IgnoreHoliday) }
func (rc *RequestConfig) SetIgnoreWeekend(b bool)  { rc.IgnoreWeekend = toPointer(b) }
func (rc *RequestConfig) GetIgnoreWeekend() bool   { return toValue(rc.IgnoreWeekend) }
func (rc *RequestConfig) SetIncludeDeletes(b bool) { rc.IncludeDeletes = toPointer(b) }
func (rc *RequestConfig) GetIncludeDeletes() bool  { return toValue(rc.IncludeDeletes) }
func (rc *RequestConfig) SetIncludeRemoved(b bool) { rc.IncludeRemoved = toPointer(b) }
func (rc *RequestConfig) GetIncludeRemoved() bool  { return toValue(rc.IncludeRemoved) }
func (rc *RequestConfig) SetIsUniversallyLocalDates(b bool) {
	rc.IsUniversallyLocalDates = toPointer(b)
}
func (rc *RequestConfig) GetIsUniversallyLocalDates() bool {
	return toValue(rc.IsUniversallyLocalDates)
}

// int
func (rc *RequestConfig) SetMaxResults(i int)      { rc.MaxResults = toPointer(i) }
func (rc *RequestConfig) GetMaxResults() int       { return toValue(rc.MaxResults) }
func (rc *RequestConfig) SetOrganizationKey(i int) { rc.OrganizationKey = toPointer(i) }
func (rc *RequestConfig) GetOrganizationKey() int  { return toValue(rc.OrganizationKey) }

// time.Time
func (rc *RequestConfig) SetEndDate(t time.Time)           { rc.EndDate = toPointer(t) }
func (rc *RequestConfig) GetEndDate() time.Time            { return toValue(rc.EndDate) }
func (rc *RequestConfig) SetRangeEndDate(t time.Time)      { rc.RangeEndDate = toPointer(t) }
func (rc *RequestConfig) GetRangeEndDate() time.Time       { return toValue(rc.RangeEndDate) }
func (rc *RequestConfig) SetRangeStartDate(t time.Time)    { rc.RangeStartDate = toPointer(t) }
func (rc *RequestConfig) GetRangeStartDate() time.Time     { return toValue(rc.RangeStartDate) }
func (rc *RequestConfig) SetScheduleEndDate(t time.Time)   { rc.ScheduleEndDate = toPointer(t) }
func (rc *RequestConfig) GetScheduleEndDate() time.Time    { return toValue(rc.ScheduleEndDate) }
func (rc *RequestConfig) SetScheduleStartDate(t time.Time) { rc.ScheduleStartDate = toPointer(t) }
func (rc *RequestConfig) GetScheduleStartDate() time.Time  { return toValue(rc.ScheduleStartDate) }
func (rc *RequestConfig) SetSinceModifiedTimestamp(t time.Time) {
	rc.SinceModifiedTimestamp = toPointer(t)
}
func (rc *RequestConfig) GetSinceModifiedTimestamp() time.Time {
	return toValue(rc.SinceModifiedTimestamp)
}
func (rc *RequestConfig) SetStartDate(t time.Time) { rc.StartDate = toPointer(t) }
func (rc *RequestConfig) GetStartDate() time.Time  { return toValue(rc.StartDate) }

func (r *Request) Parse() (*http.Request, error) {
	return http.NewRequest(r.Method, r.RequestConfig.Parse().Encode(), nil)

}

// func (rc *RequestConfig) ToRequest() *Request {
// 	return NewRequest(rc)
// }
