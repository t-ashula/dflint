package rule

import "testing"

func TestDL3000(t *testing.T) {
	name := "DL3000"

	shouldExists(name, t)
	shouldValid(name, "WORKDIR /usr/src/app", t)
	shouldInvalid(name, "WORKDIR usr/src/app", t)

	shouldValid(name, "WORKDIR ${foo}", t) // TODO: uncheckable state
}

func TestDL4000(t *testing.T) {
	name := "DL4000"

	shouldExists(name, t)
	shouldValid(name, "MAINTAINER foo@example.com", t)
	shouldInvalid(name, "MAINTAINER", t)
	shouldInvalid(name, "", t)
}
