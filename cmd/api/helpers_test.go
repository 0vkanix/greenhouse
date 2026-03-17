package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0vkanix/greenlight/internal/assert"
)

func TestWriteJSON(t *testing.T) {
	app := newTestApplication(t)

	t.Run("Custom headers", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		headers := http.Header{}
		headers.Set("X-Test", "Value")

		err := app.writeJSON(rr, r, http.StatusOK, envelope{"message": "test"}, headers)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assert.Equal(t, rr.Code, http.StatusOK)
		assert.Equal(t, rr.Header().Get("X-Test"), "Value")
		assert.Equal(t, rr.Header().Get("Content-Type"), "application/json")
	})
}

func TestReadJSON(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name          string
		body          string
		dst           any
		wantErrString string
	}{
		{
			name: "Valid JSON",
			body: `{"name": "test"}`,
			dst:  &struct{ Name string `json:"name"` }{},
		},
		{
			name:          "Malformed JSON (Syntax Error)",
			body:          `{"name": "test",}`,
			dst:           &struct{ Name string `json:"name"` }{},
			wantErrString: "body contains badly-formed JSON (at character 17)",
		},
		{
			name:          "Unexpected EOF",
			body:          `{"name": "test"`,
			dst:           &struct{ Name string `json:"name"` }{},
			wantErrString: "body contains badly formed JSON",
		},
		{
			name:          "Incorrect JSON type for field",
			body:          `{"name": 123}`,
			dst:           &struct{ Name string `json:"name"` }{},
			wantErrString: `body contains incorrect JSON type for field "name"`,
		},
		{
			name:          "Incorrect JSON type (without field name)",
			body:          `123`,
			dst:           &struct{ Name string `json:"name"` }{},
			wantErrString: `body contains incorrect JSON type (at character 3)`,
		},
		{
			name:          "Empty body",
			body:          ``,
			dst:           &struct{ Name string `json:"name"` }{},
			wantErrString: "body must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()

			err := app.readJSON(rr, r, tt.dst)

			if tt.wantErrString != "" {
				if err == nil {
					t.Fatalf("expected error %q but got nil", tt.wantErrString)
				}
				assert.Equal(t, err.Error(), tt.wantErrString)
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}

	t.Run("Panic on InvalidUnmarshalError", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("readJSON did not panic on non-pointer destination")
			}
		}()

		r, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{}`)))
		rr := httptest.NewRecorder()

		// Passing a non-pointer (struct value) should trigger json.InvalidUnmarshalError
		var dst struct{}
		app.readJSON(rr, r, dst)
	})
}
