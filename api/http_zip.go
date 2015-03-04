package zippo

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) ZipUpload(w http.ResponseWriter, r *http.Request) {
	if ok, m := postPlease(w, r); !ok {
		JSON(w, m, http.StatusMethodNotAllowed)
		return
	}

	a := NewArchive(h.cf)

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(a)
	if err != nil {
		JSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	a.RenameDuplicatePayloads()

	if a.HasCallbackURL() {
		go a.CallCallbackURL(map[string]string{})
		JSON(w, map[string]string{"message": "Request is being processed."}, http.StatusOK)
		return
	}

	l := log.WithFields(a.LogFields()).WithField("handler", "zip")

	u, err := Process(a)
	if err != nil {
		l.WithField("error", err.Error()).Warn("failure")
		internalError(w, err.Error())
	}

	JSON(w, map[string]string{"message": "OK", "url": u}, http.StatusOK)
}
