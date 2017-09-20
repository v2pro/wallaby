package util

import (
	"fmt"
	"testing"
)

func AssertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		return
	}
	message = fmt.Sprintf("%s %v == %v", message, a, b)
	t.Fatal(message)
}

func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	message = fmt.Sprintf("%s %v != %v", message, a, b)
	t.Fatal(message)
}
