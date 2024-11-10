package unit_test

import (
	"testing"

	"github.com/Samuel-Ricardo/CapybaraMQ/internal/application"
	"github.com/Samuel-Ricardo/CapybaraMQ/internal/domain/entity"
)

func TestSubscribe(t *testing.T) {
	broker := application.NewMessageBroker()
	handler := entity.NewEventHandler(func(e entity.Event) error { return nil })

	broker.Subscribe("topic1", handler)

	broker.ThreadGuard.RLock()
	defer broker.ThreadGuard.RUnlock()
	topic, exists := broker.Topics["topic1"]

	if !exists {
		t.Fatalf("Topic 'topic1' should exist")
	}

	if len(topic.Subscribers) != 1 {
		t.Fatalf("Expected 1 subscriber, got %d", len(topic.Subscribers))
	}
}
