package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	timeOutSec int
}

const (
	MaxHeaderBytes = 1 << 20 // 1 MB
	ReadTimeout    = 10 * time.Second
	WriteTimeout   = 10 * time.Second
)

// function called to create service and configure it
func NewServer(host, port string, handler http.Handler, timeOutSec int) *Server {
	server := http.Server{
		Addr:           host + ":" + port,
		Handler:        handler,
		MaxHeaderBytes: MaxHeaderBytes,
		ReadTimeout:    ReadTimeout,
		WriteTimeout:   WriteTimeout,
	}

	return &Server{
		httpServer: &server,
		timeOutSec: timeOutSec,
	}
}

// function called to run service
func (s *Server) Run() error {
	err := s.httpServer.ListenAndServe()

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(s.timeOutSec)*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctxTimeout)
}
