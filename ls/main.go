package main

import (
	"flag"
	"fmt"
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

type LS struct {
	Mode    string
	Nlink   uint16
	Owner   string
	Group   string
	Size    string
	ModTime string
	Name    string

	*LsFlags
}

func Ls(ls *LsFlags) ([]*LS, error) {
	// 現在の作業ディレクトリを取得する
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// ファイル一覧を表示する
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		return nil, err
	}

	fs := make([]*LS, 0, len(files))
	for _, file := range files {
		if !ls.showAll && file.Name()[0] == '.' {
			// -a オプションが指定されていない場合、隠しファイルを表示しない
			continue
		}
		if ls.showDetails {
			st, err := getStat(file)
			if err != nil {
				return nil, err
			}
			// ファイルのパーミッション、オーナー、グループ、モード、サイズ、更新日時を表示する
			mode := file.Mode()
			size := file.Size()
			modTime := file.ModTime().Format("1 _2 15:04")
			fs = append(fs, &LS{
				Mode:    mode.String(),
				Nlink:   st.nlink,
				Owner:   st.owner,
				Group:   st.group,
				Size:    strconv.FormatInt(size, 10),
				ModTime: modTime,
				Name:    file.Name(),
			})
		} else {
			// ファイル名だけ表示する
			fs = append(fs, &LS{
				Name: file.Name(),
			})
		}
	}
	return fs, nil
}

func main() {
	nf := NewLsFlags()
	fs, err := Ls(nf)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range fs {
		if nf.showDetails {
			fmt.Printf("%s %d %s %s %s %s %s\n", f.Mode, f.Nlink, f.Owner, f.Group, f.Size, f.ModTime, f.Name)
		} else {
			fmt.Println(f.Name)
		}
	}
}
