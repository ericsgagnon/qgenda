package main

import (
	"encoding/json"
	"fmt"
	"log"
	//"reflect"
	"time"
)

type Time time.Time

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *Time) UnmarshalJSON(data []byte) error {

	//tag := reflect.ValueOf(data).Type().Field(0).Tag.Get("json")
	//fmt.Println(tag)

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	tt, err := time.ParseInLocation(`"`+"2006-01-02T15:04:05"+`"`, string(data), location)
	*t = Time(tt)
	return err
}

type Container struct {
	Date Time `json:"Test"`
}

func main() {
	var d Container
	jd := []byte(`{ "Test": "2000-01-01T00:00:00" }`)
	err := json.Unmarshal(jd, &d)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", time.Time(d.Date))
}
