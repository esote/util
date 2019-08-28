// Package table provides file-based tabular indexing and insertion of
// fixed-width rows ordered by insertion time.
//
// Due to the indexing patterns, this has very narrow use cases.
//
// Time complexities: IndexN O(n), Insert O(1), InsertUnique O(existing rows).
//
// Space complexities: IndexN O(n), Insert and InsertUnique O(1).
package table

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/esote/util/splay"
)

// Table is a wrapper around Splay.
type Table struct {
	Splay *splay.Splay

	length int
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

// IndexN returns n table rows in order of latest first. If n == 0, all rows
// will be returned.
func (t *Table) IndexN(key string, n uint64) ([]string, error) {
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
		items = make([]string, count)
	} else {
		items = make([]string, n)
	}

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
