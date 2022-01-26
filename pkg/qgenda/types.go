package qgenda

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/jackc/pgtype"
)

// Time embeds time.Time for custom json un/marshalling
type Time struct {
	time.Time
	Valid bool
}

// Pointer simply returns a pointer to a value. It is useful
// when using literals for pointer assignments.
func Pointer[T any](t T) *T {
	return &t
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

// Date wraps jackc/pgtype's Date, adds JSON un/marshaling support,
// and allows string customization using Set/Layout(). By default, it
// supports parsing a variety of time of day style strings and will encode to
// `2006-01-02` layout.
type Date struct {
	pgtype.Date
	layout string // should be a valid go time layout https://pkg.go.dev/time#pkg-constants
}

// Layout and SetLayout can be used to customize string and marshalling representations
// on each instance. The actual Date.layout is kept private to prevent unwanted
// un/marshaling.
func (d Date) Layout() string {
	if d.layout == "" {
		return `2006-01-02`
	}
	return d.layout
}

func (d *Date) SetLayout(l string) {
	d.layout = l
}

func GuessDateLayout(s string) (string, error) {
	// attempting to handle a small variety of time of day inputs
	// `mm/DD/yyyy`
	// `DD/mm/yyyy`

	// dfr1 := regexp.MustCompile(`(?i)^(\d{1,2})(\D*)(\d{1,2})(\D*)(\d{4}|\d{2})(?:.*)`, s)
	// matches1 := dfr1.FindStringSubmatch(s)
	// match1, err := strconv.Atoi(matches1[1])
	// if err != nil {
	// 	return "", err
	// }
	// match3, err := strconv.Atoi(matches1[3])
	// if err != nil {
	// 	return "", err
	// }

	// components1 := map[string]string{}
	// for i, v := range matches1 {
	// 	if i == 0 {
	// 		continue
	// 	}
	// 	strconv.Atoi()
	// }
	// ``yyyy-mm-DDTHH:MM:SS`
	// `yyyy-mm-DD`
	// `yyyy-m-D`
	// `__yy/_m/_D`
	// rfc3339 := `(?i)^(?P<year>\d{4})(?P<ymsep>\D*)(?P<month>\d{1,2})(?P<mdsep>\D*)(?P<day>\d{1,2})((?P<dtsep>[tT ])(?P<hour>)(?P<hmsep>)(?P<minute>)(?P<mssep>)(?P<second>)(?P<subsecond>))?$`
	// `(?i)^(?P<year>\d{4})(?P<ymsep>\D*)(?P<month>\d{1,2})(?P<mdsep>\D*)(?P<day>\d{1,2})((?P<dtsep>\D*)(?P<hour>\d{1,2})(?P<hmsep>\D*)(?P<minute>\d{1,2})(?P<mssep>\D*)(?P<second>\d{1,2})(?P<subsecond>[[:punct:]]+\d*)?)?$`
	dfr := regexp.MustCompile(`(?i)^(?P<year>\d{4}|\d{2})(?P<ymsep>\D*)(?P<month>\d{1,2})(?P<mdsep>\D*)(?P<day>\d{1,2})(?P<time>(?P<dtsep>\D*)(?P<hour>\d{1,2})(?P<hmsep>\D*)(?P<minute>\d{1,2})(?P<mssep>\D*)(?P<second>\d{1,2})(?P<subsecond>[[:punct:]]+\d*)?)?$`)
	matches := dfr.FindStringSubmatch(s)
	// fmt.Printf("TimeString: %20s Matches: %s\n", ts, matches)
	components := map[string]string{}
	for _, v := range dfr.SubexpNames() {

		subMatchIndex := dfr.SubexpIndex(v)
		if subMatchIndex > 0 && subMatchIndex <= len(dfr.SubexpNames()) {
			components[v] = matches[subMatchIndex]
			// fmt.Printf("%s:\t%s\n", v, matches[subMatchIndex])

		}
	}
	switch len(components["year"]) {
	case 2:
		components["year"] = "06"
	case 4:
		components["year"] = "2006"
	}

	monthFormat := `%0` + fmt.Sprint(len(components["month"])) + `d`
	components["month"] = fmt.Sprintf(monthFormat, 1)

	dayFormat := `%0` + fmt.Sprint(len(components["day"])) + `d`
	components["day"] = fmt.Sprintf(dayFormat, 2)
	// fmt.Println(len(components["time"]))
	if len(components["time"]) > 0 {

		hourFormat := `%0` + fmt.Sprint(len(components["hour"])) + `d`
		components["hour"] = fmt.Sprintf(hourFormat, 15)
		// fmt.Println("Components - hour: ", components["hour"])
		minuteFormat := `%0` + fmt.Sprint(len(components["minute"])) + `d`
		components["minute"] = fmt.Sprintf(minuteFormat, 4)
		if len(components["second"]) > 0 {
			secondFormat := `%0` + fmt.Sprint(len(components["second"])) + `d`
			components["second"] = fmt.Sprintf(secondFormat, 5)

		}
		components["subsecond"] = regexp.
			MustCompile(`\d`).
			ReplaceAllLiteralString(components["subsecond"], "0")
	}

	format := ""
	for _, v := range dfr.SubexpNames() {
		if v != "time" {
			format = format + components[v]

		}
	}
	return format, nil
}

func ParseDateWithLayout(layout string, ds string) (Date, error) {
	d := Date{}
	t, err := time.Parse(layout, ds)
	if err != nil {
		return d, err
	}
	d.Set(t)
	d.Status = pgtype.Present
	return d, nil
}

func ParseDate(s string) (Date, error) {
	layout, err := GuessDateLayout(s)
	if err != nil {
		return Date{}, err
	}
	d, err := ParseDateWithLayout(layout, s)
	if err != nil {
		return d, err
	}
	return d, nil
}

func (d Date) String() string {
	if d.Status != pgtype.Present {
		return fmt.Sprintf("%s", time.Time{}.Format("2006-01-02"))
	}

	layout := d.layout
	if layout == "" {
		layout = "2006-01-02"
	}

	return d.Time.Format(layout)
}

func (d Date) MarshalJSON() ([]byte, error) {
	switch d.Status {
	case pgtype.Null:
		return []byte("null"), nil
	case pgtype.Undefined:
		return nil, errors.New("Cannot Date.MarshalJSON due to Date.Status == pgtype.Undefined")
	}

	if d.Status != pgtype.Present {
		return nil, errors.New("Cannot Date.MarshalJSON due to Date.Status != pgtype.Present")
	}

	var s string
	switch d.InfinityModifier {
	case pgtype.None:
		s = d.String()
	case pgtype.Infinity:
		s = "infinity"
	case pgtype.NegativeInfinity:
		s = "-infinity"
	}

	return json.Marshal(s)
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	// Ignore null, like in the main JSON package.
	if s == "null" {
		d.Status = pgtype.Null
		return nil
	}

	var layout string
	var err error
	if d != nil {
		layout = d.layout
	}
	if layout == "" {
		layout, err = GuessDateLayout(s)
		if err != nil {
			return err
		}
	}
	dt, err := ParseDateWithLayout(layout, s)
	if err != nil {
		d.Status = pgtype.Null
		return err
	}

	d.Time = dt.Time
	// d.layout = `2006-01-02`
	d.Status = pgtype.Present

	return nil
}

func (d *Date) Process() error {
	// let's make sure everything is in UTC and truncated to just the date
	if d.Time.Location() == time.UTC {
		d.Date.Time = d.Date.Time.Truncate(time.Hour * 24)
	} else {
		yr, mo, dt := d.Date.Time.Date()
		d.Date.Time = time.Date(yr, mo, dt, 0, 0, 0, 0, time.UTC)
	}
	return nil
}

// TimeOfDay wraps jackc/pgtype's Time, adds JSON un/marshaling support,
// and allows string customization using Set/Layout(). By default, it
// supports parsing a variety of time of day style strings and will encode to
// `15:04:05.999999` layout.
type TimeOfDay struct {
	pgtype.Time
	layout string // should be a valid go time layout https://pkg.go.dev/time#pkg-constants
}

// Layout and SetLayout can be used to customize string and marshalling representations
// on each instance. The actual TimeOfDay.layout is kept private to prevent unwanted
// un/marshaling.
func (t TimeOfDay) Layout() string {
	if t.layout == "" {
		return "15:04:05.999999"
	}
	return t.layout
}

func (t *TimeOfDay) SetLayout(s string) {
	t.layout = s
}

func GuessTimeOfDayLayout(s string) (string, error) {

	// attempting to handle a small variety of time of day inputs
	tfr := regexp.MustCompile(`(?i)^(?P<hour>\d{1,2})(?P<hmsep>[^\d]*)(?P<minute>\d{1,2})((?P<mssep>[^\d]*)(?P<second>\d{1,2})(?P<subsecond>[[:punct:]]\d+)?)?((?P<markprefix>[^ap]*)(?P<mark>[ap]m))?$`)
	matches := tfr.FindStringSubmatch(s)
	// fmt.Printf("TimeString: %20s Matches: %s\n", ts, matches)
	components := map[string]string{}
	for _, v := range tfr.SubexpNames() {

		subMatchIndex := tfr.SubexpIndex(v)
		if subMatchIndex > 0 && subMatchIndex <= len(tfr.SubexpNames()) {
			components[v] = matches[subMatchIndex]
			// fmt.Printf("%s:\t%s\n", v, matches[subMatchIndex])

		}
	}
	if len(components["mark"]) == 2 {
		components["mark"] = "PM"
		hourLength := fmt.Sprint(len(components["hour"]))
		hourFormat := `%0` + hourLength + `d`
		components["hour"] = fmt.Sprintf(hourFormat, 3)
	} else {
		components["hour"] = "15"
	}
	minuteFormat := `%0` + fmt.Sprint(len(components["minute"])) + `d`
	components["minute"] = fmt.Sprintf(minuteFormat, 4)
	if len(components["second"]) > 0 {
		secondFormat := `%0` + fmt.Sprint(len(components["second"])) + `d`
		components["second"] = fmt.Sprintf(secondFormat, 5)

	}
	components["subsecond"] = regexp.
		MustCompile(`\d`).
		ReplaceAllLiteralString(components["subsecond"], "9")

	format := ""
	for _, v := range tfr.SubexpNames() {
		format = format + components[v]
	}
	return format, nil
}

func ParseTimeOfDayWithLayout(layout string, ts string) (TimeOfDay, error) {
	tod := TimeOfDay{}
	t, err := time.Parse(layout, ts)
	if err != nil {
		return tod, err
	}
	tod.Set(t)
	// tod.SetLayout(`15:04:05.999999`)
	tod.Status = pgtype.Present
	return tod, nil
}

func ParseTimeOfDay(s string) (TimeOfDay, error) {
	timeFormat, err := GuessTimeOfDayLayout(s)
	if err != nil {
		return TimeOfDay{}, err
	}
	tod, err := ParseTimeOfDayWithLayout(timeFormat, s)
	if err != nil {
		return tod, err
	}
	return tod, nil
}

func (t TimeOfDay) String() string {
	if t.Status != pgtype.Present {
		return fmt.Sprintf("%s", time.Time{}.Format("15:04:05"))
	}

	layout := t.layout
	if layout == "" {
		layout = "15:04:05.999999"
	}

	td := time.Duration(t.Time.Microseconds) * time.Microsecond
	ts := time.Time{}.Add(td).Format(layout)
	return ts
}

func (t TimeOfDay) MarshalJSON() ([]byte, error) {
	switch t.Status {
	case pgtype.Null:
		return []byte("null"), nil
	case pgtype.Undefined:
		return nil, errors.New("Cannot TimeOfDay.MarshalJSON due to TimeOfDay.Status == pgtype.Undefined")
	}
	v, _ := t.Time.Value()
	if v == nil {
	}

	return json.Marshal(t.String())
	// b := []byte(`"` + t.String() + `"`)
	// return b, nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *TimeOfDay) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	// Ignore null, like in the main JSON package.
	if s == "null" {
		t.Status = pgtype.Null
		return nil
	}
	var layout string
	var err error
	if t.layout == "" {
		layout, err = GuessTimeOfDayLayout(s)
		if err != nil {
			return err
		}
	}
	tod, err := ParseTimeOfDayWithLayout(layout, s)
	if err != nil {
		return err
	}

	t.Time = tod.Time
	// t.layout = `15:04:05.999999`
	t.Status = pgtype.Present

	return nil
}

func (t *TimeOfDay) Process() error {

	// tod := *t
	// a qgenda time of day with value `00:03:00` is actually null
	// look near https://restapi.qgenda.com/#cecf3cdc-b32b-4b2c-8604-bd6fb59d7655
	todNull, err := ParseTimeOfDay(`00:03:00`)
	if err != nil {
		log.Fatalln(err)
	}

	if t.Time == todNull.Time {
		t.Time.Set(nil)
		t.Status = pgtype.Null
	}

	return nil
}

/////////////////////////////////
/////////////////////////////////
/////////////////////////////////
/////////////////////////////////
/////////////////////////////////
/////////////////////////////////
//
