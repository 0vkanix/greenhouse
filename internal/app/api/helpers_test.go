package api

import (
	"bytes"
	"errors"
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
			dst: &struct {
				Name string `json:"name"`
			}{},
		},
		{
			name: "Malformed JSON (Syntax Error)",
			body: `{"name": "test",}`,
			dst: &struct {
				Name string `json:"name"`
			}{},
			wantErrString: "body contains badly-formed JSON (at character 17)",
		},
		{
			name: "Unexpected EOF",
			body: `{"name": "test"`,
			dst: &struct {
				Name string `json:"name"`
			}{},
			wantErrString: "body contains badly formed JSON",
		},
		{
			name: "Incorrect JSON type for field",
			body: `{"name": 123}`,
			dst: &struct {
				Name string `json:"name"`
			}{},
			wantErrString: `body contains incorrect JSON type for field "name"`,
		},
		{
			name: "Incorrect JSON type (without field name)",
			body: `123`,
			dst: &struct {
				Name string `json:"name"`
			}{},
			wantErrString: `body contains incorrect JSON type (at character 3)`,
		},
		{
			name: "Empty body",
			body: ``,
			dst: &struct {
				Name string `json:"name"`
			}{},
			wantErrString: "body must not be empty",
		},
		{
			name: "Unknown field",
			body: `{"name": "test", "unknown": "field"}`,
			dst: &struct {
				Name string `json:"name"`
			}{},
			wantErrString: `body contains unknown key "unknown"`,
		},
		{
			name: "Multiple JSON values",
			body: `{"name": "test"}{"name": "test2"}`,
			dst: &struct {
				Name string `json:"name"`
			}{},
			wantErrString: "body must only contain a single JSON value",
		},
		{
			name: "Body contains extra data",
			body: `{"name": "test"} extra`,
			dst: &struct {
				Name string `json:"name"`
			}{},
			wantErrString: "body must only contain a single JSON value",
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

	t.Run("Body too large", func(t *testing.T) {
		body := `{"name": "` + strings.Repeat("a", 1_048_576) + `"}`
		r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		rr := httptest.NewRecorder()

		var dst struct {
			Name string `json:"name"`
		}
		err := app.readJSON(rr, r, &dst)

		if err == nil {
			t.Fatal("expected error but got nil")
		}
		assert.StringContains(t, err.Error(), "body must not be larger than")
	})

	t.Run("Panic on InvalidUnmarshalError", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("readJSON did not panic on non-pointer destination")
			}
		}()

		r, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{}`)))
		rr := httptest.NewRecorder()

		var dst struct{}
		err := app.readJSON(rr, r, dst)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Default error case", func(t *testing.T) {
		errReader := &errorReader{}
		r, _ := http.NewRequest(http.MethodPost, "/", errReader)
		rr := httptest.NewRecorder()

		var dst struct{}
		err := app.readJSON(rr, r, &dst)

		if err == nil {
			t.Fatal("expected error but got nil")
		}
		assert.Equal(t, err.Error(), "read error")
	})
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}
