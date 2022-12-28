package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ericsgagnon/qgenda/pkg/qgenda"
	"github.com/jmoiron/sqlx"
)

func ptr[T any](value T) *T {
	return &value
}

func val[T any](pointer *T) T {
	return *pointer
}

func main() {
	ctx := context.Background()
	var b []byte
	// cachefile cfg ---------------------------------------------------------
	cfcfg, err := qgenda.NewCacheConfig("qgenda-exporter")
	if err != nil {
		log.Println(err)
	}

	// qgenda client ---------------------------------------------------------
	qcc := &qgenda.ClientConfig{
		CompanyKey: os.Getenv("QGENDA_COMPANY_KEY"),
		Email:      os.Getenv("QGENDA_EMAIL"),
		Password:   os.Getenv("QGENDA_PASSWORD"),
	}
	c, err := qgenda.NewClient(qcc)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("preauth")
	c.Auth()
	fmt.Println("Company Key: ", c.ClientConfig.CompanyKey)
	// pg client -------------------------------------------------------------
	db, err := sqlx.Open("postgres", os.Getenv("PG_CONNECTION_STRING"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	if db.Ping() == nil {
		fmt.Printf("Connected to %s: %+v\n", db.DriverName(), db.Stats())
	}
	// var result sql.Result

	// schedules -----------------------------------------------------------
	rqf := &qgenda.RequestConfig{}
	rqf.SetStartDate(time.Now().UTC().Add(time.Hour * 24 * -120))
	rqf.SetEndDate(time.Now().UTC().Add(time.Hour * 24 * -1))
	if rqf.SinceModifiedTimestamp == nil {
		rqf.SetSinceModifiedTimestamp(time.Now().UTC().Add(time.Hour * 24 * -30))
	}

	// tm, _ := time.Parse(time.RFC3339, "2012-01-01T00:00:00Z")
	tm, _ := time.Parse(time.RFC3339, "2021-10-01T00:00:00Z")
	rqf.SetStartDate(tm.Add(time.Hour * 24 * 365))
	rqf.SetEndDate(rqf.GetStartDate().Add(time.Hour * 24 * 31))
	rqf.SetSinceModifiedTimestamp(rqf.GetStartDate())
	// c.Do()
	s := qgenda.Schedules{}
	rc := qgenda.DefaultScheduleRequestConfig()
	if err := s.Get(ctx, c, rc); err != nil {
		log.Println(err)

	}
	fmt.Sprint(s)
	// fmt.Print(s)
	// result, err = s.EPL(ctx, c, rqf, db, "qgenda", "schedule", false)
	// if err != nil {
	// 	log.Println(err)
	// }
	// if result != nil {
	// 	ra, err := result.RowsAffected()
	// 	fmt.Printf("RowsAffected: %d\tErrors: %v\n", ra, err)
	// }
	// ss := qgenda.Schedules{}
	// ss, err = ss.Extract(ctx, c, rqf)
	// if err != nil {
	// 	log.Println(err)
	// }

	// b, err = json.MarshalIndent(ss, "", "\t")
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println("Length of ss schedules: ", len(ss))
	// // fmt.Println(string(b))
	cf, err := qgenda.NewCacheFile("schedules-sample.json", "", cfcfg)
	if err != nil {
		log.Println(err)
	}
	// if err := cf.Create(); err != nil {
	// 	log.Println(err)
	// }
	// if err := cf.Write(b); err != nil {
	// 	log.Println(err)
	// }

	cfp, err := qgenda.NewCacheFile("xschedules.json", "", cfcfg)
	if err != nil {
		log.Println(err)
	}
	// xs := []qgenda.XSchedule{}
	var xs []qgenda.Schedule
	// fmt.Println(cf.String())
	b, err = cf.Read()
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(string(b))
	// if err := xs.LoadFile(cf.String()); err != nil {
	// 	log.Println(err)
	// }
	// xs.Process()
	if err := json.Unmarshal(b, &xs); err != nil {
		log.Println(err)
	}
	fmt.Println("Length of XSchedules", len(xs))

	for i, _ := range xs {
		if err := xs[i].Process(); err != nil {
			log.Println(err)
		}
	}
	// for _, v := range xs[0:9] {
	// 	if len(v.TaskTags) > 0 {
	// 		for _, vv := range v.TaskTags {
	// 			if len(vv.Tags) > 0 {
	// 				fmt.Println(*v.ScheduleKey)
	// 				fmt.Printf("%#v\n", vv)
	// 				// for _, vvv := range vv.Tags {
	// 				// 	fmt.Println(*(v.ScheduleKey))
	// 				// 	fmt.Println(*(vvv.CategoryName))
	// 				// }
	// 			}
	// 		}
	// 	}
	// }
	b, err = json.MarshalIndent(xs, "", "\t")
	if err != nil {
		log.Println(err)
	}
	// b, err = xs.MarshalJSON()
	// if err != nil {
	// 	log.Println(err)
	// }
	if err := cfp.Create(); err != nil {
		log.Println(err)
	}
	if err := cfp.Write(b); err != nil {
		log.Println(err)
	}

	// staffmembers -----------------------------------------------------
	sm := qgenda.StaffMembers{}
	if err := sm.LoadFile("../.cache/staffmember.json"); err != nil {
		log.Println(err)
	}
	fmt.Println("staffmember length: ", len(sm))
	_, err = sm.Process()
	if err != nil {
		log.Println(err)
	}
	now := qgenda.NewTime(nil)
	for i, _ := range sm {
		if sm[i].ExtractDateTime == nil {
			sm[i].ExtractDateTime = &now

		}
	}

	//---------------------------------------------------------------
	// fmt.Println(xs[0].PGQuery(true, "temptest", "xschedule"))
	// fmt.Println(xs[0].PGQuery(false, "test", "xschedule"))
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	// res, err := xs[0].PGCreateTable(ctx, tx, "slurpydurpy", "", true, "")
	// if err != nil {
	// 	log.Println(err)
	// }
	// xss := qgenda.XSchedules(xs[0:1000])

	xss := qgenda.Schedules(xs)
	res, err := xss.PGInsertRows(ctx, tx, "slurpydurpy", "", "")
	if err != nil {
		log.Println(err)
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
	}
	fmt.Println(res)

	sx := qgenda.Schedules{}
	fmt.Println(rqf.GetStartDate(), rqf.GetEndDate(), rqf.GetSinceModifiedTimestamp())
	if err := sx.Get(ctx, c, rqf); err != nil {
		log.Println(err)
	}
	if err := sx.Process(); err != nil {
		log.Println(err)
	}
	tx, err = db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	sres, err := sx.PGInsertRows(ctx, tx, "newschema", "", "")
	if err != nil {
		log.Println(err)
	} else {
		tx.Commit()
	}
	fmt.Println(sres)
	rqfn, err := qgenda.Schedule{}.PGQueryConstraints(ctx, db, "newschema", "schedule")
	if err != nil {
		log.Println(err)
	}

	sxx := qgenda.Schedules{}
	if err := sxx.Get(ctx, c, rqfn); err != nil {
		log.Println(err)
	}
	fmt.Printf("length of Schedules:\t%d\n", len(sxx))
	if err := sxx.Process(); err != nil {
		log.Println(err)
	}

	tx, err = db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	sres, err = sxx.PGInsertRows(ctx, tx, "newschema", "", "")
	if err != nil {
		log.Println(err)
	} else {
		tx.Commit()
	}
	fmt.Println(sres)

	// var rowsAffected int64
	// schema := "sm1"
	// result, err = sm[0:9].InsertToPG(ctx, db, schema, "")
	// rowsAffected, _ = result.RowsAffected()
	// fmt.Printf("Insert []StaffMember to %s.staffmember: %d\t%s\n", schema, rowsAffected, err)

	// schema = "sm2"
	// result, err = sm[0:19].InsertToPG(ctx, db, schema, "")
	// rowsAffected, _ = result.RowsAffected()
	// fmt.Printf("Insert []StaffMember to %s.staffmember: %d\t%s\n", schema, rowsAffected, err)

	// schema = "plonky"
	// result, err = sm[0:9].InsertToPG(ctx, db, schema, "")
	// rowsAffected, _ = result.RowsAffected()
	// fmt.Printf("Insert []StaffMember to %s.staffmember: %d\t%s\n", schema, rowsAffected, err)

	// result, err = sm[0:19].InsertToPG(ctx, db, schema, "")
	// rowsAffected, _ = result.RowsAffected()
	// fmt.Printf("Insert []StaffMember to %s.staffmember: %d\t%s\n", schema, rowsAffected, err)

	// result, err = sm[0:99].InsertToPG(ctx, db, schema, "")
	// rowsAffected, _ = result.RowsAffected()
	// fmt.Printf("Insert []StaffMember to %s.staffmember: %d\t%s\n", schema, rowsAffected, err)

	// schema = "importalicious"
	// var smss qgenda.StaffMembers
	// smss, err = smss.Extract(ctx, c, nil)
	// if err != nil {
	// 	log.Println(err)
	// }
	// if _, err := smss.Process(); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println("Length of staffmembers: ", len(smss))
	// smss.WriteFile(fmt.Sprintf("../.cache/staffmember-%s.json", smss[0].ExtractDateTime.Time.Format("20060102150405")))
	// result, err = smss.InsertToPG(ctx, db, schema, "")
	// rowsAffected, _ = result.RowsAffected()
	// fmt.Printf("Insert []StaffMember to %s.staffmember: %d\t%s\n", schema, rowsAffected, err)
	// tvals := map[string]bool{}
	// for _, v := range smss {
	// 	tstring := fmt.Sprint(v.ExtractDateTime.Time.Format("2006-01-02T15:04:05.99999999"))
	// 	tvals[tstring] = true
	// }
	// for k, _ := range tvals {
	// 	fmt.Println(k)
	// }

	// req := qgenda.NewStaffMemberRequest(nil)
	// fmt.Println(req.ToURL().String())

	// resp, err := c.Do(ctx, req)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// data, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// if err := os.WriteFile("../.cache/test.json", data, 0644); err != nil {
	// 	log.Println(err)
	// }
	// rqf := &qgenda.RequestConfig{}
	// rqf = &qgenda.RequestConfig{}
	// tasks -
	// rqf.SetIncludes("Tags,TaskShifts,Profiles")
	// fmt.Println(rqf.GetIncludes())
	// req = qgenda.NewRequestWithQueryField(
	// 	"task",
	// 	[]string{"Includes"},
	// 	rqf,
	// )

	// tags
	// rqf.SetIncludes("")
	// fmt.Println(rqf.GetIncludes())
	// req = qgenda.NewRequestWithQueryField(
	// 	"tags",
	// 	nil,
	// 	rqf,
	// )

	// scheduleauditlog
	// rqf.SetIncludes("Location")
	// rqf.SetScheduleStartDate(time.Now().Add(-1 * time.Hour * 24 * 14))
	// rqf.SetScheduleEndDate(time.Now().Add(-1 * time.Hour * 24))
	// fmt.Println(rqf.GetIncludes())
	// req = qgenda.NewRequestWithQueryField(
	// 	"schedule/auditLog",
	// 	[]string{"ScheduleStartDate", "ScheduleEndDate", "DateFormat"},
	// 	rqf,
	// )
	// fmt.Println(req.ToURL().String())
	// resp, err = c.Do(ctx, req)
	// if err != nil {
	// 	log.Println(err)
	// }
	// data, err = io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// if err := os.WriteFile("../.cache/sample.json", data, 0644); err != nil {
	// 	log.Println(err)
	// }

	// dsdf := qgenda.StructToFields(qgenda.DatasetDev{})
	// for i := range dsdf {
	// 	fmt.Println(dsdf[i])
	// }

	// unrelated --------------------------------------------------------
	// v := qgenda.Test{}
	// fmt.Println(v)
	// if err := v.Process(50); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(v)
	// ------------------------------------------------------------------

	// if err := json.Unmarshal(data, &sms); err != nil {
	// 	return nil, err
	// }

	// for _, v := range sm[0:9] {
	// 	if v.Tags != nil {
	// 		for _, vt := range v.Tags {
	// 			if vt.Tags != nil {

	// 				for _, vtt := range vt.Tags {
	// 					if vtt.Key != nil && vtt.Name != nil {

	// 						fmt.Printf(`ExtractDateTime: %s StaffKey: %s LastModifiedDateUTC: %s TagKey: %s TagName %s\n`, *vt.ExtractDateTime, *vt.StaffKey, *vt.LastModifiedDateUTC, *vtt.Key, *vtt.Name)
	// 					}
	// 				}

	// 			}
	// 		}
	// 	}
	// }

	// fmt.Println(db.ExecContext(ctx, "DROP SCHEMA IF EXISTS plonky CASCADE"))
	// fmt.Println(db.ExecContext(ctx, "DROP SCHEMA IF EXISTS importalicious CASCADE"))
	// result, err = db.ExecContext(ctx, "DROP SCHEMA IF EXISTS plonky CASCADE")
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(result)
	// ----------------------------------------------------------------------
	// fmt.Printf("Length of []StaffMember: %d\n", len(sm))
	// var x int
	// for i, v := range sm {
	// 	vi := v
	// 	if vi.ExtractDateTime == nil {
	// 		x = x + 1
	// 		now := qgenda.NewTime(time.Now().UTC())
	// 		vi.ExtractDateTime = &now
	// 	}
	// 	sm[i] = vi
	// }
	// fmt.Println("Number of StaffMembers with nil ExtractDateTime: ", x)

	// smtable := qgenda.StructToTable(qgenda.StaffMember{}, "staff", "qgenda", false, nil, nil)

	// smsql := qgenda.PGTableStatement(
	// 	smtable,
	// 	qgenda.PGCreateTableDevTpl,
	// 	nil,
	// )
	// pgc := &qgenda.PGClient{
	// 	DB:     db,
	// 	Config: qgenda.DBClientConfig{},
	// }
	// pgc.CreateSchema(ctx, "staffmembertest", "")
	// fmt.Println(smsql)
	// fmt.Println(db.Exec(smsql))
	// for _, field := range smtable.Fields {
	// 	fmt.Printf("field: %s\n", field.Name)
	// 	fmt.Printf("Type: %s\n", field.Type)
	// 	fmt.Printf("Kind: %s\n", field.Kind)
	// }
	// result, err = qgenda.PGInsertRowsDev(ctx, db, smtable, sm)
	// if err != nil {
	// 	log.Println(err)
	// }

	// fmt.Println(result)

	// fmt.Println(reflect.ValueOf(sm[0]).Type().Name())

	// smsql02 := qgenda.PGTableStatement(smtable, qgenda.PGInsertChangesOnlyDevTpl, nil)
	// smtable = qgenda.StructToTable(qgenda.StaffMember{}, "staffmember", "qgenda", true, nil, nil)
	// smtable.Temporary = true
	// smsql02 := qgenda.PGTableStatement(smtable, qgenda.PGInsertRowsDevTpl, nil)
	// fmt.Sprint(smsql02)
	// // fmt.Println(smsql02)
	// smsql02 = qgenda.PGTableStatement(smtable, qgenda.PGInsertChangesOnlyDevTpl, nil)
	// // fmt.Println(smsql02)
	// pgc.Tx, err = pgc.BeginTxx(ctx, nil)
	// if err != nil {
	// 	log.Println(err)
	// }
	// result, err = pgc.CreateTable(ctx, smtable, true)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(result)

	// sftrs := sm[0].FlatTags()
	// if len(sftrs) > 0 {
	// 	fmt.Println(sftrs[0])
	// }

	// smttable := qgenda.StructToTable(qgenda.FlatStaffTag{}, "staffmembertag", "qgenda", true, nil, nil)
	// smtsql := qgenda.PGTableStatement(smttable, qgenda.PGCreateTableDevTpl, nil)
	// fmt.Println(smtsql)

	// tempScheduleTable := qgenda.StructToTable(qgenda.Schedule{}, "schedule", "qgenda", true, nil, nil)
	// tempScheduleTableSQL := qgenda.PGTableStatement(tempScheduleTable, qgenda.PGCreateTableDevTpl, nil)
	// fmt.Println(tempScheduleTableSQL)

	// for _, f := range smtable.Fields {
	// 	if f.Name == "Tags" {
	// 		fmt.Println(f.StructField.Type.Kind())
	// 		fmt.Println(reflect.Slice)
	// 		fmt.Println(qgenda.PGOmit(f))
	// 	}
	// }
	// smincludes := qgenda.PGIncludeFields(smtable.Fields)
	// for _, f := range smincludes {
	// 	fmt.Println(f.Name)
	// }

}
