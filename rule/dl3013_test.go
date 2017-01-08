package rule

import "testing"

func TestDL3013(t *testing.T) {
	name := "DL3013"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "RUN pip install django==1.9"
	shouldValid(name, validSource, t)

	invalidSource := "RUN pip install django"
	shouldInvalid(name, invalidSource, t)
}
