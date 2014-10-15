package zippo

import (
	"encoding/json"
	"fmt"
	"github.com/ncw/swift"
	"net/http"
	"os"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "Zippo!")
}

func ZipHandler(w http.ResponseWriter, r *http.Request, cf swift.Connection) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	a := &Archive{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = cf.Authenticate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u, err := a.DownloadURL(cf)
	if err == nil {
		fmt.Fprint(w, u)
		return
	}

	err = a.Build()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ob, err := a.Upload(cf, os.Getenv("SWIFT_CONTAINER"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, ob)
}
