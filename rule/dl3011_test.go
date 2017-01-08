package rule

import "testing"

func TestDL3011(t *testing.T) {
	name := "DL3011"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "EXPOSE 65535"
	shouldValid(name, validSource, t)

	invalidSource := "EXPOSE 80000"
	shouldInvalid(name, invalidSource, t)

	invalidSource = "EXPOSE -2000"
	shouldInvalid(name, invalidSource, t)

	invalidSource = "EXPOSE "
	shouldInvalid(name, invalidSource, t)
}
