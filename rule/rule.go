package rule

import "github.com/docker/docker/builder/dockerfile/parser"

// Validator validate node,true -> good, false -> bad
type Validator func(root *parser.Node) bool

// Severity stand check type
type Severity struct {
	Name string
}

// Rule
type Rule struct {
	Name        string
	Severity    Severity
	Description string
	Validate    Validator
}

type Rules []*Rule

var (
	SeverityError   = Severity{Name: "ErrorC"}
	SeverityWarning = Severity{Name: "WarningC"}
	SeverityInfo    = Severity{Name: "InfoC"}
	SeverityStyle   = Severity{Name: "StyleC"}
)

var rules Rules

func RegisterRule(rule *Rule) {
	rules = append(rules, rule)
}

func GetRules() Rules {
	return rules
}
