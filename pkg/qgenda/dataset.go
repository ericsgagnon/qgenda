package qgenda

// let's try an explicit approach
type Dataset interface {
	[]Schedule | []StaffMember
}

// some dev space here...

type Schedules []Schedule

func (sch Schedules) Extract() error {
	return nil
}

func (sch Schedules) Process() error {
	return nil
}

func (sch Schedules) Load() error {
	return nil
}
