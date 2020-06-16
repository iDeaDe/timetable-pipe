package server

import (
	"io"
	"log"
	"net/http"
	"os"
)

func getLogger() *log.Logger {
	var writer io.Writer = os.Stdout
	var logger *log.Logger

	logger = log.New(writer, "[FCGI] ", log.LstdFlags)

	return logger
}

type FCGIHandler struct {
	Handlers map[string]http.HandlerFunc
}

func (fh *FCGIHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	getLogger().Printf("%s from %s", request.URL.Path, request.RemoteAddr)
	reqPath := request.URL.Path

	if fh.Handlers[reqPath] != nil {
		fh.Handlers[reqPath](writer, request)
	} else {
		writer.WriteHeader(404)
	}
}
