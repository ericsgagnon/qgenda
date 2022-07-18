package qgenda

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/jackc/pgtype"
)

// DefaultTimeLocation can be used as a sneaky, backchannel way to
// to parse Time strings that don't have time zone info. By default,
// it is nil. If neither the Time struct's location is set nor DefaultTimeLocation,
// then un/marshalling is done in UTC.
var DefaultTimeLocation *time.Location

// DefaultTimeLayout is an even sneakier and more backchannely way of parsing times
// from qgenda. Qgenda seems to prefer timestamps without time zone and just assume
// local times.
var DefaultTimeLayout string = "2006-01-02T15:04:05.999999"

// Time wraps jackc/pgtype's Timestamptz (which wraps time.Time), adds JSON un/marshaling support,
// and allows string customization using Set/Layout(). By default, it will encode to
// `2006-01-02T15:04:05.999999Z` layout.
// Note that layout and location are primarily intended to assist in unmarshalling and processing
type Time struct {
	pgtype.Timestamptz
	source         string         // optional to hold the original string used to parse time
	sourceTZ       bool           // indicate if source string has location in it
	sourceLocation *time.Location // location derived directly from parsing source
	layout         string         // should be a valid go time layout https://pkg.go.dev/time#pkg-constants
	location       *time.Location // nil by default - can also use DefaultTimeLocation
}

func (t *Time) Set(src any) error {
	return t.Timestamptz.Set(src)
}

func NewTime(src any) Time {
	t := Time{}
	if src == nil {
		t.Set(time.Now())
		return t
	}
	switch src.(type) {
	case time.Time, Time:
		t.Set(src)
	case string:
		tm, err := ParseTime(src.(string))
		if err != nil {
			panic(err)
		}
		t.Set(tm)
	}
	return t
}

// Layout and SetLayout can be used to customize string and marshalling representations
// on each instance. The actual TimeOfDay.layout is kept private to prevent unwanted
// un/marshaling.
func (t Time) Layout() string {
	if t.layout == "" {
		return "2006-01-02T15:04:05.999999Z"
	}
	return t.layout
}

func (t *Time) SetLayout(s string) {
	t.layout = s
}

func (t *Time) SetLocation(loc *time.Location) {
	t.location = loc
}

// Location will return the set location or nil, if not set.
// note that a nil *time.Location's String() method will return UTC (or local if set?)
func (t Time) Location() *time.Location {
	return t.location
}

// Source returns the source string, if one exists. Note that, if Time wasn't created from a
// string, then this will be the empty string
func (t Time) Source() string {
	return t.source
}

func (t Time) SourceHasLocation() bool {
	return (t.source != "" && t.sourceTZ)
}

// SourceLocation returns the location derived from the source string. It only applies if
// a location could be determined from the source string.
func (t Time) SourceLocation() *time.Location {
	if t.source != "" && t.sourceTZ && t.sourceLocation != nil {
		return t.sourceLocation
	}
	return nil
}

// ParseTime uses DefaultTimeLayout and DefaultTimeLocation
// to parse from only a string. If you need to control layout and location,
// used ParseTimeWithLayout or ParseTimeInLocation
func ParseTime(value string) (Time, error) {

	layout, err := MatchTimeLayout(value, DefaultTimeLayout)
	// log.Printf("Parsing %s to Time - layout: %s\n", value, layout)
	if err != nil {
		return Time{}, err
	}
	// log.Printf("ParseTime - DefaultTimeLocation is nil: %t %#v\n", DefaultTimeLocation == nil, DefaultTimeLocation)
	t, err := ParseTimeInLocation(layout, value, DefaultTimeLocation)
	// t, err := ParseTimeInLocation(layout, value, (*time.Location)(nil))
	// log.Printf("Parsing successful: %s\n", t.Time)
	// t.ChangeTimeLocation(DefaultTimeLocation)
	// log.Printf("Changed Location to %s - new time: %s\n", t.Time.Location(), t)
	// tt, err := time.Parse(layout, value)
	// log.Printf("Parsing as time.Time successful: %s\n", tt)
	if err != nil {
		return t, err
	}
	return t, nil

}

// ParseTimeWithLayout will parse in the DefaultTimeLocation, if it is set.
// this function is now fairly extraneous and may be removed.
func ParseTimeWithLayout(layout string, value string) (Time, error) {
	return ParseTimeInLocation(layout, value, DefaultTimeLocation)
}

// ParseTimeInLocation expects explicit values to be passed to it, it will not
// use DefaultTimeLayout or DefaultTimeLocation. Like time.ParseInLocation (which it wraps)
// it will only use loc if there is no time zone info in the string being parsed.
// sourceLocation is only set to the Location from the parsed time.Time if the string has
// time zone in it or loc is a valid *time.Location.
func ParseTimeInLocation(layout string, value string, loc *time.Location) (Time, error) {
	t := Time{
		source: value,
	}
	// test := (*time.Location)(nil)
	// test := &time.Location{}
	// log.Printf("loc is nil: %t %#v\n", loc == nil, loc)
	// log.Printf("test *time.Location{} is nil: %t %#v\n", test == nil, test)
	// tm, err := time.ParseInLocation(layout, value, loc)
	location := &time.Location{}
	if loc != nil {
		location = loc
	} else {
		location = time.UTC
	}
	// log.Printf("loc is nil: %t %#v\n", loc == nil, loc)
	// log.Printf("layout %s\tvalue %s\tlocation %#v\n", layout, value, location)
	tm, err := time.ParseInLocation(layout, value, location)
	// tm, err := time.ParseInLocation(layout, value, &time.Location{})
	if err != nil {
		t.Status = pgtype.Null
		return t, err
	}
	// check if layout has a location in it
	// valid go formats use one of the following to ID time zone
	// MST +-0700 +-07:00 +-07 Z0700 Z07:00 Z07
	// given these, it should be enough to check for MST or 7 (???)
	t.sourceTZ = (strings.Contains(layout, "7") || strings.Contains(layout, "MST"))
	if t.sourceTZ || (loc != nil && loc == tm.Location()) {
		t.sourceLocation = tm.Location()
	}

	t.Set(tm.UTC())
	t.Status = pgtype.Present
	return t, nil
}

// ChangeTimeLocation returns a new Time with the same calendar/wall time in the new location.
// If ignoreSource == false and the sourceLocation == loc, it will return the original Time.
func ChangeTimeLocation(t Time, loc *time.Location, ignoreSourceLocation bool) (Time, error) {
	if (!ignoreSourceLocation && t.sourceLocation == loc) || loc == nil {
		return t, nil
	}

	layout := "2006-01-02T15:04:05.999999999"
	ts := t.Time.Format(layout)
	return ParseTimeInLocation(layout, ts, loc)
}

// ChangeTimeZone is just a friendly wrapper for ChangeTimeLocation
func ChangeTimeZone(t Time, loc *time.Location, ignoreSourceLocation bool) (Time, error) {
	return ChangeTimeLocation(t, loc, ignoreSourceLocation)
}

// ChangeTimeLocation in method form
func (t *Time) ChangeLocation(loc *time.Location, ignoreSourceLocation bool) error {
	tm, err := ChangeTimeLocation(*t, loc, ignoreSourceLocation)
	if err != nil {
		return err
	}
	// log.Printf("(t *Time) ChangeLocation: %s\n", tm)
	*t = tm
	return nil
}

// ChangeTimeZone in method form
func (t *Time) ChangeTimeZone(loc *time.Location, ignoreSourceLocation bool) error {
	return t.ChangeLocation(loc, ignoreSourceLocation)
}

// ProcessTime is the underlying function used to handle any 'in-flight' processing
// before loading to a destination. It should be wrapped to satisfy any Processor type
// interfaces.
func ProcessTime(t Time) (Time, error) {
	out := Time{}
	var err error
	// reject out of range times
	if t.Time.Year() < 1800 || t.Time.Year() > 2099 {
		err = fmt.Errorf("warn qgenda.ProcessTime() year out of range (1800-2099): %s\n", t.String())
		out.Status = pgtype.Null
		return out, err
	}

	// set
	out, err = ChangeTimeLocation(t, DefaultTimeLocation, false)
	if err != nil {
		return out, err
	}

	return out, nil
}

// Process satisfies the Processor interface
func (t *Time) Process() error {
	tm, err := ProcessTime(*t)
	*t = tm
	// *t = Time{}
	return err
}

// String uses Time's layout for formatting
// Note: String does not use DefaultTimeLayout or DefaultTimeLocation -
// if you want to change the string representation, you have to t.SetLayout and t.SetLocation
func (t Time) String() string {
	layout := t.layout
	if layout == "" {
		layout = time.RFC3339Nano
	}

	if t.Status != pgtype.Present {
		return time.Time{}.Format(layout)
	}
	loc := t.location
	if loc == nil {
		loc = time.UTC
	}
	return t.Time.In(loc).Format(layout)
}

func (t Time) MarshalJSON() ([]byte, error) {
	// switch t.Status {
	// case pgtype.Null:
	// 	return []byte("null"), nil
	// case pgtype.Undefined:
	// 	// return nil, errors.New("Cannot Time.MarshalJSON due to Time.Status == pgtype.Undefined")
	// 	return []byte("null"), nil
	// }

	if t.Status != pgtype.Present {
		// return nil, errors.New("Cannot Time.MarshalJSON due to Time.Status != pgtype.Present")
		return []byte("null"), nil
	}

	var s string
	switch t.InfinityModifier {
	case pgtype.None:
		s = t.String()
	case pgtype.Infinity:
		s = "infinity"
	case pgtype.NegativeInfinity:
		s = "-infinity"
	}

	return json.Marshal(s)
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (t *Time) UnmarshalJSON(data []byte) error {
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
	if t != nil {
		layout = t.layout
	}
	layout, err = MatchTimeLayout(s, layout, DefaultTimeLayout)
	if err != nil {
		return err
	}

	var loc *time.Location
	if t != nil {
		if t.location != nil {
			loc = t.location
		} else {
			loc = DefaultTimeLocation
		}
	}
	tm, err := ParseTimeInLocation(layout, s, loc)
	if err != nil {
		t.Status = pgtype.Null
		return err
	}
	if err := (&tm).Process(); err != nil {
		log.Println("Time.UnmarshalJSON: *Time.Process error: ", err)
	}
	*t = tm

	if err := t.Process(); err != nil {
		log.Println(err)
	}
	return nil
}

// MatchTimeLayout matches the value string with the first format string
// that parses into a valid go time. If none of those match, it will attempt
// a number of other formats, including most formats listed in the time package.
func MatchTimeLayout(value string, formats ...string) (string, error) {

	// possibleTimeLayouts := []string{}
	// if DefaultTimeLayout != "" {
	// 	possibleTimeLayouts = append(possibleTimeLayouts, DefaultTimeLayout)
	// }

	standardTimeLayouts := []string{
		time.RFC3339Nano,
		"2006-01-02t15:04:05.999999999Z07:00", // RFC3339Nano with t instead of T
		"2006-01-02 15:04:05.999999999Z07:00", // RFC3339Nano with space instead of T
		time.RFC3339,
		"2006-01-02t15:04:05Z07:00",     // RFC3339 with t instead of T
		"2006-01-02 15:04:05Z07:00",     // RFC3339 with space instead of T
		"2006-01-02T15:04:05.999999999", // ~RFC3339 without time zone... necessary evil
		"2006-01-02t15:04:05.999999999", // ~RFC3339
		"2006-01-02 15:04:05.999999999", // ~RFC3339
		time.RFC1123Z,
		time.RFC1123,
		"Mon, 02 Jan 2006 15:04:05.999999999", // ~RFC1123 without time zone
		time.RFC850,
		"Monday, 02-Jan-06 15:04:05.999999999", // ~RFC850 without time zone
		time.RFC822Z,
		time.RFC822,
		"02 Jan 06 15:04.999999999", // ~RFC822 without time zone
		time.RubyDate,
		time.UnixDate,
		time.ANSIC,
		time.StampNano,
		time.StampMicro,
		time.StampMilli,
		time.Stamp,
		time.Kitchen,
		"15:04:05.999999999",
		"15:04:05",
		"15:04",
		time.Layout, // go's reference time, lol
	}
	formats = append(formats, standardTimeLayouts...)
	for _, v := range formats {
		_, err := time.Parse(v, value)
		if err == nil {
			return v, nil
		}
	}
	return "", errors.New("MatchTimeLayout couldn't find an acceptable format.")
}

func MatchTimeRegex(s string) (string, error) {
	dow :=
		strings.Join(
			[]string{
				"Mon(day)?",
				"Tue(s?day)?",
				"Wed(nesday)?",
				"Thu?r?s?(day)?",
				"Fri(day)?",
				"Sat(urday)?",
				"Sun(day)?",
			},
			"|")
	dow = fmt.Sprintf("(?P<dow>%s)?", dow)
	// (?i)^
	// (?P<year>\d{4}|\d{2})
	//(?P<ymsep>\D*)
	//(?P<month>\d{1,2})
	//(?P<mdsep>\D*)
	//(?P<day>\d{1,2})
	//(?P<time>
	//(?P<dtsep>\D*)
	//(?P<hour>\d{1,2})
	//(?P<hmsep>\D*)
	//(?P<minute>\d{1,2})
	//(?P<mssep>\D*)
	//(?P<second>\d{1,2})
	//(?P<subsecond>[[:punct:]]+\d*)?
	//)?
	//$

	// dfr := regexp.MustCompile(`(?i)^(?P<year>\d{4}|\d{2})(?P<ymsep>\D*)(?P<month>\d{1,2})(?P<mdsep>\D*)(?P<day>\d{1,2})(?P<time>(?P<dtsep>\D*)(?P<hour>\d{1,2})(?P<hmsep>\D*)(?P<minute>\d{1,2})(?P<mssep>\D*)(?P<second>\d{1,2})(?P<subsecond>[[:punct:]]+\d*)?)?$`)
	// matches := dfr.FindStringSubmatch(s)
	// // fmt.Printf("TimeString: %20s Matches: %s\n", ts, matches)
	// components := map[string]string{}
	// for _, v := range dfr.SubexpNames() {

	// 	subMatchIndex := dfr.SubexpIndex(v)
	// 	if subMatchIndex > 0 && subMatchIndex <= len(dfr.SubexpNames()) {
	// 		components[v] = matches[subMatchIndex]
	// 		// fmt.Printf("%s:\t%s\n", v, matches[subMatchIndex])

	// 	}
	// }
	// switch len(components["year"]) {
	// case 2:
	// 	components["year"] = "06"
	// case 4:
	// 	components["year"] = "2006"
	// }

	// monthFormat := `%0` + fmt.Sprint(len(components["month"])) + `d`
	// components["month"] = fmt.Sprintf(monthFormat, 1)

	// dayFormat := `%0` + fmt.Sprint(len(components["day"])) + `d`
	// components["day"] = fmt.Sprintf(dayFormat, 2)
	// // fmt.Println(len(components["time"]))
	// if len(components["time"]) > 0 {

	// 	hourFormat := `%0` + fmt.Sprint(len(components["hour"])) + `d`
	// 	components["hour"] = fmt.Sprintf(hourFormat, 15)
	// 	// fmt.Println("Components - hour: ", components["hour"])
	// 	minuteFormat := `%0` + fmt.Sprint(len(components["minute"])) + `d`
	// 	components["minute"] = fmt.Sprintf(minuteFormat, 4)
	// 	if len(components["second"]) > 0 {
	// 		secondFormat := `%0` + fmt.Sprint(len(components["second"])) + `d`
	// 		components["second"] = fmt.Sprintf(secondFormat, 5)

	// 	}
	// 	components["subsecond"] = regexp.
	// 		MustCompile(`\d`).
	// 		ReplaceAllLiteralString(components["subsecond"], "0")
	// }

	// format := ""
	// for _, v := range dfr.SubexpNames() {
	// 	if v != "time" {
	// 		format = format + components[v]

	// 	}
	// }
	// return format, nil
	return "", errors.New("MatchTimeRegex isn't ready, don't use it")
}

// func (src Time) Value() (driver.Value, error) {
// 	// var t sql.NullTime
// 	// t.Time = src.Timestamptz.Time
// 	// fmt.Println("I'm a teapot")
// 	// fmt.Println(t.Time)
// 	// return t.Time.String(), nil
// 	return src.Timestamptz.Value()
// }

// func (src Timestamptz) Value() (driver.Value, error) {
// 	switch src.Status {
// 	case Present:
// 		if src.InfinityModifier != None {
// 			return src.InfinityModifier.String(), nil
// 		}
// 		return src.Time, nil
// 	case Null:
// 		return nil, nil
// 	default:
// 		return nil, errUndefined
// 	}
// }
