package internal

import (
	"unsafe"
)

type CObjWrapper struct {
	CObjPointer unsafe.Pointer
}

func Wrap[T any](w *T, p unsafe.Pointer) {
	var w1 = (*CObjWrapper)(unsafe.Pointer(w))
	w1.CObjPointer = p
}

func WrapNew[T any](p unsafe.Pointer) *T {
	var w = new(T)
	Wrap(w, p)
	return w
}

func Unwrap[T any](w *T) unsafe.Pointer {
	var w1 = (*CObjWrapper)(unsafe.Pointer(w))
	return w1.CObjPointer
}
