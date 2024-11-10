package entity

type Event interface {
	Name() string
	Data() string
}

type SampleEvent struct {
	Message string
}

func (e SampleEvent) Name() string {
	return "SampleEvent"
}

func (e SampleEvent) Data() string {
	return e.Message
}
