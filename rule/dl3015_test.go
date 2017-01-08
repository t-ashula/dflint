package rule

import "testing"

func TestDL3015(t *testing.T) {
	name := "DL3015"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "RUN apt-get install -y --no-install-recommends python=2.7"
	shouldValid(name, validSource, t)

	invalidSource := "RUN apt-get install -y python=2.7"
	shouldInvalid(name, invalidSource, t)
}
