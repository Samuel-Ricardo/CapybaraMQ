package main

import (
	"log"
	"strconv"
	"time"

	"github.com/Samuel-Ricardo/CapybaraMQ/internal/application"
	"github.com/Samuel-Ricardo/CapybaraMQ/internal/application/middleware"
	"github.com/Samuel-Ricardo/CapybaraMQ/internal/domain/entity"
	"github.com/Samuel-Ricardo/CapybaraMQ/internal/infra"
)

func main() {
	cfg := infra.LoadConfig()
	log.Println("Config: ", cfg)

	broker := application.NewMessageBroker(middleware.WithLogging(), middleware.WithRetries(cfg.RetryCount))

	broker.Subscribe("SampleEvent", entity.NewEventHandler(func(e entity.Event) error {
		log.Println("Event received: ", e.Name())
		return nil
	}))

	for i := 0; i < 100; i++ {
		broker.Publish(entity.SampleEvent{Message: "Test Event: " + strconv.Itoa(i)})
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
