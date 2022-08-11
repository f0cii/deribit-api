package common

import (
	"errors"
)

type EventQueue struct {
	queue []interface{} // queue
	start int           // start index
}

func NewEventQueue(length int) *EventQueue {
	return &EventQueue{
		queue: make([]interface{}, length),
		start: 0,
	}
}

func (q *EventQueue) Insert(data interface{}, offset int) error {
	if offset < 0 || offset >= len(q.queue) {
		return errors.New("offset out of range")
	}

	index := q.correctIndex(offset + q.start)
	q.queue[index] = data
	return nil
}

func (q *EventQueue) GetOffset() int {
	return q.start
}

func (q *EventQueue) GetEvent() interface{} {
	return q.queue[q.start]
}

func (q *EventQueue) Next() {
	q.queue[q.start] = nil
	q.start = q.correctIndex(q.start + 1)
}

func (q *EventQueue) correctIndex(index int) int {
	if index < len(q.queue) {
		return index
	}
	return index % len(q.queue)
}
