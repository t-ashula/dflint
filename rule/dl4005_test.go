package rule

import "testing"

func TestDL4005(t *testing.T) {
	name := "DL4005"

	shouldExists(name, t)
	shouldValid(name, `SHELL ["/bin/bash", "-c"]`, t)
	shouldInvalid(name, "RUN ln -sfv /bin/bash /bin/sh", t)
	shouldInvalid(name, "RUN ln -s /bin/bash /bin/sh", t)
	shouldInvalid(name, "RUN ln -vfs /bin/bash /bin/sh", t)
	shouldInvalid(name, "RUN ln -vsf /bin/bash /bin/sh", t)
	shouldInvalid(name, "RUN ln -svf /bin/bash /bin/sh", t)
	shouldInvalid(name, "RUN ln --symbolic /bin/bash /bin/sh", t)
	// wrong option should not invalid
	shouldValid(name, "RUN ln -suffix /bin/bash /bin/sh", t)
}
