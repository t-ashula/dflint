package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func init() {
	rule := &Rule{
		Name:        "DL3012",
		Severity:    SeverityError,
		Description: "Provide an email adress or URL as maintainer.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true

		mailPattern := regexp.MustCompile(`.+@.+`)
		urlPattern := regexp.MustCompile(`https?://.+`)
		for _, child := range root.Children {
			if child.Value != "maintainer" {
				continue
			}

			if child.Next == nil {
				valid = false
				AppendResult(rule, child)
				continue
			}

			name := child.Next.Value
			if !mailPattern.MatchString(name) && !urlPattern.MatchString(name) {
				valid = false
				AppendResult(rule, child)
				continue
			}
		}

		return valid
	}
	RegisterRule(rule)
}
