package rule

import "testing"

func TestDL5000(t *testing.T) {
	name := "DL5000"

	shouldExists(name, t)
	shouldValid(name, "", t)
	shouldInvalid(name, "MAINTAINER foo@example.com", t)
}
