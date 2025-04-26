package response

import (
	"net/http"
)

type ErrResponseWriter struct {
	http.ResponseWriter
	Err error
}

func (ew *ErrResponseWriter) Header() http.Header {
	return ew.ResponseWriter.Header()
}

func (ew *ErrResponseWriter) Write(b []byte) (int, error) {
	return ew.ResponseWriter.Write(b)
}

func (ew *ErrResponseWriter) WriteHeader(statusCode int) {
	ew.ResponseWriter.WriteHeader(statusCode)
}
