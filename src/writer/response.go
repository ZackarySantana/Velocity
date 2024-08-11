package writer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	http.ResponseWriter
	bodyWritten   bool
	headerWritten bool

	writers []io.Writer

	statusCode int
}

func (w *Response) WrittenTo() bool {
	return w.bodyWritten || w.headerWritten
}

func (w *Response) StatusCode() int {
	return w.statusCode
}

func (w *Response) Write(data []byte) (int, error) {
	w.bodyWritten = true
	for _, writer := range w.writers {
		_, err := writer.Write([]byte(fmt.Sprintf("body: %v\n", data)))
		if err != nil {
			return -1, fmt.Errorf("writing to additional writers '%v'", err)
		}
	}
	return w.ResponseWriter.Write(data)
}

func (w *Response) WriteHeader(statusCode int) {
	w.headerWritten = true
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
	for _, writer := range w.writers {
		writer.Write([]byte(fmt.Sprintf("status code: %d\n", statusCode)))
	}
}

func (w *Response) WriteJSON(data any) (int, error) {
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 0, nil
	}
	w.Header().Set("Content-Type", "application/json")
	return w.Write(b)
}

type errorResponse struct {
	Error string `json:"error"`
}

func (w *Response) WriteError(err error, code int) {
	if err == nil {
		w.WriteHeader(code)
		return
	}

	b, err := json.Marshal(errorResponse{err.Error()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func (w *Response) AddOutput(ioW io.Writer) {
	w.writers = append(w.writers, ioW)
}
