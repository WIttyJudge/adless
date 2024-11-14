package fsutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyFile(t *testing.T) {
	td, err := os.MkdirTemp("", "barrier-fsutil")
	require.NoError(t, err)
	defer os.RemoveAll(td)

	t.Run("return error since there is no source file", func(t *testing.T) {
		src := td + "_nofile"
		dst := src + "_dst"
		err = CopyFile(src, dst)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "file not found")
	})

	t.Run("returns error if unable to create destination file", func(t *testing.T) {
		srcFile, err := os.CreateTemp(td, "src_")
		require.NoError(t, err)
		_, err = srcFile.WriteString("some content")
		require.NoError(t, err)

		dst := "/invalid/destination.txt"
		err = CopyFile(srcFile.Name(), dst)
		assert.Error(t, err)
	})

	t.Run("returns no errors and succefully makes copy", func(t *testing.T) {
		srcFile, err := os.CreateTemp(td, "src_")
		require.NoError(t, err)

		_, err = srcFile.WriteString("some content")
		require.NoError(t, err)

		dst := srcFile.Name() + "_dst"

		err = CopyFile(srcFile.Name(), dst)
		require.FileExists(t, dst)
		require.NoError(t, err)

		dstContent, err := os.ReadFile(dst)
		require.NoError(t, err)
		srcContent, err := os.ReadFile(srcFile.Name())
		require.NoError(t, err)

		assert.EqualValues(t, dstContent, srcContent)
	})
}
