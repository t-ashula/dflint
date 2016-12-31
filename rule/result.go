package rule

import "github.com/docker/docker/builder/dockerfile/parser"

type Result struct {
	Rule        *Rule
	StartLine   int
	EndLine     int
	StartColumn int
	EndColumn   int
}

var results []*Result

func GetResults() []*Result {
	return results
}

func AppendResult(rule *Rule, node *parser.Node) {
	res := &Result{
		Rule: rule,
	}
	if node != nil {
		res.StartLine = node.StartLine
		res.EndLine = node.EndLine
	}
	results = append(results, res)
}
