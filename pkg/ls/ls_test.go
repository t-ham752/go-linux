package ls

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewLsFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want *LsFlags
	}{
		{
			name: "フラグなし",
			args: []string{},
			want: &LsFlags{},
		},
		{
			name: "'-a'を渡す",
			args: []string{"-a"},
			want: &LsFlags{ShowAll: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			got, err := newLsFlags(tt.args)
			if err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(tt.want, got); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
		})
	}
}

func getFixedModTime(c clocker) string {
	return "4 24 19:37"
}
func TestRun(t *testing.T) {
	tests := []struct {
		name string
		fs   *LsFlags
		want string
	}{
		{
			name: "フラグなし",
			fs:   &LsFlags{},
			want: `test_file.go
test_file.txt
`,
		},
		{
			name: "'-a'を渡して隠しファイルも表示する",
			fs:   &LsFlags{ShowAll: true},
			want: `.secret
test_file.go
test_file.txt
`,
		},
		{
			name: "'-l'を渡してファイルの詳細を表示する",
			fs:   &LsFlags{ShowDetails: true},
			want: `-rw-r--r-- 1 hamoro staff 0 4 24 19:37 test_file.go
-rw-r--r-- 1 hamoro staff 0 4 24 19:37 test_file.txt
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Chdir("testdata")
			t.Cleanup(func() {
				os.Chdir("..")
			})
			stdout := bytes.Buffer{}
			err := run(tt.fs, &stdout, getFixedModTime)
			if err != nil {
				t.Fatal(err)
			}
			got := stdout.String()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetStat(t *testing.T) {
	tests := []struct {
		name string
		fs   *LsFlags
		want *StatT
	}{
		{
			name: "フラグなし",
			fs:   &LsFlags{},
			want: &StatT{
				Nlink: 1,
				Group: "staff",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Chdir("testdata")
			t.Cleanup(func() {
				os.Chdir("..")
			})
			files, err := os.ReadDir(".")
			if err != nil {
				t.Fatal(err)
			}
			fs, err := files[0].Info()
			if err != nil {
				t.Fatal(err)
			}
			st, err := getStat(fs)
			if err != nil {
				t.Fatal(err)
			}
			opt := cmpopts.IgnoreFields(*tt.want, "Owner")
			if d := cmp.Diff(tt.want, st, opt); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
		})
	}
}

// func TestLs(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		args []string
// 	}{
// 		// {
// 		// 	name: "フラグなし",
// 		// 	args: []string{},
// 		// },
// 		{
// 			name: "'-a'を渡す",
// 			args: []string{"-a"},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
// 			stdout := bytes.Buffer{}
// 			err := Ls(&stdout)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 		})
// 	}
// }
