package rule

import (
	"strings"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/mvdan/sh/syntax"
)

func init() {
	rule := &Rule{
		Name:        "DL3009",
		Severity:    SeverityError,
		Description: "Delete the apt-get lists after installing something.",
	}
	rule.Validate = func(root *parser.Node) bool {
		valid := true
		scripts := []string{""}
		for _, child := range root.Children {
			if child.Value != "run" {
				continue
			}
			src := child.Next
			if src == nil || src.Value == "" {
				continue
			}
			scripts = append(scripts, src.Value)
		}

		installed, _ := scanInstallCommand(scripts)
		if installed {
			valid = false
			AppendResult(rule, nil)
		}

		return valid
	}
	RegisterRule(rule)
}

type pkgman struct {
	Name      string
	IsInstall func([]string) (bool, error)
	IsCleanup func([]string) (bool, error)
}

var installers = []pkgman{{
	Name: "apt-get",
	IsInstall: func(scripts []string) (bool, error) {
		isInstall := false
		for _, script := range scripts {
			ast, err := syntax.Parse(strings.NewReader(script), "", 0)
			if err != nil {
				return false, err
			}
			syntax.Walk(ast, func(node syntax.Node) bool {
				expr, isCall := node.(*syntax.CallExpr)
				if !isCall {
					return true
				}

				fname, err := funcName(expr)
				if err == nil {
					if fname == "apt-get" {
						args, err := funcArgs(expr)
						if err == nil {
							if contains(args, "update") {
								isInstall = true
							}
							if contains(args, "install") {
								isInstall = true
							}
						}
					}
				}
				// no need to digging
				return false
			})
		}
		return isInstall, nil
	},
	IsCleanup: func(scripts []string) (bool, error) {
		hasClean, hasRemove := false, false
		for _, script := range scripts {
			ast, err := syntax.Parse(strings.NewReader(script), "", 0)
			if err != nil {
				return false, err
			}
			syntax.Walk(ast, func(node syntax.Node) bool {
				expr, isCall := node.(*syntax.CallExpr)
				if !isCall {
					return true
				}

				fname, err := funcName(expr)
				if err == nil {
					if fname == "apt-get" {
						args, err := funcArgs(expr)
						if err == nil {
							if contains(args, "clean") {
								hasClean = true
							}
						}
					} else if fname == "rm" {
						args, err := funcArgs(expr)
						if err == nil {
							if contains(args, "/var/lib/apt/lists/*") {
								hasRemove = true
							}
						}
					}
				}
				// no need to digging
				return false
			})
		}
		return hasClean && hasRemove, nil
	},
}}

func scanInstallCommand(scripts []string) (bool, error) {
	for _, pm := range installers {
		ins, err := pm.IsInstall(scripts)
		if err != nil {
			return false, err
		}
		del, err := pm.IsCleanup(scripts)
		if err != nil {
			return false, err
		}
		if ins && !del {
			return true, nil
		}
	}

	return false, nil
}
