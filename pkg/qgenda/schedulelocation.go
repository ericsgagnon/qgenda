package qgenda

import "database/sql"

type XScheduleLocation struct {
	CompanyKey  *string       `json:"CompanyKey,omitempty"`
	LocationKey *int64        `json:"LocationKey,omitempty"`
	ID          *string       `json:"Id,omitempty"`
	Name        *string       `json:"Name,omitempty"`
	Address     *string       `json:"Address,omitempty"`
	Abbrev      *string       `json:"Abbrev,omitempty"`
	Notes       *string       `json:"Notes,omitempty"`
	TimeZone    *string       `json:"TimeZone,omitempty"`
	Tags        []TagCategory `json:"Tags,omitempty"`
}

func (p *XScheduleLocation) UnmarshalJSON(b []byte) error {
	// TODO
	return nil
}

func (v XScheduleLocation) MarshalJSON() ([]byte, error) {
	// TODO
	return nil, nil
}

func (p *XScheduleLocation) Process() error {
	// TODO
	return nil
}

func (v XScheduleLocation) DBCreate() (sql.Result, error) {
	// TODO
	// signature is incomplete
	return nil, nil
}



type XScheduleLocations struct {
	ExtractDateTime     *Time   `json:"-" db:"_extract_date_time"`
	ScheduleKey         *string `json:"ScheduleKey,omitempty"`
	LastModifiedDateUTC *Time   `json:"LastModifiedDateUTC,omitempty"`
	Locations           []Location
}

func (sl *XScheduleLocations) Process() error {
	// TODO
	return nil
}
func ProcessScheduleLocations(sl XScheduleLocations) XScheduleLocations {
	return sl
}
