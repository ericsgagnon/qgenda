package qgenda

import (
	"sort"
)

type XScheduleTag struct {
	ExtractDateTime     *Time   `json:"-" db:"_extract_date_time"`
	LastModifiedDateUTC *Time   `json:"-" db:"lastmodifieddateutc"`
	ScheduleKey         *string `json:"-" db:"schedulekey"`
	CategoryKey         *int64  `json:"-" db:"categorykey"`
	CategoryName        *string `json:"-" db:"categoryname"`
	TagKey              *int64  `json:"Key" db:"tagkey"`
	TagName             *string `json:"Name" db:"tagname"`
}

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

func (sts *XScheduleTags) Process() error {
	if len(sts.Tags) > 0 {
		// sts.setTagMetaData()
		for i, _ := range sts.Tags {
			sts.Tags[i].ExtractDateTime = sts.ExtractDateTime
			sts.Tags[i].ScheduleKey = sts.ScheduleKey
			sts.Tags[i].LastModifiedDateUTC = sts.LastModifiedDateUTC
			sts.Tags[i].CategoryKey = sts.CategoryKey
			sts.Tags[i].CategoryName = sts.CategoryName
			if err := sts.Tags[i].Process(); err != nil {
				return err
			}
		}
		sort.SliceStable(sts.Tags, func(i, j int) bool {
			return *(sts.Tags[i].TagKey) < *(sts.Tags[j].TagKey)
		})

	}
	return nil
}

// func processScheduleTagsSlice(st []XScheduleTags) error {
// 	if len(st) > 0 {
// 		for i, _ := range st {
// 			if err := st[i].Process(); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func setScheduleTagsMetaData(s *XSchedule, sts []XScheduleTags) error {
// 	if len(sts) > 0 {
// 		for i, _ := range sts {
// 			sts[i].ExtractDateTime = s.ExtractDateTime
// 			sts[i].ScheduleKey = s.ScheduleKey
// 			sts[i].LastModifiedDateUTC = s.LastModifiedDateUTC
// 		}
// 	}
// 	return nil
// }

// // sortScheduleTagsSlice sorts based on CategoryKey
// func sortScheduleTagsSlice(st []XScheduleTags) error {
// 	sort.SliceStable(st, func(i, j int) bool {
// 		return *(st[i].CategoryKey) < *(st[j].CategoryKey)
// 	})
// 	return nil
// }

// func (sts *XScheduleTags) setTagMetaData() {
// 	for i, _ := range sts.Tags {
// 		sts.Tags[i].ExtractDateTime = sts.ExtractDateTime
// 		sts.Tags[i].ScheduleKey = sts.ScheduleKey
// 		sts.Tags[i].LastModifiedDateUTC = sts.LastModifiedDateUTC
// 	}
// }

// 	tags := []XScheduleTag{}
// 	for _, v := range sts.Tags {
// 		// p := &v
// 		if err := v.Process(); err != nil {
// 			return err
// 		}
// 		v.ExtractDateTime = sts.ExtractDateTime
// 		v.ScheduleKey = sts.ScheduleKey
// 		v.LastModifiedDateUTC = sts.LastModifiedDateUTC

// 		tags = append(tags, v)
// 	}
// 	sort.SliceStable(tags, func(i, j int) bool {
// 		return *(tags[i].TagKey) < *(tags[j].TagKey)
// 	})
// 	sts.Tags = tags
// }
