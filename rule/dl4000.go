package rule

import "github.com/docker/docker/builder/dockerfile/parser"

func init() {
	dl4000 := &Rule{
		Name:        "DL4000",
		Severity:    SeverityInfo,
		Description: "Specify a maintainer of the Dockerfile",
	}
	dl4000.Validate = func(root *parser.Node) bool {
		for _, child := range root.Children {
			if child.Value == "maintainer" && child.Next != nil && child.Next.Value != "" {
				return true
			}
		}

		AppendResult(dl4000, nil) // XXX:
		return false
	}
	// DL4000 is deprecated
	// RegisterRule(dl4000)
}
