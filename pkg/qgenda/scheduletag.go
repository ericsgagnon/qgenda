package qgenda

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/exiledavatar/gotoolkit/meta"
)

type ScheduleTag struct {
	ExtractDateTime *Time   `json:"-" db:"_extract_date_time" pgtype:"timestamp with time zone"`
	ScheduleIDHash  *string `json:"-" db:"_schedule_id_hash" pgtype:"text" parentprimarykey:"true"`
	IDHash          *string `json:"-" db:"_id_hash" pgtype:"text" primarykey:"true"` // hash of identifying fields: schedulekey-lastmodifieddateutc (rfc3339nano)
	// --
	LastModifiedDateUTC *Time   `json:"-" db:"lastmodifieddateutc" idhash:"true" pgtype:"timestamp with time zone"`
	ScheduleKey         *string `json:"-" db:"schedulekey" idhash:"true" pgtype:"text"`
	CategoryKey         *int64  `json:"-" db:"categorykey" idhash:"true" pgtype:"bigint"`
	CategoryName        *string `json:"-" db:"categoryname" pgtype:"text"`
	TagKey              *int64  `json:"Key" db:"tagkey" idhash:"true" pgtype:"bigint"`
	TagName             *string `json:"Name" db:"tagname" pgtype:"text"`
}

// Process should be run as part of Schedule.Process and after
// necessary fields from Schedule have been copied to ScheduleTag
func (st *ScheduleTag) Process() error {

	// fmt.Println(st.ScheduleKey == nil)
	if //st.ExtractDateTime == nil ||
	st.ScheduleKey == nil ||
		st.LastModifiedDateUTC == nil ||
		st.ScheduleIDHash == nil {
		// return fmt.Errorf("cannot process ScheduleTag until ScheduleKey (%s), LastModifiedDateUTC (%s), ScheduleIDHash (%s) are set", *(st.ScheduleKey), *(st.LastModifiedDateUTC), *(st.ScheduleIDHash))
		return fmt.Errorf("cannot process ScheduleTag until ScheduleKey, LastModifiedDateUTC, ScheduleIDHash are set")
	}

	if err := ProcessStruct(st); err != nil {
		return err
	}

	idh := meta.ToValueMap(*st, "idhash").Hash()
	st.IDHash = &idh
	return nil
}

type ScheduleTags []ScheduleTag

func (s *ScheduleTags) Process() error {
	sv := *s
	for i, _ := range sv {
		if err := sv[i].Process(); err != nil {
			return err
		}
	}
	sort.SliceStable(sv, func(i, j int) bool {
		cki := *(sv[i].CategoryKey)
		ckj := *(sv[j].CategoryKey)
		tki := *(sv[i].TagKey)
		tkj := *(sv[j].TagKey)
		return (cki < ckj) || (cki <= ckj && tki < tkj)
	})
	*s = sv
	return nil
}

type scheduleTagCategory struct {
	ExtractDateTime     *Time         `json:"-"`
	ScheduleKey         *string       `json:"-"`
	LastModifiedDateUTC *Time         `json:"-"`
	CategoryKey         *int64        `json:"CategoryKey"`
	CategoryName        *string       `json:"CategoryName"`
	Tags                []ScheduleTag `json:"Tags,omitempty"`
}

func (s *ScheduleTags) UnmarshalJSON(data []byte) error {
	tagCats := []scheduleTagCategory{}
	if err := json.Unmarshal(data, &tagCats); err != nil {
		return err
	}
	dest := []ScheduleTag{}
	for _, tc := range tagCats {
		for _, tag := range tc.Tags {
			st := ScheduleTag{
				CategoryKey:  tc.CategoryKey,
				CategoryName: tc.CategoryName,
				TagKey:       tag.TagKey,
				TagName:      tag.TagName,
			}
			dest = append(dest, st)
		}
	}
	*s = dest

	return nil
}

func (s ScheduleTags) MarshalJSON() ([]byte, error) {
	groupedTags := map[int64]scheduleTagCategory{}
	for _, v := range s {
		switch tc, ok := groupedTags[*v.CategoryKey]; ok {
		case true:
			tc.Tags = append(tc.Tags, v)
			groupedTags[*v.CategoryKey] = tc

		case false:
			groupedTags[*v.CategoryKey] = scheduleTagCategory{
				CategoryKey:  tc.CategoryKey,
				CategoryName: tc.CategoryName,
				Tags:         []ScheduleTag{v},
			}
		}

	}

	tagCatSlice := []scheduleTagCategory{}
	for _, tc := range groupedTags {
		tagCatSlice = append(tagCatSlice, tc)
	}

	return json.Marshal(tagCatSlice)
}

// type ScheduleTags struct {
// 	ExtractDateTime     *Time         `json:"-"`
// 	ScheduleKey         *string       `json:"-"`
// 	LastModifiedDateUTC *Time         `json:"-"`
// 	CategoryKey         *int64        `json:"CategoryKey"`
// 	CategoryName        *string       `json:"CategoryName"`
// 	Tags                []ScheduleTag `json:"Tags,omitempty"`
// }

// func (sts *ScheduleTags) Process() error {
// 	if len(sts.Tags) > 0 {
// 		// sts.setTagMetaData()
// 		for i, _ := range sts.Tags {
// 			sts.Tags[i].ExtractDateTime = sts.ExtractDateTime
// 			sts.Tags[i].ScheduleKey = sts.ScheduleKey
// 			sts.Tags[i].LastModifiedDateUTC = sts.LastModifiedDateUTC
// 			sts.Tags[i].CategoryKey = sts.CategoryKey
// 			sts.Tags[i].CategoryName = sts.CategoryName
// 			if err := sts.Tags[i].Process(); err != nil {
// 				return err
// 			}
// 		}
// 		sort.SliceStable(sts.Tags, func(i, j int) bool {
// 			return *(sts.Tags[i].TagKey) < *(sts.Tags[j].TagKey)
// 		})

// 	}
// 	return nil
// }
