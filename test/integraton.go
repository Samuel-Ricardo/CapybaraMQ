package test

import (
	"testing"

	"github.com/Samuel-Ricardo/CapybaraMQ/internal/application"
	"github.com/Samuel-Ricardo/CapybaraMQ/internal/domain/entity"
)

func TestMessageBroker(t *testing.T) {
	broker := application.NewMessageBroker()

	handler := entity.NewEventHandler(func(e entity.Event) error {
		t.Logf("Event received: %s", e.Name())
		return nil
	})

	broker.Subscribe("SampleEvent", handler)
	broker.Publish(entity.SampleEvent{Message: "Test Event"})
}
