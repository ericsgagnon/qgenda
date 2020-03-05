package main

// TaskLocationRequestConfig is intended to be used as inputs to
// api requests to the Task endpoints
type TaskLocationRequestConfig struct {
	Resource string `resource:"-"`
	Route    string `path:"-"`
	TaskID   string `path:"taskid"`
	Select   string `query:"$select"`
	Filter   string `query:"$filter"`
	OrderBy  string `query:"$orderby"`
	Expand   string `query:"$expand"`
}

// NewTaskLocationRequestConfig returns a pointer to a TaskLocationRequestConfig with default values
func NewTaskLocationRequestConfig(rc *TaskLocationRequestConfig) *TaskLocationRequestConfig {
	if rc == nil {
		rc = &TaskLocationRequestConfig{}
	}

	r := &TaskLocationRequestConfig{
		Resource: "TaskLocation",
		Route:    "/task/{{.TaskID}}/location",
		TaskID:   "",
		// Select:   "",
		// Filter:   "",
		// OrderBy:  "",
		// Expand:   "",
	}

	fillDefaults(rc, r)
	return rc
}

// NewTaskLocationRequestResponse returns a pointer to a ScheduleRequestConfig with default values
func NewTaskLocationRequestResponse(rc *TaskLocationRequestConfig) *RequestResponse {
	rr := NewRequestResponse()
	rr.RequestConfig = NewTaskLocationRequestConfig(rc)
	return rr
}

// Parse parses the RequestConfig into one or more Requests
func (rc TaskLocationRequestConfig) Parse() ([]Request, error) {
	var req []Request
	reqi, err := parseRequestConfig(rc)
	if err != nil {
		return []Request{}, err
	}
	req = append(req, reqi)
	return req, nil
}
