package fcmp

import (
	"io"
	"os"
)

// FCmp compares files quickly.
func FCmp(x, y *os.File) (bool, error) {
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

	const l = 1000

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

		if nx != ny {
			return false, nil
		}

		for i := 0; i < l; i++ {
			if bufx[i] != bufy[i] {
				return false, nil
			}
		}
	}
}

// FCmpPath compares files quickly, handling file open and close.
func FCmpPath(x, y string) (bool, error) {
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

	same, err := FCmp(fx, fy)

	// when FCmp's error is nil, care about close errors
	if errx := fx.Close(); err == nil {
		err = errx
	}

	if erry := fy.Close(); err == nil {
		err = erry
	}

	return same, err
}
