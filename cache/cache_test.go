package cache

import (
	"net/http"
	"reflect"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	resource := "/example"
	resp := &response{
		header: http.Header{"Content-Type": []string{"application/json"}},
		code:   http.StatusOK,
		body:   []byte(`{"message":"Hello, world!"}`),
	}

	set(resource, resp)

	got := get(resource)
	if !reflect.DeepEqual(resp, got) {
		t.Errorf("Expected response %+v, got %+v", resp, got)
	}

	// Test deleting entry
	set(resource, nil)

	got = get(resource)
	if got != nil {
		t.Errorf("Expected nil response, got %+v", got)
	}
}

func TestCopyHeader(t *testing.T) {
	dst := make(http.Header)
	src := make(http.Header)
	src.Add("Content-Type", "application/json")
	src.Add("Cache-Control", "no-cache")

	copyHeader(dst, src)

	expected := http.Header{
		"Content-Type":  []string{"application/json"},
		"Cache-Control": []string{"no-cache"},
	}
	if !reflect.DeepEqual(dst, expected) {
		t.Errorf("Expected header %+v, got %+v", expected, dst)
	}
}

func TestMakeResource(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com/path?param=value", nil)
	expected := "/path"

	got := MakeResource(req)

	if got != expected {
		t.Errorf("Expected resource %s, got %s", expected, got)
	}

	got = MakeResource(nil)
	if got != "" {
		t.Errorf("Expected empty resource, got %s", got)
	}
}

func TestClean(t *testing.T) {
	// Add some entries to the cache
	set("/resource1", &response{})
	set("/resource2", &response{})
	set("/resource3", &response{})

	Clean()

	if len(cache.data) != 0 {
		t.Errorf("Expected cache to be empty, got %d entries", len(cache.data))
	}
}

func TestDrop(t *testing.T) {
	resource := "/resource"
	resp := &response{}

	set(resource, resp)

	Drop(resource)

	got := get(resource)
	if got != nil {
		t.Errorf("Expected nil response, got %+v", got)
	}
}

func TestServe(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com/resource", nil)
	w := &dummyResponseWriter{}
	resp := &response{
		header: http.Header{"Content-Type": []string{"application/json"}},
		code:   http.StatusOK,
		body:   []byte(`{"message":"Hello, world!"}`),
	}
	set("/resource", resp)

	found := Serve(w, req)

	if !found {
		t.Error("Expected response to be served from cache, but not found")
	}
	if !reflect.DeepEqual(w.header, resp.header) {
		t.Errorf("Expected response header %+v, got %+v", resp.header, w.header)
	}
	if w.statusCode != resp.code {
		t.Errorf("Expected response status code %d, got %d", resp.code, w.statusCode)
	}
	if !reflect.DeepEqual(w.body, resp.body) {
		t.Errorf("Expected response body %s, got %s", resp.body, w.body)
	}

	// Test no-cache header
	req.Header.Set("Cache-Control", "no-cache")
	found = Serve(w, req)

	if found {
		t.Error("Expected response not to be served from cache, but found")
	}
}

// dummyResponseWriter is a helper struct implementing http.ResponseWriter for testing
type dummyResponseWriter struct {
	header     http.Header
	statusCode int
	body       []byte
}

func (w *dummyResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = http.Header{}
	}
	return w.header
}

func (w *dummyResponseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return len(b), nil
}

func (w *dummyResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}
