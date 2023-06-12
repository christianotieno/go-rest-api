package cache

import "net/http"

type Writer struct {
	writer   http.ResponseWriter
	response response
	resource string
}

// interface implementation check
var (
	_ http.ResponseWriter = (*Writer)(nil)
)

// NewWriter returns a cache writer
func NewWriter(w http.ResponseWriter, r *http.Request) *Writer {
	return &Writer{
		writer:   w,
		resource: MakeResource(r),
		response: response{
			header: http.Header{},
		},
	}
}

// Header returns the header map that will be sent by WriteHeader.
func (w *Writer) Header() http.Header {
	return w.response.header
}

// WriteHeader writes the data to the connection as part of an HTTP reply.
func (w *Writer) WriteHeader(code int) {
	copyHeader(w.response.header, w.writer.Header())
	w.response.code = code
	w.writer.WriteHeader(code)
}

// Write writes the data to the connection as part of an HTTP reply.
func (w *Writer) Write(b []byte) (int, error) {
	w.response.body = make([]byte, len(b))
	for k, v := range b {
		w.response.body[k] = v
	}
	copyHeader(w.Header(), w.writer.Header())
	set(w.resource, &w.response)
	return w.writer.Write(b)
}
