package middleware

import (
	"errors"
	"time"

	"github.com/Samuel-Ricardo/CapybaraMQ/internal/domain/entity"
)

func WithRetries(maxRetries int) func(entity.Subscriber, entity.Event) error {
	return func(s entity.Subscriber, e entity.Event) error {
		var attempt int

		for {
			err := s.HandleEvent(e)

			if err == nil {
				return nil
			}

			attempt++
			if attempt >= maxRetries {
				return errors.New("max retries reached")
			}

			time.Sleep(2 * time.Second)
		}
	}
}
