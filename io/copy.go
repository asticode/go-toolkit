package io

import (
	"context"
	"errors"
	"io"
	"sync"
)

// Const
const (
	bufferSize = 32 * 1024
)

// Vars
var (
	ErrShortWrite = errors.New("Short write")
)

// bufferPool is a pool of reusable buffers
var bufferPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, bufferSize)
	},
}

// newBuffer creates a new buffer
func newBuffer() []byte {
	return bufferPool.Get().([]byte)
}

// putBuffer puts an buffer back in the pool
func putBuffer(buf []byte) {
	bufferPool.Put(buf)
}

// Copy represents a cancellable copy
func Copy(ctx context.Context, src io.Reader, dst io.Writer) (written int64, err error) {
	var buf = newBuffer()
	defer putBuffer(buf)
	for {
		// Check cancellation
		if err = ctx.Err(); err != nil {
			return
		}

		// Read
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return
}
