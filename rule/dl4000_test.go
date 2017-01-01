package rule

import "testing"

func TestDL4000(t *testing.T) {
	name := "DL4000"

	shouldExists(name, t)
	shouldValid(name, "MAINTAINER foo@example.com", t)
	shouldInvalid(name, "MAINTAINER", t)
	shouldInvalid(name, "", t)
}
