package lazy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/snark87/fluentllm/internal/testutils"
)

func TestNew(t *testing.T) {
	var counter atomic.Int32
	lazyString := New(func() string {
		counter.Add(1)
		return "Hello, World!"
	})

	if lazyString.Get() != "Hello, World!" {
		t.Fatalf("Expected 'Hello, World!', got '%s'", lazyString.Get())
	}

	if counter.Load() != 1 {
		t.Fatalf("Expected counter to be 1, got %d", counter.Load())
	}

	if lazyString.Get() != "Hello, World!" {
		t.Fatalf("Expected 'Hello, World!', got '%s'", lazyString.Get())
	}

	if counter.Load() != 1 {
		t.Fatalf("Expected counter to be 1, got %d", counter.Load())
	}
}

func TestLazyWithError_WhenNoError(t *testing.T) {
	ctx := context.Background()
	var counter atomic.Int32
	lazyString := NewWithError(func(context.Context) (string, error) {
		counter.Add(1)
		return "Hello, World!", nil
	})
	value := testutils.Must1(lazyString.Get(ctx))(t)
	if value != "Hello, World!" {
		t.Fatalf("Expected 'Hello, World!', got '%s'", value)
	}

	if counter.Load() != 1 {
		t.Fatalf("Expected counter to be 1, got %d", counter.Load())
	}
	value = testutils.Must1(lazyString.Get(ctx))(t)
	if value != "Hello, World!" {
		t.Fatalf("Expected 'Hello, World!', got '%s'", value)
	}

	if counter.Load() != 1 {
		t.Fatalf("Expected counter to be 1, got %d", counter.Load())
	}
}
