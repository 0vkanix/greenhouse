package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func NotEqual[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual == expected {
		t.Errorf("got: %v; expected NOT equal to: %v", actual, expected)
	}
}

func NilError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("got error: %v; want: nil", err)
	}
}

func StringContains(t *testing.T, actual, expectedSubString string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubString) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubString)
	}
}
