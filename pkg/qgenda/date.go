package qgenda

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/jackc/pgtype"
)

// Date wraps jackc/pgtype's Date, adds JSON un/marshaling support,
// and allows string customization using Set/Layout(). By default, it
// supports parsing a variety of date style strings and will encode to
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
	components := map[string]string{}
	for _, v := range dfr.SubexpNames() {

		subMatchIndex := dfr.SubexpIndex(v)
		if subMatchIndex > 0 && subMatchIndex <= len(dfr.SubexpNames()) {
			components[v] = matches[subMatchIndex]

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
		return time.Time{}.Format("2006-01-02")
	}

	layout := d.layout
	if layout == "" {
		layout = "2006-01-02"
	}

	return d.Time.Format(layout)
}

func (d Date) MarshalJSON() ([]byte, error) {
	// switch d.Status {
	// case pgtype.Null:
	// 	return []byte("null"), nil
	// case pgtype.Undefined:
	// 	//errors.New("Cannot Date.MarshalJSON due to Date.Status == pgtype.Undefined")
	// 	return []byte("null"), nil
	// }
	// if d.Time.Year() > 2300 {
	// 	log.Println(d, d.Status)

	// }
	if d.Status != pgtype.Present || d.Date.Time.IsZero() {
		// log.Println(d, d.Status)
		// errors.New("Cannot Date.MarshalJSON due to Date.Status != pgtype.Present")
		return []byte("null"), nil
	}

	var s string
	switch d.InfinityModifier {
	case pgtype.None:
		s = d.String()
	case pgtype.Infinity:
		log.Println("Infinity? ", d)
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
	// pd, err := ProcessDate(*d)
	// if err != nil {
	// 	log.Println(err)
	// }
	// *d = pd
	err = d.Process()
	if err != nil {
		log.Println(err)
	}
	return nil
}

// ProcessDate is the underlying function used to handle any 'in-flight' processing
// before loading to a destination. It should be wrapped to satisfy any Processor type
// interfaces.
func ProcessDate(d Date) (Date, error) {
	out := &Date{}
	status := d.Status
	// truncate
	yr, mo, dt := d.Date.Time.Date()
	t := time.Date(yr, mo, dt, 0, 0, 0, 0, time.UTC)
	// if yr > 3000 {
	// 	fmt.Println(t)
	// }
	if yr < 1800 || yr > 2099 {
		t = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
		status = pgtype.Null
	}
	if err := out.Set(t); err != nil {
		log.Println(err)
	}
	out.Status = status
	return *out, nil

}

// Process satisfies the Processor interface
func (d *Date) Process() error {
	dt, err := ProcessDate(*d)
	*d = dt
	return err
}
