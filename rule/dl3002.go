package rule

import "github.com/docker/docker/builder/dockerfile/parser"

func init() {
	rule := &Rule{
		Name:        "DL3002",
		Severity:    SeverityError,
		Description: "Do not switch to root USER.",
	}
	rule.Validate = func(root *parser.Node) bool {
		for _, child := range root.Children {
			if child.Value == "user" {
				arg := child.Next
				if arg == nil {
					continue
				}
				if arg.Value == "root" {
					AppendResult(rule, child)
					return false
				}
			}
		}
		return true
	}
	RegisterRule(rule)
}
