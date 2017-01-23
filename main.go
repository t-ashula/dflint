package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/t-ashula/dflint/rule"
	"github.com/urfave/cli"
)

type lintOption struct {
	IgnoreRules    rule.Rules
	ShellCheckPath string
}

func main() {
	app := setupApp()
	app.Run(os.Args)
}

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Name = "dflint"
	app.Version = "0.0.3"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "ignore, i",
			Usage: "ignore rules",
		},
		cli.StringFlag{
			Name:  "shellcheck",
			Usage: "set ShellCheck Path",
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

	knownRules := rule.GetRules()
	ignoreNames := c.StringSlice("ignore")
	ignoreRules, errorIds := checkIgnoreRule(knownRules, ignoreNames)
	if errorIds != nil {
		for _, id := range errorIds {
			fmt.Fprintf(os.Stderr, "Ignored Unknown Rule ID: %s\n", id)
		}
	}

	scpath := c.String("shellcheck")
	if scpath == "" {
		autopath, err := shellCheckPath()
		if err == nil {
			scpath = autopath
		}
	}

	opts := &lintOption{
		IgnoreRules:    ignoreRules,
		ShellCheckPath: scpath,
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

func checkIgnoreRule(rules rule.Rules, names []string) (rule.Rules, []string) {
	uniqNames := map[string]bool{}
	for _, s := range names {
		uniqNames[s] = true
	}
	var ignores rule.Rules
	var unknowns []string
	for s := range uniqNames {
		r := rules.Find(rule.IsNameMatch(s))
		if r != nil {
			ignores = append(ignores, r)
		} else {
			unknowns = append(unknowns, s)
		}
	}

	return ignores, unknowns
}

func lint(path string, opts *lintOption) (results rule.Results, err error) {
	_, err = os.Stat(path)
	if err != nil {
		return nil, err
	}

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer fp.Close()
	ast, err := buildAST(fp)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return nil, err
	}

	rules := rule.GetRules()
	if opts.ShellCheckPath != "" {
		rule.RegisterRule(shellCheckRule(opts.ShellCheckPath))
	}
	results, err = check(ast, rules)
	if err != nil {
		return nil, err
	}

	results, err = filterResults(results, opts.IgnoreRules)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func buildAST(fp io.Reader) (*parser.Node, error) {
	d := parser.Directive{
		EscapeSeen:           false,
		LookingForDirectives: true,
	}
	parser.SetEscapeToken(parser.DefaultEscapeToken, &d)
	return parser.Parse(fp, &d)
}

func check(ast *parser.Node, rules rule.Rules) (rule.Results, error) {
	for _, r := range rules {
		r.Validate(ast)
	}
	return rule.GetResults(), nil
}

func filterResults(results rule.Results, ignores rule.Rules) (rule.Results, error) {
	var tmp rule.Results
	for _, result := range results {
		r := ignores.Find(rule.IsNameMatch(result.Rule.Name))
		if r == nil {
			tmp = append(tmp, result)
		}
	}
	return tmp, nil
}

func report(results rule.Results) {
	for i, res := range results {
		fmt.Printf("%d: %d,%d,%d,%d: %s: %s\n", i, res.StartLine, res.StartColumn, res.EndLine, res.EndColumn, res.Rule.Name, res.Rule.Description)
	}
}
