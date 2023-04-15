package ls

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"sort"
	"strconv"
	"syscall"
)

type StatT struct {
	// exportしないとcmpで比較できない
	Nlink uint16
	Owner string
	Group string
}

func getStat(fs fs.FileInfo) (*StatT, error) {
	var owner, group string
	stat, ok := fs.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("syscall.Stat_t is not *syscall.Stat_t")
	}

	// ユーザID
	uid := strconv.Itoa(int(stat.Uid))
	u, err := user.LookupId(uid)
	if err != nil {
		return nil, err
	} else {
		owner = u.Username
	}

	// グループID
	gid := strconv.Itoa(int(stat.Gid))
	g, err := user.LookupGroupId(gid)
	if err != nil {
		return nil, err
	} else {
		group = g.Name
	}

	return &StatT{
		Nlink: stat.Nlink,
		Owner: owner,
		Group: group,
	}, nil
}

var (
	commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
)

type LsFlags struct {
	// exportしないとcmpで比較できない
	ShowDetails    bool
	ShowAll        bool
	OrderBySizeAsc bool
	Reverse        bool
}

func newLsFlags(args []string) (*LsFlags, error) {
	// オプションを受け取るためのフラグを定義する
	showDetails := commandLine.Bool("l", false, "show details")
	showAll := commandLine.Bool("a", false, "show all")
	orderBySizeDesc := commandLine.Bool("S", false, "sort by size descending")
	reverse := commandLine.Bool("r", false, "reverse order")
	if err := commandLine.Parse(args); err != nil {
		return nil, err
	}

	return &LsFlags{
		ShowDetails:    *showDetails,
		ShowAll:        *showAll,
		OrderBySizeAsc: *orderBySizeDesc,
		Reverse:        *reverse,
	}, nil
}

type lsOutput struct {
	mode    string
	nlink   uint16
	owner   string
	group   string
	size    int64
	modTime string
	name    string
}

func run(fs *LsFlags, w io.Writer) error {
	// 現在の作業ディレクトリを取得する
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// ファイル一覧を表示する
	files, err := os.ReadDir(wd)
	if err != nil {
		return err
	}

	ls := make([]*lsOutput, 0, len(files))
	for _, file := range files {
		if !fs.ShowAll && file.Name()[0] == '.' {
			// -a オプションが指定されていない場合、隠しファイルを表示しない
			continue
		}
		if fs.ShowDetails {
			fi, err := file.Info()
			if err != nil {
				return err
			}
			st, err := getStat(fi)
			if err != nil {
				return err
			}

			// ファイルのパーミッション、オーナー、グループ、モード、サイズ、更新日時を表示する
			ls = append(ls, &lsOutput{
				mode:    fi.Mode().String(),
				nlink:   st.Nlink,
				owner:   st.Owner,
				group:   st.Group,
				size:    fi.Size(),
				modTime: fi.ModTime().Format("1 _2 15:04"),
				name:    fi.Name(),
			})
		} else {
			// ファイル名だけ表示する
			ls = append(ls, &lsOutput{
				name: file.Name(),
			})
		}
	}

	// 昇順にソートする
	if fs.OrderBySizeAsc {
		sort.SliceStable(ls, func(i, j int) bool {
			return ls[i].size > ls[j].size
		})
	}

	// 表示順を反対にする
	if fs.Reverse {
		for i, j := 0, len(ls)-1; i < j; i, j = i+1, j-1 {
			ls[i], ls[j] = ls[j], ls[i]
		}
	}

	// 出力する
	for _, l := range ls {
		if fs.ShowDetails {
			fmt.Fprintf(w, "%s %d %s %s %s %s %s\n", l.mode, l.nlink, l.owner, l.group, strconv.FormatInt(l.size, 10), l.modTime, l.name)
		} else {
			fmt.Fprintln(w, l.name)
		}
	}

	return nil
}

func Ls() error {
	fs, err := newLsFlags(os.Args[1:])
	if err != nil {
		return err
	}

	err = run(fs, os.Stdout)
	if err != nil {
		return err
	}

	return nil
}
