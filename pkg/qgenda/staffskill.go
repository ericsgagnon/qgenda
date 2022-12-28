package qgenda

type StaffSkill struct {
	ExtractDateTime   *Time   `json:"-" db:"_extract_date_time"`
	IDHash            *string `json:"-" db:"_id_hash"` // hash of identifying fields: for staff, this is the processed message hash
	StaffKey          *string `json:"-" db:"staffkey"`
	StaffFirstName    *string `json:"StaffFirstName,omitempty"`
	StaffLastName     *string `json:"StaffLastName,omitempty"`
	StaffAbbreviation *string `json:"StaffAbbrev,omitempty"`
	StaffId           *string `json:"StaffId,omitempty"`
	TaskName          *string `json:"TaskName,omitempty"`
	TaskAbbreviation  *string `json:"TaskAbbrev,omitempty"`
	TaskId            *string `json:"TaskId,omitempty"`
	IsSkilledMon      *bool   `json:"IsSkilledMon,omitempty"`
	MonOccurrence     *string `json:"MonOccurrence,omitempty"`
	IsSkilledTue      *bool   `json:"IsSkilledTue,omitempty"`
	TueOccurrence     *string `json:"TueOccurrence,omitempty"`
	IsSkilledWed      *bool   `json:"IsSkilledWed,omitempty"`
	WedOccurrence     *string `json:"WedOccurrence,omitempty"`
	IsSkilledThu      *bool   `json:"IsSkilledThu,omitempty"`
	ThuOccurrence     *string `json:"ThuOccurrence,omitempty"`
	IsSkilledFri      *bool   `json:"IsSkilledFri,omitempty"`
	FriOccurrence     *string `json:"FriOccurrence,omitempty"`
	IsSkilledSat      *bool   `json:"IsSkilledSat,omitempty"`
	SatOccurrence     *string `json:"SatOccurrence,omitempty"`
	IsSkilledSun      *bool   `json:"IsSkilledSun,omitempty"`
	SunOccurrence     *string `json:"SunOccurrence,omitempty"`
}

func (s *StaffSkill) Process() error {
	return ProcessStruct(s)
}
