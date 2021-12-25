package qgenda

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Time embeds time.Time for custom json un/marshalling
type Time struct {
	time.Time
	Valid bool
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *Time) UnmarshalJSON(data []byte) error {

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		return err
	}
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	tt, err := time.ParseInLocation(`"`+"2006-01-02T15:04:05"+`"`, string(data), location)
	if err != nil {
		return err
	}
	t.Time = tt
	if t.IsZero() {
		t.Valid = false
	} else {
		t.Valid = true
	}
	return err
}

// TimeUTC embeds time.Time for custom json un/marshalling
type TimeUTC struct {
	time.Time
	Valid bool
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *TimeUTC) UnmarshalJSON(data []byte) error {

	location, err := time.LoadLocation("UTC")
	if err != nil {
		return err
	}
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	tt, err := time.ParseInLocation(`"`+"2006-01-02T15:04:05"+`"`, string(data), location)
	if err != nil {
		return err
	}
	t.Time = tt
	if t.IsZero() {
		t.Valid = false
	} else {
		t.Valid = true
	}
	return err
}

//  embeds  for custom json un/marshalling
type NullTime sql.NullTime
type NullString sql.NullString
type NullBool sql.NullBool
type NullByte sql.NullByte
type NullFloat64 sql.NullFloat64
type NullInt16 sql.NullInt16
type NullInt32 sql.NullInt32
type NullInt64 sql.NullInt64

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *NullTime) UnmarshalJSON(data []byte) error {

	location, err := time.LoadLocation("UTC")
	if err != nil {
		return err
	}
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	tt, err := time.ParseInLocation(`"`+"2006-01-02T15:04:05"+`"`, string(data), location)
	if err != nil {
		return err
	}
	t.Time = tt
	if t.Time.IsZero() {
		t.Valid = false
	} else {
		t.Valid = true
	}
	return err
}

func (t NullTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() || !t.Valid {
		return []byte(""), nil
	}
	return []byte(t.Time.Format(time.RFC3339)), nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *NullString) UnmarshalJSON(data []byte) error {

	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &t.String); err != nil {
		return nil
	}
	// t.Valid = true
	return nil
}

func (t NullString) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String)
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *NullBool) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &t.Bool); err != nil {
		return nil
	}
	// t.Valid = true
	return nil
}

func (t NullBool) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return json.Marshal(t.Bool)
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *NullByte) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &t.Byte); err != nil {
		return nil
	}
	// t.Valid = true
	return nil
}

func (t NullByte) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return json.Marshal(t.Byte)
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *NullFloat64) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &t.Float64); err != nil {
		return nil
	}
	// t.Valid = true
	return nil
}

func (t NullFloat64) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return json.Marshal(t.Float64)
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *NullInt16) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &t.Int16); err != nil {
		return nil
	}
	// t.Valid = true
	return nil
}

func (t NullInt16) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return json.Marshal(t.Int16)
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *NullInt32) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &t.Int32); err != nil {
		return nil
	}
	// t.Valid = true
	return nil
}

func (t NullInt32) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return json.Marshal(t.Int32)
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *NullInt64) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == "" || string(data) == "~" {
		return nil
	}

	// var v string
	if err := json.Unmarshal(data, &t.Int64); err != nil {
		return nil
	}
	// t.Valid = true
	return nil
}

func (t NullInt64) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return json.Marshal(t.Int64)
}

// // UnmarshalJSON satisfies the json.Unmarshaler interface
// func (t *NullBool) UnmarshalJSON(data []byte) error {

// 	// Ignore null, like in the main JSON package.
// 	if string(data) == "null" {
// 		return nil
// 	}

// 	t.Bool = json.UnmarshalJSON
// 	t.Valid = true
// 	return nil
// }

// func (a *Animal) UnmarshalJSON(b []byte) error {
// 	var s string
// 	if err := json.Unmarshal(b, &s); err != nil {
// 		return err
// 	}
// 	switch strings.ToLower(s) {
// 	default:
// 		*a = Unknown
// 	case "gopher":
// 		*a = Gopher
// 	case "zebra":
// 		*a = Zebra
// 	}

// 	return nil
// }
