package os

import (
	"os"

	"context"

	"github.com/asticode/go-toolkit/io"
)

// Copy is a cross partitions cancellable copy
func Copy(src, dst string, ctx context.Context) (err error) {
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}
	defer dstFile.Close()

	// Copy the content
	_, err = io.Copy(srcFile, dstFile, ctx)
	return
}
