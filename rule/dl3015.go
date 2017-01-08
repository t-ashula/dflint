package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func init() {
	rule := &Rule{
		Name:        "DL3015",
		Severity:    SeverityError,
		Description: "Avoid additional packages by specifying --no-install-recommends.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true

		// TODO: parse shell one liner
		pattern := regexp.MustCompile(`\s*apt-get.+--no-install-recommends\s+.+\s*`)
		for _, child := range root.Children {
			if child.Value != "run" {
				continue
			}

			if child.Next == nil {
				continue
			}

			cmd := child.Next.Value
			if !pattern.MatchString(cmd) {
				valid = false
				AppendResult(rule, child)
			}
		}

		return valid
	}
	RegisterRule(rule)
}
