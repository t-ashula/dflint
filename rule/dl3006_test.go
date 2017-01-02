package rule

import "testing"

func TestDL3006(t *testing.T) {
	name := "DL3006"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "FROM debian:jessie"
	shouldValid(name, validSource, t)

	validSource = "FROM debian@abcdefg"
	shouldValid(name, validSource, t)

	invalidSource := "FROM debian"
	shouldInvalid(name, invalidSource, t)

	invalidSource = "FROM debian:"
	shouldInvalid(name, invalidSource, t)

	invalidSource = "FROM debian@"
	shouldInvalid(name, invalidSource, t)
}
