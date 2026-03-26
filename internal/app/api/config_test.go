package api

import (
	"testing"
	"time"

	"github.com/0vkanix/greenlight/internal/assert"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Config
		wantErr bool
	}{
		{
			name: "Valid flags",
			args: []string{"cmd", "-port", "4000", "-db-dsn", "postgres://foo@bar"},
			want: Config{
				Port: 4000,
				Env:  "development",
				DB: struct {
					DSN          string
					MaxOpenConns int
					MaxIdleConns int
					MaxIdleTime  time.Duration
				}{
					DSN: "postgres://foo@bar",
				},
			},
			wantErr: false,
		},
		{
			name: "Custom values",
			args: []string{
				"cmd",
				"-port", "8080",
				"-env", "production",
				"-db-dsn", "postgres://prod@db",
				"-db-max-open-conns", "50",
				"-db-max-idle-conns", "10",
				"-db-max-idle-time", "5m",
			},
			want: Config{
				Port: 8080,
				Env:  "production",
				DB: struct {
					DSN          string
					MaxOpenConns int
					MaxIdleConns int
					MaxIdleTime  time.Duration
				}{
					DSN:          "postgres://prod@db",
					MaxOpenConns: 50,
					MaxIdleConns: 10,
					MaxIdleTime:  5 * time.Minute,
				},
			},
			wantErr: false,
		},
		{
			name:    "Missing port",
			args:    []string{"cmd", "-db-dsn", "postgres://foo@bar"},
			wantErr: true,
		},
		{
			name:    "Missing DSN",
			args:    []string{"cmd", "-port", "4000"},
			wantErr: true,
		},
		{
			name:    "Invalid port type",
			args:    []string{"cmd", "-port", "abc", "-db-dsn", "postgres://foo@bar"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFlags(tt.args)

			if tt.wantErr {
				if err == nil {
					t.Error("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.Equal(t, got.Port, tt.want.Port)
			assert.Equal(t, got.Env, tt.want.Env)
			assert.Equal(t, got.DB.DSN, tt.want.DB.DSN)
			assert.Equal(t, got.DB.MaxOpenConns, tt.want.DB.MaxOpenConns)
			assert.Equal(t, got.DB.MaxIdleConns, tt.want.DB.MaxIdleConns)
			assert.Equal(t, got.DB.MaxIdleTime, tt.want.DB.MaxIdleTime)
		})
	}
}
