package assert

import (
	"testing"
)

func TestEqual(t *testing.T) {
	// 1. Success case
	t.Run("Values are equal", func(t *testing.T) {
		Equal(t, 1, 1)
	})

	// 2. Failure case (we use a subtest to simulate failure)
	t.Run("Values are not equal", func(t *testing.T) {
		mockT := new(testing.T)
		Equal(mockT, 1, 2)
		if !mockT.Failed() {
			t.Error("expected Equal to fail for unequal values")
		}
	})
}

func TestStringContains(t *testing.T) {
	// 1. Success case
	t.Run("String contains substring", func(t *testing.T) {
		StringContains(t, "Hello World", "World")
	})

	// 2. Failure case
	t.Run("String does not contain substring", func(t *testing.T) {
		mockT := new(testing.T)
		StringContains(mockT, "Hello World", "Goodbye")
		if !mockT.Failed() {
			t.Error("expected StringContains to fail when substring is missing")
		}
	})
}
