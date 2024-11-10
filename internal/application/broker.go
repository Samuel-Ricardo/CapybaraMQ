package application

import (
	"log"
	"sync"

	"github.com/Samuel-Ricardo/CapybaraMQ/internal/domain/entity"
)

type MessageBroker struct {
	topics      map[string]*entity.Topic
	threadGuard sync.RWMutex
	middlewares []func(entity.Subscriber, entity.Event) error
}

func NewMessageBroker(middlewares ...func(entity.Subscriber, entity.Event) error) *MessageBroker {
	return &MessageBroker{
		topics:      make(map[string]*entity.Topic),
		middlewares: middlewares,
	}
}

func (broker *MessageBroker) Subscribe(topicName string, subscriber entity.Subscriber) {
	broker.threadGuard.Lock()
	defer broker.threadGuard.Unlock()

	topic, exists := broker.topics[topicName]
	if !exists {
		topic = entity.NewTopic(topicName)
		broker.topics[topicName] = topic
	}

	topic.Subscribe(subscriber)
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
