package rule

import (
	"regexp"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func init() {
	rule := &Rule{
		Name:        "DL3010",
		Severity:    SeverityError,
		Description: "Delete the apt-get lists after installing something.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true
		// TODO: which extention will auto-extraction?
		isTar := regexp.MustCompile(`(?:\.tar\.gz|tar\.bz|\.tar.xz|\.tgz|\.tbz)$`)
		for _, child := range root.Children {
			if child.Value != "copy" {
				continue
			}
			src := child.Next
			if src == nil || src.Value == "" {
				continue
			}

			srcs := []string{src.Value}
			src = src.Next
			for src != nil {
				srcs = append(srcs, src.Value)
				src = src.Next
			}

			for _, src := range srcs {
				if isTar.MatchString(src) {
					valid = false
					AppendResult(rule, child)
				}
			}
		}
		return valid
	}
	RegisterRule(rule)
}
