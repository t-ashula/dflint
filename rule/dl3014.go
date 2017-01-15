package rule

import (
	"regexp"
	"strings"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/mvdan/sh/syntax"
)

func init() {
	rule := &Rule{
		Name:        "DL3014",
		Severity:    SeverityError,
		Description: "Use the -y switch.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true

		for _, child := range root.Children {
			if child.Value != "run" {
				continue
			}

			if child.Next == nil {
				continue
			}

			cmd := child.Next.Value
			if !hasYesOption(cmd) {
				valid = false
				AppendResult(rule, child)
			}
		}

		return valid
	}
	RegisterRule(rule)
}

func hasYesOption(script string) bool {
	ast, err := syntax.Parse(strings.NewReader(script), "", 0)
	if err != nil {
		return false
	}
	has := false
	yesPattern := regexp.MustCompile(`^-[^-]*y.*$`)
	syntax.Walk(ast, func(node syntax.Node) bool {
		expr, isCall := node.(*syntax.CallExpr)
		if !isCall {
			return true
		}

		fname, err := funcName(expr)

		if err != nil || fname != "apt-get" {
			return false
		}

		args, err := funcArgs(expr)
		if err != nil {
			return false
		}

		ins, yes := false, false
		for _, arg := range args {
			if arg == "install" {
				ins = true
			}
			if arg == "-y" || arg == "--yes" || arg == "--assume-yes" ||
				yesPattern.MatchString(arg) {
				yes = true
			}
			if ins && yes {
				has = true
				break
			}
		}
		// no need to digging
		return false
	})
	return has
}
