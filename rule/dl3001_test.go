package rule

import (
	"fmt"
	"testing"
)

func TestDL3001(t *testing.T) {
	name := "DL3001"

	shouldExists(name, t)
	shouldValid(name, "", t)
	shouldValid(name, "RUN apt-get", t)
	// TOOD: how to maintain NG command list
	ngcmds := []string{"shutdown", "service", "ps", "free", "top", "kill", "mount", "ifconfig", "nano", "vim", "emacs"}
	for _, cmd := range ngcmds {
		shouldInvalid(name, fmt.Sprintf("RUN %s", cmd), t)
	}
	// shouldInvalid(name, "RUN apt-get && nano", t)
}
