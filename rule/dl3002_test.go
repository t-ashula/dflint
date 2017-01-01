package rule

import "testing"

func TestDL3002(t *testing.T) {
	name := "DL3002"

	shouldExists(name, t)
	shouldValid(name, "", t)
	shouldValid(name, "USER foo", t)
	shouldInvalid(name, "USER root", t)
}
