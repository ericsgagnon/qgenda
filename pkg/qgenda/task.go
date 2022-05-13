package qgenda

type TaskShift struct {
	DayOfWeek     *string         `json:"DayOfWeek,omitempty"`
	StartTime     *TimeOfDay      `json:"StartTime,omitempty"`
	EndTime       *TimeOfDay      `json:"EndTime,omitempty"`
	MaxStaff      *int64          `json:"MaxStaff,omitempty"`
	MinStaff      *int64          `json:"MinStaff,omitempty"`
	OffAfter      *int64          `json:"OffAfter,omitempty"`
	StatCredit    *float64        `json:"StatCredit,omitempty"`
	StaffCount    *int64          `json:"StaffCount,omitempty"`
	EffectiveDate []EffectiveDate `json:"EffectiveDate:,omitempty"`
	IsActive      *bool           `json:"IsActive,omitempty"`
}

type EffectiveDate struct {
	EffectiveFromDate *Date      `json:"EffectiveFromDate,omitempty"`
	EffectiveToDate   *Date      `json:"EffectiveToDate,omitempty"`
	StartTime         *TimeOfDay `json:"StartTime,omitempty"`
	EndTime           *TimeOfDay `json:"EndTime,omitempty"`
	MaxStaff          *int64     `json:"MaxStaff,omitempty"`
	MinStaff          *int64     `json:"MinStaff,omitempty"`
	OffAfter          *int64     `json:"OffAfter,omitempty"`
	StatCredit        *float64   `json:"StatCredit,omitempty"`
}
