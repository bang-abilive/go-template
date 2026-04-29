package main

import (
	"io"
	"testing"
)

func TestParseArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		want    bool
		wantErr bool
	}{
		{name: "default no seed", args: nil, want: false, wantErr: false},
		{name: "empty no seed", args: []string{}, want: false, wantErr: false},
		{name: "seed short", args: []string{"-seed"}, want: true, wantErr: false},
		{name: "seed long form", args: []string{"-seed=true"}, want: true, wantErr: false},
		{name: "seed false explicit", args: []string{"-seed=false"}, want: false, wantErr: false},
		{name: "unknown flag", args: []string{"-unknown"}, want: false, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseArgs(tt.args, io.Discard)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
