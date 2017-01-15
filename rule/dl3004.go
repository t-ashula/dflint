package rule

import (
	"strings"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/mvdan/sh/syntax"
)

func init() {
	rule := &Rule{
		Name:        "DL3004",
		Severity:    SeverityError,
		Description: "Do not use sudo as it leads to unpredictable behavior. Use a tool like gosu to enforce root.",
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

			if hasSudo(arg.Value) {
				valid = false
				AppendResult(rule, child)
			}

		}
		return valid
	}
	RegisterRule(rule)
}

func hasSudo(script string) bool {
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

		if fname == "sudo" {
			has = true
			return false
		}

		return false
	})
	return has
}
