package application

import (
	"container/list"
	"sync"
)

type EventQueue struct {
	queue       *list.List
	threadGuard sync.Mutex
}

func NewEventQueue() *EventQueue {
	return &EventQueue{queue: list.New()}
}

func (eq *EventQueue) Enqueue(event interface{}) {
	eq.threadGuard.Lock()
	defer eq.threadGuard.Unlock()

	eq.queue.PushBack(event)
}

func (eq *EventQueue) Dequeue() (interface{}, bool) {
	eq.threadGuard.Lock()
	defer eq.threadGuard.Unlock()

	elem := eq.queue.Front()

	if elem != nil {
		eq.queue.Remove(elem)
		return elem.Value, true
	}

	return nil, false
}
