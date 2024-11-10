package unit_test

import (
	"testing"
	"time"

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

func TestPublishEventToTopic(t *testing.T) {
	broker := application.NewMessageBroker()
	var received bool

	handler := entity.NewEventHandler(func(e entity.Event) error {
		received = true
		return nil
	})

	broker.Subscribe("topic1", handler)
	broker.Publish("topic1", entity.SampleEvent{Message: "Test Event"})

	broker.ThreadGuard.RLock()
	defer broker.ThreadGuard.RUnlock()

	time.Sleep(2 * time.Second)

	if !received {
		t.Fatal("Expected event to be received by subscriber")
	}
}
