package application

import (
	"log"
	"sync"

	"github.com/Samuel-Ricardo/CapybaraMQ/internal/domain/entity"
)

type MessageBroker struct {
	Topics      map[string]*entity.Topic
	Middlewares []func(entity.Subscriber, entity.Event) error
	ThreadGuard sync.RWMutex
}

func NewMessageBroker(middlewares ...func(entity.Subscriber, entity.Event) error) *MessageBroker {
	return &MessageBroker{
		Topics:      make(map[string]*entity.Topic),
		Middlewares: middlewares,
	}
}

func (broker *MessageBroker) Subscribe(topicName string, subscriber entity.Subscriber) {
	broker.ThreadGuard.Lock()
	defer broker.ThreadGuard.Unlock()

	topic, exists := broker.Topics[topicName]
	if !exists {
		topic = entity.NewTopic(topicName)
		broker.Topics[topicName] = topic
	}

	topic.Subscribe(subscriber)
}

func (broker *MessageBroker) Publish(topicName string, event entity.Event) {
	broker.ThreadGuard.RLock()
	defer broker.ThreadGuard.RUnlock()

	topic, exists := broker.Topics[topicName]
	if !exists {
		log.Println("Topic not found: ", topicName)
		return
	}

	topic.Publish(event)

	go func() {
		for _, subscriber := range topic.Subscribers {
			go func(subscriber entity.Subscriber) {
				for _, middleware := range broker.Middlewares {
					if err := middleware(subscriber, event); err != nil {
						log.Println("Middleware failed: ", err)
						return
					}
				}

				if err := subscriber.HandleEvent(event); err != nil {
					log.Println("Error on handling event: ", err)
					return
				}
			}(subscriber)
		}
	}()
}

func (broker *MessageBroker) StartConsumer(topicName string) {
	broker.ThreadGuard.RLock()
	defer broker.ThreadGuard.RUnlock()

	topic, exists := broker.Topics[topicName]
	if !exists {
		log.Println("Topic not found: ", topicName)
		return
	}

	go func() {
		for event := range topic.Queue {
			for _, subscriber := range topic.Subscribers {
				go subscriber.HandleEvent(event)
			}
		}
	}()
}
