package quickhttp

import (
	"net/http"
	"strconv"
)

// status line 保留 50字节
func countDigits(num int) int {
	if num == 0 {
		return 1
	}

	count := 0
	if num < 0 {
		num = -num // 取绝对值处理
		count++    // 负号占一位
	}
	for num > 0 {
		num /= 10
		count++
	}
	return count
}

func countFormatStatusLine(protocol string, statusCode int, statusText string) int {
	n := len(protocol)
	n++ // ' '
	if statusCode >= 0 && statusCode <= 999 {
		n += 3
	} else {
		n += countDigits(statusCode)
	}
	n++ // ' '
	if len(statusText) == 0 {
		n += len(http.StatusText(statusCode))
	} else {
		n += len(statusText)
	}
	return n
}

func formatStatusLine(dst *[]byte, protocol string, statusCode int, statusText string) {
	*dst = append(*dst, protocol...)
	*dst = append(*dst, ' ')
	*dst = strconv.AppendInt(*dst, int64(statusCode), 10)
	*dst = append(*dst, ' ')
	if len(statusText) == 0 {
		*dst = append(*dst, str2bytes(http.StatusText(statusCode))...)
	} else {
		*dst = append(*dst, statusText...)
	}
	*dst = append(*dst, bytesCRLF...)
	return
}
