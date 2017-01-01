package main

import (
	"testing"

	"github.com/t-ashula/dflint/rule"
)

func TestFilterResult(t *testing.T) {
	results := rule.Results{
		&rule.Result{Rule: &rule.Rule{Name: "R0"}},
		&rule.Result{Rule: &rule.Rule{Name: "R1"}},
		&rule.Result{Rule: &rule.Rule{Name: "R1"}},
		&rule.Result{Rule: &rule.Rule{Name: "R2"}},
	}

	ignores := rule.Rules{
		&rule.Rule{Name: "R1"},
		&rule.Rule{Name: "R2"},
	}

	filtered, err := filterResults(results, ignores)
	if err != nil {
		t.Errorf("should no error, %s", err)
	}

	if len(filtered) != 1 {
		t.Errorf("shold filtered, len:%d, %#v", len(filtered), filtered)
	}
}

func TestCheckIgnoreRule(t *testing.T) {
	names := []string{"R1", "R2", "R3", "E0", "R1", "E0"}
	rules := rule.Rules{
		&rule.Rule{Name: "R1"},
		&rule.Rule{Name: "R2"},
		&rule.Rule{Name: "R3"},
		&rule.Rule{Name: "R4"},
	}
	ignores, unknowns := checkIgnoreRule(rules, names)
	if len(ignores) != 3 {
		t.Errorf("ignores: len:%d, %#v\n", len(ignores), ignores)
	}
	if ignores.Find(rule.IsNameMatch("R1")) == nil {
		t.Errorf("R1 should ignore: %s\n", ignores)
	}
	if ignores.Find(rule.IsNameMatch("R2")) == nil {
		t.Errorf("R2 should ignore: %s\n", ignores)
	}
	if ignores.Find(rule.IsNameMatch("R3")) == nil {
		t.Errorf("R3 should ignore: %s\n", ignores)
	}
	if ignores.Find(rule.IsNameMatch("R4")) != nil {
		t.Errorf("R4 should not ignore: %s\n", ignores)
	}
	if len(unknowns) != 1 {
		t.Errorf("unkowns: len:%d, %#v\n", len(unknowns), unknowns)
	}
	if unknowns[0] != "E0" {
		t.Errorf("unkowns: [0]:%s\n", unknowns[0])
	}
}
