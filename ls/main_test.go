package main

import (
	"bytes"
	"io/fs"
	"testing"
)

type fileDetails struct {
	mode    fs.FileMode
	size    int64
	modTime string
	name    string
}

func TestExec(t *testing.T) {
	tests := []struct {
		name  string
		flags *LsFlags
		want  string
	}{
		{
			name:  "フラグなしで実行",
			flags: &LsFlags{},
			want:  "main.go\nmain_test.go\n",
		},
		{
			name:  "'-a'フラグを渡して実行",
			flags: &LsFlags{showAll: true},
			want:  ".secret\nmain.go\nmain_test.go\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			Exec(w, tt.flags)
			if got := w.String(); got != tt.want {
				t.Errorf("Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}
