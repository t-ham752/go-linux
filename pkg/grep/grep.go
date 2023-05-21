package grep

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/urfave/cli/v2"
)

type GrepFlags struct {
	IgnoreCase bool
}

func run(args []string, flags *GrepFlags) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments")
	}
	target := args[0]
	input := args[1]

	if flags.IgnoreCase {
		target = `(?i)` + target
	}
	re := regexp.MustCompile(target)

	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			fmt.Println(line)
		}
	}

	return nil
}

func Grep() error {
	var doesIgnore bool

	app := &cli.App{
		Name:  "go-grep",
		Usage: "searches input files, selecting lines that match patterns.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "ignore-case",
				Aliases:     []string{"i"},
				Usage:       "ignore case distinctions",
				Destination: &doesIgnore,
			},
		},
		Action: func(c *cli.Context) error {
			err := run(c.Args().Slice(), &GrepFlags{
				IgnoreCase: doesIgnore,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	return nil
}
