package zippo

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/ncw/swift"
	"net/http"
	"os"
)

var (
	container   = os.Getenv("SWIFT_CONTAINER")
	metaTempKey = os.Getenv("SWIFT_META_TEMP")
	log         = logrus.New()
)

type Parker interface {
	String() string
	DownloadURL() (string, error)
	ExpirationSec() int64
}

type Transformer interface {
	Authenticate() error
	Container() string
	String() string
	Build() error
	Upload() (ob swift.Object, h swift.Headers, err error)
	DownloadURL() (string, error)
	HasCallbackURL() bool
	CallCallbackURL(data Response) error
}

func Process(t Transformer) (string, error) {
	err := t.Authenticate()
	if err != nil {
		return "", errors.New("Unable to authenticate to Rackspace Cloud Files")
	}

	u, err := t.DownloadURL()
	if err == nil {
		return u, nil
	}

	err = t.Build()
	if err != nil {
		return "", fmt.Errorf("Unable to build the file: %q", err.Error())
	}

	_, _, err = t.Upload()
	if err != nil {
		return "", fmt.Errorf("Unable to upload file to Rackspace Cloud Files: %q", err.Error())
	}

	u, err = t.DownloadURL()
	if err != nil {
		return "", fmt.Errorf("Unable to get download url for %s: %q", t.String(), err.Error())
	}

	return u, nil
}

func ProcessWithCallback(t Transformer) error {
	var d Response

	u, err := Process(t)
	if err != nil {
		d = Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	} else {
		d = Response{
			Status:  http.StatusOK,
			Message: "OK",
			URL:     u,
		}
	}

	return t.CallCallbackURL(d)
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	URL     string `json:"url,omitempty"`
}
