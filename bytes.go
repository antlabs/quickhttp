package quickhttp

var (
	bytesCRLF       = []byte("\r\n")
	bytesHTTP11     = []byte("HTTP/1.1")
	bytesColonSpace = []byte(": ")
	bytesServer     = []byte(HeaderServer)
)
