package quickhttp

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
)

const (
	rChar = byte('\r')
	nChar = byte('\n')
)

var ErrBadTrailer = errors.New("contain forbidden trailer")

type headerState int32

const (
	findHost headerState = 1 << iota
	findUserAgent
)

func (hs *headerState) set(v headerState) {
	*hs |= v
}

func (hs *headerState) clear(v headerState) {
	*hs &= ^v
}

func (hs *headerState) is(v headerState) bool {
	return (*hs & v) > 0
}

type RequestHeader struct {
	parent     *Request
	method     []byte
	host       []byte
	userAgent  []byte
	hState     headerState
	requestURI buf[int32] //uri有非常长的情况，所以使用buf结构存储
}

func (h *RequestHeader) setParent(r *Request) {
	h.parent = r
}

func (h *RequestHeader) SetMethod(method string) {
	h.method = append(h.method[:0], method...)
}

func (h *RequestHeader) SetMethodBytes(method []byte) {
	h.method = append(h.method[:0], method...)
}

func (h *RequestHeader) Reset() {
	h.resetSkipNormalize()
}

func (h *RequestHeader) RequestURI() []byte {
	return h.requestURI.getBytes(*h.parent.headerAndBody.buf)
}

func (h *RequestHeader) resetSkipNormalize() {
	h.method = h.method[:0]
}

type ResponseHeader struct {
	noCopy

	disableNormalizing   bool
	noHTTP11             bool
	connectionClose      bool
	noDefaultContentType bool
	noDefaultDate        bool

	statusCode         int
	statusMessage      string
	server             []byte
	protocol           string
	contentLength      int
	contentLengthBytes []byte
	contentEncoding    []byte
	contentType        []byte
	h                  []argsKV
	trailer            []argsKV
	bufKV              argsKV
	cookies            []argsKV
}

func (h *ResponseHeader) Server() []byte {
	return h.server
}

func (h *ResponseHeader) SetServer(server string) {
	h.server = append(h.server[:0], server...)
}

func (h *ResponseHeader) StatusCode() int {
	if h.statusCode == 0 {
		return http.StatusOK
	}
	return h.statusCode
}

func (h *ResponseHeader) SetStatusCode(statusCode int) {
	h.statusCode = statusCode
}

func (h *ResponseHeader) StatusMessage() string {
	return h.statusMessage
}

func (h *ResponseHeader) SetStatusMessage(statusMessage string) {
	h.statusMessage = statusMessage
}

func (h *ResponseHeader) SetContentEncodingBytes(contentEncoding []byte) {
	h.contentEncoding = append(h.contentEncoding[:0], contentEncoding...)
}

func (h *ResponseHeader) Protocol() string {
	if len(h.protocol) > 0 {
		return h.protocol
	}
	return strHTTP11
}

func (h *ResponseHeader) appendStatusLine(dst *[]byte) {
	statusCode := h.StatusCode()
	if statusCode < 0 {
		statusCode = http.StatusOK
	}
	formatStatusLine(dst, h.Protocol(), statusCode, h.StatusMessage())
}

func (h *ResponseHeader) SetContentLength(contentLength int) {
	h.contentLength = contentLength
	if cap(h.contentLengthBytes) == 0 {
		h.contentLengthBytes = make([]byte, 20)
	}
	h.contentLengthBytes = h.contentLengthBytes[:0]
	if contentLength >= 0 && contentLength <= 9 {
		h.contentLengthBytes = append(h.contentLengthBytes, '0'+byte(contentLength))
	} else {
		h.contentLengthBytes = strconv.AppendInt(h.contentLengthBytes, int64(contentLength), 10)
	}
}

func (h *ResponseHeader) AppendBytes(dst *[]byte) {
	h.appendStatusLine(dst)

	server := h.Server()
	if len(server) != 0 {
		appendHeaderLine(dst, bytesServer, server)
	}

	if len(h.contentLengthBytes) > 0 {
		appendHeaderLine(dst, bytesContentLength, h.contentLengthBytes)
	}

	for i, n := 0, len(h.h); i < n; i++ {
		kv := &h.h[i]

		// Exclude trailer from header
		exclude := false
		for _, t := range h.trailer {
			if bytes.Equal(kv.key, t.key) {
				exclude = true
				break
			}
		}
		if !exclude && (h.noDefaultDate || !bytes.Equal(kv.key, bytesDate)) {
			appendHeaderLine(dst, kv.key, kv.value)
		}
	}
	// if len(h.trailer) > 0 {
	// 	appendHeaderLine(dst, bytesTrailer, appendArgsKeyBytes(nil, h.trailer, strCommaSpace))
	// }

	*dst = append(*dst, bytesCRLF...)
}

func (h *ResponseHeader) ResetConnectionClose() {
	if h.connectionClose {
		h.connectionClose = false
		h.h = delAllArgsBytes(h.h, bytesConnection)
	}
}

func (h *ResponseHeader) SetContentTypeBytes(contentType []byte) {
	h.contentType = append(h.contentType[:0], contentType...)
}

func (h *ResponseHeader) SetConnectionClose() {
	h.connectionClose = true
}

func (h *ResponseHeader) setNonSpecial(key, value []byte) {
	h.h = setArgBytes(h.h, key, value, argsHasValue)
}

func (h *ResponseHeader) SetServerBytes(server []byte) {
	h.server = append(h.server[:0], server...)
}

func (h *ResponseHeader) SetTrailerBytes(trailer []byte) error {
	h.trailer = h.trailer[:0]
	return h.AddTrailerBytes(trailer)
}

func (h *ResponseHeader) AddTrailerBytes(trailer []byte) error {
	var err error
	for i := -1; i+1 < len(trailer); {
		trailer = trailer[i+1:]
		i = bytes.IndexByte(trailer, ',')
		if i < 0 {
			i = len(trailer)
		}
		key := trailer[:i]
		for len(key) > 0 && key[0] == ' ' {
			key = key[1:]
		}
		for len(key) > 0 && key[len(key)-1] == ' ' {
			key = key[:len(key)-1]
		}
		// Forbidden by RFC 7230, section 4.1.2
		if isBadTrailer(key) {
			err = ErrBadTrailer
			continue
		}
		h.bufKV.key = append(h.bufKV.key[:0], key...)
		normalizeHeaderKey(h.bufKV.key, h.disableNormalizing)
		h.trailer = appendArgBytes(h.trailer, h.bufKV.key, nil, argsNoValue)
	}

	return err
}

func isBadTrailer(key []byte) bool {
	if len(key) == 0 {
		return true
	}

	switch key[0] | 0x20 {
	case 'a':
		return caseInsensitiveCompare(key, bytesAuthorization)
	case 'c':
		if len(key) > len(HeaderContentType) && caseInsensitiveCompare(key[:8], bytesContentType[:8]) {
			// skip compare prefix 'Content-'
			return caseInsensitiveCompare(key[8:], bytesContentEncoding[8:]) ||
				caseInsensitiveCompare(key[8:], bytesContentLength[8:]) ||
				caseInsensitiveCompare(key[8:], bytesContentType[8:]) ||
				caseInsensitiveCompare(key[8:], bytesContentRange[8:])
		}
		return caseInsensitiveCompare(key, bytesConnection)
	case 'e':
		return caseInsensitiveCompare(key, bytesExpect)
	case 'h':
		return caseInsensitiveCompare(key, bytesHost)
	case 'k':
		return caseInsensitiveCompare(key, bytesKeepAlive)
	case 'm':
		return caseInsensitiveCompare(key, bytesMaxForwards)
	case 'p':
		if len(key) > len(HeaderProxyConnection) && caseInsensitiveCompare(key[:6], bytesProxyConnection[:6]) {
			// skip compare prefix 'Proxy-'
			return caseInsensitiveCompare(key[6:], bytesProxyConnection[6:]) ||
				caseInsensitiveCompare(key[6:], bytesProxyAuthenticate[6:]) ||
				caseInsensitiveCompare(key[6:], bytesProxyAuthorization[6:])
		}
	case 'r':
		return caseInsensitiveCompare(key, bytesRange)
	case 't':
		return caseInsensitiveCompare(key, bytesTE) ||
			caseInsensitiveCompare(key, bytesTrailer) ||
			caseInsensitiveCompare(key, bytesTransferEncoding)
	case 'w':
		return caseInsensitiveCompare(key, bytesWWWAuthenticate)
	}
	return false
}

func (h *ResponseHeader) setSpecialHeader(key, value []byte) bool {
	if len(key) == 0 {
		return false
	}

	switch key[0] | 0x20 {
	case 'c':
		switch {
		case caseInsensitiveCompare(bytesContentType, key):
			h.SetContentTypeBytes(value)
			return true
		case caseInsensitiveCompare(bytesContentLength, key):
			if contentLength, err := parseContentLength(value); err == nil {
				h.contentLength = contentLength
				h.contentLengthBytes = append(h.contentLengthBytes[:0], value...)
			}
			return true
		case caseInsensitiveCompare(bytesContentEncoding, key):
			h.SetContentEncodingBytes(value)
			return true
		case caseInsensitiveCompare(bytesConnection, key):
			if bytes.Equal(bytesClose, value) {
				h.SetConnectionClose()
			} else {
				h.ResetConnectionClose()
				h.setNonSpecial(key, value)
			}
			return true
		}
	case 's':
		if caseInsensitiveCompare(strServer, key) {
			h.SetServerBytes(value)
			return true
		} else if caseInsensitiveCompare(bytesSetCookie, key) {
			var kv *argsKV
			h.cookies, kv = allocArg(h.cookies)
			kv.key = getCookieKey(kv.key, value)
			kv.value = append(kv.value[:0], value...)
			return true
		}
	case 't':
		if caseInsensitiveCompare(bytesTransferEncoding, key) {
			// Transfer-Encoding is managed automatically.
			return true
		} else if caseInsensitiveCompare(bytesTrailer, key) {
			_ = h.SetTrailerBytes(value)
			return true
		}
	case 'd':
		if caseInsensitiveCompare(bytesDate, key) {
			// Date is managed automatically.
			return true
		}
	}

	return false
}

func normalizeHeaderKey(b []byte, disableNormalizing bool) {
	if disableNormalizing {
		return
	}

	n := len(b)
	if n == 0 {
		return
	}

	b[0] = toUpperTable[b[0]]
	for i := 1; i < n; i++ {
		p := &b[i]
		if *p == '-' {
			i++
			if i < n {
				b[i] = toUpperTable[b[i]]
			}
			continue
		}
		*p = toLowerTable[*p]
	}
}

func getHeaderKeyBytes(kv *argsKV, key string, disableNormalizing bool) []byte {
	kv.key = append(kv.key[:0], key...)
	normalizeHeaderKey(kv.key, disableNormalizing)
	return kv.key
}

func initHeaderKV(kv *argsKV, key, value string, disableNormalizing bool) {
	kv.key = getHeaderKeyBytes(kv, key, disableNormalizing)
	// https://tools.ietf.org/html/rfc7230#section-3.2.4
	kv.value = append(kv.value[:0], value...)
	kv.value = removeNewLines(kv.value)
}

func (h *ResponseHeader) SetCanonical(key, value []byte) {
	if h.setSpecialHeader(key, value) {
		return
	}
	h.setNonSpecial(key, value)
}

func (h *ResponseHeader) Set(key, value string) {
	initHeaderKV(&h.bufKV, key, value, h.disableNormalizing)
	h.SetCanonical(h.bufKV.key, h.bufKV.value)
}

func appendHeaderLine(dst *[]byte, key, value []byte) {
	*dst = append(*dst, key...)
	*dst = append(*dst, bytesColonSpace...)
	*dst = append(*dst, value...)
	*dst = append(*dst, bytesCRLF...)
}

func parseContentLength(b []byte) (int, error) {
	return strconv.Atoi(b2s(b))
}

func removeNewLines(raw []byte) []byte {
	// check if a `\r` is present and save the position.
	// if no `\r` is found, check if a `\n` is present.
	foundR := bytes.IndexByte(raw, rChar)
	foundN := bytes.IndexByte(raw, nChar)
	start := 0

	switch {
	case foundN != -1:
		if foundR > foundN {
			start = foundN
		} else if foundR != -1 {
			start = foundR
		}
	case foundR != -1:
		start = foundR
	default:
		return raw
	}

	for i := start; i < len(raw); i++ {
		switch raw[i] {
		case rChar, nChar:
			raw[i] = ' '
		default:
			continue
		}
	}
	return raw
}
