package utils

import "net/http"

type ResponseWriter struct {
	Headers    map[string]string
	StatusCode int
	Body       string
}

func (rw *ResponseWriter) Header() http.Header {
	return http.Header{}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
}

func (rw *ResponseWriter) Write(body []byte) (int, error) {
	rw.Body = string(body)
	return len(body), nil
}
