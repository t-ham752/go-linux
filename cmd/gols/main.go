package main

import (
	"log"
	"os"

	"github.com/t-ham752/go-linux/pkg/ls"
)

func main() {
	err := ls.Ls(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
