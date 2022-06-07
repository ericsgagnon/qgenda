package qgenda

// let's try an explicit approach
type Dataset interface {
	[]Schedule | []StaffMember
	// CreatePGTable()
}

// some dev space here...
