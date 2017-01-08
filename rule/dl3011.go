package rule

import (
	"strconv"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func init() {
	rule := &Rule{
		Name:        "DL3011",
		Severity:    SeverityError,
		Description: "Valid UNIX ports range from 0 to 65535.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true

		for _, child := range root.Children {
			if child.Value != "expose" {
				continue
			}

			port := child.Next
			if port == nil {
				valid = false
				AppendResult(rule, child)
				continue
			}

			portNum, err := strconv.Atoi(port.Value)
			if err != nil || portNum < 0 || portNum > 65535 {
				valid = false
				AppendResult(rule, child)
			}
		}

		return valid
	}
	RegisterRule(rule)
}
