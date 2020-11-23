// Package splay implements a splayed file tree.
//
// For example in the directory "images" with cutoff of 3, inserting the files
// "transistor", "speaker", and "speech" would create a file tree as such:
//
//	images/
//		spe/
//			ech
//			aker
//		tra/
//			nsistor
//
// This is useful to keep the number of files within a directory manageable.
// The same concept is used with Git object hashes: see ".git/objects/".
package splay

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	// ErrInvalidCutoff only happens when cutoff is zero, in which case
	// Splay should not be used at all.
	ErrInvalidCutoff = errors.New("cutoff is zero")

	// ErrInvalidName means the file name cannot fit the splay because its
	// length is less than cutoff.
	ErrInvalidName = errors.New("file name cannot fit the splay scheme")
)

// Splay represents a file tree where names are splayed according to a cutoff
// index.
type Splay struct {
	cutoff uint64
	dir    string
}

// NewSplay creates a new splay. Names will be splayed at cutoff.
func NewSplay(dir string, cutoff uint64) (*Splay, error) {
	if cutoff == 0 {
		return nil, ErrInvalidCutoff
	}

	dir = filepath.Clean(dir)

	if dir == "." {
		return nil, errors.New("cannot use working directory")
	}

	if err := mkdirExists(dir); err != nil {
		return nil, err
	}

	return &Splay{
		cutoff: cutoff,
		dir:    dir,
	}, nil
}

// Exists checks if a file exists in the splay.
func (s *Splay) Exists(name string) bool {
	_, err := s.Stat(name)

	// Even when os.IsExist(err), if err != nil then the file cannot be used
	// and therefore does not exist from the splay's point of view. This can
	// be due to EACCES, ENAMETOOLONG, etc.
	return err == nil
}

// Open splay file for reading. For simple reading operations prefer Read. One
// good use case is seeking while reading.
//
// The file must be closed by the caller.
func (s *Splay) Open(name string) (*os.File, error) {
	_, file, err := s.parts(name)

	if err != nil {
		return nil, err
	}

	return os.Open(file)
}

// OpenFile opens splay file. For simple reading and writing operations prefer
// Read and Write. Good use cases include writing to a file without overwriting
// it, writing to a file in append-only mode, or seeking while writing.
//
// The file must be closed by the caller.
func (s *Splay) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	dir, file, err := s.parts(name)

	if err != nil {
		return nil, err
	}

	if (flag & os.O_CREATE) == os.O_CREATE {
		// File is being created
		if err = mkdirExists(dir); err != nil {
			return nil, err
		}
	}

	return os.OpenFile(file, flag, perm)
}

// Read the entirety of a file from the splay.
func (s *Splay) Read(name string) ([]byte, error) {
	_, file, err := s.parts(name)

	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(file)
}

// Remove a file from the splay.
func (s *Splay) Remove(name string) error {
	dir, file, err := s.parts(name)

	if err != nil {
		return err
	}

	if err = os.Remove(file); err != nil {
		return err
	}

	return removeEmpty(dir)
}

// RemoveAll splay contents including the splay directory itself.
func (s *Splay) RemoveAll() error {
	return os.RemoveAll(s.dir)
}

// Stat returns os.FileInfo describing the splay file. For checking file
// existence prefer Exists.
func (s *Splay) Stat(name string) (os.FileInfo, error) {
	_, file, err := s.parts(name)

	if err != nil {
		return nil, err
	}

	return os.Stat(file)
}

// Write a file to the splay. Overwrites existing files.
func (s *Splay) Write(name string, data []byte) error {
	dir, file, err := s.parts(name)

	if err != nil {
		return err
	}

	if err = mkdirExists(dir); err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0600)
}

func (s *Splay) parts(name string) (dir string, file string, err error) {
	if uint64(len(name)) <= s.cutoff {
		err = ErrInvalidName
		return
	}

	dir = filepath.Join(s.dir, name[:s.cutoff])
	file = filepath.Join(dir, name[s.cutoff:])
	return
}

func mkdirExists(dir string) error {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return os.Mkdir(dir, 0700)
	}

	return err
}

func removeEmpty(dir string) error {
	d, err := os.Open(dir)

	if err != nil {
		return err
	}

	list, err := d.Readdirnames(1)

	if err != io.EOF || len(list) != 0 {
		return d.Close()
	}

	return os.Remove(dir)
}
