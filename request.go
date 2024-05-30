package quickhttp

type Request struct {
	Header    RequestHeader
	bodyStart int
	bodyEnd   int
}
