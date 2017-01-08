package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func init() {
	rule := &Rule{
		Name:        "DL3013",
		Severity:    SeverityError,
		Description: "Pin versions in pip.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true

		// TODO: parse shell one liner
		pipPattern := regexp.MustCompile(`\s*pip\s+install\s+.+==.+\s*`)
		for _, child := range root.Children {
			if child.Value != "run" {
				continue
			}

			if child.Next == nil {
				continue
			}

			cmd := child.Next.Value
			if !pipPattern.MatchString(cmd) {
				valid = false
				AppendResult(rule, child)
			}
		}

		return valid
	}
	RegisterRule(rule)
}
