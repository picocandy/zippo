package main

import (
	"flag"
	"fmt"
	"github.com/ncw/swift"
	zippo "github.com/picocandy/zippo/api"
	"log"
	"net/http"
	"runtime"
)

var cf swift.Connection

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
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

	h := zippo.NewHandler(cf)
	mux := http.NewServeMux()
	mux.HandleFunc("/", zippo.HomeHandler)
	mux.HandleFunc("/version", zippo.VersionHandler)
	mux.HandleFunc("/z", h.ZipUpload)
	mux.HandleFunc("/u", h.Upload)

	bind := fmt.Sprintf("%s:%s", *address, *port)

	if err := http.ListenAndServe(bind, zippo.LogHandler(mux)); err != nil {
		log.Fatal(err)
	}
}
