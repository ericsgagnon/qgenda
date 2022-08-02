package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

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

	// qgenda client ---------------------------------------------------------
	qcc := &qgenda.ClientConfig{
		Email:    os.Getenv("QGENDA_EMAIL"),
		Password: os.Getenv("QGENDA_PASSWORD"),
	}
	c, err := qgenda.NewClient(qcc)
	if err != nil {
		log.Fatalln(err)
	}
	c.Auth()

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
	var result sql.Result

	// schedules -----------------------------------------------------------
	// rqf := &qgenda.RequestQueryFields{}
	// rqf.SetStartDate(time.Now().UTC().Add(time.Hour * 24 * -120))
	// rqf.SetEndDate(time.Now().UTC().Add(time.Hour * 24 * -1))
	// if rqf.SinceModifiedTimestamp == nil {
	// 	rqf.SetSinceModifiedTimestamp(time.Now().UTC().Add(time.Hour * 24 * -30))
	// }

	// tm, _ := time.Parse(time.RFC3339, "2012-01-01T00:00:00Z")
	// rqf.SetStartDate(tm.Add(time.Hour * 24 * 365))
	// rqf.SetEndDate(rqf.GetStartDate().Add(time.Hour * 24 * 31))
	// rqf.SetSinceModifiedTimestamp(rqf.GetStartDate())
	// s := qgenda.Schedules{}
	// result, err = s.EPL(ctx, c, rqf, db, "qgenda", "schedule", false)
	// if err != nil {
	// 	log.Println(err)
	// }
	// if result != nil {
	// 	ra, err := result.RowsAffected()
	// 	fmt.Printf("RowsAffected: %d\tErrors: %v\n", ra, err)
	// }

	// staffmembers -----------------------------------------------------
	sm := qgenda.StaffMembers{}
	if err := sm.LoadFile("../.cache/staffmember.json"); err != nil {
		log.Println(err)
	}
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
	var rowsAffected int64
	schema := "sm1"
	result, err = sm[0:9].InsertToPG(ctx, db, schema, "")
	rowsAffected, _ = result.RowsAffected()
	fmt.Printf("Insert []StaffMember to %s.staffmember: %d\t%s\n", schema, rowsAffected, err)

	schema = "sm2"
	result, err = sm[0:19].InsertToPG(ctx, db, schema, "")
	rowsAffected, _ = result.RowsAffected()
	fmt.Printf("Insert []StaffMember to %s.staffmember: %d\t%s\n", schema, rowsAffected, err)

	schema = "plonky"
	result, err = sm[0:9].InsertToPG(ctx, db, schema, "")
	rowsAffected, _ = result.RowsAffected()
	fmt.Printf("Insert []StaffMember to %s.staffmember: %d\t%s\n", schema, rowsAffected, err)

	result, err = sm[0:19].InsertToPG(ctx, db, schema, "")
	rowsAffected, _ = result.RowsAffected()
	fmt.Printf("Insert []StaffMember to %s.staffmember: %d\t%s\n", schema, rowsAffected, err)

	result, err = sm[0:99].InsertToPG(ctx, db, schema, "")
	rowsAffected, _ = result.RowsAffected()
	fmt.Printf("Insert []StaffMember to %s.staffmember: %d\t%s\n", schema, rowsAffected, err)

	schema = "importalicious"
	var smss qgenda.StaffMembers
	smss, err = smss.Extract(ctx, c, nil)
	if err != nil {
		log.Println(err)
	}
	if _, err := smss.Process(); err != nil {
		log.Println(err)
	}

	smss.WriteFile(fmt.Sprintf("../.cache/staffmember-%s.json", smss[0].ExtractDateTime.Time.Format("20060102150405")))
	result, err = smss.InsertToPG(ctx, db, schema, "")
	rowsAffected, _ = result.RowsAffected()
	fmt.Printf("Insert []StaffMember to %s.staffmember: %d\t%s\n", schema, rowsAffected, err)
	tvals := map[string]bool{}
	for _, v := range smss {
		tstring := fmt.Sprint(v.ExtractDateTime.Time.Format("2006-01-02T15:04:05.99999999"))
		tvals[tstring] = true
	}
	for k, _ := range tvals {
		fmt.Println(k)
	}

	req := qgenda.NewStaffMemberRequest(nil)
	fmt.Println(req.ToURL().String())

	resp, err := c.Do(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	if err := os.WriteFile("../.cache/test.json", data, 0644); err != nil {
		log.Println(err)
	}
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
