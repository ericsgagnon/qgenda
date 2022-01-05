package qgenda

import (
	"path"
	"time"
)

type ScheduleRequest struct {
	Request
}

func NewScheduleRequest(rqf *RequestQueryFields) *ScheduleRequest {

	r := NewRequest()
	r.Path = path.Join(r.Path, "schedule")
	r.SetIncludes("StaffTags,TaskTags,LocationTags")
	r.SetStartDate(time.Now().AddDate(0, 0, -14).UTC())
	r.SetEndDate(time.Now().UTC())
	if rqf != nil {
		if rqf.CompanyKey != nil {
			r.SetCompanyKey(rqf.GetCompanyKey())
		}
		if rqf.StartDate != nil {
			r.SetStartDate(rqf.GetStartDate())
		}
		if rqf.EndDate != nil {
			r.SetEndDate(rqf.GetEndDate())
		}
		if rqf.IncludeDeletes != nil {
			r.SetIncludeDeletes(rqf.GetIncludeDeletes())
		}
		if rqf.SinceModifiedTimestamp != nil {
			r.SetSinceModifiedTimestamp(rqf.GetSinceModifiedTimestamp())
		}
		if rqf.DateFormat != nil {
			r.SetDateFormat(rqf.GetDateFormat())
		}
		if rqf.Includes != nil {
			r.SetIncludes(rqf.GetIncludes())
		}
		if rqf.Select != nil {
			r.SetSelect(rqf.GetSelect())
		}
		if rqf.Filter != nil {
			r.SetFilter(rqf.GetFilter())
		}
		if rqf.Orderby != nil {
			r.SetOrderby(rqf.GetOrderby())
		}
		if rqf.Expand != nil {
			r.SetExpand(rqf.GetExpand())
		}
	}
	s := ScheduleRequest{}
	s.Request = *r
	return &s
}

type Schedule struct {
	ScheduleKey *string  `json:"ScheduleKey,omitempty"`
	CallRole    *string  `json:"CallRole,omitempty"`
	CompKey     *string  `json:"CompKey,omitempty"`
	Credit      *float64 `json:"Credit,omitempty"`
	// Date                   *time.Time `json:"Date,omitempty"`
	// StartDateUTC           *time.Time `json:"StartDateUTC,omitempty"`
	// EndDateUTC             *time.Time `json:"EndDateUTC,omitempty"`
	// EndDate                *time.Time `json:"EndDate,omitempty"`
	// EndTime                *time.Time `json:"EndTime,omitempty"`
	Date              *string `json:"Date,omitempty"`
	StartDateUTC      *string `json:"StartDateUTC,omitempty"`
	EndDateUTC        *string `json:"EndDateUTC,omitempty"`
	EndDate           *string `json:"EndDate,omitempty"`
	EndTime           *string `json:"EndTime,omitempty"`
	IsCred            *bool   `json:"IsCred,omitempty"`
	IsPublished       *bool   `json:"IsPublished,omitempty"`
	IsLocked          *bool   `json:"IsLocked,omitempty"`
	IsStruck          *bool   `json:"IsStruck,omitempty"`
	Notes             *string `json:"Notes,omitempty"`
	IsNotePrivate     *bool   `json:"IsNotePrivate,omitempty"`
	StaffAbbrev       *string `json:"StaffAbbrev,omitempty"`
	StaffBillSysId    *string `json:"StaffBillSysId,omitempty"`
	StaffEmail        *string `json:"StaffEmail,omitempty"`
	StaffEmrId        *string `json:"StaffEmrId,omitempty"`
	StaffErpId        *string `json:"StaffErpId,omitempty"`
	StaffInternalId   *string `json:"StaffInternalId,omitempty"`
	StaffExtCallSysId *string `json:"StaffExtCallSysId,omitempty"`
	StaffFName        *string `json:"StaffFName,omitempty"`
	StaffId           *string `json:"StaffId,omitempty"`
	StaffKey          *string `json:"StaffKey,omitempty"`
	StaffLName        *string `json:"StaffLName,omitempty"`
	StaffMobilePhone  *string `json:"StaffMobilePhone,omitempty"`
	StaffNpi          *string `json:"StaffNpi,omitempty"`
	StaffPager        *string `json:"StaffPager,omitempty"`
	StaffPayrollId    *string `json:"StaffPayrollId,omitempty"`
	StaffTags         []any   `json:"StaffTags,omitempty"`
	// StartDate              *time.Time `json:"StartDate,omitempty"`
	// StartTime              *time.Time `json:"StartTime,omitempty"`
	StartDate              *string `json:"StartDate,omitempty"`
	StartTime              *string `json:"StartTime,omitempty"`
	TaskAbbrev             *string `json:"TaskAbbrev,omitempty"`
	TaskBillSysId          *string `json:"TaskBillSysId,omitempty"`
	TaskContactInformation *string `json:"TaskContactInformation,omitempty"`
	TaskExtCallSysId       *string `json:"TaskExtCallSysId,omitempty"`
	TaskId                 *string `json:"TaskId,omitempty"`
	TaskKey                *string `json:"TaskKey,omitempty"`
	TaskName               *string `json:"TaskName,omitempty"`
	TaskPayrollId          *string `json:"TaskPayrollId,omitempty"`
	TaskEmrId              *string `json:"TaskEmrId,omitempty"`
	TaskCallPriority       *string `json:"TaskCallPriority,omitempty"`
	TaskDepartmentId       *string `json:"TaskDepartmentId,omitempty"`
	TaskIsPrintEnd         *bool   `json:"TaskIsPrintEnd,omitempty"`
	TaskIsPrintStart       *bool   `json:"TaskIsPrintStart,omitempty"`
	TaskShiftKey           *string `json:"TaskShiftKey,omitempty"`
	TaskType               *string `json:"TaskType,omitempty"`
	TaskTags               []any   `json:"TaskTags,omitempty"`
	LocationName           *string `json:"LocationName,omitempty"`
	LocationAbbrev         *string `json:"LocationAbbrev,omitempty"`
	LocationID             *string `json:"LocationID,omitempty"`
	LocationAddress        *string `json:"LocationAddress,omitempty"`
	TimeZone               *string `json:"TimeZone,omitempty"`
	// LastModifiedDateUTC    *time.Time `json:"LastModifiedDateUTC,omitempty"`
	LastModifiedDateUTC *string `json:"LastModifiedDateUTC,omitempty"`
	LocationTags        []any   `json:"LocationTags,omitempty"`
	IsRotationTask      *bool   `json:"IsRotationTask"`
}

// func Schedule(rqf *RequestQueryFields) *Request {

// 	r := NewRequest()
// 	r.Path = path.Join(r.Path, "schedule")
// 	r.SetIncludes("StaffTags,TaskTags,LocationTags")
// 	r.SetStartDate(time.Now().AddDate(0, 0, -14).UTC())
// 	r.SetEndDate(time.Now().UTC())
// 	if rqf != nil {
// 		if rqf.CompanyKey != nil {
// 			r.SetCompanyKey(rqf.GetCompanyKey())
// 		}
// 		if rqf.StartDate != nil {
// 			r.SetStartDate(rqf.GetStartDate())
// 		}
// 		if rqf.EndDate != nil {
// 			r.SetEndDate(rqf.GetEndDate())
// 		}
// 		if rqf.IncludeDeletes != nil {
// 			r.SetIncludeDeletes(rqf.GetIncludeDeletes())
// 		}
// 		if rqf.SinceModifiedTimestamp != nil {
// 			r.SetSinceModifiedTimestamp(rqf.GetSinceModifiedTimestamp())
// 		}
// 		if rqf.DateFormat != nil {
// 			r.SetDateFormat(rqf.GetDateFormat())
// 		}
// 		if rqf.Includes != nil {
// 			r.SetIncludes(rqf.GetIncludes())
// 		}
// 		if rqf.Select != nil {
// 			r.SetSelect(rqf.GetSelect())
// 		}
// 		if rqf.Filter != nil {
// 			r.SetFilter(rqf.GetFilter())
// 		}
// 		if rqf.Orderby != nil {
// 			r.SetOrderby(rqf.GetOrderby())
// 		}
// 		if rqf.Expand != nil {
// 			r.SetExpand(rqf.GetExpand())
// 		}
// 	}

// 	return r
// }
