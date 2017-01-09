package rule

import "testing"

func TestDL4001(t *testing.T) {
	name := "DL4001"

	shouldExists(name, t)
	shouldValid(name, "RUN curl http://google.com\nRUN curl http://bing.com", t)
	shouldInvalid(name, "RUN curl http://google.com\nRUN wget http://bing.com", t)
}
