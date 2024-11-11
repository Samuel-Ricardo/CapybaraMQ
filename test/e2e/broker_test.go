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
	bus := application.NewMessageBroker(middleware.WithLogging())

	var wg sync.WaitGroup
	wg.Add(3)

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

	bus.Subscribe("topic1", handler1)
	bus.Subscribe("topic1", handler2)
	bus.Subscribe("topic2", handler3)

	bus.StartConsumer("topic1")
	bus.StartConsumer("topic2")

	bus.Publish("topic1", entity.SampleEvent{Message: "E2E Test Event for Topic 1"})
	bus.Publish("topic2", entity.SampleEvent{Message: "E2E Test Event for Topic 2"})

	wg.Wait()
}
