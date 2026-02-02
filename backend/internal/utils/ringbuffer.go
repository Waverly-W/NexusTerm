package utils

import (
	"sync"
)

// RingBuffer is a thread-safe circular buffer
type RingBuffer struct {
	data   []byte
	size   int
	start  int
	length int
	mu     sync.RWMutex
}

// NewRingBuffer creates a new RingBuffer with the given size
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]byte, size),
		size: size,
	}
}

// Write writes data to the buffer, overwriting old data if necessary
func (rb *RingBuffer) Write(p []byte) (n int, err error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	n = len(p)
	if n > rb.size {
		// If data is larger than buffer, only keep the last size bytes
		p = p[n-rb.size:]
		n = rb.size
	}

	// Case 1: Data fits in the remaining space at the end
	// Case 2: Data wraps around
	
	// We implementing a simple overwrite logic
	// The 'start' pointer moves if we overwrite
	
	for i := 0; i < n; i++ {
		pos := (rb.start + rb.length) % rb.size
		rb.data[pos] = p[i]
		if rb.length < rb.size {
			rb.length++
		} else {
			rb.start = (rb.start + 1) % rb.size
		}
	}
	
	return len(p), nil
}

// Bytes returns all available data in the buffer
func (rb *RingBuffer) Bytes() []byte {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	out := make([]byte, rb.length)
	if rb.length == 0 {
		return out
	}

	if rb.start+rb.length <= rb.size {
		copy(out, rb.data[rb.start:rb.start+rb.length])
	} else {
		// Wrapped around
		firstPart := rb.size - rb.start
		copy(out, rb.data[rb.start:])
		copy(out[firstPart:], rb.data[:rb.length-firstPart])
	}

	return out
}

// String returns the buffer content as string
func (rb *RingBuffer) String() string {
	return string(rb.Bytes())
}
