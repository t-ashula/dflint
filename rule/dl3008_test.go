package rule

import "testing"

func TestDL3008(t *testing.T) {
	name := "DL3008"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "RUN apt-get install python=2.7 emacs=24"
	shouldValid(name, validSource, t)

	invalidSource := "RUN apt-get install python emacs"
	shouldInvalid(name, invalidSource, t)

	invalidSource = "RUN apt-get install python emacs=24"
	shouldInvalid(name, invalidSource, t)
}
