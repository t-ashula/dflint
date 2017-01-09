package rule

import "github.com/docker/docker/builder/dockerfile/parser"

func init() {
	rule := &Rule{
		Name:        "DL4003",
		Severity:    SeverityInfo,
		Description: "Multiple CMD instructions found. Only the first CMD instruction will take effect.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true
		cmdCount := 0
		for _, child := range root.Children {
			if child.Value != "cmd" {
				continue
			}

			cmdCount++

			if cmdCount > 1 {
				valid = false
				AppendResult(rule, child)
			}
		}

		return valid
	}
	RegisterRule(rule)
}
