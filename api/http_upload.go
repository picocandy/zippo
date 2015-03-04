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
		JSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if p.HasCallbackURL() {
		go p.CallCallbackURL(map[string]string{})
		JSON(w, map[string]string{"message": "Request is being processed."}, http.StatusOK)
		return
	}

	l := log.WithFields(p.LogFields()).WithField("handler", "upload")

	u, err := Process(p)
	if err != nil {
		l.WithField("error", err.Error()).Warn("failure")
		internalError(w, err.Error())
	}

	JSON(w, map[string]string{"message": "OK", "url": u}, http.StatusOK)
}
