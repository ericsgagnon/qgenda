package main

import (
	"context"
	"database/sql"
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
	if err := sm.LoadFromFile("../.cache/staffmember.json"); err != nil {
		log.Println(err)
	}

	// b, err := os.ReadFile("../.cache/staffmember.json")
	// if err != nil {
	// 	log.Println(err)
	// }
	// if err := json.Unmarshal(b, &sm); err != nil {
	// 	log.Println(err)
	// }
	fmt.Printf("Length of []StaffMember: %d\n", len(sm))
	for i, v := range sm {
		vi := v
		if vi.ExtractDateTime == nil {
			now := qgenda.NewTime(time.Now().UTC())
			vi.ExtractDateTime = &now
		}
		sm[i] = vi
	}

	smtable := qgenda.StructToTable(qgenda.StaffMember{}, "staff", "qgenda", false, nil, nil)

	smsql := qgenda.PGTableStatement(
		smtable,
		qgenda.PGCreateTableDevTpl,
		nil,
	)
	pgc := &qgenda.PGClient{
		DB:     db,
		Config: qgenda.DBClientConfig{},
	}
	pgc.CreateSchema(ctx, "staffmembertest", "")
	// fmt.Println(smsql)
	fmt.Println(db.Exec(smsql))
	// for _, field := range smtable.Fields {
	// 	fmt.Printf("field: %s\n", field.Name)
	// 	fmt.Printf("Type: %s\n", field.Type)
	// 	fmt.Printf("Kind: %s\n", field.Kind)
	// }
	result, err = qgenda.PGInsertRowsDev(ctx, db, smtable, sm)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(result)

	// fmt.Println(reflect.ValueOf(sm[0]).Type().Name())

	// smsql02 := qgenda.PGTableStatement(smtable, qgenda.PGInsertChangesOnlyDevTpl, nil)
	smtable = qgenda.StructToTable(qgenda.StaffMember{}, "staffmember", "qgenda", true, nil, nil)
	smtable.Temporary = true
	smsql02 := qgenda.PGTableStatement(smtable, qgenda.PGInsertRowsDevTpl, nil)
	fmt.Sprint(smsql02)
	// fmt.Println(smsql02)
	smsql02 = qgenda.PGTableStatement(smtable, qgenda.PGInsertChangesOnlyDevTpl, nil)
	// fmt.Println(smsql02)
	pgc.Tx, err = pgc.BeginTxx(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	result, err = pgc.CreateTable(ctx, smtable, true)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)
	now := qgenda.NewTime(nil)
	for i, _ := range sm {
		if sm[i].ExtractDateTime == nil {
			sm[i].ExtractDateTime = &now

		}
	}
	result, err = sm.InsertToPG(ctx, db, "plonky", "")
	fmt.Printf("Insert []StaffMember to plonky.staffmember: %s\t%s\n", result, err)

	sftrs := sm[0].FlatTags()
	if len(sftrs) > 0 {
		fmt.Println(sftrs[0])
	}

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
	result, err = sm[0:9].InsertToPG(ctx, db, "plonky", "")
	fmt.Println(result, err)

	smss, err := sm.Extract(ctx, c, nil)
	if err != nil {
		log.Println(err)
	}
	if _, err := smss.Process(); err != nil {
		log.Println(err)
	}
	result, err = smss.InsertToPG(ctx, db, "plonky", "")
	fmt.Println(result, err)
	// result, err = db.ExecContext(ctx, "DROP SCHEMA IF EXISTS plonky CASCADE")
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(result)
}
