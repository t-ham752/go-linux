package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/user"
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
	nlink := stat.Nlink

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
		nlink: nlink,
		owner: owner,
		group: group,
	}, nil
}

type LsFlags struct {
	showDetails bool
	showAll     bool
}

func NewLsFlags() *LsFlags {
	// -l オプションと -a オプションを受け取るためのフラグを定義する
	showDetails := flag.Bool("l", false, "show details")
	showAll := flag.Bool("a", false, "show all")
	flag.Parse()

	return &LsFlags{
		showDetails: *showDetails,
		showAll:     *showAll,
	}
}

func Exec(w io.Writer, lf *LsFlags) {
	// 現在の作業ディレクトリを取得する
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// ファイル一覧を表示する
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !lf.showAll && file.Name()[0] == '.' {
			// -a オプションが指定されていない場合、隠しファイルを表示しない
			continue
		}
		if lf.showDetails {
			st, err := getStat(file)
			if err != nil {
				log.Fatal(err)
			}
			// ファイルのパーミッション、オーナー、グループ、モード、サイズ、更新日時を表示する
			mode := file.Mode()
			size := file.Size()
			modTime := file.ModTime().Format("1 _2 15:04")
			fmt.Fprintf(w, "%s %d %s %s %s %s %s\n", mode.String(), st.nlink, st.owner, st.group, strconv.FormatInt(size, 10), modTime, file.Name())
		} else {
			// ファイル名だけ表示する
			fmt.Fprintln(w, file.Name())
		}
	}
}

func main() {
	lf := NewLsFlags()
	Exec(os.Stdout, lf)
}
