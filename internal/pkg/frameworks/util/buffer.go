package util

import (
	"errors"
	"strings"
)

// CircularBuffer represents a fixed-size buffer of strings.
type CircularBuffer struct {
	buffer []string
	size   int
	index  int
	full   bool
}

// NewCircularBuffer initializes a new circular buffer with the specified size.
func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		buffer: make([]string, size),
		size:   size,
	}
}

// Add appends a string to the circular buffer. If the buffer is full, it overwrites the oldest entry.
func (cb *CircularBuffer) Add(s string) {
	cb.buffer[cb.index] = s
	cb.index = (cb.index + 1) % cb.size
	if cb.index == 0 {
		cb.full = true
	}
}

// Get returns all strings in the buffer, starting from the oldest and moving to the newest.
func (cb *CircularBuffer) Get() ([]string, error) {
	if !cb.full {
		return nil, errors.New("buffer not yet full")
	}
	result := make([]string, 0, cb.size)
	start := cb.index - len(cb.buffer)
	if start < 0 {
		start += cb.size
	}
	for i := 0; i < len(cb.buffer); i++ {
		result = append(result, cb.buffer[(start+i)%cb.size])
	}
	return result, nil
}

// Join returns a single string by joining the strings in the buffer with spaces.
func (cb *CircularBuffer) Join() (string, error) {
	stringSlice, err := cb.Get()
	if err != nil {
		return "", err
	}

	return strings.Join(stringSlice, " "), nil
}
