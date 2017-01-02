package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
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
			if child.Value == "run" {
				arg := child.Next
				if arg == nil || arg.Value == "" {
					continue
				}

				// TODO: parse shell one liner
				cd := regexp.MustCompile(`[\s;(|&]*cd[\s)|&]*`)
				if cd.MatchString(arg.Value) {
					AppendResult(rule, child)
					valid = false
				}
			}
		}
		return valid
	}
	RegisterRule(rule)
}
