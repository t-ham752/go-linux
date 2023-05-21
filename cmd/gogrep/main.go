package main

import (
	"log"

	"github.com/t-ham752/go-linux/pkg/grep"
)

func main() {
	err := grep.Grep()
	if err != nil {
		log.Fatal(err)
	}
}
