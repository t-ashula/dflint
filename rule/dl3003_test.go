package rule

import "testing"

func TestDL3003(t *testing.T) {
	name := "DL3003"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "RUN git clone git@github.com:lukasmartinelli/hadolint.git /usr/src/app"
	shouldValid(name, validSource, t)
	invalidSource := "RUN cd /usr/src/app && git clone git@github.com:lukasmartinelli/hadolint.git"
	shouldInvalid(name, invalidSource, t)

	invalidSource = "RUN cd&&cd"
	shouldInvalid(name, invalidSource, t)
}
