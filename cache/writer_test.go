package cache

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type mockWriter response

func newMockWriter() *mockWriter {
	return &mockWriter{
		body:   []byte{},
		header: http.Header{},
	}
}

func (mw *mockWriter) Write(b []byte) (int, error) {
	mw.body = make([]byte, len(b))
	copy(mw.body, b)
	return len(b), nil
}

func (mw *mockWriter) WriteHeader(code int) { mw.code = code }
func (mw *mockWriter) Header() http.Header  { return mw.header }

func TestWriter(t *testing.T) {
	mw := newMockWriter()

	res := "/test/url?with=params"
	u, err := url.Parse(res)
	if err != nil {
		t.Fatal("Invalid URL")
	}
	req := &http.Request{
		URL: u,
	}

	t.Log("Testing NewWriter")
	w := NewWriter(mw, req)
	if w.resource != res {
		t.Errorf("Expected resource %s, got %s", res, w.resource)
	}
	if w.writer != mw {
		t.Errorf("Expected writer %+v, got %+v", mw, w.writer)
	}

	t.Log("Testing Header")
	h := w.Header()
	h.Add("Content-Type", "application/json")
	h2 := w.response.header
	if h2.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be application/json, got %s", h.Get("Content-Type"))
	}

	t.Log("Testing WriteHeader")
	c := 201
	w.WriteHeader(c)
	if w.response.code != c {
		t.Errorf("Expected code %d, got %d", c, w.response.code)
	}
	if mw.code != c {
		t.Errorf("Expected code %d, got %d", c, mw.code)
	}
	h2 = w.response.header
	if h2.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be application/json, got %s", h.Get("Content-Type"))
	}

	t.Log("Testing Write")
	bd := []byte{1, 2, 3, 4, 5}
	n, err := w.Write(bd)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if n != len(bd) {
		t.Errorf("Expected %d bytes written, got %d", len(bd), n)
	}
	if &w.response.body == &bd {
		t.Errorf("Expected body to be a copy, got %p", &w.response.body)
	}
	if !reflect.DeepEqual(w.response.body, bd) {
		t.Errorf("Expected body %v, got %v", bd, w.response.body)
	}

	if !reflect.DeepEqual(mw.body, bd) {
		t.Errorf("Expected body %v, got %v", bd, mw.body)
	}
}
