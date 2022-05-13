package qgenda

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/jackc/pgtype"
)

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

// ProcessTimeOfDay is the underlying function used to handle any 'in-flight' processing
// before loading to a destinatino. It should be wrapped to satisfy any Processor type
// interfaces.
func ProcessTimeOfDay(t TimeOfDay) (TimeOfDay, error) {
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

	return t, nil
}

// ProcessValue satisfies the ValueProcessor interface
func (t TimeOfDay) ProcessValue() (TimeOfDay, error) {
	return ProcessTimeOfDay(t)
}

// Process satisfies the Processor interface
func (t *TimeOfDay) Process() error {
	tod, err := ProcessTimeOfDay(*t)
	*t = tod
	return err
}
