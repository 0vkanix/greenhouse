package main

import (
	"io"
	"log/slog"
	"testing"
)

func TestRun(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "Default Config",
			args:    []string{"cmd", "-env", "test-config"},
			wantErr: false,
		},
		{
			name:    "Custom Port",
			args:    []string{"cmd", "-port", "8080", "-env", "test-config"},
			wantErr: false,
		},
		{
			name:    "Invalid Port Flag",
			args:    []string{"cmd", "-port", "abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.args, io.Discard, logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
