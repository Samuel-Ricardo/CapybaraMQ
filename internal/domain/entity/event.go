package entity

type Event interface {
	Name() string
}

type SampleEvent struct {
	Message string
}

func (e SampleEvent) Name() string {
	return "SampleEvent"
}
