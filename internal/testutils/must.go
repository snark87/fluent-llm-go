package testutils

import "testing"

func Must0(err error) func(t *testing.T) {
	return func(t *testing.T) {
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	}
}

func Must1[T any](value T, err error) func(t *testing.T) T {
	return func(t *testing.T) T {
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		return value
	}
}
