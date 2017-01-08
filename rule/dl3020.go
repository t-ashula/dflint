package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func init() {
	rule := &Rule{
		Name:        "DL3020",
		Severity:    SeverityError,
		Description: "Use COPY instead of ADD for files and folders.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true

		// TODO: whici extention will auto-extraction?
		tarPattern := regexp.MustCompile(`(\.tar\.xz|\.tar\.gz|\.txz|\.tgz)$`)
		for _, child := range root.Children {
			if child.Value != "add" {
				continue
			}

			if child.Next == nil {
				continue
			}

			src := child.Next
			for src.Next != nil {
				if !tarPattern.MatchString(src.Value) {
					valid = false
					AppendResult(rule, child)
					break
				}
				src = src.Next
			}
		}

		return valid
	}
	RegisterRule(rule)
}
