// Package table provides file-based tabular deletion, indexing, and insertion
// of fixed-width rows ordered by insertion time.
//
// Due to the indexing and memory patterns, this has very narrow use cases.
//
// Time complexities: Delete O(n), IndexN O(n), Insert O(1), InsertUnique O(n).
// Space complexities: Delete O(n), IndexN O(n), Insert and InsertUnique O(1).
package table

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/esote/util/splay"
)

// Table is a wrapper around Splay.
type Table struct {
	Splay *splay.Splay

	length int
	mu     sync.Mutex
}

// NewTable constructs a splay and creates a new Table.
func NewTable(dir string, cutoff uint64, rowLength int) (t *Table, err error) {
	if rowLength <= 0 {
		return nil, errors.New("invalid rowLength")
	}

	t = &Table{
		length: rowLength,
	}

	t.Splay, err = splay.NewSplay(dir, cutoff)
	return
}

// Delete row from table. Shifts all others rows down.
func (t *Table) Delete(key, row string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.Splay.Exists(key) {
		return errors.New("no such key")
	}

	if len(row) != t.length {
		return errors.New("row length invalid")
	}

	f, err := t.Splay.OpenFile(key, os.O_RDWR, 0)

	if err != nil {
		return err
	}

	defer f.Close()

	count, err := readCount(f)

	if err != nil {
		return err
	}

	buf := make([]byte, t.length)
	var index int64

	// Search for row.
	for {
		if _, err = f.Read(buf); err == io.EOF {
			// Row not found.
			return nil
		} else if err != nil {
			return err
		}

		index++

		if string(buf) == row {
			break
		}

	}

	// Shift rows down.
	rat := 8 + index*int64(t.length)
	wat := rat - int64(t.length)

	size := (int64(count) - index) * int64(t.length)
	buf = make([]byte, size)

	if _, err = f.ReadAt(buf, rat); err != nil {
		return err
	}

	if _, err = f.WriteAt(buf, wat); err != nil {
		return err
	}

	// Truncate last row.
	size = 8 + int64(count-1)*int64(t.length)
	if err = f.Truncate(size); err != nil {
		return err
	}

	_, err = f.WriteAt(encodeCount(count-1), 0)
	return err
}

// IndexN returns n table rows in order of latest first. If n == 0, all rows
// will be returned.
func (t *Table) IndexN(key string, n uint64) ([]string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.Splay.Exists(key) {
		return nil, errors.New("no such key")
	}

	f, err := t.Splay.Open(key)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	count, err := readCount(f)

	if err != nil {
		return nil, err
	}

	var items []string

	if n == 0 || n >= count {
		n = count
	}

	items = make([]string, n)

	offset := int64(len(items) * t.length)

	if _, err = f.Seek(-offset, io.SeekEnd); err != nil {
		return nil, err
	}

	buf := make([]byte, t.length)

	for i := len(items) - 1; i >= 0; i-- {
		if _, err = f.Read(buf); err != nil {
			return nil, err
		}

		items[i] = string(buf)
	}

	return items, nil
}

// Insert row into table.
func (t *Table) Insert(key, row string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(row) != t.length {
		return errors.New("row length invalid")
	}

	if t.Splay.Exists(key) {
		return t.appendRow(key, row)
	}

	return t.create(key, row)
}

// InsertUnique row into table. If the row already exists, no change is made to
// the key file.
func (t *Table) InsertUnique(key, row string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(row) != t.length {
		return errors.New("row length invalid")
	}

	if t.Splay.Exists(key) {
		return t.appendUniqueRow(key, row)
	}

	return t.create(key, row)
}

func (t *Table) create(key, row string) error {
	var b bytes.Buffer
	b.Grow(8 + t.length)

	_, _ = b.Write(encodeCount(1))
	_, _ = b.WriteString(row)

	return t.Splay.Write(key, b.Bytes())
}

func (t *Table) appendRow(key, row string) error {
	f, err := t.Splay.OpenFile(key, os.O_RDWR, 0)

	if err != nil {
		return err
	}

	defer f.Close()

	count, err := readCount(f)

	if err != nil {
		return err
	}

	if _, err = f.WriteAt(encodeCount(count+1), 0); err != nil {
		return err
	}

	if _, err = f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	_, err = f.WriteString(row)
	return err
}

func (t *Table) appendUniqueRow(key, row string) error {
	f, err := t.Splay.OpenFile(key, os.O_RDWR, 0)

	if err != nil {
		return err
	}

	defer f.Close()

	count, err := readCount(f)

	if err != nil {
		return err
	}

	buf := make([]byte, t.length)

	for {
		if _, err = f.Read(buf); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if string(buf) == row {
			return nil
		}
	}

	// At end of file: write row and increment count.

	if _, err = f.WriteString(row); err != nil {
		return err
	}

	_, err = f.WriteAt(encodeCount(count+1), 0)
	return err
}

func encodeCount(count uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, count)
	return b
}

// Read and decode count. Assumes seek offset 0.
func readCount(f *os.File) (out uint64, err error) {
	b := make([]byte, 8)
	if _, err = f.Read(b); err != nil {
		return
	}
	out = binary.LittleEndian.Uint64(b)
	return
}
