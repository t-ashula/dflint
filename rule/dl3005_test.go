package rule

import "testing"

func TestDL3005(t *testing.T) {
	name := "DL3005"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "RUN apt-get update"
	shouldValid(name, validSource, t)

	invalidSource := "RUN apt-get update && apt-get upgrade"
	shouldInvalid(name, invalidSource, t)

	invalidSource = "RUN apt-get upgrade"
	shouldInvalid(name, invalidSource, t)

	invalidSource = "RUN apt-get dist-upgrade"
	shouldInvalid(name, invalidSource, t)
}
