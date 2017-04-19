package rule

import "github.com/docker/docker/builder/dockerfile/parser"

func init() {
	rule := &Rule{
		Name:        "DL3012",
		Severity:    SeverityError,
		Description: "Provide an email adress or URL as maintainer.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true
		return valid
	}
	// DL3012 is deprecated
	// RegisterRule(dl4000)
}
