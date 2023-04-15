package main

import (
	"log"

	"github.com/t-ham752/go-linux/ls"
)

func main() {
	err := ls.Ls()
	if err != nil {
		log.Fatal(err)
	}
}
