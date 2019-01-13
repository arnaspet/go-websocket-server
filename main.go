package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 3001, "Port on which http server is run")
	flag.Parse()

	logger := log.New(os.Stdout, "teso_task ", log.LstdFlags | log.Lshortfile)

	logger.Printf("Starting server on %v\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)

	if err != nil {
		logger.Fatalf("Cannot run server: %v", err)
	}
}
