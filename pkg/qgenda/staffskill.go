package qgenda

import (
	"fmt"
	"sort"

	"github.com/exiledavatar/gotoolkit/meta"
	"golang.org/x/exp/slices"
)

type StaffSkill struct {
	ExtractDateTime *Time   `json:"-" db:"_extract_date_time" pgtype:"timestamp with time zone"`
	StaffIDHash     *string `json:"-" db:"_staff_id_hash" pgtype:"text" parentprimarykey:"true" idhash:"true"`
	IDHash          *string `json:"-" db:"_id_hash" pgtype:"text" primarykey:"true"` // hash of identifying fields
	// --
	StaffKey          *string `json:"-" db:"staffkey" pgtype:"text" idhash:"true"`
	StaffFirstName    *string `json:"StaffFirstName,omitempty" db:"stafffirstname" pgtype:"text"`
	StaffLastName     *string `json:"StaffLastName,omitempty" db:"stafflastname" pgtype:"text"`
	StaffAbbreviation *string `json:"StaffAbbrev,omitempty" db:"staffabbrev" pgtype:"text"`
	StaffId           *string `json:"StaffId,omitempty" db:"staffid" pgtype:"text"`
	TaskName          *string `json:"TaskName,omitempty" db:"taskname" pgtype:"text" idhash:"true"`
	TaskAbbreviation  *string `json:"TaskAbbrev,omitempty" db:"taskabbrev" pgtype:"text" idhash:"true"`
	TaskId            *string `json:"TaskId,omitempty" db:"taskid" pgtype:"text" idhash:"true"`
	IsSkilledMon      *bool   `json:"IsSkilledMon,omitempty" db:"isskilledmon" pgtype:"boolean" idhash:"true"`
	MonOccurrence     *string `json:"MonOccurrence,omitempty" db:"monoccurrence" pgtype:"text" idhash:"true"`
	IsSkilledTue      *bool   `json:"IsSkilledTue,omitempty" db:"isskilledtue" pgtype:"boolean" idhash:"true"`
	TueOccurrence     *string `json:"TueOccurrence,omitempty" db:"tueoccurrence" pgtype:"text" idhash:"true"`
	IsSkilledWed      *bool   `json:"IsSkilledWed,omitempty" db:"isskilledwed" pgtype:"boolean" idhash:"true"`
	WedOccurrence     *string `json:"WedOccurrence,omitempty" db:"wedoccurrence" pgtype:"text" idhash:"true"`
	IsSkilledThu      *bool   `json:"IsSkilledThu,omitempty" db:"isskilledthu" pgtype:"boolean" idhash:"true"`
	ThuOccurrence     *string `json:"ThuOccurrence,omitempty" db:"thuoccurrence" pgtype:"text" idhash:"true"`
	IsSkilledFri      *bool   `json:"IsSkilledFri,omitempty" db:"isskilledfri" pgtype:"boolean" idhash:"true"`
	FriOccurrence     *string `json:"FriOccurrence,omitempty" db:"frioccurrence" pgtype:"text" idhash:"true"`
	IsSkilledSat      *bool   `json:"IsSkilledSat,omitempty" db:"isskilledsat" pgtype:"boolean" idhash:"true"`
	SatOccurrence     *string `json:"SatOccurrence,omitempty" db:"satoccurrence" pgtype:"text" idhash:"true"`
	IsSkilledSun      *bool   `json:"IsSkilledSun,omitempty" db:"isskilledsun" pgtype:"boolean" idhash:"true"`
	SunOccurrence     *string `json:"SunOccurrence,omitempty" db:"sunoccurrence" pgtype:"text" idhash:"true"`
}

func (s *StaffSkill) Process() error {
	if s.StaffIDHash == nil {
		return fmt.Errorf("cannot process StaffSkill until StaffIDHash is set")
	}

	if err := meta.ProcessStruct(s); err != nil {
		return err
	}

	idh := meta.ToValueMap(*s, "idhash").Hash()
	s.IDHash = &idh
	return nil

	return meta.ProcessStruct(s)
}

type StaffSkills []StaffSkill

func (s *StaffSkills) Sort() *StaffSkills {
	ss := *s
	sort.SliceStable(ss, func(i, j int) bool {
		ski := *(ss[i].StaffKey)
		skj := *(ss[j].StaffKey)
		tni := *(ss[i].TaskName)
		tnj := *(ss[j].TaskName)
		tai := *(ss[i].TaskAbbreviation)
		taj := *(ss[j].TaskAbbreviation)
		tii := *(ss[i].TaskId)
		tij := *(ss[j].TaskId)
		return (ski < skj) ||
			(ski <= skj && tni < tnj) ||
			(tni <= tnj && tai < taj) ||
			(tai <= taj && tii < tij)
	})
	*s = ss
	return s
}

func (s *StaffSkills) Process() error {

	ss := *s
	for i, _ := range ss {
		if err := ss[i].Process(); err != nil {
			return err
		}
	}
	sort.SliceStable(ss, func(i, j int) bool {
		ski := *(ss[i].StaffKey)
		skj := *(ss[j].StaffKey)
		tni := *(ss[i].TaskName)
		tnj := *(ss[j].TaskName)
		tai := *(ss[i].TaskAbbreviation)
		taj := *(ss[j].TaskAbbreviation)
		tii := *(ss[i].TaskId)
		tij := *(ss[j].TaskId)
		return (ski < skj) ||
			(ski <= skj && tni < tnj) ||
			(tni <= tnj && tai < taj) ||
			(tai <= taj && tii < tij)
	})
	ss = slices.CompactFunc(ss, func(s1, s2 StaffSkill) bool {
		return *(s1.IDHash) == *(s2.IDHash) && *(s1.StaffIDHash) == *(s2.StaffIDHash)
	})
	*s = ss

	return nil
}
