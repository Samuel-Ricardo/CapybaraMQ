package middleware

import (
	"log"

	"github.com/Samuel-Ricardo/CapybaraMQ/internal/domain/entity"
)

func WithLogging() func(entity.Subscriber, entity.Event) error {
	return func(s entity.Subscriber, e entity.Event) error {
		log.Printf("[Middleware] - Subscriber: %s, Event: %s", s, e)
		return nil
	}
}
