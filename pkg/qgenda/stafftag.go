package qgenda

import "sort"

type XStaffTag struct {
	ExtractDateTime *Time   `json:"-" db:"_extract_date_time"`
	IDHash          *string `json:"-" db:"_id_hash"` // hash of identifying fields: for staff, this is the processed message hash
	StaffKey        *string `json:"-" db:"staffkey"`
	CategoryKey     *int64  `json:"-" db:"categorykey"`
	CategoryName    *string `json:"-" db:"categoryname"`
	TagKey          *int64  `json:"Key" db:"tagkey"`
	TagName         *string `json:"Name" db:"tagname"`
}

func (s *XStaffTag) Process() error {
	return ProcessStruct(s)
}

type XStaffTags struct {
	ExtractDateTime *Time       `json:"-" db:"_extract_date_time"`
	IDHash          *string     `json:"-" db:"_id_hash"` // hash of identifying fields: for staff, this is the processed message hash
	StaffKey        *string     `json:"-" db:"staffkey"`
	CategoryKey     *int64      `json:"CategoryKey"`
	CategoryName    *string     `json:"CategoryName"`
	Tags            []XStaffTag `json:"Tags,omitempty"`
}

func (s *XStaffTags) Process() error {
	for i, _ := range s.Tags {
		s.Tags[i].ExtractDateTime = s.ExtractDateTime
		s.Tags[i].IDHash = s.IDHash
		s.Tags[i].StaffKey = s.StaffKey
		if err := s.Tags[i].Process(); err != nil {
			return err
		}
	}
	sort.SliceStable(s.Tags, func(i, j int) bool {
		return *(s.Tags[i].TagKey) < *(s.Tags[j].TagKey)
	})
	return nil
}
