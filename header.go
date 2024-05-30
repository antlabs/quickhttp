package quickhttp

import "net/http"

type RequestHeader struct {
	method []byte
}

type ResponseHeader struct {
	statusCode    int
	statusMessage string
	protocol      string
	Request       Request
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

func (h *ResponseHeader) Protocol() string {
	if len(h.protocol) > 0 {
		return h.protocol
	}
	return strHTTP11
}

func (h *ResponseHeader) appendStatusLine(dst []byte) []byte {
	statusCode := h.StatusCode()
	if statusCode < 0 {
		statusCode = http.StatusOK
	}
	return formatStatusLine(dst, h.Protocol(), statusCode, h.StatusMessage())
}
