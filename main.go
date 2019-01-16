package main

import (
	"flag"
	"teso_task/http"
)

func main() {
	port := flag.Uint("port", 3001, "Port on which http server is run")
	flag.Parse()

	s := http.New()
	s.Start(uint16(*port))
}
