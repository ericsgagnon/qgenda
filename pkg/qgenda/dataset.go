package qgenda

import (
	"crypto/sha1"
	"fmt"
)

// let's try an explicit approach
type Dataset interface {
	// []Schedule | []StaffMember
	// CreatePGTable()
}

// some dev space here...
type DatasetDev struct {
	Config
	MetaData
	Data []any
}

type DatasetConfig struct {
	Name string
	// EndpointName
}

type DatasetDevInterfaceType interface {
	DefaultRequestConfig(*RequestConfig) *RequestConfig
	NewRequest(*RequestConfig) *Request
}

type MetaData struct {
	ExtractDateTime  Time   `db:"_extract_date_time"`
	RawMessage       string `db:"_raw_message"`
	ProcessedMessage string `db:"_processed_message"`
	HashID           string `db:"_hash_id"`
}

func Hash[V []byte | string](b V) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(b)))
}
