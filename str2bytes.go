package quickhttp

import "unsafe"

func str2bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
