package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgtype"
)

func main() {
	t := time.Now().UTC()
	var snt NullTime
	snt.Time = t
	snt.Valid = true
	sal := ScheduleAuditLog{
		StaffFirstName:            stringPointer("MyFirstName"),
		TaskKey:                   stringPointer("MyTaskKey"),
		UserLastName:              stringPointer("MyLastName"),
		ScheduleEntryStartTimeUTC: &t,
		ScheduleEntryEndTimeUTC:   &snt,
	}
	fn := FieldNames(sal)
	fmt.Println(fn)
	fmt.Println([]byte(""))
	fmt.Printf("is it nil: %t\n", []byte("") == nil)
	var x []byte
	x = nil
	fmt.Println(x)
	fmt.Printf("is it nil: %t\n", x == nil)
	fmt.Println(nil)

	b, err := json.MarshalIndent(sal, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))

	b, err = json.MarshalIndent(snt, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))

	var sntt NullTime
	b, err = json.MarshalIndent(sntt, "", "\t")
	if err != nil {
		log.Fatalln(err)

	}
	fmt.Println(string(b))

	// let's check custom marshalling on value type vs pointer type
	tc := TimeContainer{
		NullTimePointer: &snt,
	}
	b, err = json.MarshalIndent(tc, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
	tc.NullTimePointer = nil
	b, err = json.MarshalIndent(tc, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))

	// regexp
	ts := "00:03:00"
	tf, err := regexp.MatchString(`^\d{2}:\d{2}:\d{2}$`, ts)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%t\n", tf)
	fmt.Println(time.Kitchen)

	tmt, err := time.Parse(`15:04:05`, ts)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(tmt)
	fmt.Println(strings.ToUpper("03:04pm"))

	// date manipulation
	now := time.Now().UTC()
	nowDate := now.Truncate(time.Hour * 24)
	todd := now.Sub(nowDate)
	tod := time.Time{}
	tod = tod.Add(todd)
	fmt.Println(tod)

	pgtod := TimeOfDay{}
	pgtod.Set(time.Now().UTC())

	b, err = pgtod.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(string(b))
	// fmt.Println(pgtod)

	// fmt.Println(pgtod)
	// todString, err := pgtod.Value()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Printf("%s\n", todString)

	// b, err = json.Marshal(pgtod)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(string(b))
	// fmt.Println(len("03"))

	// tfReg := regexp.MustCompile(`(?i)^(?P<hour>\d{1,2})(?P<hmsep>[^\d]*)(?P<minute>\d{1,2})(?P<mssep>[^\d]*)(?P<second>\d{1,2})(?P<subsecond>\.\d*)?(?P<markprefix>.*)(?P<mark>[ap]m)?$`)
	// hourLength := len(tfReg.ReplaceAllString(`12:30:19`, `${hour}`))
	// fmt.Println(hourLength)
	// fmt.Println(tfReg.ReplaceAllString(`12:30:19AM`, `${hour}`))
	// if err := pgtod.Set("09:03:01"); err != nil {
	// 	log.Println(err)
	// }

	// fmt.Println(pgtod.String())
	// fmt.Println(pgtod)

	// var pgtodnull TimeOfDay
	// fmt.Println(pgtodnull)
	// xyz := fmt.Sprintf(`"%s"`, pgtod)
	// fmt.Println(xyz)
	// var todTest TODTest
	// todTest.TOD.Set(time.Now().UTC())
	// todJSON, err := json.Marshal(todTest)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(string(todJSON))
	// q, err := GuessTimeOfDayFormat(`12:17:25.00402 AM`)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(q)
	// fmt.Println(time.Kitchen)
	// q, err = GuessTimeOfDayFormat(time.Kitchen)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(q)

	samples := []string{
		`12:17:25.00402 AM`,
		time.Kitchen,
		`4:45:23`,
		`00:03:23`,
	}
	for _, v := range samples {
		q, err := GuessTimeOfDayFormat(v)
		if err != nil {
			log.Println(err)
		}
		// fmt.Println(q)
		tod, err := ParseTOD(v)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("%25s%25s%25s\n", v, q, tod)
	}
	fmt.Println(TimeOfDay{})
}

type TODTest struct {
	TOD TimeOfDay
}

func GuessTimeOfDayFormat(ts string) (string, error) {

	// attempting to handle a small variety of time of day inputs
	tfr := regexp.MustCompile(`(?i)^(?P<hour>\d{1,2})(?P<hmsep>[^\d]*)(?P<minute>\d{1,2})((?P<mssep>[^\d]*)(?P<second>\d{1,2})(?P<subsecond>\.\d+)?)?((?P<markprefix>[^ap]*)(?P<mark>[ap]m))?$`)
	matches := tfr.FindStringSubmatch(ts)
	// fmt.Printf("Matches: %q\n", matches)
	// fmt.Printf("Number of matches: %d\n", len(matches))
	// fmt.Printf("%q\n", matches[])
	// fmt.Printf("Hour: %s\n", matches[tfr.SubexpIndex("hour")])
	components := map[string]string{}
	for _, v := range tfr.SubexpNames() {
		subMatchIndex := tfr.SubexpIndex(v)
		if subMatchIndex > 0 && subMatchIndex <= len(tfr.SubexpNames()) {
			components[v] = matches[subMatchIndex]
			// fmt.Printf("%s:\t%s\n", v, matches[subMatchIndex])

		}
	}
	// fmt.Println(components)
	if len(components["mark"]) == 2 {
		components["mark"] = "PM"
		hourLength := fmt.Sprint(len(components["hour"]))
		hourFormat := `%0` + hourLength + `d`
		components["hour"] = fmt.Sprintf(hourFormat, 3)
		// fmt.Println(components["hour"])
	} else {
		components["hour"] = "15"
	}
	minuteFormat := `%0` + fmt.Sprint(len(components["minute"])) + `d`
	components["minute"] = fmt.Sprintf(minuteFormat, 4)
	if len(components["second"]) > 0 {
		// fmt.Printf("Length of second component: %d\n", len(components["second"]))
		secondFormat := `%0` + fmt.Sprint(len(components["second"])) + `d`
		components["second"] = fmt.Sprintf(secondFormat, 5)

	}
	components["subsecond"] = regexp.
		MustCompile(`\d`).
		ReplaceAllLiteralString(components["subsecond"], "0")

	// `.` + strings.Repeat("0", len(components["subsecond"]))

	format := ""
	for _, v := range tfr.SubexpNames() {
		// fmt.Printf("%s:\t%s\n", v, components[v])
		format = format + components[v]
	}
	// fmt.Println(format)
	return format, nil
	// fmt.Printf("%s%s%s\n", components["hour"], components["hmsep"], components["minute"])
	// fmt.Println(tfRegex.ReplaceAllString(`12:30:19AM`, `${hour}`))
	// ampmMark := regexp.MustCompile(``)

	// Try reasonable 'wall clock' time formats
	// there's probably an easier way to do this...
	// var timeFormat string
	// if tf, _ := regexp.MatchString(`^\d{2}:\d{2}:\d{2}$`, ts); tf {
	// 	timeFormat = `15:04:05`
	// } else if tf, _ := regexp.MatchString(`^\d{2}:\d{2}$`, ts); tf {
	// 	timeFormat = `15:04`
	// } else if tf, _ := regexp.MatchString(`(?i)^\d{2}:\d{2}:\d{2}[AP]M$`, ts); tf {
	// 	timeFormat = `03:04:05PM`
	// } else if tf, _ := regexp.MatchString(`(?i)^\d{2}:\d{2}:\d{2} [AP]M$`, ts); tf {
	// 	timeFormat = `03:04:05 PM`
	// } else if tf, _ := regexp.MatchString(`(?i)^\d{2}:\d{2}\.\d{2}[AP]M$`, ts); tf {
	// 	timeFormat = `03:04.05PM`
	// } else if tf, _ := regexp.MatchString(`(?i)^\d{2}:\d{2}\.\d{2} [AP]M$`, ts); tf {
	// 	timeFormat = `03:04.05 PM`
	// } else if tf, _ := regexp.MatchString(`(?i)^\d{2}:\d{2}[AP]M$`, ts); tf {
	// 	timeFormat = `03:04PM`
	// } else if tf, _ := regexp.MatchString(`(?i)^\d{2}:\d{2} [AP]M$`, ts); tf {
	// 	timeFormat = `03:04 PM`
	// } else if tf, _ := regexp.MatchString(`(?i)^\d:\d{2}:\d{2}[AP]M$`, ts); tf {
	// 	timeFormat = `3:04:05PM`
	// } else if tf, _ := regexp.MatchString(`(?i)^\d:\d{2}:\d{2} [AP]M$`, ts); tf {
	// 	timeFormat = `3:04:05 PM`
	// } else if tf, _ := regexp.MatchString(`(?i)^\d:\d{2}\.\d{2}[AP]M$`, ts); tf {
	// 	timeFormat = `3:04.05PM`
	// } else if tf, _ := regexp.MatchString(`(?i)^\d:\d{2}\.\d{2} [AP]M$`, ts); tf {
	// 	timeFormat = `3:04.05 PM`
	// } else if tf, _ := regexp.MatchString(`(?i)^\d:\d{2}[AP]M$`, ts); tf {
	// 	timeFormat = `3:04PM`
	// } else if tf, _ := regexp.MatchString(`(?i)^\d:\d{2} [AP]M$`, ts); tf {
	// 	timeFormat = `3:04 PM`
	// } else {
	// 	// just give up
	// 	return "", errors.New("Unable to guess TimeOfDay format.")
	// }
	// return timeFormat, nil
}

func ParseTODWithLayout(layout string, ts string) (TimeOfDay, error) {
	tod := TimeOfDay{}
	t, err := time.Parse(layout, ts)
	if err != nil {
		return tod, err
	}
	tod.Set(t)
	return tod, nil
}

func ParseTOD(s string) (TimeOfDay, error) {
	// tod := TimeOfDay{}
	timeFormat, err := GuessTimeOfDayFormat(s)
	if err != nil {
		return TimeOfDay{}, err
	}
	tod, err := ParseTODWithLayout(timeFormat, s)
	// t, err := time.Parse(timeFormat, s)
	if err != nil {
		return tod, err
	}
	// tod.Set(t)
	return tod, nil
}

type TimeOfDay struct {
	pgtype.Time
}

func (tod TimeOfDay) String() string {
	v, _ := tod.Time.Value()
	return fmt.Sprint(v)
}

func (tod *TimeOfDay) UnmarshalJSON(data []byte) error {
	if err := tod.Set(data); err != nil {
		return err
	}
	return nil
}

func (tod TimeOfDay) MarshalJSON() ([]byte, error) {
	v, err := tod.Time.Value()
	if err != nil {
		return nil, err
	}
	b := []byte(`"` + fmt.Sprint(v) + `"`)
	return b, nil
}

// v, err := tod.Value()
// if err != nil {
// 	return err
// }
func stringPointer(s string) *string {
	return &s
}

type ScheduleAuditLog struct {
	StaffFirstName            *string       `json:"StaffFirstName,omitempty"`
	StaffLastName             *string       `json:"StaffLastName,omitempty"`
	StaffAbbreviation         *string       `json:"StaffAbbreviation,omitempty"`
	StaffKey                  *string       `json:"StaffKey,omitempty"`
	TaskName                  *string       `json:"TaskName,omitempty"`
	TaskAbbreviation          *string       `json:"TaskAbbreviation,omitempty"`
	TaskKey                   *string       `json:"TaskKey,omitempty"`
	ScheduleEntryDate         *string       `json:"ScheduleEntryDate,omitempty"`
	ScheduleEntryStartTimeUTC *time.Time    `json:"ScheduleEntryStartTimeUTC,omitempty"`
	ScheduleEntryStartTime    *string       `json:"ScheduleEntryStartTime,omitempty"`
	ScheduleEntryEndTimeUTC   *NullTime     `json:"ScheduleEntryEndTimeUTC,omitempty"`
	ScheduleEntryEndTime      *string       `json:"ScheduleEntryEndTime,omitempty"`
	ScheduleEntryKey          *string       `json:"ScheduleEntryKey,omitempty"`
	ActivityType              *string       `json:"ActivityType,omitempty"`
	SourceType                *string       `json:"SourceType,omitempty"`
	UserFirstName             *string       `json:"UserFirstName,omitempty"`
	UserLastName              *string       `json:"UserLastName,omitempty"`
	UserKey                   *string       `json:"UserKey,omitempty"`
	TimestampUTC              *string       `json:"TimestampUTC,omitempty"`
	Timestamp                 *string       `json:"Timestamp,omitempty"`
	AdditionalInformation     *string       `json:"AdditionalInformation,omitempty"`
	Locations                 []interface{} `json:"Locations,omitempty"`
	IPAddress                 *string       `json:"IPAddress,omitempty"`
}

func FieldNames(a any) []string {
	allowableFiels := []string{
		"ActivityType",
		"UserFirstName",
		"TaskAbbreviation",
		"ScheduleEntryStartTimeUTC",
		"ScheduleEntryEndTimeUTC",
	}

	afMap := map[string]struct{}{}
	for _, v := range allowableFiels {
		afMap[v] = struct{}{}
	}

	v := reflect.ValueOf(a)
	fmt.Printf("ValueOf: %s\n", v)
	tt := reflect.TypeOf(a)
	fmt.Printf("TypeOf: %s\n", tt)

	fmt.Printf("ValueOf.Type.Name: %s\n", v.Type().Name())
	fmt.Printf("ValueOf.Type.Kind: %s\n", v.Type().Kind())
	fmt.Printf("ValueOf.Type.NumField: %d\n", v.Type().NumField())
	fmt.Printf("ValueOf.Type.NumField: %v\n", v.Type().Field(0))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		// fv := reflect.ValueOf(f)
		if _, ok := afMap[f.Name]; ok {
			fmt.Printf("%2d:\t%s\t%t\n", i, f.Name, v.Field(i).IsNil())
		}
	}

	return []string{}
}

type NullTime struct {
	sql.NullTime
}

// MarshalJSON satisfies the json.Marshaler interface
func (t NullTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() || !t.Valid {
		return []byte("null"), nil
		// return nil, nil
	}
	return []byte(`"` + t.Time.Format("15:04:05") + `"`), nil
}

type TimeContainer struct {
	NullTime           NullTime  `json:"NullTime,omitempty"`
	NullTimePointer    *NullTime `json:"NullTimePointer"`
	SQLNullTimePointer *sql.NullTime
}
