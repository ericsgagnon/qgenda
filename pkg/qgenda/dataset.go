package qgenda

// let's try an explicit approach
type Dataset interface {
	// []Schedule | []StaffMember
	// CreatePGTable()
}

// some dev space here...
type DatasetDev struct {
	Data []any
}

type DatasetDevInterfaceType interface {
	DefaultRequestQueryFields(*RequestQueryFields) *RequestQueryFields
	NewRequest(*RequestQueryFields) *Request
	
}
