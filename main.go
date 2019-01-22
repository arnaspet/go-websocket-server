package main

import (
	"flag"
	"log"
	"os"
	"teso_task/http"
)

func main() {
	port := flag.Uint("port", 3001, "Port on which http server is run")
	flag.Parse()
	logger := log.New(os.Stdout, "http_server ", log.LstdFlags | log.Lshortfile)

	s := http.NewServer(logger)
	s.Start(uint16(*port))
}
