package fsutil

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file %s not found", src)
	}
	if err != nil {
		return err
	}

	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
