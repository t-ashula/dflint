package rule

import "testing"

func TestDL4005(t *testing.T) {
	name := "DL4005"

	shouldExists(name, t)
	shouldValid(name, `SHELL ["/bin/bash", "-c"]`, t)
	shouldInvalid(name, "RUN ln -sfv /bin/bash /bin/sh", t)
}
