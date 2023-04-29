package cat

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "one arg",
			args: []string{"cmd", "testdata/sample1.txt"},
			want: "okane",
		},
		{
			name: "two args",
			args: []string{"cmd", "testdata/sample1.txt", "testdata/sample2.txt"},
			want: "okanehoshii",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := bytes.Buffer{}
			err := run(tt.args, &stdout)
			if err != nil {
				t.Errorf("run() error = %v", err)
			}
			if got := stdout.String(); got != tt.want {
				t.Errorf("run() = %v, want %v", got, tt.want)
			}
		})
	}
}
