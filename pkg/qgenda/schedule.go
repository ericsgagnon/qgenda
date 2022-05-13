package qgenda

type Schedule struct {
	// RawMessage        *string    `json:"-" db:"_raw_message"`
	// ExtractDateTime   *Time      `json:"-" db:"_extract_date_time"`
	ScheduleKey            *string       `json:"ScheduleKey,omitempty" primarykey:"true"`
	CallRole               *string       `json:"CallRole,omitempty"`
	CompKey                *string       `json:"CompKey,omitempty"`
	Credit                 *float64      `json:"Credit,omitempty"`
	Date                   *Date         `json:"Date,omitempty"`
	StartDateUTC           *Time         `json:"StartDateUTC,omitempty"`
	EndDateUTC             *Time         `json:"EndDateUTC,omitempty"`
	EndDate                *Date         `json:"EndDate,omitempty"`
	EndTime                *TimeOfDay    `json:"EndTime,omitempty"`
	IsCred                 *bool         `json:"IsCred,omitempty"`
	IsPublished            *bool         `json:"IsPublished,omitempty"`
	IsLocked               *bool         `json:"IsLocked,omitempty"`
	IsStruck               *bool         `json:"IsStruck,omitempty"`
	Notes                  *string       `json:"Notes,omitempty"`
	IsNotePrivate          *bool         `json:"IsNotePrivate,omitempty"`
	StaffAbbrev            *string       `json:"StaffAbbrev,omitempty"`
	StaffBillSysId         *string       `json:"StaffBillSysId,omitempty"`
	StaffEmail             *string       `json:"StaffEmail,omitempty"`
	StaffEmrId             *string       `json:"StaffEmrId,omitempty"`
	StaffErpId             *string       `json:"StaffErpId,omitempty"`
	StaffInternalId        *string       `json:"StaffInternalId,omitempty"`
	StaffExtCallSysId      *string       `json:"StaffExtCallSysId,omitempty"`
	StaffFName             *string       `json:"StaffFName,omitempty"`
	StaffId                *string       `json:"StaffId,omitempty"`
	StaffKey               *string       `json:"StaffKey,omitempty"`
	StaffLName             *string       `json:"StaffLName,omitempty"`
	StaffMobilePhone       *string       `json:"StaffMobilePhone,omitempty"`
	StaffNpi               *string       `json:"StaffNpi,omitempty"`
	StaffPager             *string       `json:"StaffPager,omitempty"`
	StaffPayrollId         *string       `json:"StaffPayrollId,omitempty"`
	StaffTags              []ScheduleTag `json:"StaffTags,omitempty"`
	StartDate              *Date         `json:"StartDate,omitempty"`
	StartTime              *TimeOfDay    `json:"StartTime,omitempty"`
	TaskAbbrev             *string       `json:"TaskAbbrev,omitempty"`
	TaskBillSysId          *string       `json:"TaskBillSysId,omitempty"`
	TaskContactInformation *string       `json:"TaskContactInformation,omitempty"`
	TaskExtCallSysId       *string       `json:"TaskExtCallSysId,omitempty"`
	TaskId                 *string       `json:"TaskId,omitempty"`
	TaskKey                *string       `json:"TaskKey,omitempty"`
	TaskName               *string       `json:"TaskName,omitempty"`
	TaskPayrollId          *string       `json:"TaskPayrollId,omitempty"`
	TaskEmrId              *string       `json:"TaskEmrId,omitempty"`
	TaskCallPriority       *string       `json:"TaskCallPriority,omitempty"`
	TaskDepartmentId       *string       `json:"TaskDepartmentId,omitempty"`
	TaskIsPrintEnd         *bool         `json:"TaskIsPrintEnd,omitempty"`
	TaskIsPrintStart       *bool         `json:"TaskIsPrintStart,omitempty"`
	TaskShiftKey           *string       `json:"TaskShiftKey,omitempty"`
	TaskType               *string       `json:"TaskType,omitempty"`
	TaskTags               []ScheduleTag `json:"TaskTags,omitempty"`
	LocationName           *string       `json:"LocationName,omitempty"`
	LocationAbbrev         *string       `json:"LocationAbbrev,omitempty"`
	LocationID             *string       `json:"LocationID,omitempty"`
	LocationAddress        *string       `json:"LocationAddress,omitempty"`
	TimeZone               *string       `json:"TimeZone,omitempty"`
	LastModifiedDateUTC    *Time         `json:"LastModifiedDateUTC,omitempty" primarykey:"true" querycondition:"ge" qf:"SinceModifiedTimestamp"`
	LocationTags           []Location    `json:"LocationTags,omitempty"`
	IsRotationTask         *bool         `json:"IsRotationTask"`
}

type ScheduleTag struct {
	ScheduleKey         *string `json:"ScheduleKey,omitempty" nullable:"false"`
	LastModifiedDateUTC *Time   `json:"LastModifiedDateUTC,omitempty" nullable:"false"`
	CategoryKey         *int64  `json:"CategoryKey" nullable:"false"`
	CategoryName        *string `json:"CategoryName" nullable:"false"`
	Tags                []struct {
		Key  *int64  `json:"Key" db:"tagkey" nullable:"false"`
		Name *string `json:"Name" db:"tagname" nullable:"false"`
	}
}

func NewScheduleRequest(rqf *RequestQueryFields) *Request {
	requestPath := "schedule"
	queryFields := []string{
		"CompanyKey",
		"StartDate",
		"EndDate",
		"IncludeDeletes",
		"SinceModifiedTimestamp",
		"DateFormat",
		"Includes",
		"Select",
		"Filter",
		"OrderBy",
		"Expand",
	}
	if rqf != nil {
		if rqf.Includes == nil {
			rqf.SetIncludes("StaffTags,TaskTags,LocationTags")
		}
	}
	r := NewRequestWithQueryField(requestPath, queryFields, rqf)
	return r
}

func PGCreateScheduleTableStatement(schema, table string) string {
	return PGCreateTableStatement(Schedule{}, schema, table)
}

func PGCreateScheduleTagStatement(schema, table string) string {
	return PGCreateTableStatement(ScheduleTag{}, schema, table)
}

func PGInsertScheduleStatement(schema, table string) string {
	return PGInsertStatement(Schedule{}, schema, table)
}

func PGInsertScheduleTagStatement(schema, table string) string {
	return PGInsertStatement(ScheduleTag{}, schema, table)
}

func PGQueryScheduleConstraintStatement(schema, table string) string {
	return PGQueryConstraintsStatement(Schedule{}, schema, table)
}

func PGQueryScheduleTagConstraintStatement(schema, table string) string {
	return PGQueryConstraintsStatement(ScheduleTag{}, schema, table)
}

func (sch Schedule) PGCreateTable() {

}

// func xxPGCreateScheduleTagTableStatement(schema, table string) string {
// 	if table == "" {
// 		table = "scheduletag"
// 	}

// 	tplValues := struct {
// 		Schema     string
// 		Table      string
// 		Fields     []string
// 		PrimaryKey []string
// 	}{
// 		Schema: schema,
// 		Table:  table,
// 		// Fields:     sqlFieldDefs,
// 		// PrimaryKey: primaryKey,
// 	}

// 	tpl := `
// 	CREATE TABLE IF NOT EXISTS {{ .Schema -}}{{- .Table }} (
//         schedulekey text  not null,
//         lastmodifieddateutc timestamp with time zone not null,
//         categorykey text,
//         categoryname text,
// 		tagkey text,
// 		tagname text
// 	)
// 	`
// 	var buf bytes.Buffer

// 	if err := template.Must(template.
// 		New("").
// 		Funcs(template.FuncMap{"join": strings.Join}).
// 		Parse(tpl)).
// 		Execute(&buf, tplValues); err != nil {
// 		log.Println(err)
// 		panic(err)
// 	}
// 	return buf.String()

// }

// func PGCreateScheduleTable(ctx context.Context, db sqlx.DB, schema string, table string) (sql.Result, error) {
// 	return db.ExecContext(ctx, PGCreateScheduleTableStatement(schema, table))
// }

// func xxPGCreateScheduleTagTable(ctx context.Context, db sqlx.DB, schema string, table string) (sql.Result, error) {
// 	return db.ExecContext(ctx, xxPGCreateScheduleTagTableStatement(schema, table))
// }

// //
// func xxPGInsertScheduleStatement(schema, table string) string {
// 	if table == "" {
// 		table = "scheduletag"
// 	}

// 	tplValues := struct {
// 		Schema     string
// 		Table      string
// 		Fields     []string
// 		PrimaryKey []string
// 	}{
// 		Schema: schema,
// 		Table:  table,
// 		// Fields:     sqlFieldDefs,
// 		// PrimaryKey: primaryKey,
// 	}
// 	tpl := `
// 	INSERT INTO {{ .Schema -}}{{- .Table }} (
// 	) VALUES (
// 		:first_name,
// 		:last_name,
// 		:email
// 	)
// 	CREATE TABLE IF NOT EXISTS {{ .Schema -}}{{- .Table }} (
// 		schedulekey text  not null,
//         lastmodifieddateutc timestamp with time zone not null,
//         categorykey text,
//         categoryname text,
// 		tagkey text,
// 		tagname text
// 	)
// 	`
// 	var buf bytes.Buffer

// 	if err := template.Must(template.
// 		New("").
// 		Funcs(template.FuncMap{"join": strings.Join}).
// 		Parse(tpl)).
// 		Execute(&buf, tplValues); err != nil {
// 		log.Println(err)
// 		panic(err)
// 	}
// 	return buf.String()
// }

// var pgSQLColumnSpecTpl = `
// {{- range  $index, $field := .Fields -}}
// {{- if ne $index 0 -}},{{- end }}
// 	{{ $field -}}
// {{- end -}}
// {{- $primarykey := join .PrimaryKey  ", " -}}
// {{ if ne $primarykey "" }},
// 	PRIMARY KEY ( {{ $primarykey }} ) {{ end }}
// `

// func XXPGCreateTableStatement[T any](schema, table string) string {
// 	return createSQLStatement[T](schema, table, pgCreateTableTpl)
// }

// func SQLTest[T any](schema string, table string) string {

// 	return createSQLStatement[T](schema, table, pgSQLColumnSpecTpl)
// }

// func XCreateScheduleTableSQL() string {
// 	// CreateTableSQL[Schedule]()
// 	return ""
// }

// func InsertScheduleRowSQL(schema, table string) string {
// 	return createSQLStatement[Schedule](schema, table, pgInsertRowsTpl)
// }

// func InsertScheduleTaskTagRowSQL(schema, table string)

//
// func InsertSchedulePG(db sqlx.DB, sch []Schedule) (sql.Result, error) {

// 	result, err := db.NamedExec(PGInsertScheduleStatement("", ""), sch)
// 	if err != nil {
// 		return result, err
// 	}
// 	type scheduleTagForSQL struct {
// 		ScheduleKey         *string `json:"ScheduleKey,omitempty" sql:",not null"`
// 		LastModifiedDateUTC *Time   `json:"LastModifiedDateUTC,omitempty" sql:",not null"`
// 		CategoryKey         *string
// 		CategoryName        *string
// 		TagKey              *string
// 		TagName             *string
// 	}

// 	var scht []scheduleTagForSQL
// 	for _, sc := range sch {
// 		schedulekey := sc.ScheduleKey
// 		lastmodifieddateutc := sc.LastModifiedDateUTC
// 		for _, ttcat := range sc.TaskTags {
// 			for _, tt := range ttcat.Tags {
// 				var stfs scheduleTagForSQL
// 				stfs.ScheduleKey = schedulekey
// 				stfs.LastModifiedDateUTC = lastmodifieddateutc
// 				stfs.CategoryKey = ttcat.CategoryKey
// 				stfs.CategoryName = ttcat.CategoryName
// 				stfs.TagKey = tt.Key
// 				stfs.TagName = tt.Name
// 				scht = append(scht, stfs)
// 			}

// 		}
// 	}
// 	result, err = db.NamedExec(PGInsertScheduleTagStatement("", ""), scht)

// 	return result, err
// }

// `
// INSERT INTO (
// _raw_message,
// _extract_date_time,
// schedulekey,
// callrole,
// compkey,
// credit,
// date,
// startdateutc,
// enddateutc,
// enddate,
// endtime,
// iscred,
// ispublished,
// islocked,
// isstruck,
// notes,
// isnoteprivate,
// staffabbrev,
// staffbillsysid,
// staffemail,
// staffemrid,
// stafferpid,
// staffinternalid,
// staffextcallsysid,
// stafffname,
// staffid,
// staffkey,
// stafflname,
// staffmobilephone,
// staffnpi,
// staffpager,
// staffpayrollid,
// startdate,
// starttime,
// taskabbrev,
// taskbillsysid,
// taskcontactinformation,
// taskextcallsysid,
// taskid,
// taskkey,
// taskname,
// taskpayrollid,
// taskemrid,
// taskcallpriority,
// taskdepartmentid,
// taskisprintend,
// taskisprintstart,
// taskshiftkey,
// tasktype,
// locationname,
// locationabbrev,
// locationid,
// locationaddress,
// timezone,
// lastmodifieddateutc,
// isrotationtask ,
// `

// `
// _raw_message,
// _extract_date_time,
// schedulekey,
// callrole,
// compkey,
// credit,
// date,
// startdateutc,
// enddateutc,
// enddate,
// endtime,
// iscred,
// ispublished,
// islocked,
// isstruck,
// notes,
// isnoteprivate,
// staffabbrev,
// staffbillsysid,
// staffemail,
// staffemrid,
// stafferpid,
// staffinternalid,
// staffextcallsysid,
// stafffname,
// staffid,
// staffkey,
// stafflname,
// staffmobilephone,
// staffnpi,
// staffpager,
// staffpayrollid,
// startdate,
// starttime,
// taskabbrev,
// taskbillsysid,
// taskcontactinformation,
// taskextcallsysid,
// taskid,
// taskkey,
// taskname,
// taskpayrollid,
// taskemrid,
// taskcallpriority,
// taskdepartmentid,
// taskisprintend,
// taskisprintstart,
// taskshiftkey,
// tasktype,
// locationname,
// locationabbrev,
// locationid,
// locationaddress,
// timezone,
// lastmodifieddateutc,
// isrotationtask
// `

// func CreateScheduleTableSQL(schema, table string) string {
// }

// func CreateScheduleTagTableSQL(schema, table string) string {
// func PGCreateScheduleTableStatement(schema, table string) string {
// 	if table == "" {
// 		table = "schedule"
// 	}
// 	return CreateTableSQL[Schedule](schema, table)

// 	// if table == "" {
// 	// 	table = "schedule"
// 	// }

// 	// tplValues := struct {
// 	// 	Schema     string
// 	// 	Table      string
// 	// 	Fields     []string
// 	// 	PrimaryKey []string
// 	// }{
// 	// 	Schema: schema,
// 	// 	Table:  table,
// 	// 	// Fields:     sqlFieldDefs,
// 	// 	// PrimaryKey: primaryKey,
// 	// }

// 	// tpl := `
// 	// CREATE TABLE IF NOT EXISTS {{ .Schema -}}{{- .Table }} (
// 	//     schedulekey text  not null,
// 	//     lastmodifieddateutc timestamp with time zone not null,
// 	//     categorykey text,
// 	//     categoryname text,
// 	// 	tagkey text,
// 	// 	tagname text
// 	// )
// 	// `
// 	// var buf bytes.Buffer

// 	// if err := template.Must(template.
// 	// 	New("").
// 	// 	Funcs(template.FuncMap{"join": strings.Join}).
// 	// 	Parse(tpl)).
// 	// 	Execute(&buf, tplValues); err != nil {
// 	// 	log.Println(err)
// 	// 	panic(err)
// 	// }
// 	// return buf.String()

// }
