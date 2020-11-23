package e2e

import (
	"os/exec"
	"testing"

	"github.com/onsi/gomega/gexec"
	"github.com/stretchr/testify/assert"
)

func TestE2E(t *testing.T) {
	mcrawl, err := gexec.Build("../../cmd/mcrawl/")
	assert.NoError(t, err)

	server := startMockServer()
	defer server.Close()

	cmd := exec.Command(mcrawl, "--workers", "2", server.URL)
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Contains(t, string(output), "Unique URLs crawled: 3")

	gexec.CleanupBuildArtifacts()
}
