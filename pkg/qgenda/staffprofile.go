package qgenda

type XStaffProfile struct {
	ExtractDateTime *Time   `json:"-" db:"_extract_date_time"`
	IDHash          *string `json:"-" db:"_id_hash"` // hash of identifying fields: for staff, this is the processed message hash
	StaffKey        *string `json:"-" db:"staffkey"`
	Name            *string `json:"Name,omitempty"`
	ProfileKey      *string `json:"ProfileKey,omitempty"`
	IsViewable      *bool   `json:"IsViewable,omitempty"`
	IsSchedulable   *bool   `json:"IsSchedulable,omitempty"`
}

func (s *XStaffProfile) Process() error {
	return ProcessStruct(s)
}
