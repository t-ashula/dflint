package rule

import (
	"errors"
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

func contains(values []string, target string) bool {
	for _, e := range values {
		if e == target {
			return true
		}
	}
	return false
}

func funcName(expr *syntax.CallExpr) (string, error) {
	if expr == nil {
		return "", errors.New("expr nil")
	}
	if expr.Args == nil || len(expr.Args) < 1 {
		return "", errors.New("no commands")
	}
	// CallExpr:{ Args  []Word{}},
	// Word    :{ Parts []WordPart{} }
	a0 := expr.Args[0].Parts[0]
	name := ""
	switch node := a0.(type) {
	case *syntax.Lit:
		name = node.Value
	}

	return name, nil
}

func funcArgs(expr *syntax.CallExpr) ([]string, error) {
	if expr == nil {
		return nil, errors.New("expr nil")
	}
	if expr.Args == nil || len(expr.Args) < 1 {
		return nil, errors.New("no command")
	}
	args := make([]string, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = wordAsString(arg)
	}
	return args, nil
}

func partAsString(part syntax.WordPart) string {
	s := ""
	switch v := part.(type) {
	case *syntax.Lit:
		s = v.Value
	case *syntax.SglQuoted:
		s = v.Value
	case *syntax.DblQuoted:
		for _, p := range v.Parts {
			s += partAsString(p)
		}
	default:
		s = ""
	}
	return s
}

func wordAsString(word *syntax.Word) string {
	parts := make([]string, len(word.Parts))
	for i, part := range word.Parts {
		parts[i] = partAsString(part)
	}
	return strings.Join(parts, "")
}

/*
cabal install --dependencies-only --enable-tests
CallExpr:
  Args:
    Word
      Lit: # WordPart
        Value:
          cabal
      Lit: # WordPart
        Value:
          install
      Lit: # WordPart
        Value:
          --dependencies-only
      Lit: # WordPart
        Value:
           --enable-tests
*/
