package middleware

import (
	"errors"
	"log"
	"time"

	"github.com/Samuel-Ricardo/CapybaraMQ/internal/domain/entity"
)

func WithRetries(maxRetries int) func(entity.Subscriber, entity.Event) error {
	return func(s entity.Subscriber, e entity.Event) error {
		var attempt int

		for {
			err := s.HandleEvent(e)

			log.Printf("Verify error: %v.", err)
			log.Printf("Attempt %d of %d", attempt, maxRetries)

			if err == nil {
				return nil
			}

			attempt++
			if attempt >= maxRetries {
				return errors.New("max retries reached")
			}

			log.Printf("Error on handling event: %v. Retrying...", err)

			time.Sleep(2 * time.Second)
		}
	}
}
