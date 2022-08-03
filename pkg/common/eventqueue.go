package common

import (
	"errors"
)

type EventQueue struct {
	q []interface{} // queue
	s int           // start index
}

func NewEventQueue(length int) *EventQueue {
	return &EventQueue{
		q: make([]interface{}, length),
	}
}

func (eq *EventQueue) Insert(e interface{}, offset int) error {
	if offset < 0 || offset >= len(eq.q) {
		return errors.New("offset out of range")
	}

	index := eq.correctIndex(offset + eq.s)
	eq.q[index] = e
	return nil
}

func (eq *EventQueue) GetOffset() int {
	return eq.s
}

func (eq *EventQueue) GetEvent() interface{} {
	return eq.q[eq.s]
}

func (eq *EventQueue) Next() {
	eq.q[eq.s] = nil
	eq.s = eq.correctIndex(eq.s + 1)
}

func (eq *EventQueue) correctIndex(index int) int {
	if index < len(eq.q) {
		return index
	}
	return index % len(eq.q)
}
