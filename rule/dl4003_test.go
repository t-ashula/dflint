package rule

import "testing"

func TestDL4003(t *testing.T) {
	name := "DL4003"

	shouldExists(name, t)
	shouldValid(name, "CMD /bin/true", t)
	shouldInvalid(name, "CMD /bin/true\nCMD /bin/false", t)
}
