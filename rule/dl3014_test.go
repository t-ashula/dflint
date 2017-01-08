package rule

import "testing"

func TestDL3014(t *testing.T) {
	name := "DL3014"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "RUN apt-get -y install python=2.7"
	shouldValid(name, validSource, t)

	invalidSource := "RUN apt-get install python=2.7"
	shouldInvalid(name, invalidSource, t)
}
