package qgenda

import (
	"net/http"
	"time"
)

type Response struct {
	Time time.Time
	data []byte
}

func Read(r *http.Response) *Response {
	var resp Response
	return &resp
}
