package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "dflint"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "ignore, i",
			Usage: "ignore rules",
		},
	}
	app.Usage = "lint Dockerfile"
	app.UsageText = "dflint [global options] Dockerfile [Dockerfile,...]"
	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 {
			return errors.New("at least one Dockerfile required")
		}
		files := c.Args()
		for i, file := range files {
			fmt.Printf("target file[%03d]: %s\n", i, file)
		}

		ignoreRules := c.StringSlice("ignore")
		for i, rule := range ignoreRules {
			fmt.Printf("ignore rule[%03d]: %s\n", i, rule)
		}
		return nil
	}

	app.Run(os.Args)
}
