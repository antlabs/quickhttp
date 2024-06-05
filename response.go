package quickhttp

import (
	"io"
)

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

func (r *Response) reset() {
	if r.headerRaw != nil {
		PutBytes(r.headerRaw)
		r.headerRaw = nil
	}

	if r.backHeaderRaw != nil {
		r.backHeaderRaw = nil
	}
	if r.bodyRaw != nil {
		PutBytes(r.bodyRaw)
		r.bodyRaw = nil
	}
}

func (r *Response) init() {
	if r.headerRaw == nil {

		r.headerRaw = GetAndResetBytes(1024)
	}
	r.backHeaderRaw = r.headerRaw
}

func (r *Response) Reset() {
	r.reset()
}

func (r *Response) Write(w io.Writer) error {
	contentLength := 0
	if r.bodyRaw != nil && len(*r.bodyRaw) > 0 {
		contentLength = len(*r.bodyRaw)
	}
	r.Header.SetContentLength(contentLength)

	r.Header.AppendBytes(r.headerRaw)

	if r.bodyRaw == nil || len(*r.bodyRaw) == 0 {
		// 只有header数据
		w.Write(*r.headerRaw)
	} else {
		// 有header数据和body数据
		mergeBody := GetAndResetBytes(len(*r.bodyRaw) + len(*r.headerRaw) + 2)
		*mergeBody = append(*mergeBody, *r.headerRaw...)
		*mergeBody = append(*mergeBody, *r.bodyRaw...)

		w.Write(*mergeBody)
		PutBytes(mergeBody)
	}

	return nil
}
