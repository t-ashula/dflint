package rule

import "testing"

func TestDL4004(t *testing.T) {
	name := "DL4004"

	shouldExists(name, t)
	shouldValid(name, "ENTRYPOINT /bin/true", t)
	shouldInvalid(name, "ENTRYPOINT /bin/true\nENTRYPOINT /bin/false", t)
}
