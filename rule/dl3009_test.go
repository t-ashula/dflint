package rule

import "testing"

func TestDL3009(t *testing.T) {
	name := "DL300"

	shouldExists(name, t)
	t.Skip("not supported. yet")
	// 	shouldValid(name, "", t)

	// 	validSource :=
	// 		`RUN apt-get update && apt-get install -y python \
	//  && apt-get clean \
	//  && rm -rf /var/lib/apt/lists/*
	// `
	// 	shouldValid(name, validSource, t)

	// 	invalidSource := "RUN apt-get update && apt-get install -y python"
	// 	shouldInvalid(name, invalidSource, t)
}
