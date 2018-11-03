package main

import (
	"errors"
)

// Queue is the Go implementation of Queue
type Queue []interface{}

// Offer adds an element to the back of this queue.
func (queue *Queue) Offer(element interface{}) {
	*queue = append(*queue, element)
}

// Poll removes the head element of this queue. If queue is empty, it returns
// -1 and an error.
func (queue *Queue) Poll() (interface{}, error) {
	if len(*queue) > 0 {
		polled := (*queue)[0]
		*queue = (*queue)[1:]
		return polled, nil
	}
	return -1, errors.New("queue is empty")
}

// Peek returns the head element of this queue. If queue is empty, it returns
// -1 and an error.
func (queue *Queue) Peek() (interface{}, error) {
	if len(*queue) > 0 {
		return (*queue)[0], nil
	}
	return -1, errors.New("queue is empty")
}
