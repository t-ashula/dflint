package rule

import "github.com/docker/docker/builder/dockerfile/parser"

func init() {
	dl3000 := &Rule{
		Name:        "DL3000",
		Severity:    SeverityError,
		Description: "Use absolute WORKDIR",
	}
	dl3000.Validate = func(root *parser.Node) bool {
		for _, child := range root.Children {
			if child.Value == "workdir" {
			}
		}
		return true
	}
	RegisterRule(dl3000)

	dl4000 := &Rule{
		Name:        "DL4000",
		Severity:    SeverityInfo,
		Description: "Specify a maintainer of the Dockerfile",
	}
	dl4000.Validate = func(root *parser.Node) bool {
		for _, child := range root.Children {
			if child.Value == "maintainer" {
				return true
			}
		}

		AppendResult(dl4000, nil) // XXX:
		return false
	}
	RegisterRule(dl4000)
}
