package qgenda

import (
	"fmt"
	"sort"

	"github.com/exiledavatar/gotoolkit/meta"
	"golang.org/x/exp/slices"
)

type StaffProfile struct {
	ExtractDateTime *Time   `json:"-" db:"_extract_date_time" pgtype:"timestamp with time zone"`
	StaffIDHash     *string `json:"-" db:"_staff_id_hash" pgtype:"text" parentprimarykey:"true" idhash:"true"`
	IDHash          *string `json:"-" db:"_id_hash" pgtype:"text" primarykey:"true"` // hash of identifying fields
	// --
	StaffKey      *string `json:"-" db:"staffkey" pgtype:"text" idhash:"true"`
	Name          *string `json:"Name,omitempty" db:"name" pgtype:"text" idhash:"true"`
	ProfileKey    *string `json:"ProfileKey,omitempty" db:"profilekey" pgtype:"text" idhash:"true"`
	IsViewable    *bool   `json:"IsViewable,omitempty" db:"isviewable" pgtype:"boolean" idhash:"true"`
	IsSchedulable *bool   `json:"IsSchedulable,omitempty" db:"isschedulable" pgtype:"boolean" idhash:"true"`
}

func (s *StaffProfile) Process() error {
	if s.StaffIDHash == nil {
		return fmt.Errorf("cannot process StaffTag until StaffIDHash is set")
	}

	if err := meta.ProcessStruct(s); err != nil {
		return err
	}

	idh := meta.ToValueMap(*s, "idhash").Hash()
	s.IDHash = &idh
	return meta.ProcessStruct(s)
}

type StaffProfiles []StaffProfile

func (s *StaffProfiles) Sort() *StaffProfiles {
	ss := *s
	sort.SliceStable(ss, func(i, j int) bool {
		ski := *(ss[i].StaffKey)
		skj := *(ss[j].StaffKey)
		pki := *(ss[i].ProfileKey)
		pkj := *(ss[j].ProfileKey)
		return (ski < skj) ||
			(ski <= skj && pki < pkj)
	})
	*s = ss
	return s
}

func (s *StaffProfiles) Process() error {

	ss := *s
	for i, _ := range ss {
		if err := ss[i].Process(); err != nil {
			return err
		}
	}
	sort.SliceStable(ss, func(i, j int) bool {
		ski := *(ss[i].StaffKey)
		skj := *(ss[j].StaffKey)
		pki := *(ss[i].ProfileKey)
		pkj := *(ss[j].ProfileKey)
		return (ski < skj) ||
			(ski <= skj && pki < pkj)
	})
	ss = slices.CompactFunc(ss, func(s1, s2 StaffProfile) bool {
		return *(s1.IDHash) == *(s2.IDHash) && *(s1.StaffIDHash) == *(s2.StaffIDHash)
	})

	*s = ss

	return nil
}
