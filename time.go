package main

import (
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
