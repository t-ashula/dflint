package rule

import "github.com/docker/docker/builder/dockerfile/parser"

func init() {
	rule := &Rule{
		Name:        "DL3001",
		Severity:    SeverityError,
		Description: "For some bash commands it makes no sense running them in a Docker container like ssh, vim, shutdown, service, ps, free, top, kill, mount, ifconfig.",
	}
	rule.Validate = func(root *parser.Node) bool {
		ngcmds := []string{"shutdown", "service", "ps", "free", "top", "kill", "mount", "ifconfig", "nano", "vim", "emacs"}
		for _, child := range root.Children {
			if child.Value == "run" {
				arg := child.Next
				if arg == nil {
					continue
				}
				for _, ng := range ngcmds {
					if arg.Value == ng {
						AppendResult(rule, child)
						return false
					}
				}
			}
		}
		return true
	}
	RegisterRule(rule)
}
