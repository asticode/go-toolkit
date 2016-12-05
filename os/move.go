package os

import "os"

// Move is a cross partitions cancellable move even if files are on different partitions
func Move(src, dst string, channelCancel chan bool) (err error) {
	// Copy
	if err = Copy(src, dst, channelCancel); err != nil {
		return
	}

	// Delete
	err = os.Remove(src)
	return
}
