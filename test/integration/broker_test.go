package test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Samuel-Ricardo/CapybaraMQ/internal/application"
	"github.com/Samuel-Ricardo/CapybaraMQ/internal/application/middleware"
	"github.com/Samuel-Ricardo/CapybaraMQ/internal/domain/entity"
)

func TestMessageBrokerWithMultipleSubscribers(t *testing.T) {
	broker := application.NewMessageBroker(middleware.WithLogging(), middleware.WithRetries(3))

	var wg sync.WaitGroup
	wg.Add(6)

	handler1 := entity.NewEventHandler(func(e entity.Event) error {
		fmt.Printf("Handled by subscriber 1: %s\n", e.(entity.SampleEvent).Message)
		wg.Done()
		return nil
	})

	handler2 := entity.NewEventHandler(func(e entity.Event) error {
		fmt.Printf("Handled by subscriber 2: %s\n", e.(entity.SampleEvent).Message)
		wg.Done()
		return nil
	})

	broker.Subscribe("topic1", handler1)
	broker.Subscribe("topic1", handler2)

	broker.StartConsumer("topic1")
	broker.Publish("topic1", entity.SampleEvent{Message: "Integration Test Event"})

	wg.Wait()
}

func TestMessageBrokerWithRetriesMiddleware(t *testing.T) {
	broker := application.NewMessageBroker(middleware.WithRetries(2))
	var attempts int

	handler := entity.NewEventHandler(func(e entity.Event) error {
		attempts++
		if attempts < 2 {
			return fmt.Errorf("test attempts: %d", attempts)
		}
		return nil
	})

	broker.Subscribe("topic_with_retries", handler)
	//	broker.StartConsumer("topic_with_retries")
	broker.Publish("topic_with_retries", entity.SampleEvent{Message: "Event with retries"})

	time.Sleep(2 * time.Second)

	if attempts != 2 {
		t.Fatalf("Expected 2 attempts, got %d", attempts)
	}
}
