// Package fcmp implements functions to compare files quickly.
package fcmp

import (
	"bytes"
	"io"
	"os"
)

// Files compares files quickly.
//
// Begins comparison at the current file (seek) offset, but preserves that
// offset to reload after comparison finishes. Despite this, Files will return
// false if the files differ in size regardless of offset. Use Bare if you wish
// to disregard file size and offset reloading.
func Files(x, y *os.File) (bool, error) {
	statx, err := x.Stat()

	if err != nil {
		return false, err
	}

	staty, err := y.Stat()

	if err != nil {
		return false, err
	}

	if statx.Size() != staty.Size() {
		return false, nil
	}

	// Get current file offsets to reload later.
	offsetx, err := x.Seek(0, os.SEEK_CUR)

	if err != nil {
		return false, err
	}

	offsety, err := y.Seek(0, os.SEEK_CUR)

	if err != nil {
		return false, err
	}

	equal, err := Bare(x, y)

	// When Bare's error is nil, care about offset reload errors.
	if _, errx := x.Seek(offsetx, os.SEEK_SET); err == nil {
		err = errx
	}

	if _, erry := y.Seek(offsety, os.SEEK_SET); err == nil {
		err = erry
	}

	return equal, err

}

// Paths compares files quickly, handling file open and close as a wrapper
// for Files.
func Paths(x, y string) (bool, error) {
	if x == y {
		return true, nil
	}

	fx, err := os.Open(x)

	if err != nil {
		return false, err
	}

	fy, err := os.Open(y)

	if err != nil {
		return false, err
	}

	same, err := Files(fx, fy)

	// When Files' error is nil, care about close errors.
	if errx := fx.Close(); err == nil {
		err = errx
	}

	if erry := fy.Close(); err == nil {
		err = erry
	}

	return same, err
}

// Bare compares files quickly, disregarding file size.
//
// Bare will not preserve the file offset: for this look to Files.
func Bare(x, y *os.File) (bool, error) {
	const l = 4096

	bufx := make([]byte, l)
	bufy := make([]byte, l)

	for {
		nx, errx := x.Read(bufx)
		ny, erry := y.Read(bufy)

		if errx == io.EOF && erry == io.EOF {
			return true, nil
		} else if errx != nil {
			return false, errx
		} else if erry != nil {
			return false, erry
		}

		if nx != ny || !bytes.Equal(bufx, bufy) {
			return false, nil
		}
	}
}
