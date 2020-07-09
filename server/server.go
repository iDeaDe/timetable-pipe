package server

import (
	"net"
	"net/http"
)

type Server struct {
	Address string
	Handler http.HandlerFunc
}

func (s *Server) Run(errCh chan error, handlers map[string]http.HandlerFunc) {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		errCh <- err
		return
	}

	handler := new(FCGIHandler)
	handler.Handlers = handlers

	getLogger().Printf("Starting on %s\n", s.Address)
	if err = http.Serve(listener, handler); err != nil {
		errCh <- err
		return
	}
}
