package util

import (
	"fmt"
	"testing"
)

// AssertNotEqual test not equal or fail
func AssertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		return
	}
	message = fmt.Sprintf("%s %v == %v", message, a, b)
	t.Fatal(message)
}

// AssertEqual test equal or fail
func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	message = fmt.Sprintf("%s %v != %v", message, a, b)
	t.Fatal(message)
}
