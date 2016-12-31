package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/t-ashula/dflint/rule"
	"github.com/urfave/cli"
)

type lintOption struct {
	IgnoreRules []*rule.Rule
}

func main() {
	app := setupApp()
	app.Run(os.Args)
}

func setupApp() *cli.App {
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
	app.Action = realMain
	return app
}

func realMain(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("at least one Dockerfile required")
	}
	ignoreRules, errorIds := checkIgnoreRule(c.StringSlice("ignore"))
	if errorIds != nil {
		for _, id := range errorIds {
			fmt.Fprintf(os.Stderr, "Ignored Unknown Rule ID: %s\n", id)
		}
	}

	opts := &lintOption{
		IgnoreRules: ignoreRules,
	}

	files := c.Args()
	for _, file := range files {
		result, err := lint(file, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s, %s\n", file, err)
			continue
		}
		report(result)
	}

	return nil
}

func checkIgnoreRule(names []string) ([]*rule.Rule, []string) {
	var rules []*rule.Rule
	return rules, nil
}

func lint(path string, opts *lintOption) (results []*rule.Result, err error) {
	_, err = os.Stat(path)
	if err != nil {
		return nil, err
	}

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer fp.Close()

	d := parser.Directive{
		EscapeSeen:           false,
		LookingForDirectives: true,
	}
	parser.SetEscapeToken(parser.DefaultEscapeToken, &d)
	ast, err := parser.Parse(fp, &d)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return nil, err
	}

	rules := rule.GetRules()
	results, err = check(ast, rules)
	if err != nil {
		return nil, err
	}

	err = filter(&results, opts.IgnoreRules)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func check(ast *parser.Node, rules []*rule.Rule) ([]*rule.Result, error) {
	for _, r := range rules {
		r.Validate(ast)
	}
	return rule.GetResults(), nil
}

func filter(results *[]*rule.Result, ignores []*rule.Rule) error {
	return nil
}

func report(results []*rule.Result) {
	for i, res := range results {
		fmt.Printf("%d: %d,%d,%d,%d: %s: %s\n", i, res.StartLine, res.StartColumn, res.EndLine, res.EndColumn, res.Rule.Name, res.Rule.Description)
	}
}
