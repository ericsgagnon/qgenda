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
	Data     *[]byte `yaml:"-"`
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

// CacheFile parses the input filename or default cache file, creates the directory if
// necessary and returns the complete path/filename as a string
func CacheFile(filename string, defaultDir string, defaultFileName string) (string, error) {
	var err error
	// parse path to cache directory or default to user's cache directory/qgenda/auth
	p := strings.TrimSuffix(filename, filepath.Base(filename))
	if p == "" || p == "." {
		if p, err = os.UserCacheDir(); err != nil {
			log.Printf("Error retrieving cache directory: %v", err)
			return "", err
		}
		p = p + defaultDir
	}
	// make cache directory
	if err := os.MkdirAll(p, 0777); err != nil {
		log.Printf("Error making directory %v: %#v", p, err)
		return "", err
	}

	// parse filename or default to authtoken.json
	f := filepath.Base(filename)
	f = strings.ToLower(f)
	if f == "" || f == "*" || f == "." {
		f = defaultFileName
	}
	// compile absolute path + file name
	f = filepath.Join(p, f)
	return f, nil
}

// ToJSONFile writes the raw json bytes to a formatted json text file at
// filename, or uses *Response.Metadata.Name to form a default
func (r *Response) ToJSONFile(filename string) error {
	f, err := CacheFile(filename, "/qgenda/data/in", strings.ToLower(r.Metadata.Name)+".json")
	if err != nil {
		log.Printf("Error %v", err)
		return err
	}
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
