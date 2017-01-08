package rule

import "testing"

func TestDL3020(t *testing.T) {
	name := "DL3020"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "COPY requirements.txt /usr/src/app/"
	shouldValid(name, validSource, t)

	invalidSource := "ADD requirements.txt /usr/src/app/"
	shouldInvalid(name, invalidSource, t)
}
