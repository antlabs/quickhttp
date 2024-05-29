package quickhttp

import (
	"net/http"
	"strconv"
)

func formatStatusLine(dst, protocol []byte, statusCode int, statusText string) []byte {
	dst = append(dst, protocol...)
	dst = append(dst, ' ')
	dst = strconv.AppendInt(dst, int64(statusCode), 10)
	dst = append(dst, ' ')
	if len(statusText) == 0 {
		dst = append(dst, str2bytes(http.StatusText(statusCode))...)
	} else {
		dst = append(dst, statusText...)
	}
	return append(dst, bytesCRLF...)
}
