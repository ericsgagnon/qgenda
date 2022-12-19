package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	_ "time/tzdata"

	_ "github.com/lib/pq"

	"github.com/ericsgagnon/qgenda/pkg/qgenda"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	ctx := context.Background()
	var sch []qgenda.Schedule
	b, err := os.ReadFile("../.cache/rawschedule.json")
	if err != nil {
		log.Println(err)
	}
	if err := json.Unmarshal(b, &sch); err != nil {
		log.Println(err)
	}
	if err := qgenda.ProcessRecursively(sch); err != nil {
		log.Println(err)
	}
	jsonOut, err := json.MarshalIndent(sch, "", "\t")
	if err != nil {
		log.Println(err)
	}
	os.WriteFile("schedule-structured.json", jsonOut, 0644)

	// var staff []qgenda.StaffMember
	// b, err = os.ReadFile("../.cache/rawstaffmember.json")
	// if err != nil {
	// 	log.Println(err)
	// }
	// if err := json.Unmarshal(b, &staff); err != nil {
	// 	log.Println(err)
	// }
	// if err := qgenda.ProcessRecursively(staff); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(qgenda.CreateTableSQL[qgenda.Schedule]("", "schedule"))

	// dev pg inserts
	// db, err := sqlx.Open("pgx", os.Getenv("PG_CONNECTION_STRING"))
	db, err := sqlx.Open("postgres", os.Getenv("PG_CONNECTION_STRING"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	// fmt.Println(qgenda.PGCreateTableStatement(qgenda.Schedule{}, "", "schedule"))
	// fmt.Println(qgenda.PGInsertStatement(qgenda.Schedule{}, "", "schedule"))
	if dbDropResult, err := db.ExecContext(context.Background(), "DROP TABLE IF EXISTS schedule"); err != nil {
		fmt.Println(dbDropResult, err)
	}
	if dbDropResult, err := db.ExecContext(context.Background(), "DROP TABLE IF EXISTS schedulestafftag"); err != nil {
		fmt.Println(dbDropResult, err)
	}
	if dbDropResult, err := db.ExecContext(context.Background(), "DROP TABLE IF EXISTS scheduletasktag"); err != nil {
		fmt.Println(dbDropResult, err)
	}

	// pipelineResult, err := qgenda.LoadSchedulesToPG(ctx, db, sch, "", "schedule")
	// if err != nil {
	// 	log.Println(err)
	// }
	// prra, err := pipelineResult.RowsAffected()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("%d", prra)

	// dbResult, err := db.NamedExecContext(
	// 	context.Background(),
	// 	qgenda.PGCreateTableStatement(qgenda.Schedule{}, "", "schedule"),
	// 	sch,
	// )
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(dbResult)
	// dbResult, err = db.NamedExecContext(
	// 	context.Background(),
	// 	qgenda.PGInsertStatement(qgenda.Schedule{}, "", "schedule"),
	// 	sch,
	// )
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(dbResult)
	// fmt.Println(qgenda.PGCreateTableStatement(qgenda.ScheduleTag{}, "", "scheduletag"))
	// fmt.Println(qgenda.PGInsertStatement(qgenda.ScheduleTag{}, "", "scheduletag"))
	// fmt.Println(qgenda.PGQueryConstraintsStatement(qgenda.Schedule{}, "", "schedule"))
	rqf := qgenda.RequestConfig{}
	if err := db.Get(
		&rqf,
		qgenda.PGQueryConstraintsStatement(qgenda.Schedule{}, "", "schedule"),
	); err != nil {
		log.Println(err)
	}
	fmt.Println("---------------------------------------------------------------------")
	fmt.Printf("%s\n", *rqf.SinceModifiedTimestamp)
	sr := qgenda.NewScheduleRequest(&rqf)
	fmt.Println("---------------------------------------------------------------------")
	fmt.Printf("%s\n", sr.Encode())
	fmt.Println("---------------------------------------------------------------------")
	// sr.SetEndDate(time.Time{})
	// sr.StartDate = nil
	// sr.EndDate = nil

	sr.SetStartDate(time.Now().UTC().Add(-1 * 24 * 15 * time.Hour))
	sr.SetEndDate(time.Now().UTC().Add(-1 * 24 * 1 * time.Hour))

	qcc := &qgenda.ClientConfig{
		Email:    os.Getenv("QGENDA_EMAIL"),
		Password: os.Getenv("QGENDA_PASSWORD"),
	}
	c, err := qgenda.NewClient(qcc)
	if err != nil {
		log.Fatalln(err)
	}
	c.Auth()

	schedules := qgenda.Schedules{}
	// result, err := schedules.DropTable(ctx, db, "qgenda", "schedule")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	sd := time.Now().UTC().Add(-1 * 24 * 15 * time.Hour)
	ed := time.Now().UTC().Add(-1 * 24 * 1 * time.Hour)
	ts := time.Now().UTC().Add(-1 * 24 * 90 * time.Hour)
	rqf.StartDate = &sd
	rqf.EndDate = &ed
	rqf.SinceModifiedTimestamp = &ts
	// fmt.Println(rqf.Parse().Encode())

	schedules, err = schedules.Extract(ctx, c, &rqf)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("length of schedules: %d\n", len(schedules))
	if err != nil {
		log.Println(err)
	}
	schedules, err = schedules.Process()
	if err != nil {
		log.Println(err)
	}
	// var b []byte
	// bout, err := yaml.Marshal(schedules[0])
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(string(bout))

	schResult, err := schedules.InsertRows(ctx, db, "", "schedule")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(schResult)

	jsonOut, err = json.MarshalIndent(schedules, "", "\t")
	if err != nil {
		log.Println(err)
	}
	os.WriteFile("../.cache/schedules-type.json", jsonOut, 0644)

	// fmt.Printf("Length of schedules: %d\n", len(schedules))

	// updatePipelineResult, err := qgenda.LoadSchedulesToPG(ctx, db, schedules, "", "schedule")
	// if err != nil {
	// 	log.Println(err)
	// }
	// uprra, err := updatePipelineResult.RowsAffected()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("%d\n", uprra)

	// fmt.Println(sr.ToHTTPRequest().URL.String())
	// resp, err := c.Do(ctx, sr)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println("after response ##########")
	// for k, v := range resp.Header {
	// 	fmt.Printf("%s %-80s\n", k, "-")
	// 	for vi, vv := range v {
	// 		fmt.Printf("\t%3d: %40s\n", vi, vv)
	// 	}
	// }
	// respOut, err := yaml.Marshal(resp)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(string(respOut))

	// fmt.Println(resp.Status)
	// resp.Header.Get(http.CanonicalHeaderKey("Date"))

	// if err != nil {
	// 	log.Println(err)
	// }
	// data, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println("## Data ##########################################")
	// fmt.Println(string(data))
	// fmt.Println("##################################################")
	// data2 := *&data
	// var schTest []qgenda.Schedule
	// if err := json.Unmarshal(data, &schTest); err != nil {
	// 	log.Println(err)
	// }
	// qgenda.Process(schTest)

	// fmt.Println(db.DriverName())

	// ds := map[string]any{
	// 	"schedule":    qgenda.Schedule{},
	// 	"staffmember": qgenda.StaffMember{},
	// }
	// for k, v := range ds {
	// 	fmt.Printf("%s\t%T\n", k, v)
	// }

	// pingResult, err := sqlx.Connect()
	// process data

	// out, err := yaml.Marshal(schTest)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(string(out))

	// os.WriteFile("scheduleTest.yaml", out, 0644)
	// dbInsertRowsResult, err := qgenda.PGInsertRows(ctx, db, schTest, "", "schedule")
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(dbInsertRowsResult)

	// fmt.Println(sch[0])
	// var greeting string
	// err = db.QueryRow("select 'Hello, world!'").Scan(&greeting)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(greeting)

	// fmt.Println(qgenda.PGInsertRowsStatement[qgenda.Schedule]("", "schedule"))
	// fmt.Println("##########################################################")
	// fmt.Println(qgenda.LetsInsertRows[qgenda.Schedule]("", "schedule"))

	// // pgxpool version
	// dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("PG_CONNECTION_STRING"))
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer dbpool.Close()

	// err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(greeting)
	// results, err := dbpool.Exec(
	// 	context.Background(),
	// 	qgenda.CreateTableSQL[qgenda.Schedule]("", "schedule"),
	// )
	// if err != nil {
	// 	log.Println("Error creating table")
	// }
	// fmt.Println(results)
	// fmt.Println(qgenda.LetsInsertRows("", "schedule", sch))
	// fmt.Println(len(sch))
	// copycount, err := dbpool.CopyFrom(
	// 	context.Background(),
	// 	pgx.Identifier{"schedule"},
	// 	[]string{"schedulekey", "lastmodifieddateutc", "compkey", "staffemail", "taskname", "date", "endtime"},
	// 	pgx.CopyFromSlice(
	// 		len(sch),
	// 		// 5,
	// 		func(i int) ([]interface{}, error) {
	// 			d := sch[i].LastModifiedDateUTC
	// 			return []interface{}{
	// 					sch[i].ScheduleKey,
	// 					d,
	// 					// time.Now().UTC(),
	// 					// sch[i].Date,
	// 					sch[i].CompKey,
	// 					sch[i].StaffEmail,
	// 					sch[i].TaskName,
	// 					// "raw_message",
	// 					// time.Now().UTC(),
	// 					// time.Now().UTC(),
	// 					// sch[i].ScheduleKey,
	// 					// sch[i].CallRole,
	// 					// sch[i].CompKey,
	// 					// sch[i].Credit,
	// 					sch[i].Date,
	// 					// sch[i].StartDateUTC,
	// 					// sch[i].EndDateUTC,
	// 					sch[i].EndTime,
	// 					// sch[i].IsCred,
	// 					// sch[i].IsPublished,
	// 					// sch[i].IsLocked,
	// 					// sch[i].IsStruck,
	// 					// sch[i].Notes,
	// 					// sch[i].IsNotePrivate,
	// 					// sch[i].StaffAbbrev,
	// 					// sch[i].StaffBillSysId,
	// 					// sch[i].StaffEmail,
	// 					// sch[i].StaffEmrId,
	// 					// sch[i].StaffErpId,
	// 					// sch[i].StaffInternalId,
	// 					// sch[i].StaffExtCallSysId,
	// 					// sch[i].StaffFName,
	// 					// sch[i].StaffLName,
	// 					// sch[i].StaffMobilePhone,
	// 					// sch[i].StaffNpi,
	// 					// sch[i].StaffPager,
	// 					// sch[i].StaffPayrollId,
	// 					// sch[i].StartDate,
	// 					// sch[i].StartTime,
	// 					// sch[i].TaskAbbrev,
	// 					// sch[i].TaskBillSysId,
	// 					// sch[i].TaskContactInformation,
	// 					// sch[i].TaskExtCallSysId,
	// 					// sch[i].TaskId,
	// 					// sch[i].TaskKey,
	// 					// sch[i].TaskName,
	// 					// sch[i].TaskPayrollId,
	// 					// sch[i].TaskEmrId,
	// 					// sch[i].TaskCallPriority,
	// 					// sch[i].TaskDepartmentId,
	// 					// sch[i].TaskIsPrintEnd,
	// 					// sch[i].TaskIsPrintStart,
	// 					// sch[i].TaskShiftKey,
	// 					// sch[i].TaskType,
	// 					// sch[i].LocationName,
	// 					// sch[i].LocationAbbrev,
	// 					// sch[i].LocationID,
	// 					// sch[i].LocationAddress,
	// 					// sch[i].TimeZone,
	// 					// sch[i].LastModifiedDateUTC,
	// 					// sch[i].IsRotationTask,
	// 				},
	// 				nil
	// 		},
	// 	),
	// )
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(copycount)
	// fmt.Println(*sch[1].CompKey)
	// type User struct {
	// 	FirstName string
	// 	LastName  string
	// 	Age       int
	// }
	// rows := []User{
	// 	{"John", "Smith", 36},
	// 	{"Jane", "Doe", 29},
	// }

	// copyCount, err := dbpool.CopyFrom(
	// 	context.Background(),
	// 	pgx.Identifier{"people"},
	// 	[]string{"first_name", "last_name", "age"},
	// 	pgx.CopyFromSlice(len(rows), func(i int) ([]interface{}, error) {
	// 		return []interface{}{rows[i].FirstName, rows[i].LastName, rows[i].Age}, nil
	// 	}),
	// )
	// fmt.Println(err)
	// fmt.Println(copyCount)

	// copyCount, err = dbpool.CopyFrom(
	// 	context.Background(),
	// 	pgx.Identifier{"people"},
	// 	[]string{"first_name", "last_name"},
	// 	pgx.CopyFromSlice(len(rows), func(i int) ([]interface{}, error) {
	// 		return []interface{}{rows[i].FirstName, rows[i].LastName}, nil
	// 	}),
	// )
	// fmt.Println(err)
	// fmt.Println(copyCount)

	// type Person struct {
	// 	FirstName string `db:"first_name"`
	// 	LastName  string `db:"last_name"`
	// 	Email     string
	// }
	// personStructs := []Person{
	// 	{FirstName: "Ardie", LastName: "Savea", Email: "asavea@ab.co.nz"},
	// 	{FirstName: "Sonny Bill", LastName: "Williams", Email: "sbw@ab.co.nz"},
	// 	{FirstName: "Ngani", LastName: "Laumape", Email: "nlaumape@ab.co.nz"},
	// }

	// result, err := db.NamedExec(`CREATE TABLE IF NOT EXISTS person ( first_name text, last_name text, email text )`, personStructs)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(result)
	// result, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
	//     VALUES (:first_name, :last_name, :email)`, personStructs)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(result)
	// fmt.Println(qgenda.PGCreateScheduleTableStatement("", ""))
	// fields := qgenda.StructToFields(qgenda.Schedule{})
	// for _, field := range fields {
	// 	fmt.Println(field)
	// }
	// out, err := yaml.Marshal(fields)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(string(out))

	// fmt.Println(qgenda.CreateTableSQL[qgenda.TagCategory]("", "scheduletasktags"))
	// fmt.Println(qgenda.SQLTest[qgenda.ScheduleTag]("", "scheduletasktags"))
	// fmt.Println(qgenda.CreateScheduleTagTableSQL("", ""))
	// fmt.Println(qgenda.PGInsertRowsStatement[qgenda.Schedule]("", ""))
	// fmt.Printf(nada)
	// for _, v := range staff {
	// 	for _, u := range (&v).TTCMTags {
	// 		if err := (&u).Process(); err != nil {
	// 			log.Println(err)
	// 		}
	// 	}
	// }
	// jsonOut, err = json.MarshalIndent(staff, "", "\t")
	// if err != nil {
	// 	log.Println(err)
	// }
	// os.WriteFile("staffmember-structured.json", jsonOut, 0644)

	// dJSON := `"9999-12-31"`
	// d := qgenda.Date{}
	// if err := json.Unmarshal([]byte(dJSON), &d); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("%#v\n", d)
	// fmt.Printf("%t\n", d.Date.Time.Year() > 2100)
	// if err := (&d).Process(); err != nil {
	// 	log.Println(err)
	// }
	// sc := sch[0]
	// sf := qgenda.StructFields(sc)
	// for _, v := range sf {
	// 	fmt.Printf("%s\n", v)
	// }

	// for k, v := range qgenda.GoToPGTypeMap {
	// 	fmt.Printf("%s: %s\n", k, v)
	// }

	// v := reflect.ValueOf(sc)
	// iv := reflect.Indirect(v)
	// fields := qgenda.StructFields(iv)
	// sqlFieldDefs := []string{}
	// for i := 0; i < iv.NumField(); i++ {
	// 	stf := fields[i]
	// 	// // fmt.Printf("ProcessStruct Field %d %s\n", i, sf.Name)
	// 	// fv := iv.Field(i)
	// 	// k := fv.Kind()
	// 	// if k == reflect.Pointer {
	// 	// 	t := reflect.TypeOf(fv.Interface()).Elem()
	// 	// 	fmt.Printf("typeof: %s\n", t)
	// 	// }
	// 	// fiv := reflect.Indirect(fv)
	// 	// fmt.Printf("Field Indirect: %s\n", fiv.Kind())
	// 	// fmt.Printf("%s:\t%s\t%s\t%s\n", stf.Name, stf.Type, stf.Tag, stf.Type.Kind())
	// 	// fmt.Printf("%s\n", fv)
	// 	dt := qgenda.DynamicType(iv.Field(i))
	// 	switch dt.Kind() {
	// 	case reflect.Slice:
	// 		continue
	// 	case reflect.Map:
	// 		continue
	// 	case reflect.Array:
	// 		continue
	// 	case reflect.Chan:
	// 		continue
	// 	}
	// 	// vt := reflect.Type{}
	// 	// if dt.Kind() == reflect.Slice {
	// 	// 	vt := dt.Elem()
	// 	// 	fmt.Printf("Kind: %-8s Type: %-25s ValueType: %-25s ValueKind: %-25s PGType: %-20s\n", dt.Kind(), dt, vt, vt.Kind(), qgenda.GoToPGTypeMap[dt.Name()])
	// 	// }
	// 	// fmt.Printf("Kind: %-8s Type: %-99s PGType: %-20s\n", dt.Kind(), dt, qgenda.GoToPGTypeMap[dt.Name()])
	// 	pgnm := strings.ToLower(stf.Name)
	// 	if stf.Tag.Get("sqlname") != "" {
	// 		pgnm = stf.Tag.Get("sqlname")
	// 	}
	// 	pgtp := qgenda.GoToPGTypeMap[dt.Name()]
	// 	pgcnstrnt := stf.Tag.Get("constraint")
	// 	sqlFieldDefs = append(sqlFieldDefs, fmt.Sprintf("%s %s %s", pgnm, pgtp, pgcnstrnt))

	// 	// stf.Tag
	// 	// stf.Name
	// 	// stf.Type
	// }
	// fmt.Println(strings.Join(sqlFieldDefs, ",\n"))
	// fmt.Println(qgenda.CreateScheduleTable(nil))

	// od, err := qgenda.ProcessDate(d)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("%#v\n", od)
	// fmt.Printf("%s\n", d)
	// fmt.Printf("%t\n", d.Date.Time.Year() > 2100)
	// tt := staff[2]
	// fmt.Println(tt.TTCMTags[0])
	// if err := qgenda.Process(&tt); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(tt.TTCMTags[0])

	// // some useful samples
	// todStrings := []string{
	// 	`00:05:00`,
	// 	`12:15:00`,
	// 	`00:03:00`,
	// 	`08:33:00`,
	// 	`22:59:00`,
	// 	`00:03:00`,
	// 	`03:05:00`,
	// 	`00:03:00`,
	// }
	// todMap := map[string]*qgenda.TimeOfDay{}
	// todMapValues := map[string]qgenda.TimeOfDay{}
	// todSlice := []*qgenda.TimeOfDay{}
	// todSliceValues := []qgenda.TimeOfDay{}
	// todStructSlice := []struct{ Time qgenda.TimeOfDay }{}
	// for _, v := range todStrings {
	// 	tod, err := qgenda.ParseTimeOfDay(v)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	todSlice = append(todSlice, &tod)
	// 	todSliceValues = append(todSliceValues, tod)
	// 	todMap[v] = &tod
	// 	todMapValues[v] = *&tod
	// 	todStructSlice = append(todStructSlice,
	// 		struct{ Time qgenda.TimeOfDay }{Time: tod})
	// }
	// todSlice2 := *(&todSlice)
	// todSliceValues2 := *(&todSliceValues)
	// todStruct := struct{ Time qgenda.TimeOfDay }{Time: *todSlice[0]}
	// fmt.Sprint(todStruct)

	// ////////////////////////////////////////////////////
	// b, err := os.ReadFile("schedule.json")
	// if err != nil {
	// 	log.Println(err)
	// }

	// var sch []qgenda.Schedule
	// if err := json.Unmarshal(b, &sch); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println("--------------------------------------------------------------------------------------------------------------------")
	// qgenda.ReflectionInfo(&todSliceValues2[0])
	// qgenda.ReflectionInfo(&sch[0])
	// qgenda.ReflectionInfo(todSliceValues2[0])
	// qgenda.ReflectionInfo(sch[0])
	// // qgenda.ReflectionInfo(&xts)
	// // qgenda.ReflectionInfo(xts)
	// // qgenda.ReflectionInfo(&todSliceValues2)
	// fmt.Println("--------------------------------------------------------------------------------------------------------------------")
	// qgenda.ReflectionInfo(&todSliceValues2)
	// qgenda.ReflectionInfo(&todSlice2)
	// qgenda.ReflectionInfo(todSliceValues2)
	// qgenda.ReflectionInfo(todSlice2)
	// fmt.Println("--------------------------------------------------------------------------------------------------------------------")
	// qgenda.ReflectionInfo(&todMap)
	// qgenda.ReflectionInfo(todMap)
	// fmt.Println("--------------------------------------------------------------------------------------------------------------------")

	// fmt.Println("Pre-Processing Slice:")
	// for i, v := range todSlice {
	// 	fmt.Printf("%2d: %s\n", i, v)
	// }
	// if err := qgenda.Process(todSlice); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println("Post-Processing Slice:")
	// for i, v := range todSlice {
	// 	fmt.Printf("%2d: %s\n", i, v)
	// }
	// if err := sch[0].StartTime.Process(); err != nil {
	// 	log.Println(err)
	// }
	// for i, v := range sch {
	// 	if err := v.StartTime.Process(); err != nil {
	// 		log.Println(err)
	// 	}
	// 	if err := v.EndTime.Process(); err != nil {
	// 		log.Println(err)
	// 	}
	// 	if err := v.StartDate.Process(); err != nil {
	// 		log.Println(err)
	// 	}
	// 	if err := v.EndDate.Process(); err != nil {
	// 		log.Println(err)
	// 	}
	// 	if err := v.Date.Process(); err != nil {
	// 		log.Println(err)
	// 	}
	// 	sch[i] = v

	// 	// fmt.Println(FieldNames(v))
	// 	// fmt.Printf("sch[%d] is a struct: %t\n", i, IsStruct(v))
	// }
	// if err := qgenda.ProcessRecursively(sch); err != nil {
	// 	log.Println(err)
	// }
	// for _, v := range sch {
	// 	if err := qgenda.ProcessRecursively(v); err != nil {
	// 		log.Println(err)
	// 	}
	// }
	// fmt.Printf("Pre-sch[0]:\n%#v\n", sch[0])
	// fmt.Printf("Pre-sch[0].EndTime: %s\n", sch[0].EndTime.String())
	// if err := qgenda.Process(sch[0]); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("Post-sch[0].EndTime: %s\n", sch[0].EndTime.String())
	// tod, err := qgenda.ParseTimeOfDay("00:03:00")
	// if err != nil {
	// 	log.Println(err)
	// }
	// todStruct.Time = tod
	// fmt.Printf("Pre Processed todStruct:\n%s\n", todStruct)
	// if err := qgenda.Process(&todStruct); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("Post Processed todStruct:\n%s\n", todStruct)
	// if err := qgenda.Process(sch); err != nil {
	// 	log.Println(err)
	// }
	// jsonOut, err := json.MarshalIndent(sch, "", "\t")
	// if err != nil {
	// 	log.Println(err)
	// }
	// os.WriteFile("schedule-test.json", jsonOut, 0644)
	// jsonOut, err = json.MarshalIndent(sch[0], "", "\t")
	// if err != nil {
	// 	log.Println(err)
	// }

	// ctx := context.Background()
	// db, err := pgx.Connect(ctx, os.Getenv("PG_CONNECTION_STRING"))
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println("connection to postgres is open: ", !db.IsClosed())
	// defer func(db *pgx.Conn) {
	// 	if err := db.Close(ctx); err != nil {
	// 		log.Println(err)
	// 	}

	// }(db)

	// fmt.Printf("%t\n", qgenda.DefaultTimeLocation == nil)
	// tzNY, err := time.LoadLocation("America/New_York")
	// if err != nil {
	// 	log.Println(err)
	// }
	// // qgenda.DefaultTimeLocation = tzNY
	// // qgenda.DefaultTimeLocation = time.UTC
	// trash, err := time.ParseInLocation(time.RFC3339Nano, "2022-01-21T14:10:46.6126755Z", nil)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("trash: %s\n", trash)
	// trash2, err := qgenda.ParseTimeInLocation(time.RFC3339Nano, "2022-01-21T14:10:46.6126755Z", nil)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("trash2: %s\n", trash2)
	// // qtt.SetLocation(tzNY)
	// timeStrings := []string{
	// 	"1901-01-01T00:00:00",
	// 	"2021-04-30T01:26:23.3569266",
	// 	"9999-12-31T00:00:00",
	// 	"2022-01-21T13:00:00Z",
	// 	"2022-01-21T14:10:46.6126755Z",
	// 	"2022-01-21T13:15:23.0767721Z",
	// }
	// for i, v := range timeStrings {
	// 	fmt.Println("---------------------------------------------------------------")
	// 	t, err := qgenda.ParseTime(v)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	fmt.Println(i, v, t.String(), t.Time.Location(), t.SourceLocation())
	// 	qgenda.DefaultTimeLocation = tzNY
	// 	if err := (&t).Process(); err != nil {
	// 		log.Println(err)
	// 	}
	// 	fmt.Println(i, v, t.String(), t.Time.Location(), t.SourceLocation())
	// 	if err := (&t).ChangeLocation(tzNY, false); err != nil {
	// 		log.Println(err)
	// 	}
	// 	fmt.Println(i, v, t.String(), t.Time.Location(), t.SourceLocation())
	// 	// if err := (&t).ChangeLocation(tzNY, false); err != nil {
	// 	// 	log.Println(err)
	// 	// }
	// 	// fmt.Println(i, v, t.String(), t.Time.Location(), t.Location())
	// 	// if err := (&t).ChangeLocation(tzNY, false); err != nil {
	// 	// 	log.Println(err)
	// 	// }
	// 	// fmt.Println(i, v, t.String(), t.Time.Location(), t.Location())
	// 	// if err := (&t).ChangeLocation(tzNY, false); err != nil {
	// 	// 	log.Println(err)
	// 	// }
	// 	// fmt.Println(i, v, t.String(), t.Time.Location(), t.Location())
	// 	// fmt.Printf("location is nil: %t\n", t.Location() == nil)
	// 	// tt, err := time.Parse("2006-01-02T15:04:05", v)
	// 	// if err != nil {
	// 	// 	log.Println(err)
	// 	// }
	// 	// fmt.Println(i, v, tt)

	// }
	// tlt, err := time.Parse("2006-01-02T15:04:05 MST", "2022-02-22T12:23:17 EST")
	// tlt, err := time.ParseInLocation("02 Jan 06 15:04 MST", "22 Feb 22 12:23 EST", tzNY)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("%s\n", tlt)
	// fmt.Printf("%s\n", tlt.Add(time.Hour*24*7*8))
	// fmt.Printf("%s\n", tlt.In(tzNY))
	// fmt.Printf("%s\n", tlt.UTC())
	// fmt.Printf("%t\n", tlt.Location() == nil)
	// tzEST, err := time.LoadLocation("EST")
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(tzEST)
	// fmt.Println(time.Now().UTC())
	// fmt.Println(time.Date(2020, 01, 01, 00, 05, 00, 00, tzNY))
	// qgenda.ChangeLocation()
	// fmt.Println(string(jsonOut))
	// fmt.Printf("%#v\n", qgenda.GenericTests(sch[0]))
	// fmt.Println(todSlice)
	// tmp, err := qgenda.SliceProcess(todSlice)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(tmp)
	// fmt.Println(todSlice)
	// ttmp, err := qgenda.SliceProcess(testStructSlice)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(ttmp)

	// fmt.Println()
	// qgenda.SliceProcess([]int{1, 2, 3, 4, 5, 6, 7})
	// fmt.Println(todSlice)
	// sv, err := qgenda.ProcessSliceValue(todSlice)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("%s\n", sv)
	// xts := &todSlice2

	// fmt.Println(todMap)
	// fmt.Println(qgenda.ToMap(todMap))
	// out, err := todMap["00:03:00"].ProcessValue()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(out)
	// outm, err := qgenda.ProcessValue(todMap)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(outm)
	// mapRV := reflect.Indirect(reflect.ValueOf(todMap))
	// fmt.Println(mapRV.MapKeys()[0])
	// var x99 interface{} = 35
	// // var x100 = 50
	// switch i := x99.(type) {
	// default:
	// 	fmt.Printf("%#v\n", i)
	// }
	// var todMapi interface{} = todMap
	// mv, ok := todMapi.(map[any]any)
	// if !ok {
	// 	fmt.Println("nope")
	// }
	// fmt.Printf("%T: %v\n", mv, mv)
	// for k, v := range todMap {
	// 	var ki interface{} = k
	// 	var vi interface{} = v
	// 	switch kt := ki.(type) {
	// 	default:
	// 		fmt.Printf("key type: %T\n", kt)
	// 	}
	// 	switch vt := vi.(type) {
	// 	default:
	// 		fmt.Printf("value type: %T\n", vt)
	// 	}
	// 	break
	// }
	// if err := qgenda.Process(todMap); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println("todMap")
	// for k, v := range todMap {
	// 	fmt.Printf("%s: %s\n", k, v)
	// }
	// fmt.Println("---------------------------------------------------------------------------------------------")
	// if err := qgenda.Process(todMapValues); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println("todMapValues")
	// for k, v := range todMapValues {
	// 	fmt.Printf("%s: %s\n", k, v)
	// }
	// fmt.Println("---------------------------------------------------------------------------------------------")
	// sliceTODMapValues := []map[string]qgenda.TimeOfDay{
	// 	todMapValues,
	// }
	// if err := qgenda.Process(&sliceTODMapValues); err != nil {
	// 	log.Println(err)
	// }
	// for _, v := range sliceTODMapValues {
	// 	fmt.Println(v)
	// }
	// sliceTODMap := []map[string]*qgenda.TimeOfDay{
	// 	todMap,
	// }
	// if err := qgenda.Process(sliceTODMap); err != nil {
	// 	log.Println(err)
	// }
	// for _, v := range sliceTODMap {
	// 	fmt.Println(v)
	// }
	// todv := reflect.ValueOf(todMap)

	// mapKeyType :=
	// qgenda.PointerProcesser(&todMap)
	// var it interface{} = todSliceValues2
	// switch x := (it).(type) {
	// case []interface{}:
	// 	fmt.Printf("%T is a slice\n", x)

	// case qgenda.Processor:
	// 	fmt.Printf("%T is a Processor\n", x.(qgenda.Processor))

	// default:
	// 	fmt.Printf("%T\n", x)
	// }
	// fmt.Println(qgenda.IsValueProcessor(qgenda.TimeOfDay{}))
	// fmt.Println(qgenda.IsValueProcessor(&qgenda.TimeOfDay{}))
	// for i, v := range todSliceValues2 {
	// 	tod, err := v.ProcessValue()
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	todSliceValues2[i] = tod
	// }
	// fmt.Printf("todSlice2: %s\n", todSlice2)
	// // fmt.Sprintf("%s\n", todSliceValues2)
	// qgenda.StructFieldProcess(todSliceValues[0])
	// qgenda.StructFieldProcess(map[string]qgenda.TimeOfDay{
	// 	"test": todSliceValues2[0],
	// })
	// out := qgenda.Mappinator(&todSlice)
	// fmt.Printf("%T %s\n", out)
	// out = qgenda.Mappinator(map[string]qgenda.TimeOfDay{})
	// qgenda.SliceTest(todSliceValues2)
	// fmt.Printf("%s\n", todSliceValues2)

	// qgenda.MapTest(todMapValues)
	// for k, v := range todMapValues {
	// 	fmt.Printf("%s: %s\n", k, v)
	// }

	// fmt.Printf("%s\n", todMap)
	// fmt.Println(todMap)
	// todSliceValues := []qgenda.TimeOfDay{}
	// ttest, err := qgenda.AnySlice(todSliceValues)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(qgenda.ProcessSlice(ttest))
	// qgenda.AnySlice(testStruct)
	// qgenda.ReflectionStuff([]qgenda.TimeOfDay{})
	// qgenda.ReflectionStuff([]*qgenda.TimeOfDay{})
	// qgenda.SliceProcess(testStructSlice)
	// fmt.Println(StructFieldNames(sch[0]))
	// fmt.Println(StructFields(sch[0]))
	// fmt.Println(StructFieldByName(sch[0], "LocationAddress"))
	// fmt.Println(StructFieldByName(sch[0], ""))
	// fmt.Println(IsInterface(sch[0], (*qgenda.Schedule)(nil)))
	// fmt.Println(IsInterface(&qgenda.TimeOfDay{}, new(qgenda.Processor)))
	// fmt.Println(IsInterface(new(qgenda.Processor), &qgenda.TimeOfDay{}))
	// fmt.Println(IsInterface(new(qgenda.Processor), sch[0]))
	// fmt.Println(checkBetter[qgenda.Schedule](sch[0]))
	// fmt.Println(ImplementsInterface[qgenda.Processor](&qgenda.TimeOfDay{}))
	// fmt.Println(ImplementsInterface[io.ReadCloser](&qgenda.TimeOfDay{}))
	// fmt.Println(qgenda.Process(sch[0]))
	// fmt.Println(qgenda.Process(&qgenda.TimeOfDay{}))

	// tod, err := qgenda.ParseTimeOfDay(`00:03:00`)
	// // tod.SetLayout(`15:04:05.99999`)
	// if err != nil {
	// 	log.Println(err)
	// }
	// // if err := tod.Process(); err != nil {
	// // 	log.Println(err)
	// // }
	// qgenda.Process(&tod)
	// // for _, v := range todSlice {
	// // 	fmt.Printf("%#v\n", v)
	// // }
	// if err := qgenda.ProcessSlice(todSlice); err != nil {
	// 	log.Println(err)
	// }
	// for _, v := range todSlice {
	// 	fmt.Println(v)
	// }

	// fmt.Printf("%T is a slice: %t\n", todSlice, qgenda.IsSlice(todSlice))
	// fmt.Printf("%T is a struct: %t\n", todSlice, qgenda.IsStruct(todSlice))
	// fmt.Printf("%T is a slice: %t\n", qgenda.IndirectReflectionValue(todSlice), qgenda.IsSlice(qgenda.IndirectReflectionValue(todSlice)))
	// fmt.Printf("%T is a struct: %t\n", qgenda.IndirectReflectionValue(todSlice), qgenda.IsStruct(qgenda.IndirectReflectionValue(todSlice)))
	// fmt.Printf("reflect.ValueOf(reflect.ValueOf)): %t\n", reflect.Indirect(reflect.ValueOf(reflect.ValueOf(reflect.ValueOf(todSlice)))).CanSet())

	// ReflectionStuff(qgenda.Schedule{})
	// ReflectionStuff(&qgenda.Schedule{})
	// ReflectionStuff(qgenda.IndirectReflectionValue(qgenda.Schedule{}))
	// ReflectionStuff(qgenda.IndirectReflectionValue(&qgenda.Schedule{}))
	// ReflectionStuff(reflect.ValueOf(qgenda.Schedule{}))
	// ReflectionStuff(reflect.ValueOf(&qgenda.Schedule{}))
	// ReflectionStuff(reflect.Indirect(reflect.ValueOf(qgenda.Schedule{})))
	// ReflectionStuff(reflect.Indirect(reflect.ValueOf(&qgenda.Schedule{})))
	// ReflectionStuff(&map[any]any{})
	// ReflectionStuff(&map[string]any{})
	// ReflectionStuff(qgenda.IndirectReflectionValue(&map[any]any{}))
	// ReflectionStuff(qgenda.IndirectReflectionValue(&map[string]any{}))
	// fmt.Println(qgenda.IndirectReflectionKind(reflect.ValueOf(map[string]any{})))
	// fmt.Println(qgenda.IndirectReflectionKind(qgenda.Schedule{}))
	// fmt.Println(qgenda.IsKind(&map[string]string{}, "map"))
	// fmt.Println(qgenda.IndirectReflectionValue(&map[string]any{}).CanSet())

	// // if err := qgenda.Process(todMap); err != nil {
	// // 	log.Println(err)
	// // }
	// // for k, v := range todMap {
	// // 	fmt.Printf("%s: %s\n", k, v)
	// // 	rv := reflect.ValueOf(v)
	// // 	fmt.Printf("%T %#v\n", rv, reflect.Indirect(rv))
	// // }
	// // for k, v := range todMap {
	// // 	mv := reflect.ValueOf(v)
	// // 	if mv.Kind() != reflect.Pointer {
	// // 		mv = reflect.ValueOf(&v)
	// // 		fmt.Printf("%s %s %s %s\n", k, mv.Type(), mv.Kind(), mv.Type().Kind())
	// // 		qgenda.Process(mv)
	// // 	}
	// // 	// m[k]
	// // }

	// // qgenda.MapTest(todMap)
	// qgenda.ProcessMap(todMap)
	// if err := qgenda.ProcessMap(todMap); err != nil {
	// 	log.Println(err)
	// }
	// for k, v := range todMap {
	// 	fmt.Println(k, v)
	// }
	// if err := qgenda.ProcessSlice(todMap); err != nil {
	// 	log.Println(err)
	// }
	// if err := qgenda.ProcessMap(1); err != nil {
	// 	log.Println(err)
	// }
	// var xdf *string
	// xdfv := qgenda.IndirectReflectionValue(xdf)
	// xdfvs := xdfv.Addr()
	// fmt.Printf("%s\n", xdfvs)
	// fmt.Printf("%s is settable: %t\n", xdfv, xdfv.CanAddr())
	// var xdg *string
	// xdgv := reflect.ValueOf(xdg).Elem()
	// fmt.Printf("%s is settable: %t\n", xdgv, xdgv.CanAddr())
	// v := reflect.ValueOf(&qgenda.Schedule{})
	// vv := reflect.ValueOf(v)

	// fmt.Printf("%T: %s\n", vv, vv.Type())
	// fmt.Printf("%T: %s\n", vv, vv.Kind())
	// fmt.Println(vv.Kind() == reflect.Pointer)
	// vvi := reflect.Indirect(vv)
	// fmt.Printf("%T: %s\n", vvi, vvi.Kind())

	// fmt.Printf("%T: %s\n")
	// fmt.Println(tod)
	// b, err = tod.MarshalJSON()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(string(b))
	// structFun(sch[0])
	// FieldNames(sch[0])

	// var schTest []qgenda.Schedule
	// if err := json.Unmarshal([]byte(scheduleJSON), &schTest); err != nil {
	// 	log.Println(err)
	// }
	// // fmt.Printf("scheduleJSON: %#v\n", schTest[0].StartTime.String())
	// field, ok := reflect.TypeOf(schTest[0]).FieldByName("EndTime")
	// if !ok {
	// 	fmt.Println(ok)
	// }
	// tag := string(field.Tag)
	// fmt.Println(tag)
	// fmt.Printf("%#v\n", field)

	// todTest := &qgenda.TimeOfDay{}
	// if err := todTest.UnmarshalJSON([]byte(`00:12:00`)); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(todTest)

	// tod, err := qgenda.ParseTimeOfDay(`00:12:00,045`)
	// tod.SetLayout(`15:04:05.99999`)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("tod: %25s\n", tod.String())

	// todDev := &qgenda.TimeOfDay{}
	// fmt.Println(todDev.Time.Microseconds)
	// todDev.Set(time.Now().UTC())
	// fmt.Println(todDev.Time.Microseconds)
	// fmt.Println(todDev.Time.Microseconds)

	// tDur, err := time.ParseDuration(fmt.Sprintf("%dus", todDev.Time.Microseconds))
	// if err != nil {
	// 	log.Println(err)
	// }
	// tttt := &time.Time{}
	// fmt.Println(tttt.Add(tDur).Format("15:04:05"))

	// fmt.Println((time.Time{}).Add(time.Duration(todDev.Time.Microseconds) * time.Microsecond).Format("15:04:05"))
	// // todDev.SetLayout(`03:04PM`)
	// // todDev.SetLayout(`15:04:05.99999`)
	// fmt.Println(todDev)

	// b, err = json.Marshal(todDev)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("todDev:          %20s\n", string(b))
	// fmt.Printf("Empty TimeOfDay: %20s\n", (&qgenda.TimeOfDay{}).String())

	// dateRegex := `(?i)^(?P<year>\d{4}|\d{2})(?P<ymsep>\D*)(?P<month>\d{1,2})(?P<mdsep>\D*)(?P<day>\d{1,2})(?P<time>(?P<dtsep>\D*)(?P<hour>\d{1,2})(?P<hmsep>\D*)(?P<minute>\d{1,2})(?P<mssep>\D*)(?P<second>\d{1,2})(?P<subsecond>[[:punct:]]+\d*)?)?$`
	// dr := regexp.MustCompile(dateRegex)

	// matches := dr.FindStringSubmatch("2022-01-05T00:00:00.0000")
	// matches = dr.FindStringSubmatch("2022-01-05")
	// matches = dr.FindStringSubmatch("2022/01/05T01:02:03")
	// // fmt.Printf("TimeString: %20s Matches: %s\n", ts, matches)
	// components := map[string]string{}
	// for _, v := range dr.SubexpNames() {

	// 	subMatchIndex := dr.SubexpIndex(v)
	// 	if subMatchIndex > 0 && subMatchIndex <= len(dr.SubexpNames()) {
	// 		components[v] = matches[subMatchIndex]
	// 		fmt.Printf("%s:\t%s\n", v, matches[subMatchIndex])

	// 	}
	// }
	// dateLayoutSamples := []string{
	// 	`2022/01/05T01:02:03.8348348`,
	// 	`1904-04-28t22:13:12.85`,
	// }
	// dateSamples := []qgenda.Date{}
	// for _, v := range dateLayoutSamples {
	// 	ds, err := qgenda.GuessDateLayout(v)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	fmt.Printf("Date Layout: %35s %35s\n", v, ds)
	// 	d, err := qgenda.ParseDate(v)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	// d.SetLayout(`01/02/2006`)
	// 	dateSamples = append(dateSamples, d)

	// 	fmt.Printf("%35s %35s\n", v, d)

	// }

	// dateJSON, err := json.Marshal(dateSamples)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(string(dateJSON))

	// dtLayoutSamples := `
	// [
	// 	"2022/01/05T01:02:03.8348348",
	// 	"1904-04-28t22:13:12.85"
	// ]
	// `

	// var s []string
	// if err := json.Unmarshal([]byte(dtLayoutSamples), &s); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("s is awesome: %s\n", s)
	// // var dtXDT struct {
	// // 	Date *qgenda.Date `json:"Date,omitempty"`
	// // }
	// dtSamples := []qgenda.Date{}
	// // if err := json.Unmarshal(dateJSON, &dtSamples); err != nil {
	// // 	log.Println(err)
	// // }
	// if err := json.Unmarshal([]byte(dtLayoutSamples), &dtSamples); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("dtSamples: %v\n", dtSamples)
	// var xdt *qgenda.Date
	// xdt := &qgenda.Date{}
	// fmt.Printf("*qgenda.Date.layout: %s\n", xdt.Layout())

	// if err := xdt.UnmarshalJSON([]byte(`"2022/01/05T01:02:03.8348348"`)); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("xdt: %s\n", xdt)

	// fmt.Printf("nil time of date.Layout: %#v\n", (&qgenda.TimeOfDay{}).Layout())
	// fmt.Printf("nil time of date.Layout: %#v\n", (&qgenda.TimeOfDay{}).Layout())
	// fmt.Println(&time.Time{})
	// var x *qgenda.TimeOfDay
	// fmt.Println(x)
	// // var zdt *qgenda.TimeOfDay
	// zdt := &qgenda.TimeOfDay{}
	// (zdt).SetLayout(`03:04PM`)
	// fmt.Printf("nil time of date.Layout: %#v\n", zdt.Layout())

	// regexp.MatchString(``, `12/31/1999`)

}

var scheduleJSON = `
[
	{
		"ScheduleKey": "df349b7d-6080-4c0a-ac13-4a4ec9ea957a",
		"CompKey": "8c44c075-d894-4b00-9ae7-3b3842226626",
		"Credit": 1,
		"Date": "2022-01-05T00:00:00",
		"StartDateUTC": "2022-01-05T05:00:00Z",
		"EndDateUTC": "2022-01-06T04:59:59.999Z",
		"EndDate": "2022-01-05T00:00:00",
		"EndTime": "00:03:00",
		"IsCred": false,
		"IsPublished": true,
		"IsLocked": false,
		"IsStruck": true,
		"IsNotePrivate": false,
		"StaffAbbrev": "MO",
		"StaffBillSysId": "",
		"StaffEmail": "mobryant@pennstatehealth.psu.edu",
		"StaffEmrId": "mobryant",
		"StaffErpId": "",
		"StaffInternalId": "",
		"StaffExtCallSysId": "",
		"StaffFName": "Michael",
		"StaffId": "OBryant",
		"StaffKey": "93e3f801-8ce3-476a-8bc1-db86a432311a",
		"StaffLName": "O'Bryant",
		"StaffMobilePhone": "",
		"StaffNpi": "",
		"StaffPager": "",
		"StaffPayrollId": "",
		"StaffTags": [
			{
				"CategoryKey": 10448,
				"CategoryName": "Division",
				"Tags": [
					{
						"Key": 64255,
						"Name": "Resident"
					}
				]
			},
			{
				"CategoryKey": 10449,
				"CategoryName": "Task Type",
				"Tags": [
					{
						"EffectiveFromDate": "0001-01-01T00:00:00",
						"EffectiveToDate": "9999-12-31T00:00:00",
						"Key": 46339,
						"Name": "Resident"
					}
				]
			}
		],
		"StartDate": "2022-01-05T00:00:00",
		"StartTime": "00:03:00",
		"TaskAbbrev": "RES - Moonlight",
		"TaskId": "ML",
		"TaskKey": "614ada45-ef62-42f1-afc1-7010428a3893",
		"TaskName": "RES - Moonlight",
		"TaskIsPrintEnd": false,
		"TaskIsPrintStart": false,
		"TaskShiftKey": "8c829f91-58e0-4d7f-b188-0ed502657c94",
		"TaskType": "Working",
		"TaskTags": [
			{
				"CategoryKey": 10448,
				"CategoryName": "Division",
				"Tags": [
					{
						"Key": 64255,
						"Name": "Resident"
					}
				]
			},
			{
				"CategoryKey": 19429,
				"CategoryName": "System Task Type",
				"Tags": [
					{
						"Key": 180438,
						"Name": "Working"
					}
				]
			}
		],
		"LocationName": "",
		"LocationAbbrev": "",
		"LocationID": "",
		"LocationAddress": "",
		"TimeZone": "(UTC-05:00) Eastern Time (US \u0026 Canada)",
		"LastModifiedDateUTC": "2022-01-06T16:29:56.7909196Z",
		"IsRotationTask": false
	}
]
`

// // IsStruct returns true if a's kind is a struct or a pointer to a struct
// func IsStruct(a any) bool {
// 	v := reflect.Indirect(reflect.ValueOf(a))
// 	k := v.Type().Kind()
// 	return (k.String() == "struct")
// }

// // StructFields de-references as
// func StructFields(a any) []reflect.StructField {
// 	var structFields []reflect.StructField
// 	if IsStruct(a) {

// 		v := reflect.Indirect(reflect.ValueOf(a))
// 		t := v.Type()
// 		for i := 0; i < t.NumField(); i++ {
// 			f := t.Field(i)
// 			structFields = append(structFields, f)
// 		}
// 	}
// 	return structFields
// }

// func StructFieldNames(a any) []string {
// 	var fieldNames []string
// 	if IsStruct(a) {
// 		v := reflect.Indirect(reflect.ValueOf(a))
// 		t := v.Type()
// 		for i := 0; i < t.NumField(); i++ {
// 			f := t.Field(i).Name
// 			fieldNames = append(fieldNames, f)

// 			// fv := reflect.ValueOf(f)
// 			// if _, ok := afMap[f.Name]; ok {
// 			// 	fmt.Printf("%2d:\t%s\t%t\n", i, f.Name, v.Field(i).IsNil())
// 			// }
// 		}

// 	}
// 	return fieldNames
// }

// func StructFieldByName(a any, s string) reflect.StructField {
// 	if IsStruct(a) {
// 		v := reflect.Indirect(reflect.ValueOf(a))
// 		f, ok := v.Type().FieldByName(s)
// 		if !ok {
// 			log.Printf("Type %s doesn't have a field %#v\n", v.Type().String(), s)
// 		}
// 		return f

// 	}
// 	return reflect.StructField{}
// }

// // ImplementsInterface returns true if value implements Reference interface
// // note: Reference must be passed as a type parameter
// func ImplementsInterface[Reference any](value interface{}) bool {
// 	_, ok := value.(Reference)
// 	return ok
// }

// func FieldNames(a any) ([]string, error) {
// 	afMap := map[string]struct{}{}
// 	v := reflect.ValueOf(a)
// 	// fmt.Printf("ValueOf: %s\n", v)
// 	tt := reflect.TypeOf(a)
// 	fmt.Printf("TypeOf: %s\n", tt)

// 	fmt.Printf("reflect.Indirect(ValueOf): %s\n", reflect.Indirect(v))
// 	v = reflect.Indirect(v)
// 	fmt.Printf("ValueOf.Type.Name: %s\n", v.Type().Name())
// 	fmt.Printf("ValueOf.Type.Kind: %s\n", v.Type().Kind())
// 	fmt.Printf("ValueOf.Type.NumField: %d\n", v.Type().NumField())
// 	fmt.Printf("ValueOf.Type.Field(%d): %v\n", 0, v.Type().Field(0))
// 	t := v.Type()
// 	for i := 0; i < t.NumField(); i++ {
// 		f := t.Field(i)
// 		// fv := reflect.ValueOf(f)
// 		if _, ok := afMap[f.Name]; ok {
// 			fmt.Printf("%2d:\t%s\t%t\n", i, f.Name, v.Field(i).IsNil())
// 		}
// 	}

// 	return []string{}, nil
// }

// func ReflectionStuff(a any) {
// 	fmt.Printf("%-17T: Type: %-17s Kind: %-15s Indirect: %-17s Indirect.Kind: %-25s\n",
// 		a,
// 		reflect.ValueOf(a).Type(),
// 		reflect.ValueOf(a).Kind(),
// 		reflect.Indirect(reflect.ValueOf(a)).Type(),
// 		reflect.Indirect(reflect.ValueOf(a)).Kind(),
// 		// reflect.Indirect(reflect.ValueOf(a)).Interface(),
// 	)
// }
