package zippo

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if ok, m := postPlease(w, r); !ok {
		JSON(w, m, http.StatusMethodNotAllowed)
		return
	}

	p := NewPayload(h.cf)

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(p)
	if err != nil {
		JSON(w, Response{Status: http.StatusBadRequest, Message: err.Error()}, http.StatusBadRequest)
		return
	}

	if p.HasCallbackURL() {
		go ProcessWithCallback(p)
		JSON(w, Response{Status: http.StatusAccepted, Message: "Request is being processed."}, http.StatusAccepted)
		return
	}

	l := log.WithFields(p.LogFields()).WithField("handler", "upload")

	u, err := Process(p)
	if err != nil {
		l.WithField("error", err.Error()).Warn("failure")
		internalError(w, err.Error())
		return
	}

	JSON(w, Response{Status: http.StatusOK, Message: "OK", URL: u}, http.StatusOK)
}
