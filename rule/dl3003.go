package rule

import (
	"strings"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/mvdan/sh/syntax"
)

func init() {
	rule := &Rule{
		Name:        "DL3003",
		Severity:    SeverityError,
		Description: "Use WORKDIR to switch to a directory.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true
		for _, child := range root.Children {
			if child.Value != "run" {
				continue
			}

			arg := child.Next
			if arg == nil || arg.Value == "" {
				continue
			}

			if hasCd(arg.Value) {
				AppendResult(rule, child)
				valid = false
			}

		}
		return valid
	}
	RegisterRule(rule)
}

func hasCd(script string) bool {
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

		if err != nil {
			return false
		}

		if fname == "cd" {
			has = true
			return false
		}

		return false
	})
	return has
}
