package entity

type EventHandler struct {
	handler func(event Event) error
}

func (self *EventHandler) HandleEvent(e Event) error {
	return self.handler(e)
}

func NewEventHandler(handler func(event Event) error) *EventHandler {
	return &EventHandler{handler: handler}
}
