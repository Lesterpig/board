package main

import "testing"

func Test_parseConfigString(t *testing.T) {
	tests := []struct {
		name     string
		args     string
		wantDir  string
		wantName string
	}{
		{
			args:     "/etc/config/board.yaml",
			wantDir:  "/etc/config",
			wantName: "board",
		},
		{
			args:     "../config/conf.yaml",
			wantDir:  "../config",
			wantName: "conf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			gotDir, gotName := parseConfigString(tt.args)
			if gotDir != tt.wantDir {
				t.Errorf("parseConfigString() gotDir = %v, want %v", gotDir, tt.wantDir)
			}
			if gotName != tt.wantName {
				t.Errorf("parseConfigString() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}
