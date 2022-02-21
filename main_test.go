package main

import "testing"

func TestRun(t *testing.T) {

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "success",
			args:    []string{"google.com"},
			wantErr: false,
		},
		{
			name:    "invalid-url",
			args:    []string{"invalid-url"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Tool{
				Parallel: 1,
			}
			err := s.Run(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}

}
