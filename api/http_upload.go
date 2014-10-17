package zippo

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/ncw/swift"
	"net/http"
)

func UploadHandler(w http.ResponseWriter, r *http.Request, cf swift.Connection) {
	if ok, m := postPlease(w, r); !ok {
		JSON(w, m, http.StatusMethodNotAllowed)
		return
	}

	p := &Payload{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(p)
	if err != nil {
		JSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	err = cf.Authenticate()
	if err != nil {
		internalError(w, "Unable to authenticate to Rackspace Cloud Files")
		return
	}

	l := log.WithFields(logrus.Fields{
		"handler":        "upload",
		"hash":           p.Hash(),
		"filename":       p.String(),
		"url":            p.URL,
		"content_type":   p.ContentType,
		"content_length": p.ContentLength,
	})

	u, err := p.DownloadURL(cf)
	if err == nil {
		l.WithField("container", container).Info("existing secure url")
		JSON(w, map[string]string{"message": "OK", "url": u}, http.StatusOK)
		return
	}

	err = p.Download()
	if err != nil {
		l.WithField("error", err.Error()).Warn("download error")
		internalError(w, fmt.Sprintf("Unable to download the file: '%s'", err.Error()))
		return
	}

	_, _, err = p.Upload(cf, container)
	if err != nil {
		l.WithFields(logrus.Fields{"tmp": p.TempFile, "error": err.Error()}).Warn("upload error")
		internalError(w, "Unable to upload file to Rackspace Cloud Files")
		return
	}

	u, err = p.DownloadURL(cf)
	if err != nil {
		l.WithField("error", err.Error()).Warn("generating secure url failed")
		internalError(w, fmt.Sprintf("Unable to get download url for %s", p.String()))
		return
	}

	l.WithFields(logrus.Fields{"container": container, "tmp": p.TempFile}).Info("secure url generated")
	JSON(w, map[string]string{"message": "OK", "url": u}, http.StatusOK)
}
