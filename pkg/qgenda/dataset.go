package qgenda

// let's try an explicit approach
type Dataset interface {
	// []Schedule | []StaffMember
	// CreatePGTable()
}

// some dev space here...
type DatasetDev struct {
	Data []any
	MetaData
}

type DatasetDevInterfaceType interface {
	DefaultRequestQueryFields(*RequestQueryFields) *RequestQueryFields
	NewRequest(*RequestQueryFields) *Request
}

type MetaData struct {
	ExtractDateTime  Time   `db:"_extract_date_time"`
	RawMessage       string `db:"_raw_message"`
	ProcessedMessage string `db:"_processed_message"`
	HashID           string `db:"_hash_id"`
}
