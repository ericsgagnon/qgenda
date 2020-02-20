package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Response holds select information from the qgenda api response
type Response struct {
	Metadata *Metadata
	Data     *[]byte
}

// Metadata captures relevant metadata from each response
type Metadata struct {
	APIVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	Timestamp  time.Time `json:"time"`
}

// NewResponse initializes an empty struct and returns a points
// it's primarily meant in other constructors to avoid null pointer
func NewResponse() *Response {
	r := &Response{
		Metadata: &Metadata{},
		Data:     &[]byte{},
	}
	return r
}

// FormatJSON decodes and indents the raw json bytes to a formatted json string
func (r *Response) FormatJSON() ([]byte, error) {
	var b bytes.Buffer
	if err := json.Indent(&b, *r.Data, "", "  "); err != nil {
		log.Printf("Error marshalling to json: %+v", err)
		return nil, err
	}
	return b.Bytes(), nil

}

// ToJSONFile writes the raw json bytes to a formatted json text file at
// filename, or uses *Response.Metadata.Name to form a default
func (r *Response) ToJSONFile(filename string) error {

	p := strings.TrimSuffix(filename, filepath.Base(filename))
	//p := filepath.Dir(filename)
	// fmt.Println(p)
	// create directory, or use "data/in" as default
	if p == "" || p == "." {
		p = "data/in"
	}
	// fmt.Println(p)

	if err := os.MkdirAll(p, 0777); err != nil {
		log.Printf("Error making directory %v: %#v", p, err)
		return err
	}

	f := filepath.Base(filename)
	f = strings.ToLower(f)
	// build filename if not provided
	if f == "" || f == "*" || f == "." {
		f = strings.ToLower(r.Metadata.Name) + ".json"
	}
	f = filepath.Join(p, f)
	// fmt.Println(f)
	var b bytes.Buffer
	if err := json.Indent(&b, *r.Data, "", "  "); err != nil {
		log.Printf("Error marshalling to json: %+v", err)
		return err
	}

	if err := ioutil.WriteFile(f, b.Bytes(), 0755); err != nil {
		log.Printf("Error writing file %v to disk: %v", f, err)
		return err
	}

	return nil
}

// FromJSONFile reads an itemlist from a jsonfile
// func FromJSONFile(f string, il *ItemList) error {
// 	b, err := ioutil.ReadFile(f)
// 	if err != nil {
// 		log.Printf("Error Reading file %v: %v", f, err)
// 		return err
// 	}

// 	if err := json.Unmarshal(b, il); err != nil {
// 		log.Printf("Error Unmarshaling file %v: %v", f, err)
// 		return err
// 	}
// 	return nil
// }
