package main

// TaskRequestConfig is intended to be used as inputs to
// api requests to the Task endpoints
type TaskRequestConfig struct {
	Resource string `resource:"-"`
	Route    string `path:"-"`
	Includes string `query:"includes"`
	Select   string `query:"$select"`
	Filter   string `query:"$filter"`
	OrderBy  string `query:"$orderby"`
	Expand   string `query:"$expand"`
}

// NewTaskRequestConfig returns a pointer to a TaskRequestConfig with default values
func NewTaskRequestConfig(rc *TaskRequestConfig) *TaskRequestConfig {
	if rc == nil {
		rc = &TaskRequestConfig{}
	}

	r := &TaskRequestConfig{
		Resource: "Task",
		Route:    "/task",
		Includes: "Tags,TaskShifts",
		// Select:   "",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}

	fillDefaults(rc, r)
	return rc
}

// NewTaskRequestResponse returns a pointer to a ScheduleRequestConfig with default values
func NewTaskRequestResponse(rc *TaskRequestConfig) *RequestResponse {
	rr := NewRequestResponse()
	rr.RequestConfig = NewTaskRequestConfig(rc)
	return rr
}

// Parse parses the RequestConfig into one or more Requests
func (rc TaskRequestConfig) Parse() ([]Request, error) {
	var req []Request
	reqi, err := parseRequestConfig(rc)
	if err != nil {
		return []Request{}, err
	}
	req = append(req, reqi)
	return req, nil
}
