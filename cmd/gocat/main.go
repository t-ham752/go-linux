package main

import (
	"log"

	"github.com/t-ham752/go-linux/pkg/cat"
)

func main() {
	err := cat.Cat()
	if err != nil {
		log.Fatal(err)
	}
}
