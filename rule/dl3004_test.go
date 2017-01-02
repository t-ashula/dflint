package rule

import "testing"

func TestDL3004(t *testing.T) {
	name := "DL3004"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "RUN apt-get install"
	shouldValid(name, validSource, t)
	invalidSource := "RUN sudo apt-get install"
	shouldInvalid(name, invalidSource, t)

	invalidSource = "RUN cd&&sudo apt-get install"
	shouldInvalid(name, invalidSource, t)
}
