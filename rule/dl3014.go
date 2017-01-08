package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func init() {
	rule := &Rule{
		Name:        "DL3014",
		Severity:    SeverityError,
		Description: "Use the -y switch.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true

		// TODO: parse shell one liner
		yesPattern := regexp.MustCompile(`\s*apt-get\s+-y\s+.+\s*`)
		for _, child := range root.Children {
			if child.Value != "run" {
				continue
			}

			if child.Next == nil {
				continue
			}

			cmd := child.Next.Value
			if !yesPattern.MatchString(cmd) {
				valid = false
				AppendResult(rule, child)
			}
		}

		return valid
	}
	RegisterRule(rule)
}
