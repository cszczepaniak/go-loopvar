// Package integrationtests contains tests that run the analyzer with -fix to check that they
// properly fix findings.

package integrationtests

import (
	"bytes"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixes(t *testing.T) {
	cmd := exec.Command(`go`, `build`, `-o`, `integrationtests/test_binary`)
	cmd.Dir = `..`

	t.Cleanup(func() {
		assert.NoError(t, os.Remove(`test_binary`))
	})

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	fixCmd := exec.Command(`./test_binary`, `-fix`, `./testdata`)
	out, err = fixCmd.CombinedOutput()

	gitRestorePathOnCleanup(t, `./testdata`)

	// We expect an error because the program exits with code 3 when it finds things.
	assert.Error(t, err)
	assert.Equal(t, 3, fixCmd.ProcessState.ExitCode(), `unexpected error: %v`, err)

	exps := getExpectationLoaders()
	for name, loadExpBytes := range exps {
		actual, err := os.ReadFile(filepath.Join(`./testdata`, name+`.go`))
		require.NoError(t, err)

		expBytes := loadExpBytes(t)

		assert.Equal(
			t,
			string(trimBuildTag(t, expBytes)),
			string(trimBuildTag(t, actual)),
		)
	}
}

func trimBuildTag(t *testing.T, bs []byte) []byte {
	ln, rest, ok := bytes.Cut(bs, []byte("\n"))
	require.True(t, ok, `expected to find a newline`)
	require.True(t, bytes.HasPrefix(ln, []byte("//go:build")), `expected first line to contain a build tag`)
	return rest
}

func getExpectationLoaders() map[string]func(t *testing.T) []byte {
	res := make(map[string]func(t *testing.T) []byte)
	filepath.WalkDir(`./testdata`, func(path string, d fs.DirEntry, err error) error {
		_, filename := filepath.Split(path)
		name := strings.TrimSuffix(filename, `_exp.go`)
		if name == filename {
			// It didn't have the desired suffix... skip it
			return nil
		}

		res[name] = func(t *testing.T) []byte {
			t.Helper()

			bs, err := os.ReadFile(path)
			require.NoError(t, err)

			return bs
		}

		return nil
	})

	return res
}

func gitRestorePathOnCleanup(t *testing.T, path string) {
	t.Helper()

	t.Cleanup(func() {
		out, err := exec.Command(`git`, `restore`, path).CombinedOutput()
		require.NoError(t, err, `restore failed: %v`, string(out))
	})
}
