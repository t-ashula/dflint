package rule

import (
	"regexp"
	"strings"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/mvdan/sh/syntax"
)

func init() {
	rule := &Rule{
		Name:        "DL3008",
		Severity:    SeverityError,
		Description: "Using latest is prone to errors if the image will ever update. Pin the version explicitly to a release tag.",
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

			if notAptgetPinned(arg.Value) {
				AppendResult(rule, child)
				valid = false
			}
		}

		return valid
	}
	RegisterRule(rule)
}

func notAptgetPinned(script string) bool {
	ast, err := syntax.Parse(strings.NewReader(script), "", 0)
	if err != nil {
		return false
	}
	verPattern := regexp.MustCompile(`.+=.+`)
	pinned := true
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
