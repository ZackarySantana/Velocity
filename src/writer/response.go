package writer

import (
	"net/http"
)

type Response struct {
	http.ResponseWriter
	statusCode int
}

func (w *Response) StatusCode() int {
	return w.statusCode
}

func (w *Response) Write(data []byte) (int, error) {
	w.statusCode = http.StatusOK
	return w.ResponseWriter.Write(data)
}

func (w *Response) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
