package zippo

import (
	"fmt"
	"net/http"
)

type Response struct {
	Zipname     string    `json:"zipname"`
	Payloads    []Payload `json:"payloads"`
	Length      int       `json:"length"`
	ContentType string    `json:"content_type"`
}

func ZipHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Zippo!")
}

func ServeMux() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/", ZipHandler)
	return m
}
