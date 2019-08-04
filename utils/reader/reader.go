package reader

import (
	"io"
	"sync/atomic"
)

// Reader counts the bytes read through it.
type Reader struct {
	r io.Reader
	n int64
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: r,
	}
}
func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	atomic.AddInt64(&r.n, int64(n))
	return
}

func (r *Reader) N() int64 {
	return atomic.LoadInt64(&r.n)
}
