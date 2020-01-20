package fragmented

import "io"

// Selector returns the next node.
type Selector interface {
	Next() (io.ReadWriteCloser, error)
}

// Fragmented breaks IO operations across a series of nodes to form an
// amalgamated read-writer. Can be used to construct a mesh network by having
// the selector maintain a list of network participants.
type Fragmented struct {
	curr    io.ReadWriteCloser
	s       Selector
	size, n int
}

// New constructs a new fragmented read-writer. Returns nil if size < 0.
func New(s Selector, size int) *Fragmented {
	if size < 0 {
		return nil
	}
	return &Fragmented{
		s:    s,
		size: size,
	}
}

// Reads until the size is saturated, then selects the next reader.
func (f *Fragmented) Read(p []byte) (n int, err error) {
	for {
		if f.curr == nil {
			if f.curr, err = f.s.Next(); err != nil {
				return
			}
		}

		var m int
		m, err = f.curr.Read(p[n:])
		n += m

		if m == 0 || n == len(p) || (err != nil && err != io.EOF) {
			return
		}

		if err = f.curr.Close(); err != nil {
			return
		}

		f.curr = nil
	}
}

// Writes until the size is saturated, then selects the next writer.
func (f *Fragmented) Write(p []byte) (n int, err error) {
	for {
		if f.curr == nil {
			f.n = 0
			if f.curr, err = f.s.Next(); err != nil {
				return
			}
		}

		var m int
		m, err = f.curr.Write(p[n:min(len(p), n+f.size-f.n)])
		f.n += m
		n += m

		if m == 0 || n == len(p) || err != nil {
			return
		}

		if err = f.curr.Close(); err != nil {
			return
		}

		f.curr = nil
	}
}

// Close the fragmented read-writer.
func (f *Fragmented) Close() error {
	if f.curr == nil {
		return nil
	}
	return f.curr.Close()
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
