package typehack

import "unsafe"

type Hasher func(unsafe.Pointer, uintptr) uintptr

// Keep sync with src/runtime/runtime2.go
type eface struct {
	_type *MapType
	data  unsafe.Pointer
}

// NewHasher returns a new hash function for the comparable type.
func NewHasher[T comparable]() Hasher {
	var m map[T]struct{}
	tmp := interface{}(m)
	eface := (*eface)(unsafe.Pointer(&tmp))
	return eface._type.Hasher
}
