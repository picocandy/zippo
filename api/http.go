package zippo

import (
	"encoding/json"
	"fmt"
	"github.com/ncw/swift"
	"net/http"
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
		JSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	err = cf.Authenticate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u, err := a.DownloadURL(cf)
	if err == nil {
		JSON(w, map[string]string{"message": "OK", "url": u}, http.StatusOK)
		return
	}

	err = a.Build()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = a.Upload(cf, container)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u, err = a.DownloadURL(cf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, map[string]string{"message": "OK", "url": u}, http.StatusOK)
}
