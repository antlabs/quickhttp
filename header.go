package quickhttp

import "net/http"

const (
	HeaderServer = "Server"
)

type RequestHeader struct {
	method []byte
}

type ResponseHeader struct {
	statusCode    int
	statusMessage string
	server        []byte
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

func (h *ResponseHeader) Server() []byte {
	return h.server
}

func (h *ResponseHeader) appendStatusLine(dst *[]byte) {
	statusCode := h.StatusCode()
	if statusCode < 0 {
		statusCode = http.StatusOK
	}
	formatStatusLine(dst, h.Protocol(), statusCode, h.StatusMessage())
}

func (h *ResponseHeader) AppendBytes(dst *[]byte) {
	h.appendStatusLine(dst)

	server := h.Server()
	if len(server) != 0 {
		appendHeaderLine(dst, bytesServer, server)
	}

	*dst = append(*dst, bytesCRLF...)
}

func appendHeaderLine(dst *[]byte, key, value []byte) {
	*dst = append(*dst, key...)
	*dst = append(*dst, bytesColonSpace...)
	*dst = append(*dst, value...)
	*dst = append(*dst, bytesCRLF...)
}
