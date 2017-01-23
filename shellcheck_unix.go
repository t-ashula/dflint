package main

import (
	"os"
	"strings"
)

func getEnvPaths() []string {
	envPath := os.ExpandEnv(os.Getenv("PATH"))
	if envPath == "" {
		return nil
	}
	return strings.Split(envPath, ":")
}
