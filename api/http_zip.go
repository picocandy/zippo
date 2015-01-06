package zippo

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/ncw/swift"
	"net/http"
)

func ZipHandler(w http.ResponseWriter, r *http.Request, cf swift.Connection) {
	if ok, m := postPlease(w, r); !ok {
		JSON(w, m, http.StatusMethodNotAllowed)
		return
	}

	a := &Archive{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(a)
	if err != nil {
		JSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	a.RenameDuplicatePayloads()

	err = cf.Authenticate()
	if err != nil {
		internalError(w, "Unable to authenticate to Rackspace Cloud Files")
		return
	}

	l := log.WithFields(logrus.Fields{
		"handler":      "zip",
		"hash":         a.Hash(),
		"filename":     a.String(),
		"content_type": "application/zip",
		"expiration":   a.ExpirationSec(),
	})

	u, err := a.DownloadURL(cf)
	if err == nil {
		l.WithField("container", container).Info("existing secure zip url")
		JSON(w, map[string]string{"message": "OK", "url": u}, http.StatusOK)
		return
	}

	err = a.Build()
	if err != nil {
		l.WithField("error", err.Error()).Warn("build zip error")
		internalError(w, fmt.Sprintf("Unable to build zip file: '%s'", err.Error()))
		return
	}

	_, _, err = a.Upload(cf, container)
	if err != nil {
		l.WithFields(logrus.Fields{"tmp": a.TempFile, "error": err.Error()}).Warn("upload zip error")
		internalError(w, "Unable to upload zip file to Rackspace Cloud Files")
		return
	}

	u, err = a.DownloadURL(cf)
	if err != nil {
		l.WithField("error", err.Error()).Warn("generating secure zip url failed")
		internalError(w, fmt.Sprintf("Unable to get download url for %s", a.String()))
		return
	}

	l.WithFields(logrus.Fields{"container": container, "tmp": a.TempFile}).Info("secure zip url generated")
	JSON(w, map[string]string{"message": "OK", "url": u}, http.StatusOK)
}
