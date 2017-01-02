package rule

import (
	"regexp"
	"strings"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func init() {
	rule := &Rule{
		Name:        "DL3008",
		Severity:    SeverityError,
		Description: "Using latest is prone to errors if the image will ever update. Pin the version explicitly to a release tag.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true
		for _, child := range root.Children {
			if child.Value == "run" {
				arg := child.Next
				if arg == nil || arg.Value == "" {
					continue
				}

				// TODO: parse shell one liner
				cmd := regexp.MustCompile(`[\s;(|&]*apt-get(?:[^;&|]*)?install\s+(([-a-zA-Z0-9.='"]+\s*)+)[\s)|&]*`)
				apt := cmd.FindAllStringSubmatch(arg.Value, -1)
				if len(apt) == 1 && len(apt[0]) > 1 {
					pkgs := regexp.MustCompile(`\s+`).Split(apt[0][1], -1)
					for _, pkg := range pkgs {
						if strings.Index(pkg, "=") == -1 {
							AppendResult(rule, child)
							valid = false
						}
					}
				}
			}
		}
		return valid
	}
	RegisterRule(rule)
}
