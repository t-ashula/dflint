package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
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
			if child.Value == "run" {
				arg := child.Next
				if arg == nil || arg.Value == "" {
					continue
				}

				// TODO: parse shell one liner
				cmd := regexp.MustCompile(`[\s;(|&]*sudo[\s)|&]*`)
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
