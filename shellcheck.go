package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/t-ashula/dflint/rule"
)

type SCResult struct {
	File      string `json:"file"`
	Line      int    `json:"line"`
	EndLine   int    `json:"endLine"`
	Column    int    `json:"column"`
	EndColumn int    `json:"endColumn"`
	Level     string `json:"level"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
}

func runShellCheck(path, script string) (rule.Results, error) {
	tmpfile, err := ioutil.TempFile("", "dflint")
	if err != nil {
		fmt.Printf("create temp file error:%v\n", err)
		return nil, err
	}

	content := "#!/bin/sh\n"
	content += script
	content += "\n"
	if _, err := tmpfile.WriteString(content); err != nil {
		tmpfile.Close()
		fmt.Printf("write string to %s error:%v\n", tmpfile.Name(), err)
		return nil, err
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Printf("temp file close error:%v\n", err)
		return nil, err
	}

	defer os.Remove(tmpfile.Name())
	fmt.Printf("path:%s, file:%s\n", path, tmpfile.Name())

	stdout, err := exec.Command(path, "-f", "json", "-s", "sh", tmpfile.Name()).Output()
	if err != nil {
		fmt.Printf("exec command error:%v\n", err)
	}

	var scrs []SCResult
	json.Unmarshal(stdout, &scrs)

	var ress rule.Results
	for _, scr := range scrs {
		res := &rule.Result{
			Rule: &rule.Rule{
				Name:        fmt.Sprintf("SC%d", scr.Code),
				Description: scr.Message,
				Severity:    convertSeverity(scr.Level),
			},
			StartColumn: scr.Column,
			EndColumn:   scr.EndColumn,
		}
		ress = append(ress, res)
	}

	return ress, nil
}

func convertSeverity(sclevel string) rule.Severity {
	switch sclevel {
	case "error":
		return rule.SeverityError
	case "warning":
		return rule.SeverityWarning
	case "info":
		return rule.SeverityInfo
	case "style":
		return rule.SeverityStyle
	default:
		return rule.SeverityUnknown
	}
}

func shellCheckPath() (string, error) {
	paths := getEnvPaths()
	if paths == nil {
		paths = []string{}
	}
	for _, p := range paths {
		sc := path.Join(p, "shellcheck")
		if _, err := os.Stat(sc); err == nil {
			return sc, nil
		}
	}
	return "", errors.New("Not installed")
}

func shellCheckRule(path string) *rule.Rule {
	r := &rule.Rule{
		Name: "SC", // dummy
	}
	r.Validate = func(root *parser.Node) bool {
		valid := true
		for _, child := range root.Children {
			if child.Value != "run" {
				continue
			}
			arg := child.Next
			if arg == nil || arg.Value == "" {
				continue
			}

			ress, err := runShellCheck(path, arg.Value)
			if err != nil {
				// TODO report error
				valid = false
				continue
			}
			if ress != nil {
				valid = false
				for _, res := range ress {
					rule.AppendResult(res.Rule, child)
				}
				continue
			}
		}
		return valid
	}
	return r
}
