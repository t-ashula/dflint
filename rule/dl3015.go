package rule

import (
	"strings"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/mvdan/sh/syntax"
)

func init() {
	rule := &Rule{
		Name:        "DL3015",
		Severity:    SeverityError,
		Description: "Avoid additional packages by specifying --no-install-recommends.",
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
			if !hasNoIstallRecommends(cmd) {
				valid = false
				AppendResult(rule, child)
			}
		}

		return valid
	}

	RegisterRule(rule)
}

func hasNoIstallRecommends(script string) bool {
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
			has = true
			return false
		}

		args, err := funcArgs(expr)
		if err != nil {
			has = true
			return false
		}

		ins, rec := false, false
		for _, arg := range args {
			if arg == "install" {
				ins = true
			}
			if arg == "--no-install-recommends" {
				rec = true
			}
			if ins && rec {
				has = true
				break
			}
		}
		// no need to digging
		return false
	})
	return has
}
