package rule

import (
	"regexp"
	"strings"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/mvdan/sh/syntax"
)

func init() {
	rule := &Rule{
		Name:        "DL3013",
		Severity:    SeverityError,
		Description: "Pin versions in pip.",
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

			if notPipPinned(cmd) {
				valid = false
				AppendResult(rule, child)
			}
		}

		return valid
	}
	RegisterRule(rule)
}

func notPipPinned(script string) bool {
	ast, err := syntax.Parse(strings.NewReader(script), "", 0)
	if err != nil {
		return false
	}
	verPattern := regexp.MustCompile(`.+==.+`)
	pinned := true
	syntax.Walk(ast, func(node syntax.Node) bool {
		expr, isCall := node.(*syntax.CallExpr)
		if !isCall {
			return true
		}

		fname, err := funcName(expr)

		if err != nil || fname != "pip" {
			return false
		}

		args, err := funcArgs(expr)
		if err != nil {
			return false
		}

		// TODO: refactor
		ins, ver := false, true
		for _, arg := range args {
			if strings.HasPrefix(arg, "-") {
				continue
			}
			if arg == "install" {
				ins = true
				continue
			}
			if !verPattern.MatchString(arg) {
				ver = false
			}
			if ins && !ver {
				pinned = false
				break
			}
		}
		// no need to digging
		return false
	})

	return !pinned
}
