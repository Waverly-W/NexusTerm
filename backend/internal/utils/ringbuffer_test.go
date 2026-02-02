package utils

import (
	"testing"
)

func TestRingBuffer(t *testing.T) {
	rb := NewRingBuffer(5)

	// Test 1: Write within capacity
	rb.Write([]byte("123"))
	if rb.String() != "123" {
		t.Errorf("Expected '123', got '%s'", rb.String())
	}

	// Test 2: Fill capacity
	rb.Write([]byte("45"))
	if rb.String() != "12345" {
		t.Errorf("Expected '12345', got '%s'", rb.String())
	}

	// Test 3: Overwrite (1 byte)
	rb.Write([]byte("6"))
	// Buffer size 5. Old: 12345. New: 23456.
	if rb.String() != "23456" {
		t.Errorf("Expected '23456', got '%s'", rb.String())
	}

	// Test 4: Overwrite multiple bytes
	rb.Write([]byte("78"))
	// Old: 23456. Write 78 -> 45678.
	if rb.String() != "45678" {
		t.Errorf("Expected '45678', got '%s'", rb.String())
	}

	// Test 5: Overwrite more than size
	rb.Write([]byte("ABCDEF"))
	// Should keep last 5: BCDEF
	if rb.String() != "BCDEF" {
		t.Errorf("Expected 'BCDEF', got '%s'", rb.String())
	}
}
