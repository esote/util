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
//
// For example in the directory "images" with cutoff of 3, inserting the files
// "transistor", "speaker", and "speach" would create a file tree as such:
//
//	images/
//		spe/
//			ach
//			aker
//		tra/
//			nsistor
//
// This is useful to keep the number of files within a directory manageable.
// The same concept is used with Git object hashes: see ".git/objects/".
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

// Read a file from the splay.
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

func removeEmpty(name string) error {
	f, err := os.Open(name)

	if err != nil {
		return err
	}

	list, err := f.Readdirnames(1)

	if err != io.EOF || len(list) != 0 {
		return f.Close()
	}

	if err = f.Close(); err != nil {
		return err
	}

	return os.Remove(name)
}
