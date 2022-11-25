package qgenda

import (
	"sort"
)

type XScheduleTag struct {
	ExtractDateTime     *Time   `json:"-" db:"_extract_date_time"`
	ScheduleKey         *string `json:"ScheduleKey,omitempty"`
	LastModifiedDateUTC *Time   `json:"LastModifiedDateUTC,omitempty"`
	CategoryKey         *int64  `json:"CategoryKey,omitempty"`
	CategoryName        *string `json:"CategoryName,omitempty"`
	TagKey              *int64  `json:"Key" db:"tagkey"`
	TagName             *string `json:"Name" db:"tagname"`
}

// func (s *XScheduleTag) UnmarshalJSON(b []byte) error {
// 	// alias technique to avoid infinite recursion
// 	type Alias XScheduleTag
// 	var a Alias

// 	if err := json.Unmarshal(b, &a); err != nil {
// 		return err
// 	}

// 	dest := XScheduleTag(a)
// 	var bb bytes.Buffer
// 	if err := json.Compact(&bb, b); err != nil {
// 		return err
// 	}
// 	// rawMessage := bb.String()
// 	// dest.RawMessage = &rawMessage

// 	*s = dest
// 	return nil
// }

func (st *XScheduleTag) Process() error {
	return ProcessStruct(st)
}

type XScheduleTags struct {
	ExtractDateTime     *Time          `json:"-"`
	ScheduleKey         *string        `json:"-"`
	LastModifiedDateUTC *Time          `json:"-"`
	CategoryKey         *int64         `json:"CategoryKey"`
	CategoryName        *string        `json:"CategoryName"`
	Tags                []XScheduleTag `json:"Tags,omitempty"`
}

// func (tc *XScheduleTagCategory) UnmarshalJSON(b []byte) error {
// 	return json.Unmarshal(b, tc)
// }

// func (tc *XScheduleTagCategory) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(tc.)
// }

// type XScheduleTags struct {
// 	ExtractDateTime     *Time                  `json:"-" db:"_extract_date_time"`
// 	ScheduleKey         *string                `json:"ScheduleKey,omitempty"`
// 	LastModifiedDateUTC *Time                  `json:"LastModifiedDateUTC,omitempty"`
// 	TagCategories       []XScheduleTagCategory `json:"Tags,omitempty"`
// }

// func (sts *XScheduleTags) UnmarshalJSON(b []byte) error {
// 	var tcs []XScheduleTagCategory
// 	if err := json.Unmarshal(b, &tcs); err != nil {
// 		return err
// 	}
// 	sts.TagCategories = tcs
// 	return nil
// }

// func (sts *XScheduleTags) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(sts.Tags)
// }

// func (st *XScheduleTags) Process() error {

// 	if len(st.Tags) > 0 {
// 		tags := []XScheduleTag{}
// 		for _, tag := range st.Tags {
// 			tagp := &tag
// 			if err := tagp.Process(); err != nil {
// 				return err
// 			}

// 			tagp.ExtractDateTime = st.ExtractDateTime
// 			tagp.ScheduleKey = st.ScheduleKey
// 			tagp.LastModifiedDateUTC = st.LastModifiedDateUTC
// 			tags = append(tags, tag)
// 		}
// 		sort.SliceStable(tags, func(i, j int) bool {
// 			return *(tags[i].CategoryKey) < *(tags[j].CategoryKey)
// 		})
// 		st.Tags = tags
// 	}
// 	return nil
// }

func (tc *XScheduleTags) Process() error {
	if len(tc.Tags) > 0 {
		tags := []XScheduleTag{}
		for _, v := range tc.Tags {
			p := &v
			if err := p.Process(); err != nil {
				return err
			}
			p.ExtractDateTime = tc.ExtractDateTime
			p.ScheduleKey = tc.ScheduleKey
			p.LastModifiedDateUTC = tc.LastModifiedDateUTC

			tags = append(tags, v)
		}
		sort.SliceStable(tags, func(i, j int) bool {
			return *(tags[i].TagKey) < *(tags[j].TagKey)
		})
		tc.Tags = tags
	}
	return nil
}

type Test struct {
	V int
}

func (t *Test) Process(i int) error {
	t.V = i
	return nil
}
