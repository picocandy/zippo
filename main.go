package main

import (
	"flag"
	"fmt"
	"github.com/ncw/swift"
	zippo "github.com/picocandy/zippo/api"
	"log"
	"net/http"
	"os"
)

var cf swift.Connection

func init() {
	cf = swift.Connection{
		UserName: os.Getenv("SWIFT_API_USER"),
		ApiKey:   os.Getenv("SWIFT_API_KEY"),
		AuthUrl:  os.Getenv("SWIFT_AUTH_URL"),
		Region:   os.Getenv("SWIFT_REGION"),
		TenantId: os.Getenv("SWIFT_TENANT_ID"),
	}
}

func main() {
	address := flag.String("address", "0.0.0.0", "Bind to HOST address (default: 0.0.0.0)")
	port := flag.String("port", "5020", "Use PORT (default: 5020)")

	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", zippo.HomeHandler)
	mux.HandleFunc("/z", func(w http.ResponseWriter, r *http.Request) {
		zippo.ZipHandler(w, r, cf)
	})

	bind := fmt.Sprintf("%s:%s", *address, *port)

	if err := http.ListenAndServe(bind, mux); err != nil {
		log.Fatal(err)
	}
}
