package qgenda

// Parameters is a key-value map to represent arguments
// it is generally used to pass arguments for getting or sending
// data in data models
type Parameters map[any]any
type Data any

type Model interface {
	Config() *any
	Request(p Parameters) *Request
	Data() []Data
}

type ScheduleModel struct {
	
}

func (sm *ScheduleModel) SourceConfig(p Parameters) {
	
}

