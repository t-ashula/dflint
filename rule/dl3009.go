package rule

import "github.com/docker/docker/builder/dockerfile/parser"

func init() {
	rule := &Rule{
		Name:        "DL3009",
		Severity:    SeverityError,
		Description: "Delete the apt-get lists after installing something.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true
		return valid
	}
	RegisterRule(rule)
}
