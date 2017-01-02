package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func init() {
	rule := &Rule{
		Name:        "DL3006",
		Severity:    SeverityError,
		Description: "Always tag the version of an image explicitly",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true
		for _, child := range root.Children {
			if child.Value == "from" {
				arg := child.Next
				if arg == nil || arg.Value == "" {
					continue
				}

				image := regexp.MustCompile(`.+[:@].+`)
				if !image.MatchString(arg.Value) {
					AppendResult(rule, child)
					valid = false
				}
			}
		}
		return valid
	}
	RegisterRule(rule)
}
