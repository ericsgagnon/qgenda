package qgenda

// type DBClient interface {
// 	// Open(cfg *DBClientConfig) (*sqlx.DB, error)
// 	// Ping(ctx context.Context) (bool, error)
// 	// DB() *sqlx.DB
// 	CreateSchema(ctx context.Context, schema string) (sql.Result, error)
// 	CreateTable(ctx context.Context, table Table) (sql.Result, error)
// 	DropTable(ctx context.Context, table Table) (sql.Result, error)
// 	// InsertRows(ctx context.Context, data Dataset) (sql.Result, error)
// 	// QueryConstraints(ctx context.Context, data Dataset) error
// }

// // type DClient struct {
// // 	*sqlx.DB
// // }

// // func (c *DClient) CreateTable() {
// // 	switch c.DriverName() {
// // 	case "postgres":
// // 	}
// // }

// type DBClientConfig struct {
// 	Name               string // descriptive name - used only for logs and reference
// 	Type               string // descriptive type - used only for logs and reference
// 	Driver             string // driver name - will be passed to sqlx.DB
// 	DataSourceName     string // only applicable if using DSN's
// 	ConnectionString   string // will be parsed to url
// 	Schema             string // schema to use for this client
// 	ExpandEnvVars      bool   // whether or not to interpolate env vars of the form ${ENV_VAR} in connection string and dsn
// 	ExpandFileContents bool   // whether or not to interpolate file contents of the form {file:/path/to/file} in connection string and dsn
// 	// User             string   // prefer to reference env var or file contents by ${ENV_VAR_NAME} or {file:/path/to/file}
// 	// Password         string   // prefer to reference env var or file contents by ${ENV_VAR_NAME} or {file:/path/to/file}
// 	// url *url.URL // let the program handle this
// }

// func ExampleDBClientConfig() DBClientConfig {
// 	cfg := DBClientConfig{
// 		Name:               "database",
// 		Type:               "odbc",
// 		Driver:             "odbc",
// 		ConnectionString:   "${DB_CONN_SCHEME}://${DB_USER}:${DB_PASSWORD}@${DB_HOSTNAME}:${DB_PORT}/${DB_DATABASE}?${DB_ARGUMENTS}",
// 		Schema:             "qgenda",
// 		ExpandEnvVars:      true,
// 		ExpandFileContents: true,
// 		// User:             "${DB_USER}",
// 		// Password:         "${DB_PASSWORD}",
// 	}
// 	return cfg
// }

// func (cfg DBClientConfig) String() string {
// 	s := cfg.ConnectionString
// 	if cfg.ExpandEnvVars {
// 		s = ExpandEnvVars(s)
// 	}
// 	if cfg.ExpandFileContents {
// 		s = ExpandFileContents(s)
// 	}
// 	return s
// }

// // OpenDBConnection doesn't technically 'open' a real connection, it follows
// // the go default of creating a DB struct that manages connections as needed
// func OpenDBClient(cfg *DBClientConfig) (DBClient, error) {
// 	return nil, nil
// }

// func NewDBClient(cfg *DBClientConfig) (*sqlx.DB, error) {
// 	connString := ExpandEnvVars(cfg.ConnectionString)
// 	// fmt.Printf("Driver: %s\t ConnString: %s\n", cfg.Driver, connString)
// 	return sqlx.Open(cfg.Driver, connString)
// }

// type Table struct {
// 	Name           string
// 	Type           reflect.Type
// 	Schema         string
// 	Temporary      bool
// 	UUID           string
// 	Constraints    map[string]string // these should be table constraints, not field constraints
// 	Fields         Fields
// 	Children       []Table
// 	Tags           map[string][]string
// 	UpdateStrategy string
// 	Parent         *Table
// }

// func (t Table) FullIdentifier() string {
// 	schema := strings.ToLower(t.Schema)
// 	name := strings.ToLower(t.Name)
// 	uuid := strings.ReplaceAll(t.UUID, "-", "")
// 	switch {
// 	case t.Temporary && t.UUID != "":
// 		return "_tmp_" + uuid + "_" + t.Name
// 	case t.Temporary && uuid == "":
// 		return "_tmp_" + name
// 	case t.Schema != "":
// 		return schema + "." + name
// 	default:
// 		return name
// 	}
// }

// func (t Table) PrimaryKey() string {
// 	for k, v := range t.Constraints {
// 		key := strings.ToLower(k)
// 		key = strings.ReplaceAll(key, " ", "")
// 		if key == "primarykey" {
// 			return v
// 		}
// 	}
// 	return ""
// }

// type Fields []Field

// // type UpdateStrategy int8

// // func (u UpdateStrategy) String() string {
// // 	return ""
// // }

// // const (
// //  AppendChanges UpdateStrategy = iota
// // AppendAll
// // ReplaceChanges
// // ReplaceAll
// // )

// // ReflectType converts value to a reflect.Type
// // if value is already a reflect.Type, it asserts and returns it
// func ReflectType(value any) reflect.Type {
// 	if rt, ok := value.(reflect.Type); ok {
// 		return rt
// 	}
// 	return reflect.TypeOf(value)
// 	// rt := reflect.TypeOf(value)
// 	// rtta := reflect.TypeOf(reflect.TypeOf(any("")))

// 	// if rt == rtta {
// 	// 	rt = value.(reflect.Type)
// 	// }
// 	// return rt
// }

// func NewTable(value any) Table {

// 	rt := reflect.TypeOf(value)
// 	if v, ok := value.(reflect.Type); ok {
// 		rt = v
// 	}
// 	if rt.Kind() != reflect.Struct {
// 		panic(fmt.Errorf("%T is not a struct", value))
// 	}

// 	return Table{}
// }

// func StructToTable[T any](value T, name, schema string, temporary bool, id string, constraints map[string]string, tags map[string][]string, parent *Table) Table {

// 	fields := StructToFields(value)
// 	// fields := StructToFields(*new(T))
// 	// rv := reflect.ValueOf(value)
// 	// fmt.Println("---------------------------------------------------")
// 	// fmt.Printf("\t%T\n", *new(T))
// 	// rv := reflect.ValueOf(*new(T))
// 	// rv := reflect.ValueOf(value)
// 	// fmt.Println(rv)
// 	// rt := rv.Type()
// 	rt := ReflectType(value)
// 	if rt.Kind() != reflect.Struct {
// 		panic(fmt.Errorf("%T is not a struct", value))
// 	}

// 	if name == "" {
// 		name = strings.ToLower(rt.Name())
// 	}

// 	if len(constraints) == 0 {
// 		constraints = map[string]string{}
// 		pk := strings.Join(PrimaryKey(fields), ", ")
// 		if pk != "" {
// 			constraints["primarykey"] = pk

// 		}
// 		uf := []string{}
// 		for _, v := range UniqueFields(fields) {

// 			uf = append(uf, PGName(v))
// 		}
// 		if len(uf) > 0 {
// 			constraints["unique"] = strings.Join(uf, ", ")
// 		}
// 	}
// 	if id == "" {
// 		id = strings.ReplaceAll(uuid.NewString(), "-", "")
// 	}

// 	// check https://www.postgresql.org/docs/current/limits.html for current
// 	// identifier limits. Limit is 63 at time of coding.
// 	pgIDLimit := 63
// 	idLength := len(id)
// 	permIDLength := len("_tmp_") + len(name)
// 	maxIDLength := pgIDLimit - permIDLength
// 	if idLength > maxIDLength {
// 		id = id[0:maxIDLength]
// 		log.Printf("length of _tmp_[id]_[name] is %d, exceeding postgres identifier limit of %d, truncating [id] to %d characters: %s", (permIDLength + idLength), pgIDLimit, maxIDLength, id)
// 	}
// 	// for _, field := range fields {
// 	// 	if field.Kind == "slice" || field.Kind == "map" {
// 	// 		field.StructField.Interface()
// 	// 	}
// 	// }
// 	return Table{
// 		Name:        name,
// 		Type:        rt,
// 		Schema:      schema,
// 		Temporary:   temporary,
// 		UUID:        id,
// 		Constraints: constraints,
// 		Fields:      fields,
// 		Tags:        tags,
// 		Parent:      parent,
// 	}
// }

// type Field struct {
// 	Name        string
// 	Type        string // the type, after dereferencing
// 	Pointer     bool   // since we dereference, was it originally a pointer
// 	Client      string // client name - eg: pg, sqlserver, mysql, oracle
// 	ClientType  string // intended client data type - eg: numeric, text, char[n], etc
// 	PrimaryKey  bool   // is it part of the table's primary key
// 	Unique      bool   // does it have a unique constraint
// 	Nullable    bool   // nullable follows the sql standard of defaulting to true
// 	Constraints []string
// 	// Tags        map[string][]string
// 	StructField reflect.StructField // should this be embedded?
// }

// func StructToFields[T any](value T) []Field {

// 	rt := reflect.TypeOf(value)
// 	structfields := StructFields(value)
// 	fields := []Field{}
// 	for i := 0; i < rt.NumField(); i++ {
// 		sf := structfields[i]
// 		tf := Field{StructField: sf}
// 		field := Field{
// 			Name:        sf.Name,
// 			Type:        fmt.Sprint(tf.UnderlyingType()),
// 			Pointer:     sf.Type.Kind() == reflect.Pointer,
// 			PrimaryKey:  tf.TagKeyIsPositive("primarykey"),
// 			Unique:      tf.TagKeyIsPositive("unique"),
// 			Nullable:    tf.TagKeyIsPositive("nullable"),
// 			Constraints: tf.Tags()["constraints"],
// 			StructField: sf,
// 		}

// 		fields = append(fields, field)
// 	}

// 	return fields
// }

// func (f Field) Kind() reflect.Kind {
// 	sf := f.StructField
// 	kind := sf.Type.Kind()
// 	switch {
// 	case kind == reflect.Pointer:
// 		// this isn't really any more readable when you expand it, is there an easier way?
// 		kind = reflect.Indirect(reflect.New(sf.Type.Elem())).Type().Kind()
// 	}
// 	return kind
// }

// func (f Field) UnderlyingType() reflect.Type {
// 	ft := f.StructField.Type
// 	if ft.Kind() == reflect.Pointer {
// 		return reflect.Indirect(reflect.New(ft.Elem())).Type()
// 	}
// 	return ft
// }

// // ElementType returns the string representation of the element type of slices and maps
// func (f Field) ElementType() reflect.Type {
// 	ft := f.StructField.Type
// 	if f.Kind() == reflect.Slice || f.Kind() == reflect.Map {
// 		return reflect.Indirect(reflect.New(ft.Elem())).Type()
// 	}
// 	return nil
// }

// func (f Field) Tags() map[string][]string {
// 	tagString := string(f.StructField.Tag)
// 	pattern := regexp.MustCompile(`(?m)(?P<key>\w+):\"(?P<value>[^"]+)\"`)
// 	matches := pattern.FindAllStringSubmatch(tagString, -1)
// 	var tkv = map[string][]string{}
// 	for _, match := range matches {
// 		tkv[match[1]] = strings.Split(match[2], ",")
// 	}
// 	return tkv
// }

// // TagKeyIsPostive returns false if the key is missing, or
// // if the first value matches any of false, -, or ""
// func (f Field) TagKeyIsPositive(key string) bool {
// 	values, ok := f.Tags()[key]
// 	// fmt.Println(values[0])
// 	value := ""
// 	if values != nil {
// 		value = strings.ToLower(values[0])
// 	}
// 	return ok && value != "false" && value != "-" && value != ""
// }

// // TagIsTrue returns true if the tag exists and the first value isn't empty, false, or -
// func (f Field) TagIsTrue(key string) bool {
// 	values, ok := f.Tags()[key]
// 	return values != nil && ok && values[0] != "false" && values[0] != "-" && values[0] != ""
// 	// switch values, ok := f.Tags()[key]; {
// 	// case values == nil || !ok:
// 	// 	return false
// 	// default:
// 	// 	return values[0] != "false" && values[0] != "-" && values[0] != ""
// 	// }
// }

// // TagIsFalse returns true only if the tag exists and the first value is false or -
// func (f Field) TagIsFalse(key string) bool {
// 	switch values, ok := f.Tags()[key]; {
// 	case values == nil || !ok:
// 		return false
// 	default:
// 		return values[0] == "false" || values[0] == "-"
// 	}
// }

// func (f Field) TagKeyValueExists(key string, value string) bool {
// 	tvs, ok := f.Tags()[key]
// 	for i, v := range tvs {
// 		tvs[i] = strings.ToLower(v)
// 	}
// 	return ok && slices.Contains(tvs, strings.ToLower(value))
// }

// func TagKeyValues(s string) map[string][]string {

// 	pattern := regexp.MustCompile(`(?m)(?P<key>\w+):\"(?P<value>[^"]+)\"`)
// 	matches := pattern.FindAllStringSubmatch(s, -1)
// 	var out = map[string][]string{}
// 	for _, match := range matches {
// 		out[match[1]] = strings.Split(match[2], ",")
// 	}
// 	return out
// }

// func PrimaryKey(fields []Field) []string {
// 	pk := []string{}
// 	for _, field := range fields {
// 		if field.PrimaryKey {
// 			pk = append(pk, PGName(field))
// 		}
// 	}
// 	return pk
// }

// func QueryFieldName(field Field) string {
// 	if nametags, ok := field.Tags()["qf"]; ok {
// 		return nametags[0]
// 	}
// 	return ""
// }

// func FieldNames(fields []Field) []string {
// 	var fn []string
// 	for _, field := range fields {
// 		fn = append(fn, field.Name)
// 	}
// 	return fn
// }

// func JoinStringSlice(sep string, s []string) string {
// 	return strings.Join(s, sep)
// }

// // func JoinStringSlices(sep string, s ...[]string) string {
// // 	ss := [][]string{}

// // 	for _, sj := range s {
// // 		for _, sji := range sj {

// // 		}
// // 	}
// // }

// // SQLResult combines any number of sql.Result's
// // note that this effectively collapses multiple results
// func SQLResult(res ...sql.Result) Result {
// 	var lis, ras int64
// 	var lies, raes error
// 	for _, r := range res {
// 		if r == nil {
// 			continue
// 		}
// 		li, lie := r.LastInsertId()
// 		ra, rae := r.RowsAffected()
// 		lis = li
// 		ras = ras + ra
// 		lies = fmt.Errorf("[%v]: [%w]", lie, lies)
// 		raes = fmt.Errorf("[%v]: [%w]", rae, raes)
// 	}
// 	return Result{
// 		lastInsertID:      lis,
// 		lastInsertIDError: lies,
// 		rowsAffected:      ras,
// 		rowsAffectedError: raes,
// 	}
// }

// // Result is used to satisfy the sql.Result interface and enable aggregating multiple sql.Results
// type Result struct {
// 	lastInsertID      int64
// 	lastInsertIDError error
// 	rowsAffected      int64
// 	rowsAffectedError error
// }

// func (r *Result) AddResult(res ...sql.Result) {
// 	if r != nil && len(res) > 0 {
// 		rs := []sql.Result{*r}
// 		rs = append(rs, res...)
// 		out := SQLResult(rs...)
// 		*r = out
// 	}
// }

// func (r Result) LastInsertId() (int64, error) {
// 	return r.lastInsertID, r.lastInsertIDError
// }

// func (r Result) RowsAffected() (int64, error) {
// 	return r.rowsAffected, r.rowsAffectedError
// }

// // TableStatement is primarily intended for creating SQL statements from a Table
// // and a template. It includes a handful of funcs, but accepts a funcmap that can
// // overwrite the included funcs
// func TableStatement(table Table, tpl string, funcs template.FuncMap) string {

// 	var buf bytes.Buffer

// 	if err := template.Must(template.
// 		New("").
// 		Option("missingkey=zero").
// 		Funcs(template.FuncMap{
// 			"join":               strings.Join,
// 			"joinss":             JoinStringSlice,
// 			"qfname":             QueryFieldName,
// 			"uniquefields":       UniqueFields,
// 			"fieldswithtagvalue": FieldsWithTagValue,
// 			"fieldnames":         FieldNames,
// 		}).
// 		Funcs(funcs).
// 		Parse(tpl)).
// 		Execute(&buf, table); err != nil {
// 		log.Println(err)
// 		panic(err)
// 	}
// 	return buf.String()
// }

// // UniqueFields is intended to be used in templates and is included in the
// // default TableStatement funcmap as uniquefields
// func UniqueFields(fields []Field) []Field {
// 	f := []Field{}
// 	for _, field := range fields {
// 		if field.Unique {
// 			f = append(f, field)
// 		}
// 	}
// 	return f
// }

// // FieldsWithTagValue returns only those fields with the given key-value pair
// // it is included in the TableStatement funcmap as fieldswithtagvalue
// func FieldsWithTagValue(fields []Field, key, value string) []Field {
// 	f := []Field{}
// 	for _, field := range fields {
// 		tagSlice, ok := field.Tags()[key]
// 		if ok && len(tagSlice) > 0 {
// 			for _, tagi := range tagSlice {
// 				if tagi == value {
// 					f = append(f, field)
// 				}
// 			}
// 		}
// 	}
// 	return f
// }

// // FieldsWithoutTagValue returns only those fields with the given key-value pair
// // it is included in the TableStatement funcmap as fieldswithtagvalue
// func FieldsWithoutTagValue(fields []Field, key, value string) []Field {
// 	f := []Field{}

// 	for _, field := range fields {
// 		tagSlice, ok := field.Tags()[key]
// 		if ok && len(tagSlice) > 0 {
// 			for _, tagi := range tagSlice {
// 				if tagi == value {
// 					f = append(f, field)
// 				}
// 			}
// 		}
// 	}
// 	return f
// }

// // FieldHasTagValue returns true if the key: value exists in the given tag
// // it is included in the TableStatement funcmap as fieldhastagvalue
// func FieldHasTagValue(field Field, key, value string) bool {

// 	tag, ok := field.Tags()[key]
// 	if !ok {
// 		return false
// 	}
// 	for _, tv := range tag {
// 		if tv == value {
// 			return true
// 		}
// 	}
// 	return false
// }

// // Field.HasTagValue is a direct wrap of FieldHasTagValue
// func (f Field) HasTagValue(key, value string) bool {
// 	return FieldHasTagValue(f, key, value)
// }

// // WithTagValue returns fields that test true for Field.HasTagValue
// func (f Fields) WithTagValue(key, value string) Fields {
// 	ff := Fields{}
// 	for _, field := range f {
// 		if field.HasTagValue(key, value) {
// 			ff = append(ff, field)
// 		}
// 	}
// 	return ff
// }

// // WithoutTagValue returns fields that test false for Field.HasTagValue
// func (f Fields) WithoutTagValue(key, value string) Fields {
// 	ff := Fields{}
// 	for _, field := range f {
// 		if !field.HasTagValue(key, value) {
// 			ff = append(ff, field)
// 		}
// 	}

// 	return ff
// }
