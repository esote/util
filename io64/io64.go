// Package io64 implements IO operations for uint64 lengths.
package io64

import (
	"io"
	"math"
)

// LimitedReader is io.LimitedReader using uint64 instead of int64.
type LimitedReader struct {
	R io.Reader
	N uint64
}

// LimitReader constructs LimitedReader.
func LimitReader(r io.Reader, n uint64) io.Reader {
	return &LimitedReader{
		R: r,
		N: n,
	}
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, io.EOF
	}
	if uint64(len(p)) > l.N {
		p = p[:l.N]
	}
	n, err = l.R.Read(p)
	l.N -= uint64(n)
	return
}

// CopyN is a wrapper of io.CopyN for handling copy length > math.MaxInt64.
func CopyN(dst io.Writer, src io.Reader, l uint64) (written uint64, err error) {
	var n int64
	// Loops twice if l is math.MaxUint64.
	for ; l > uint64(math.MaxInt64); l -= uint64(math.MaxInt64) {
		n, err = io.CopyN(dst, src, math.MaxInt64)
		written += uint64(n)
		if err != nil {
			return
		}
	}
	n, err = io.CopyN(dst, src, int64(l))
	written += uint64(n)
	return
}
