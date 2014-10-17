package main

import (
	"flag"
	"fmt"
	"github.com/ncw/swift"
	zippo "github.com/picocandy/zippo/api"
	"log"
	"net/http"
)

var cf swift.Connection

func init() {
	cf = zippo.NewConnection()
}

func main() {
	address := flag.String("address", "0.0.0.0", "Bind to HOST address (default: 0.0.0.0)")
	port := flag.String("port", "5020", "Use PORT (default: 5020)")

	flag.Parse()

	err := zippo.UpdateAccountMetaTempURL(cf)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", zippo.HomeHandler)
	mux.HandleFunc("/z", func(w http.ResponseWriter, r *http.Request) {
		zippo.ZipHandler(w, r, cf)
	})
	mux.HandleFunc("/u", func(w http.ResponseWriter, r *http.Request) {
		zippo.UploadHandler(w, r, cf)
	})

	bind := fmt.Sprintf("%s:%s", *address, *port)

	if err := http.ListenAndServe(bind, zippo.LogHandler(mux)); err != nil {
		log.Fatal(err)
	}
}
