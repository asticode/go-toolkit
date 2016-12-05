package os

import (
	"context"
	"os"
)

// Move is a cross partitions cancellable move even if files are on different partitions
func Move(src, dst string, ctx context.Context) (err error) {
	// Copy
	if err = Copy(src, dst, ctx); err != nil {
		return
	}

	// Delete
	err = os.Remove(src)
	return
}
