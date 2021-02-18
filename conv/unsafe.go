package conv

import (
	"reflect"
	"unsafe"
)

func UnsafeBytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func UnsafeStrToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Len:  len(s),
		Cap:  len(s),
		Data: (*reflect.StringHeader)(unsafe.Pointer(&s)).Data,
	}))
}
