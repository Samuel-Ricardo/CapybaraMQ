package application

import (
	"log"
	"sync"

	"github.com/Samuel-Ricardo/CapybaraMQ/internal/domain/entity"
)

type MessageBroker struct {
	subscribers map[string][]entity.Subscriber
	threadGuard sync.RWMutex
	middlewares []func(entity.Subscriber, entity.Event) error
}

func NewMessageBroker(middlewares ...func(entity.Subscriber, entity.Event) error) *MessageBroker {
	return &MessageBroker{
		subscribers: make(map[string][]entity.Subscriber),
		middlewares: middlewares,
	}
}

func (broker *MessageBroker) Subscribe(eventType string, subscriber entity.Subscriber) {
	broker.threadGuard.Lock()
	defer broker.threadGuard.Unlock()

	broker.subscribers[eventType] = append(broker.subscribers[eventType], subscriber)
}

func (broker *MessageBroker) Publish(event entity.Event) {
	broker.threadGuard.RLock()
	defer broker.threadGuard.RUnlock()

	if subscribers, ok := broker.subscribers[event.Name()]; ok {
		for _, subscriber := range subscribers {
			go func(subscriber entity.Subscriber, event entity.Event) {
				for _, middleware := range broker.middlewares {
					if err := middleware(subscriber, event); err != nil {
						log.Println("Middleware failed: ", err)
						return
					}
				}

				if err := subscriber.HandleEvent(event); err != nil {
					log.Println("Error on handling event: ", err)
				}
			}(subscriber, event)
		}
	}
}
