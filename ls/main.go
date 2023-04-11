package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/user"
	"sort"
	"strconv"
	"syscall"
)

type statT struct {
	nlink uint16
	owner string
	group string
}

func getStat(fs fs.FileInfo) (*statT, error) {
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

	return &statT{
		nlink: stat.Nlink,
		owner: owner,
		group: group,
	}, nil
}

var (
	commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
)

type LsFlags struct {
	ShowDetails    bool
	ShowAll        bool
	OrderBySizeAsc bool
	Reverse        bool
}

func NewLsFlags(args []string) *LsFlags {
	// オプションを受け取るためのフラグを定義する
	showDetails := commandLine.Bool("l", false, "show details")
	showAll := commandLine.Bool("a", false, "show all")
	orderBySizeDesc := commandLine.Bool("S", false, "sort by size descending")
	reverse := commandLine.Bool("r", false, "reverse order")
	if err := commandLine.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	return &LsFlags{
		ShowDetails:    *showDetails,
		ShowAll:        *showAll,
		OrderBySizeAsc: *orderBySizeDesc,
		Reverse:        *reverse,
	}
}

type LS struct {
	Mode    string
	Nlink   uint16
	Owner   string
	Group   string
	Size    int64
	ModTime string
	Name    string
}

func Ls(ls *LsFlags) ([]*LS, error) {
	// 現在の作業ディレクトリを取得する
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// ファイル一覧を表示する
	files, err := os.ReadDir(wd)
	if err != nil {
		return nil, err
	}

	fs := make([]*LS, 0, len(files))
	for _, file := range files {
		if !ls.ShowAll && file.Name()[0] == '.' {
			// -a オプションが指定されていない場合、隠しファイルを表示しない
			continue
		}
		if ls.ShowDetails {
			fi, err := file.Info()
			if err != nil {
				return nil, err
			}
			st, err := getStat(fi)
			if err != nil {
				return nil, err
			}

			// ファイルのパーミッション、オーナー、グループ、モード、サイズ、更新日時を表示する
			fs = append(fs, &LS{
				Mode:    fi.Mode().String(),
				Nlink:   st.nlink,
				Owner:   st.owner,
				Group:   st.group,
				Size:    fi.Size(),
				ModTime: fi.ModTime().Format("1 _2 15:04"),
				Name:    fi.Name(),
			})
		} else {
			// ファイル名だけ表示する
			fs = append(fs, &LS{
				Name: file.Name(),
			})
		}
	}

	// 昇順にソートする
	if ls.OrderBySizeAsc {
		sort.SliceStable(fs, func(i, j int) bool {
			return fs[i].Size > fs[j].Size
		})
	}

	// 表示順を反対にする
	if ls.Reverse {
		for i, j := 0, len(fs)-1; i < j; i, j = i+1, j-1 {
			fs[i], fs[j] = fs[j], fs[i]
		}
	}

	return fs, nil
}

func main() {
	nf := NewLsFlags(os.Args[1:])
	fs, err := Ls(nf)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range fs {
		if nf.ShowDetails {
			fmt.Printf("%s %d %s %s %s %s %s\n", f.Mode, f.Nlink, f.Owner, f.Group, strconv.FormatInt(f.Size, 10), f.ModTime, f.Name)
		} else {
			fmt.Println(f.Name)
		}
	}
}
