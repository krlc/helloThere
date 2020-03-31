package helloThere

import (
	"net/http"
	"strings"
	"testing"
)

func TestNewServer(t *testing.T) {
	const addr = "test"

	if server := NewServer("", nil); server == nil {
		t.Fatal("NewServer() got nil instance")
	}

	if server := NewServer(addr, nil); server.srv.Addr != addr {
		t.Errorf("NewServer() got %v, want %v", server.srv.Addr, addr)
	}
}

func TestServer_Close(t *testing.T) {
	if err := NewServer("", nil).Close(); err != nil {
		t.Error(err)
	}
}

type ResponseWriterMock struct {
	header http.Header
}

func (w *ResponseWriterMock) Header() http.Header {
	return w.header
}

func (w *ResponseWriterMock) Write(_ []byte) (int, error) {
	return 0, nil
}

func (w *ResponseWriterMock) WriteHeader(_ int) {}

func TestSetCookie(t *testing.T) {
	const cname = "testCookie"

	w := &ResponseWriterMock{header: http.Header{}}
	SetCookie(w, cname)

	if c := w.Header().Get("Set-Cookie"); !strings.Contains(c, cname) {
		t.Errorf("SetCookie() got %v, want %v", c, cname)
	}
}
