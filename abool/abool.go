// Package abool implements an atomic boolean.
package abool

import "sync/atomic"

const (
	f = iota
	t
)

// Bool is an atomic boolean. Must be used as a pointer.
type Bool int32

// New returns an unset boolean.
func New() *Bool {
	return new(Bool)
}

// Set the boolean to true, returns whether the boolean was false.
func (b *Bool) Set() bool {
	return atomic.CompareAndSwapInt32((*int32)(b), f, t)
}

// Unset the boolean to false, returns whether the boolean was true.
func (b *Bool) Unset() bool {
	return atomic.CompareAndSwapInt32((*int32)(b), t, f)
}

// IsSet returns the value of the boolean.
func (b *Bool) IsSet() bool {
	return atomic.LoadInt32((*int32)(b)) == t
}
