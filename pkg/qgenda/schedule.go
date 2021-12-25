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
