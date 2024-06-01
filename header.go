package quickhttp

import (
	"net/http"
	"strconv"
)

const (
	HeaderServer = "Server"
)

type RequestHeader struct {
	method []byte
}

type ResponseHeader struct {
	statusCode         int
	statusMessage      string
	server             []byte
	protocol           string
	contentLength      int
	contentLengthBytes []byte
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
	h.contentLengthBytes = strconv.AppendInt(h.contentLengthBytes, int64(contentLength), 10)
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

	*dst = append(*dst, bytesCRLF...)
}

func appendHeaderLine(dst *[]byte, key, value []byte) {
	*dst = append(*dst, key...)
	*dst = append(*dst, bytesColonSpace...)
	*dst = append(*dst, value...)
	*dst = append(*dst, bytesCRLF...)
}
