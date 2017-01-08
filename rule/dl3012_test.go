package rule

import "testing"

func TestDL3012(t *testing.T) {
	name := "DL3012"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "MAINTAINER Lukas Martinelli <me@lukasmartinelli.ch>"
	shouldValid(name, validSource, t)

	validSource = "MAINTAINER Lukas Martinelli http://lukasmartinelli.ch"
	shouldValid(name, validSource, t)

	invalidSource := "MAINTAINER Lukas"
	shouldInvalid(name, invalidSource, t)
}
