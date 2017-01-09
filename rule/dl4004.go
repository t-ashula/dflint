package rule

import "github.com/docker/docker/builder/dockerfile/parser"

func init() {
	rule := &Rule{
		Name:        "DL4004",
		Severity:    SeverityInfo,
		Description: "Multiple ENTRYPOINT instructions found. Only the first ENTRYPOINT instruction will take effect.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true
		cmdCount := 0
		for _, child := range root.Children {
			if child.Value != "entrypoint" {
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
