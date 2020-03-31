package helloThere

import (
	"context"
	"net/http"
	"time"
)

const (
	srvClosedTimeout = 10 * time.Second
	readTimeout      = 5 * time.Second
	writeTimeout     = 5 * time.Second

	defaultCookieLifetime = 24 * time.Hour
)

type Server struct {
	srv *http.Server
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), srvClosedTimeout)
	defer cancel()

	s.srv.SetKeepAlivesEnabled(false)
	return s.srv.Shutdown(ctx)
}

func NewServer(endpoint string, handler *http.HandlerFunc) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         endpoint,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			Handler:      handler,
		},
	}
}

func SetCookie(w http.ResponseWriter, cookieName string) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Expires: time.Now().Add(defaultCookieLifetime),
	})
}
