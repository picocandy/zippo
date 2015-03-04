package zippo

import (
	"github.com/Sirupsen/logrus"
	"github.com/ncw/swift"
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

func Log() *logrus.Logger {
	return log
}

type Expiration struct {
	Duration int64 `json:"expiration,omitempty"`
}

func (e *Expiration) ExpirationSec() int64 {
	if e.Duration == 0 {
		return 600
	}

	return e.Duration
}

type Callback struct {
	CallbackURL string `json:"callback_url,omitempty"`
}

func (c *Callback) HasCallbackURL() bool {
	return c.CallbackURL != ""
}

func (c *Callback) CallCallbackURL(data interface{}) error {
	return PostJSON(c.CallbackURL, data)
}

type CloudFile struct {
	cf swift.Connection
}

func (c *CloudFile) SetConnection(conn swift.Connection) {
	c.cf = conn
}

func (c *CloudFile) Authenticate() error {
	return c.cf.Authenticate()
}
