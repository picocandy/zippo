package zippo

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"net/http"
)

func LogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr := r.Header.Get("X-Real-IP")
		if addr == "" {
			addr = r.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = r.RemoteAddr
			}
		}

		log.WithFields(logrus.Fields{"method": r.Method, "path": r.URL.Path, "remote": addr}).Info("request")
		handler.ServeHTTP(w, r)
	})
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "zippo!")
}

func postPlease(w http.ResponseWriter, r *http.Request) (bool, map[string]string) {
	if r.Method != "POST" {
		return false, map[string]string{"error": fmt.Sprintf("Method %s is not allowed. Please use 'POST' instead.", r.Method)}
	}

	return true, map[string]string{}
}

func internalError(w http.ResponseWriter, msg string) {
	JSON(w, map[string]string{"error": msg}, http.StatusInternalServerError)
}
