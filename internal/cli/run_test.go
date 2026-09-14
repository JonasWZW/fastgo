package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{
			name:       "default command",
			wantCode:   0,
			wantStdout: "fastgo: ready",
		},
		{
			name:       "version",
			args:       []string{"-version"},
			wantCode:   0,
			wantStdout: "version=dev commit=none built=unknown",
		},
		{
			name:       "unknown flag",
			args:       []string{"-does-not-exist"},
			wantCode:   2,
			wantStderr: "flag provided but not defined",
		},
		{
			name:       "positional argument",
			args:       []string{"extra"},
			wantCode:   2,
			wantStderr: "unexpected arguments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			gotCode := Run(tt.args, &stdout, &stderr)

			if gotCode != tt.wantCode {
				t.Errorf("Run() code = %d, want %d", gotCode, tt.wantCode)
			}
			if !strings.Contains(stdout.String(), tt.wantStdout) {
				t.Errorf("stdout = %q, want it to contain %q", stdout.String(), tt.wantStdout)
			}
			if !strings.Contains(stderr.String(), tt.wantStderr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}
