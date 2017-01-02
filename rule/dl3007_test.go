package rule

import "testing"

func TestDL3007(t *testing.T) {
	name := "DL3007"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "FROM debian:jessie"
	shouldValid(name, validSource, t)

	validSource = "FROM debian@abcd"
	shouldValid(name, validSource, t)

	invalidSource := "FROM debian:latest"
	shouldInvalid(name, invalidSource, t)
}
