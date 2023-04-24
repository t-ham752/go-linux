package cat

import (
	"fmt"
	"io"
	"os"
)

func run(r io.Reader, w io.Writer) error {
	if len(os.Args) < 2 {
		return fmt.Errorf("too few arguments")
	}
	for _, arg := range os.Args[1:] {
		file, err := os.Open(arg)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, file)
	}
	return nil
}
func Cat() error {
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		return err
	}
	return nil
}
