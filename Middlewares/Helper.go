package middlewares

import "net/http"

type responseWriterWithCode struct {
	http.ResponseWriter
	StatusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriterWithCode {
	return &responseWriterWithCode{w, http.StatusOK}
}

func (obj *responseWriterWithCode) WriteHeader(code int) {
	obj.StatusCode = code
	obj.ResponseWriter.WriteHeader(code)
}
