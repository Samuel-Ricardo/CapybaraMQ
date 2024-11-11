package test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Samuel-Ricardo/CapybaraMQ/internal/application"
	"github.com/Samuel-Ricardo/CapybaraMQ/internal/application/middleware"
	"github.com/Samuel-Ricardo/CapybaraMQ/internal/domain/entity"
)

func TestEndToEndEventProcessing(t *testing.T) {
	broker := application.NewMessageBroker(middleware.WithLogging())

	var wg sync.WaitGroup
	wg.Add(6)

	handler1 := entity.NewEventHandler(func(e entity.Event) error {
		fmt.Printf("Subscriber 1 handled event: %s\n", e.(entity.SampleEvent).Message)
		wg.Done()
		return nil
	})

	handler2 := entity.NewEventHandler(func(e entity.Event) error {
		fmt.Printf("Subscriber 2 handled event: %s\n", e.(entity.SampleEvent).Message)
		wg.Done()
		return nil
	})

	handler3 := entity.NewEventHandler(func(e entity.Event) error {
		fmt.Printf("Subscriber 3 handled event: %s\n", e.(entity.SampleEvent).Message)
		wg.Done()
		return nil
	})

	broker.Subscribe("topic1", handler1)
	broker.Subscribe("topic1", handler2)
	broker.Subscribe("topic2", handler3)

	broker.StartConsumer("topic1")
	broker.StartConsumer("topic2")

	broker.Publish("topic1", entity.SampleEvent{Message: "E2E Test Event for Topic 1"})
	broker.Publish("topic2", entity.SampleEvent{Message: "E2E Test Event for Topic 2"})

	wg.Wait()
}

func TestEndToEndWithRetriesAndLogging(t *testing.T) {
	broker := application.NewMessageBroker(middleware.WithLogging(), middleware.WithRetries(3))

	var attempts int
	var wg sync.WaitGroup
	wg.Add(2)

	handler := entity.NewEventHandler(func(e entity.Event) error {
		attempts++
		if attempts < 2 {
			fmt.Println("Simulated error in handler")
			return fmt.Errorf("simulated error")
		}
		fmt.Printf("Event processed successfully after %d attempts\n", attempts)
		wg.Done()
		return nil
	})

	broker.Subscribe("topic_with_retries", handler)
	broker.StartConsumer("topic_with_retries")

	broker.Publish("topic_with_retries", entity.SampleEvent{Message: "Event needing retries"})

	wg.Wait()

	if attempts != 3 {
		t.Fatalf("Expected 3 attempts, got %d", attempts)
	}
}
