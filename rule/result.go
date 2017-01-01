package rule

import "github.com/docker/docker/builder/dockerfile/parser"

type Result struct {
	Rule        *Rule
	StartLine   int
	EndLine     int
	StartColumn int
	EndColumn   int
}

type Results []*Result

var results Results

func GetResults() Results {
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
