package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func init() {
	rule := &Rule{
		Name:        "DL4005",
		Severity:    SeverityInfo,
		Description: "Use SHELL to change the default shell.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true
		// TODO: parse shell one liner
		pattern := regexp.MustCompile(`ln\s+-[sfv]+\s+/bin/bash\s*/bin/sh`)
		for _, child := range root.Children {
			if child.Value != "run" {
				continue
			}

			if child.Next == nil {
				continue
			}

			cmd := child.Next.Value
			if pattern.MatchString(cmd) {
				valid = false
				AppendResult(rule, child)
			}
		}

		return valid
	}
	RegisterRule(rule)
}
