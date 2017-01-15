package rule

import (
	"errors"
	"strings"

	"github.com/mvdan/sh/syntax"
)

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
