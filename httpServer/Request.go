package httpServer

type Request struct {
	headers map[string]string
	body string
}

func (req *Request) Get (h string) string {
	return req.headers[h]
}