package quickhttp

// copy from fasthttp

func addLeadingSlash(dst, src []byte) []byte {
	// zero length 、"C:/" and "a" case
	isDisk := len(src) > 2 && src[1] == ':'
	if len(src) == 0 || (!isDisk && src[0] != '/') {
		dst = append(dst, '/')
	}

	return dst
}
