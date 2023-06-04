package internal

// #include <stdlib.h>
// typedef void * _cgo_vulkan_void_ptr;
import "C"

import (
	"unsafe"
)

type CHandleWrapper[T any] struct {
	CHandle T
}

func Wrap[T any](w unsafe.Pointer, p *T) {
	var w1 = (*CHandleWrapper[T])(w)
	w1.CHandle = *p
}

func Unwrap[T any](w unsafe.Pointer) *T {
	var w1 = (*CHandleWrapper[T])(w)
	return &w1.CHandle
}

func CStringArray(a []string) (unsafe.Pointer, []func()) {

	if nil == a || 0 == len(a) {
		return nil, nil
	}

	var r []func()

	var cnt = len(a)

	// Assume pointer size = sizeof(uint)
	// const sizeof_p = C.sizeof_uint
	const sizeof_p = C.sizeof__cgo_vulkan_void_ptr
	var p = C.malloc(C.ulonglong(cnt * sizeof_p))

	r = append(r, func() {
		C.free(p)
	})

	for i := 0; i < cnt; i++ {

		var s = C.CString(a[i])
		r = append(r, func() {
			C.free(unsafe.Pointer(s))
		})

		var p1 = (**C.char)(unsafe.Pointer(uintptr(p) + uintptr(i*sizeof_p)))
		*p1 = s
	} // for

	return p, r
}

func CallAll(a []func()) {

	if nil == a {
		return
	}

	for _, f := range a {
		if nil != f {
			f()
		}
	}
}
