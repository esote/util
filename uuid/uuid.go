// Package uuid implements functions for generating cryptographically-random
// UUIDs or MegaUUIDs.
package uuid

import (
	"crypto/rand"
)

const (
	// LenUUID is the byte length of a UUID (128 bits).
	LenUUID = 16

	// LenMegaUUID is the byte length of a MegaUUID (1024 bits).
	LenMegaUUID = 128
)

// NewUUID generates a new UUID.
func NewUUID() ([]byte, error) {
	u := make([]byte, LenUUID)
	_, err := rand.Read(u)
	return u, err
}

// NewMegaUUID generates a new MegaUUID.
func NewMegaUUID() ([]byte, error) {
	u := make([]byte, LenMegaUUID)
	_, err := rand.Read(u)
	return u, err
}
