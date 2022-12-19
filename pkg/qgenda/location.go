package qgenda

func NewLocationRequest(rqf *RequestConfig) *Request {
	requestPath := "location"
	queryFields := []string{
		"CompanyKey",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
		"Includes",
	}
	if rqf != nil {
		if rqf.Includes == nil {
			rqf.SetIncludes("Tags")
		}
	}

	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}

type Location struct {
	CompanyKey  *string       `json:"CompanyKey,omitempty"`
	LocationKey *int64        `json:"LocationKey,omitempty"`
	ID          *string       `json:"Id,omitempty"`
	Name        *string       `json:"Name,omitempty"`
	Address     *string       `json:"Address,omitempty"`
	Abbrev      *string       `json:"Abbrev,omitempty"`
	Notes       *string       `json:"Notes,omitempty"`
	TimeZone    *string       `json:"TimeZone,omitempty"`
	Tags        []TagCategory `json:"Tags,omitempty"`
}

// func NewLocationStaffRequest(rqf *RequestConfig) *Request {
// 	requestPath := "location/:locationId/staff"
// 	queryFields := []string{
// 		"CompanyKey",
// 		"Select",
// 		"Filter",
// 		"OrderBy",
// 		"Expand",
// 	}

// 	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
// 	return r
// }
// func NewLocationTasksRequest(rqf *RequestConfig) *Request {
// 	requestPath := "location/:locationId/tasks"
// 	queryFields := []string{
// 		"CompanyKey",
// 		"Select",
// 		"Filter",
// 		"OrderBy",
// 		"Expand",
// 	}

// 	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
// 	return r
// }
