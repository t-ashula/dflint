package rule

import (
	"strings"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/mvdan/sh/syntax"
)

func init() {
	rule := &Rule{
		Name:        "DL3005",
		Severity:    SeverityError,
		Description: "Do not use apt-get upgrade or dist-upgrade.",
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

			if aptUpgrade(arg.Value) {
				AppendResult(rule, child)
				valid = false
			}
		}

		return valid
	}
	RegisterRule(rule)
}

func aptUpgrade(script string) bool {
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

		if err != nil || fname != "apt-get" {
			return false
		}

		args, err := funcArgs(expr)
		if err != nil {
			return false
		}

		for _, arg := range args {
			if arg == "upgrade" || arg == "dist-upgrade" {
				has = true
				break
			}
		}
		// no need to digging
		return false
	})
	return has
}
