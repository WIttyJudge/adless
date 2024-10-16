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
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// CopyFileIsExist copies a file from source to destination only in case if
// destination file exists.
func CopyFileIfExist(src, dst string) error {
	if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file %s not found", src)
	}

	return CopyFile(src, dst)
}
