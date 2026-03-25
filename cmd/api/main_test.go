package main

// func TestRun(t *testing.T) {
// 	logger := slog.New(slog.DiscardHandler)
//
// 	tests := []struct {
// 		name    string
// 		args    []string
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Missing Flags Error",
// 			args:    []string{"cmd"},
// 			wantErr: true,
// 		},
// 		{
// 			name:    "Valid Config",
// 			args:    []string{"cmd", "-port", "4000", "-env", "test", "-db-dsn", "postgres://foo@bar"},
// 			wantErr: false,
// 		},
// 		{
// 			name:    "Invalid Port Flag",
// 			args:    []string{"cmd", "-port", "abc", "-env", "test", "-db-dsn", "postgres://foo@bar"},
// 			wantErr: true,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := run(tt.args, io.Discard, logger)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
