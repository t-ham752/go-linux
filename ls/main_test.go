package main

import (
	"io/fs"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type fileDetails struct {
	mode    fs.FileMode
	size    int64
	modTime string
	name    string
}

func TestLs(t *testing.T) {
	tests := []struct {
		name  string
		flags *LsFlags
		want  []*LS
	}{
		{
			name:  "フラグなしで実行",
			flags: &LsFlags{},
			want: []*LS{
				{Name: "test_file.go"},
			},
		},
		{
			name:  "'-a'フラグを渡して実行",
			flags: &LsFlags{showAll: true},
			want: []*LS{
				{Name: ".secret"},
				{Name: "test_file.go"},
			},
		},
		{
			name:  "'-l'フラグを渡して実行",
			flags: &LsFlags{showDetails: true},
			want: []*LS{
				{
					Name:  "test_file.go",
					Nlink: 1,
					Owner: "hamoro",
					Group: "staff",
					Size:  5,
					Mode:  "-rw-r--r--",
				},
			},
		},
		{
			name:  "'-l'と'-a'フラグを渡して実行",
			flags: &LsFlags{showDetails: true, showAll: true},
			want: []*LS{
				{
					Name:  ".secret",
					Nlink: 1,
					Owner: "hamoro",
					Group: "staff",
					Size:  0,
					Mode:  "-rw-r--r--",
				},
				{
					Name:  "test_file.go",
					Nlink: 1,
					Owner: "hamoro",
					Group: "staff",
					Size:  5,
					Mode:  "-rw-r--r--",
				},
			},
		},
		{
			name:  "'-S'と'-l'と'-r'と'-a'フラグを渡して実行",
			flags: &LsFlags{orderBySizeAsc: true, showDetails: true, showAll: true, reverse: true},
			want: []*LS{
				{
					Name:  ".secret",
					Nlink: 1,
					Owner: "hamoro",
					Group: "staff",
					Size:  0,
					Mode:  "-rw-r--r--",
				},
				{
					Name:  "test_file.go",
					Nlink: 1,
					Owner: "hamoro",
					Group: "staff",
					Size:  5,
					Mode:  "-rw-r--r--",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			// ディレクトリを作成
			if err := os.Mkdir("test", 0777); err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll("test")

			// ファイルを作成
			f, err := os.Create("test/test_file.go")
			if err != nil {
				t.Fatal(err)
			}
			_, err = os.Create("test/.secret")
			if err != nil {
				t.Fatal(err)
			}

			// ファイルへの書き込み
			data := []byte("hello")
			_, err = f.Write(data)
			if err != nil {
				t.Fatal(err)
			}

			// test に移動
			if err := os.Chdir("test"); err != nil {
				t.Fatal(err)
			}
			defer os.Chdir("..")

			fs, err := Ls(tt.flags)
			if err != nil {
				t.Fatal(err)
			}

			for i, f := range fs {
				if d := cmp.Diff(f, tt.want[i], cmpopts.IgnoreFields(*f, "ModTime")); len(d) != 0 {
					t.Errorf("differs: (-got +want)\n%s", d)
				}
			}
		})
	}
}
