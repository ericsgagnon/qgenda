package qgenda

type TagCompany struct {
	CompanyName   *string       `json:"CompanyName"`
	CompanyKey    *string       `json:"CompanyKey"`
	TagCategories []TagCategory `json:"Tags"`
}

type TagCategory struct {
	CategoryKey                    *int64  `json:"CategoryKey"`
	CategoryName                   *string `json:"CategoryName"`
	CategoryDateCreated            *Time   `json:"CategoryDateCreated,omitempty"`
	CategoryDateLastModified       *Time   `json:"CategoryDateLastModified,omitempty"`
	IsAvailableForCreditAllocation *bool   `json:"IsAvailableForCreditAllocation,omitempty"`
	IsAvailableForLDailySum        *bool   `json:"IsAvailableForLDailySum,omitempty"`
	IsAvailableForHoliday          *bool   `json:"IsAvailableForHoliday,omitempty"`
	IsAvailableForLocation         *bool   `json:"IsAvailableForLocation,omitempty"`
	IsAvailableForProfile          *bool   `json:"IsAvailableForProfile,omitempty"`
	IsAvailableForSeries           *bool   `json:"IsAvailableForSeries,omitempty"`
	IsAvailableForStaff            *bool   `json:"IsAvailableForStaff,omitempty"`
	IsAvailableForStaffLocation    *bool   `json:"IsAvailableForStaffLocation,omitempty"`
	IsAvailableForStaffTarget      *bool   `json:"IsAvailableForStaffTarget,omitempty"`
	IsAvailableForRequestLimit     *bool   `json:"IsAvailableForRequestLimit,omitempty"`
	IsAvailableForTask             *bool   `json:"IsAvailableForTask,omitempty"`
	IsTTCMCategory                 *bool   `json:"IsTTCMCategory,omitempty"`
	IsSingleTaggingOnly            *bool   `json:"IsSingleTaggingOnly,omitempty"`
	IsPermissionCategory           *bool   `json:"IsPermissionCategory,omitempty"`
	IsUsedForStats                 *bool   `json:"IsUsedForStats,omitempty"`
	CategoryBackgroundColor        *string `json:"CategoryBackgroundColor,omitempty"`
	CategoryTextColor              *string `json:"CategoryTextColor,omitempty"`
	Tags                           []Tag   `json:"Tags"`
}

type Tag struct { 
	Key               *int64  `json:"Key"`
	Name              *string `json:"Name"`
	DateCreated       *Time   `json:"DateCreated,omitempty"`
	DateLastModified  *Time   `json:"DateLastModified,omitempty"`
	BackgroundColor   *string `json:"BackgroundColor,omitempty"`
	TextColor         *string `json:"TextColor,omitempty"`
	EffectiveFromDate *Date   `json:"EffectiveFromDate,omitempty"` // only applies to links to other entities
	EffectiveToDate   *Date   `json:"EffectiveToDate,omitempty"`   // only applies to links to other entities
}

type TagDetail struct {
	Key                            *int64  `json:"Key"`
	Name                           *string `json:"Name"`
	DateCreated                    *Time   `json:"DateCreated,omitempty"`
	DateLastModified               *Time   `json:"DateLastModified,omitempty"`
	BackgroundColor                *string `json:"BackgroundColor,omitempty"`
	TextColor                      *string `json:"TextColor,omitempty"`
	CompanyName                    *string `json:"CompanyName"`
	CompanyKey                     *string `json:"CompanyKey"`
	CategoryKey                    *int64  `json:"CategoryKey"`
	CategoryName                   *string `json:"CategoryName"`
	CategoryDateCreated            *Time   `json:"CategoryDateCreated,omitempty"`
	CategoryDateLastModified       *Time   `json:"CategoryDateLastModified,omitempty"`
	IsAvailableForCreditAllocation *bool   `json:"IsAvailableForCreditAllocation,omitempty"`
	IsAvailableForLDailySum        *bool   `json:"IsAvailableForLDailySum,omitempty"`
	IsAvailableForHoliday          *bool   `json:"IsAvailableForHoliday,omitempty"`
	IsAvailableForLocation         *bool   `json:"IsAvailableForLocation,omitempty"`
	IsAvailableForProfile          *bool   `json:"IsAvailableForProfile,omitempty"`
	IsAvailableForSeries           *bool   `json:"IsAvailableForSeries,omitempty"`
	IsAvailableForStaff            *bool   `json:"IsAvailableForStaff,omitempty"`
	IsAvailableForStaffLocation    *bool   `json:"IsAvailableForStaffLocation,omitempty"`
	IsAvailableForStaffTarget      *bool   `json:"IsAvailableForStaffTarget,omitempty"`
	IsAvailableForRequestLimit     *bool   `json:"IsAvailableForRequestLimit,omitempty"`
	IsAvailableForTask             *bool   `json:"IsAvailableForTask,omitempty"`
	IsTTCMCategory                 *bool   `json:"IsTTCMCategory,omitempty"`
	IsSingleTaggingOnly            *bool   `json:"IsSingleTaggingOnly,omitempty"`
	IsPermissionCategory           *bool   `json:"IsPermissionCategory,omitempty"`
	IsUsedForStats                 *bool   `json:"IsUsedForStats,omitempty"`
	CategoryBackgroundColor        *string `json:"CategoryBackgroundColor,omitempty"`
	CategoryTextColor              *string `json:"CategoryTextColor,omitempty"`
}

// func NewTagRequest(rqf *RequestQueryFields) *Request {
// 	rPath := "tags"
// 	allowableFields := []string{
// 		"CompanyKey",
// 		"OrganizationKey",
// 		"Expand",
// 		"Filter",
// 		"Orderby",
// 		"Select",
// 		"DailyConfigurationKey",
// 		"DateFormat",
// 		"EndDate",
// 		"IgnoreHoliday",
// 		"IgnoreWeekend",
// 		"IncludeDeletes",
// 		"IncludeRemoved",
// 		"Includes",
// 		"IsUniversallyLocalDates",
// 		"MaxResults",
// 		"PageToken",
// 		"RangeEndDate",
// 		"RangeStartDate",
// 		"ScheduleEndDate",
// 		"ScheduleStartDate",
// 		"SinceModifiedTimestamp",
// 		"StartDate",
// 		"SyncToken",
// 	}

// 	r := NewRequestWithQueryField(rPath, allowableFields, rqf)
// 	// r.SetIncludes("StaffTags,TaskTags,LocationTags")
// 	// r.SetStartDate(time.Now().AddDate(0, 0, -14).UTC())
// 	// r.SetEndDate(time.Now().UTC())

// 	return r
// }

func NewTagRequest(rqf *RequestQueryFields) *Request {
	requestPath := "tags"
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

func (p *TagCompany) Process() error {
	ProcessStruct(p)
	for i, _ := range p.TagCategories {
		(&p.TagCategories[i]).Process()
	}
	return nil
}

func (p *TagCategory) Process() error {
	ProcessStruct(p)
	for i, _ := range p.Tags {
		(&p.Tags[i]).Process()
	}
	return nil
}

func (p *Tag) Process() error {
	return ProcessStruct(p)
}
