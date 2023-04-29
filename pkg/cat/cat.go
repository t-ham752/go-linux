package cat

import (
	"fmt"
	"io"
	"os"
)

func run(args []string, w io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("too few arguments")
	}
	for _, arg := range args[1:] {
		file, err := os.Open(arg)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, file)
	}
	return nil
}
func Cat() error {
	err := run(os.Args, os.Stdout)
	if err != nil {
		return err
	}
	return nil
}
