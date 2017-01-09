package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func init() {
	rule := &Rule{
		Name:        "DL4001",
		Severity:    SeverityInfo,
		Description: "Either use Wget or Curl but not both.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true
		curl, wget := false, false

		curlPattern := regexp.MustCompile(`^curl\s*`)
		wgetPattern := regexp.MustCompile(`^wget\s*`)
		for _, child := range root.Children {
			if child.Value != "run" {
				continue
			}

			if child.Next == nil {
				continue
			}

			cmd := child.Next.Value
			if curlPattern.MatchString(cmd) {
				curl = true
			}
			if wgetPattern.MatchString(cmd) {
				wget = true
			}

			if curl && wget {
				valid = false
				AppendResult(rule, child)
			}
		}

		return valid
	}
	RegisterRule(rule)
}
