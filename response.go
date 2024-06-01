package quickhttp

import "io"

type Response struct {
	Header        ResponseHeader
	headerRaw     *[]byte
	backHeaderRaw *[]byte
	bodyRaw       *[]byte
}

func (r *Response) AppendBody(p []byte) {
	if r.bodyRaw == nil {
		r.bodyRaw = GetAndResetBytes(len(p))
	}

	newBodySize := len(*r.bodyRaw) + len(p)
	if cap(*r.bodyRaw) < newBodySize {
		// TODO 优化
		newBodyRaw := GetAndResetBytes(newBodySize)
		*newBodyRaw = append(*newBodyRaw, *r.bodyRaw...)
		PutBytes(r.bodyRaw)
		r.bodyRaw = newBodyRaw
	}

	*r.bodyRaw = append(*r.bodyRaw, p...)
}

func (r *Response) init() {
	r.headerRaw = GetAndResetBytes(1024)
	r.backHeaderRaw = r.backHeaderRaw
}

func (r *Response) free() {
	PutBytes(r.headerRaw)

}
func (r *Response) Write(w io.Writer) error {
	contentLength := 0
	if r.bodyRaw != nil && len(*r.bodyRaw) > 0 {
		contentLength = len(*r.bodyRaw)
	}
	r.Header.SetContentLength(contentLength)

	r.Header.AppendBytes(r.bodyRaw)

	if r.bodyRaw == nil || len(*r.bodyRaw) == 0 {
		w.Write(*r.headerRaw)
	} else {
		mergeBody := GetAndResetBytes(len(*r.bodyRaw) + len(*r.headerRaw) + 2)
		copy(*mergeBody, *r.headerRaw)
		copy(*mergeBody, *r.bodyRaw)
		*mergeBody = append(*mergeBody, bytesCRLF...)

		w.Write(*mergeBody)
		PutBytes(mergeBody)
	}

	return nil
}
