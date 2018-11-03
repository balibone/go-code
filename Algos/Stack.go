package main

import (
	"errors"
)

// Stack is the Go implementation of stack
type Stack []interface{}

// Push ...
func (s *Stack) Push(element interface{}) {
	*s = append(*s, element)
}

// Pop removes the last element of this stack. If stack is empty, it returns
// -1 and an error.
func (s *Stack) Pop() (interface{}, error) {
	if len(*s) > 0 {
		popped := (*s)[len(*s)-1]
		*s = (*s)[:len(*s)-1]
		return popped, nil
	}
	return -1, errors.New("stack is empty")
}

// Peek returns the topmost element of the stack. If stack is empty, it returns
// -1 and an error.
func (s *Stack) Peek() (interface{}, error) {
	if len(*s) > 0 {
		return (*s)[len(*s)-1], nil
	}
	return -1, errors.New("stack is empty")
}
