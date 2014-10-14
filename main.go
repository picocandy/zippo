package main

import (
	"flag"
	"fmt"
	zippo "github.com/picocandy/zippo/api"
	"log"
	"net/http"
)

func main() {
	address := flag.String("address", "0.0.0.0", "Bind to HOST address (default: 0.0.0.0)")
	port := flag.String("port", "5020", "Use PORT (default: 5020)")

	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", zippo.HomeHandler)
	mux.HandleFunc("/z", zippo.ZipHandler)

	bind := fmt.Sprintf("%s:%s", *address, *port)

	if err := http.ListenAndServe(bind, mux); err != nil {
		log.Fatal(err)
	}
}
