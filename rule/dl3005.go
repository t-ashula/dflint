package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
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
			if child.Value == "run" {
				arg := child.Next
				if arg == nil || arg.Value == "" {
					continue
				}

				// TODO: parse shell one liner
				cmd := regexp.MustCompile(`[\s;(|&]*apt-get([^;&|]*)?(\s+dist-)?upgrade[\s)|&]*`)
				if cmd.MatchString(arg.Value) {
					AppendResult(rule, child)
					valid = false
				}
			}
		}
		return valid
	}
	RegisterRule(rule)
}
