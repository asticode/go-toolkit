package os

import (
	"context"
	"os"
)

// Move is a cross partitions cancellable move even if files are on different partitions
func Move(ctx context.Context, src, dst string) (err error) {
	// Copy
	if err = Copy(ctx, src, dst); err != nil {
		return
	}

	// Delete
	err = os.Remove(src)
	return
}
