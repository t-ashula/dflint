package rule

import "github.com/docker/docker/builder/dockerfile/parser"

func init() {
	dl5000 := &Rule{
		Name:        "DL5000",
		Severity:    SeverityInfo,
		Description: "MAINTAINER was an early very limited form of LABEL which should be used instead.",
	}
	dl5000.Validate = func(root *parser.Node) bool {
		for _, child := range root.Children {
			if child.Value == "maintainer" {
				return false
			}
		}

		AppendResult(dl5000, nil) // XXX:
		return true
	}
	RegisterRule(dl5000)
}
