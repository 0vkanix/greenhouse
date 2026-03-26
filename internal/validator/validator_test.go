package validator

import (
	"regexp"
	"testing"
)

func TestMatches(t *testing.T) {
	rx := regexp.MustCompile("^[a-z]+$")

	if !Matches("abc", rx) {
		t.Error("expected 'abc' to match [a-z]+")
	}
	if Matches("123", rx) {
		t.Error("expected '123' NOT to match [a-z]+")
	}
}

func TestUnique(t *testing.T) {
	if !Unique([]string{"a", "b", "c"}) {
		t.Error("expected unique slice to be true")
	}
	if Unique([]string{"a", "b", "a"}) {
		t.Error("expected non-unique slice to be false")
	}
}

func TestAddError(t *testing.T) {
	v := New()
	v.AddError("key", "message")

	if v.Valid() {
		t.Error("expected validator to be invalid after adding error")
	}
	if v.Errors["key"] != "message" {
		t.Errorf("expected message %q, got %q", "message", v.Errors["key"])
	}

	v.AddError("key", "new message")
	if v.Errors["key"] != "message" {
		t.Error("expected first error message to be preserved")
	}
}

func TestCheck(t *testing.T) {
	v := New()
	v.Check(true, "key", "message")
	if !v.Valid() {
		t.Error("expected true check to be valid")
	}

	v.Check(false, "key", "message")
	if v.Valid() {
		t.Error("expected false check to be invalid")
	}
}

func TestNew(t *testing.T) {
	v := New()
	if !v.Valid() {
		t.Error("expected new validator to be valid")
	}
	if len(v.Errors) != 0 {
		t.Errorf("expected 0 errors, got %d", len(v.Errors))
	}
}

func TestPermittedValue(t *testing.T) {
	if !PermittedValue("a", "a", "b", "c") {
		t.Error("expected 'a' to be permitted in [a, b, c]")
	}
	if PermittedValue("x", "a", "b", "c") {
		t.Error("expected 'x' NOT to be permitted in [a, b, c]")
	}
}
