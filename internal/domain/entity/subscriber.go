package entity

type Subscriber interface {
	HandleEvent(event Event) error
}
