package rule

import (
	"regexp"
	"strings"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/mvdan/sh/syntax"
)

func init() {
	rule := &Rule{
		Name:        "DL4005",
		Severity:    SeverityInfo,
		Description: "Use SHELL to change the default shell.",
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

			script := child.Next.Value
			has := hasSymbolicLinkCommand(script)

			if has {
				valid = false
				AppendResult(rule, child)
			}
		}

		return valid
	}
	RegisterRule(rule)
}

func hasSymbolicLinkCommand(script string) bool {
	shortArg := regexp.MustCompile(`^-.*s.*$`)
	ast, err := syntax.Parse(strings.NewReader(script), "", 0)
	if err != nil {
		return false
	}
	has := false
	syntax.Walk(ast, func(node syntax.Node) bool {
		expr, isCall := node.(*syntax.CallExpr)
		if !isCall {
			return true
		}

		fname, err := funcName(expr)

		if err != nil || fname != "ln" {
			return false
		}
		args, err := funcArgs(expr)
		if err != nil {
			return false
		}

		for _, arg := range args {
			if !strings.HasPrefix(arg, "-") {
				break // for args
			}
			if arg == "--symbolic" || arg == "-s" {
				has = true
				break
			}
			if shortArg.MatchString(arg) {
				has = true
				break
			}
		}
		// no need to digging
		return false
	})
	return has
}
