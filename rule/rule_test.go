package rule

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func TestRuleFind(t *testing.T) {
	rules := Rules{
		&Rule{
			Name: "R1",
		},
		&Rule{
			Name: "R2",
		},
	}

	r0 := rules.Find(func(r *Rule) bool { return r.Name == "R0" })
	if r0 != nil {
		t.Errorf("r0 should not found")
	}

	r1 := rules.Find(func(r *Rule) bool { return r.Name == "R1" })
	if r1 == nil {
		t.Errorf("r1 should found")
	}

	r2 := rules.Find(IsNameMatch("R2"))
	if r2 == nil {
		t.Errorf("r2 +isNameMatch should found")
	}

	r0 = rules.Find(IsNameMatch("R0"))
	if r0 != nil {
		t.Errorf("r0 +isNameMatch should not found")
	}
}

// utilties for test
func buildAST(fp io.Reader) (*parser.Node, error) {
	d := parser.Directive{
		EscapeSeen:           false,
		LookingForDirectives: true,
	}
	parser.SetEscapeToken(parser.DefaultEscapeToken, &d)
	return parser.Parse(fp, &d)
}

func dumpAST(node *parser.Node) string {
	if node == nil {
		return ""
	}
	dump := "Node"
	dump += "{"
	dump += fmt.Sprintf("Value:\"%s\", ", node.Value)
	dump += fmt.Sprintf("Original:\"%s\", ", node.Original)
	dump += fmt.Sprintf("StartLine:%d, ", node.StartLine)
	dump += fmt.Sprintf("EndLine:%d, ", node.EndLine)
	dump += fmt.Sprintf("Attributes:%#v, ", node.Attributes) // TODO: more detail?
	dump += fmt.Sprintf("Flags:%#v, ", node.Flags)           // TODO: more detail?
	if node.Next == nil {
		dump += fmt.Sprintf("Next:%s, ", "(nil)")
	} else {
		dump += fmt.Sprintf("Next:%s, ", dumpAST(node.Next))
	}
	dump += "Children:["
	if node.Children != nil {
		for _, child := range node.Children {
			dump += dumpAST(child)
			dump += ", "
		}
	}
	dump += "], "

	dump += "} "
	return dump
}

func shold(name string, t *testing.T, fn func(rule *Rule, t *testing.T)) {
	rules := GetRules()
	rule := rules.Find(IsNameMatch(name))
	fn(rule, t)
}

func shouldExists(name string, t *testing.T) {
	shold(name, t, func(rule *Rule, t *testing.T) {
		if rule == nil {
			t.Errorf("rule %s not found.\n", name)
		}
	})
}

func shouldNotExists(name string, t *testing.T) {
	shold(name, t, func(rule *Rule, t *testing.T) {
		if rule != nil {
			t.Errorf("rule %s should not found.\n", name)
		}
	})
}

// TODO: refactor valid/invalid
func shouldValid(name string, source string, t *testing.T) {
	shold(name, t, func(rule *Rule, t *testing.T) {
		r := strings.NewReader(source)
		root, err := buildAST(r)
		if err != nil {
			t.Fatalf("parse failed. docker API changed? source:'%s', err:%s\n", source, err)
		}

		ok := rule.Validate(root)
		if !ok {
			t.Errorf("should Valid but Invalid. source:'%s', AST:%s\n", source, dumpAST(root))
		}
	})
}

// TODO: refactor valid/invalid
func shouldInvalid(name string, source string, t *testing.T) {
	shold(name, t, func(rule *Rule, t *testing.T) {
		r := strings.NewReader(source)
		root, err := buildAST(r)
		if err != nil {
			t.Fatalf("parse failed. docker API changed? source:'%s', err:%s\n", source, err)
		}

		ok := rule.Validate(root)
		if ok {
			t.Errorf("should Invalid but Valid. source:'%s', AST:%s\n", source, dumpAST(root))
		}
	})
}
