package main

import (
	"flag"
	"github.com/ideade/timetable-pipe/cache"
	"github.com/ideade/timetable-pipe/server"
	"github.com/ideade/timetable-pipe/timetable"
	"io"
	"log"
	"net/http"
	"os"
)

func getLogger() *log.Logger {
	var writer io.Writer = os.Stdout
	var logger *log.Logger

	logger = log.New(writer, "[Main] ", log.Ldate|log.Ltime)

	return logger
}

func main() {
	var address string
	flag.StringVar(&address, "address", "127.0.0.1:9000", "Server address")
	flag.Parse()

	logger := getLogger()

	logger.Println("Initialization of the timetable parser")
	ttHandler := timetable.NewHandler()
	ttHandler.ReqUrl()

	logger.Println("Initialization of the cache store")
	tCache := new(cache.Store)

	gHandler := new(GroupsHandler)
	tHandler := new(TimetableHandler)

	gHandler.cache = tCache
	gHandler.ttHandler = *ttHandler

	tHandler.cache = tCache
	tHandler.ttHandler = *ttHandler

	logger.Println("Initialization of the fastcgi server")
	srv := new(server.Server)
	srv.Address = address

	logger.Println("Registering handlers")
	handlers := map[string]http.HandlerFunc{
		"/api/groups":    gHandler.ServeHTTP,
		"/api/timetable": tHandler.ServeHTTP,
	}

	logger.Println("Starting the fcgi-server")
	errCh := make(chan error, 64)
	go srv.Run(errCh, handlers)
	for {
		srvErr := <-errCh
		if srvErr != nil {
			log.Fatal(srvErr)
		}
	}
}
