package main

import (
	"strings"
	"testing"

	"github.com/t-ashula/dflint/rule"
)

func TestRunShellCheck(t *testing.T) {
	res, err := runShellCheck("/usr/bin/shellcheck", "$foo is $BAR")
	if err != nil {
		t.Fatalf("err:%v\n", err)
	}
	if res == nil {
		t.Fatalf("res:%v\n", err)
	}
}

func TestShellCheckRulecheck(t *testing.T) {
	path, err := shellCheckPath()
	if err != nil || path == "" {
		t.Skip("no shellcheck found")
	}

	script := `RUN echo "$foo is $BAR"`
	fp := strings.NewReader(script)
	ast, err := buildAST(fp)
	r := shellCheckRule(path)
	res := r.Validate(ast)
	if res == true {
		report(rule.GetResults())
		t.Fatalf("failed")
	}
	if len(rule.GetResults()) != 1 {
		report(rule.GetResults())
		t.Fatalf("failed, results size should be 1")
	}
	detail := rule.GetResults()[0]
	if detail == nil {
		t.Fatalf("failed, results[0] should not be nil")
	}
	if detail.Rule == nil {
		t.Fatalf("failed, results[0] should not be nil:%v", detail)
	}
	if !strings.HasPrefix(detail.Rule.Name, "SC") {
		t.Fatalf("failed, results[0] should not be starts with 'SC':%s", detail.Rule.Name)
	}
}
