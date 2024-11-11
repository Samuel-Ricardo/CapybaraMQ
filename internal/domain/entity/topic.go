package entity

type Topic struct {
	Queue       chan Event
	Name        string
	Subscribers []Subscriber
}

func NewTopic(name string) *Topic {
	return &Topic{
		Queue:       make(chan Event, 100),
		Name:        name,
		Subscribers: make([]Subscriber, 0),
	}
}

func (t *Topic) Publish(event Event) {
	t.Queue <- event
}

func (t *Topic) Subscribe(subscriber Subscriber) {
	t.Subscribers = append(t.Subscribers, subscriber)
}
