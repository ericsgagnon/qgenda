package qgenda

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/exiledavatar/gotoolkit/meta"
	"golang.org/x/exp/slices"
)

type StaffTag struct {
	ExtractDateTime *Time   `json:"-" db:"_extract_date_time" pgtype:"timestamp with time zone"`
	StaffIDHash     *string `json:"-" db:"_staff_id_hash" pgtype:"text" parentprimarykey:"true" idhash:"true"`
	IDHash          *string `json:"-" db:"_id_hash" pgtype:"text" primarykey:"true"` // hash of identifying fields
	// --
	StaffKey     *string `json:"-" db:"staffkey" pgtype:"text" idhash:"true"`
	CategoryKey  *int64  `json:"-" db:"categorykey" pgtype:"bigint" idhash:"true"`
	CategoryName *string `json:"-" db:"categoryname" pgtype:"text" idhash:"true"`
	TagKey       *int64  `json:"Key" db:"tagkey" pgtype:"bigint" idhash:"true"`
	TagName      *string `json:"Name" db:"tagname" pgtype:"text" idhash:"true"`
}

func (s *StaffTag) Process() error {
	if s.StaffIDHash == nil {
		return fmt.Errorf("cannot process StaffTag until StaffIDHash is set")
	}

	if err := meta.ProcessStruct(s); err != nil {
		return err
	}

	idh := meta.ToValueMap(*s, "idhash").Hash()
	s.IDHash = &idh
	return nil
}

type StaffTags []StaffTag

func (s *StaffTags) Sort() *StaffTags {
	sv := *s
	sort.SliceStable(sv, func(i, j int) bool {
		ski := *(sv[i].StaffKey)
		skj := *(sv[j].StaffKey)
		cki := *(sv[i].CategoryKey)
		ckj := *(sv[j].CategoryKey)
		tki := *(sv[i].TagKey)
		tkj := *(sv[j].TagKey)
		return (ski < skj) ||
			(ski <= skj && cki < ckj) ||
			(cki <= ckj && tki < tkj)
	})
	*s = sv
	return s
}

func (s *StaffTags) Process() error {
	sv := *s
	for i, _ := range sv {
		if err := sv[i].Process(); err != nil {
			return err
		}
	}
	sort.SliceStable(sv, func(i, j int) bool {
		ski := *(sv[i].StaffKey)
		skj := *(sv[j].StaffKey)
		cki := *(sv[i].CategoryKey)
		ckj := *(sv[j].CategoryKey)
		tki := *(sv[i].TagKey)
		tkj := *(sv[j].TagKey)
		return (ski < skj) ||
			(ski <= skj && cki < ckj) ||
			(cki <= ckj && tki < tkj)
	})
	sv = slices.CompactFunc(sv, func(s1, s2 StaffTag) bool {
		return *(s1.IDHash) == *(s2.IDHash) && *(s1.StaffIDHash) == *(s2.StaffIDHash)
	})

	*s = sv
	return nil
}

type staffTagCategory struct {
	CategoryKey  *int64     `json:"CategoryKey" db:"categorykey" pgtype:"numeric"`
	CategoryName *string    `json:"CategoryName" db:"categoryname" pgtype:"text"`
	Tags         []StaffTag `json:"Tags,omitempty"`
}

func (s *StaffTags) UnmarshalJSON(data []byte) error {
	tagCats := []staffTagCategory{}
	if err := json.Unmarshal(data, &tagCats); err != nil {
		return err
	}
	dest := []StaffTag{}
	for _, tc := range tagCats {
		for _, tag := range tc.Tags {
			st := StaffTag{
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

func (s StaffTags) MarshalJSON() ([]byte, error) {
	groupedTags := map[int64]staffTagCategory{}
	for _, v := range s {
		switch tc, ok := groupedTags[*v.CategoryKey]; ok {
		case true:
			tc.Tags = append(tc.Tags, v)
			groupedTags[*v.CategoryKey] = tc

		case false:
			groupedTags[*v.CategoryKey] = staffTagCategory{
				CategoryKey:  tc.CategoryKey,
				CategoryName: tc.CategoryName,
				Tags:         []StaffTag{v},
			}
		}

	}

	tagCatSlice := []staffTagCategory{}
	for _, tc := range groupedTags {
		tagCatSlice = append(tagCatSlice, tc)
	}

	return json.Marshal(tagCatSlice)
}
