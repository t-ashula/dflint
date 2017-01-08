package rule

import "testing"

func TestDL3010(t *testing.T) {
	name := "DL3010"

	shouldExists(name, t)
	shouldValid(name, "", t)

	validSource := "COPY Gemfile /app/"
	shouldValid(name, validSource, t)

	validSource = "ADD rootfs.tar.xz /app/"
	shouldValid(name, validSource, t)

	invalidSource := "COPY rootfs.tar.xz /"
	shouldInvalid(name, invalidSource, t)

	invalidSource = "COPY rootfs.tar.xz /rootfs.tar.xz"
	shouldInvalid(name, invalidSource, t)

	invalidSource = "COPY [\"root fs.tar.xz\", \"/path\"]"
	shouldInvalid(name, invalidSource, t)
}
