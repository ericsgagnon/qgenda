package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"
)

// String represents a 'nullable' string by using  Valid member
// as a flag, in the same way sql.NullString does
type String struct {
	Data  string
	Valid bool
}

// NewString is strictly for convenience  - it sets the
// Valid member to true, even for empty strings
func NewString(s string) String {
	return String{
		Data:  s,
		Valid: true,
	}
}

func (v String) String() string {
	if v.Valid {
		return v.Data
	}
	return ""
}

// func (v String) MarshalBinary() ([]byte, error)     {}
// func (v *String) UnmarshalBinary(data []byte) error {}

func (v String) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Data)
	}
	return []byte{}, nil
}
func (v *String) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	v.Valid = true
	return nil
}

func (v *String) Scan(value interface{}) error {
	var sqlV sql.NullString
	if err := sqlV.Scan(value); err != nil {
		v.Data, v.Valid = "", false
		return err
	}
	v.Data = sqlV.String
	v.Valid = true
	return nil
}

func (v String) Value() (driver.Value, error) {
	sqlV := sql.NullString{
		String: v.Data,
		Valid:  v.Valid,
	}
	return sqlV.Value()
}

// func (v ) GoString() string {}
// func (v ) String() string {}

// func (v *) GobDecode(data []byte) error {}
// func (v ) GobEncode() ([]byte, error) {}

// func (v ) MarshalBinary() ([]byte, error) {}
// func (v *) UnmarshalBinary(data []byte) error {}

// func (v ) MarshalJSON() ([]byte, error) {}
// func (v *) UnmarshalJSON(data []byte) error {}

// func (v ) MarshalText() ([]byte, error) {}
// func (v *) UnmarshalText(data []byte) error {}

// func (v *) Scan (value interface{}) error {}
// func (v ) Value () (driver.Value, error ) {}

// func (n Nullable) GoString() string {

// 	return ()
// }

// func (n Nullable) String() string {

// }

// func (n *Nullable) GobDecode(data []byte) error {

// }

// func (n Nullable) GobEncode() ([]byte, error) {

// }

// func (n Nullable) MarshalBinary() ([]byte, error) {

// }

// func (n *Nullable) Scan(value interface{}) error {

// }

// func (n Nullable) Value() (driver.Value, error) {

// }

type Time struct {
	Value time.Time
	Valid bool
}

type Int struct {
	Value int
	Valid bool
}

type Byte struct {
	Value byte
	Valid bool
}

type Bool struct {
	Value bool
	Valid bool
}

// // just a test
// func (n Nullable) OK() bool {
// 	return n.Valid
// }

// func Test(n Nullable) bool {
// 	return n.OK()
// }

type TNullable interface {
	OK() bool
	Set(interface{}) error
}

func (s String) OK() bool {
	return s.Valid
}

// func (n Nullable) GetData() interface{} {
// 	return n.Data
// }

// func (n *Nullable) SetData(v interface{}) error {
// 	n.Data = v
// 	return nil
// }

// func (s *String) Set(i interface{}) error {
// 	s.Value = string(i)

// }

type Still struct {
	A int
	B *int
	C *int
	D *bool
	E *string
	F string
	G *time.Time `json:",omitempty"`
}

func main() {
	// x := String{}
	// fmt.Printf("%+v\n", Test(x.Nullable))

	// var y Wrapper
	// y.Valid = true
	// y.Data = "yay"
	// fmt.Printf("%+v\n", y.Maybe())
	// fmt.Printf("%+v\n", y.GetData())
	// if err := y.SetData("setting data"); err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Printf("%+t\n", y.GetData())
	// fmt.Printf("%+t\n", "test")
	// fmt.Printf("%+v\n", y)
	jd := []byte(`{
		"a": 50,
		"b": 10,
		"e": "yep",
		"f": "nope"
	}`)

	var st Still
	json.Unmarshal(jd, &st)

	fmt.Printf("%+v\n", st)
	fmt.Printf("%s\n", st)

	ts := Still{
		A: 10,
		B: intPointer(20),
		E: stringPointer("stringPointer"),
		F: "thestring",
		G: timePointer(time.Now()),
	}
	dOut, err := json.Marshal(ts)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(dOut))

	dOut2, err := json.Marshal(st)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(dOut2))
	rv := reflect.ValueOf(st)
	fmt.Printf("%#v\n", rv)
	rvi := reflect.Indirect(rv)
	fmt.Printf("%#v\n", rvi)
	rvia := rvi.FieldByName("A")
	fmt.Printf("%#v\n", rvia)
	rvib := rvi.FieldByName("B")
	fmt.Printf("%#v\n", rvib)
	rvibi := reflect.Indirect(rvib)
	fmt.Printf("%#v\n", rvibi)

	fmt.Printf("%#v\n", rvib.Type().Kind())
	// .FieldByName("Resource").Interface().(string)

}

func stringPointer(s string) *string {
	return &s
}

func intPointer(i int) *int {
	return &i
}

func timePointer(t time.Time) *time.Time {
	return &t
}

