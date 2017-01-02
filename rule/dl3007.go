package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func init() {
	rule := &Rule{
		Name:        "DL3007",
		Severity:    SeverityError,
		Description: "Using latest is prone to errors if the image will ever update. Pin the version explicitly to a release tag.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true
		for _, child := range root.Children {
			if child.Value == "from" {
				arg := child.Next
				if arg == nil || arg.Value == "" {
					continue
				}

				image := regexp.MustCompile(`^.+:latest$`)
				if image.MatchString(arg.Value) {
					AppendResult(rule, child)
					valid = false
				}
			}
		}
		return valid
	}
	RegisterRule(rule)
}
