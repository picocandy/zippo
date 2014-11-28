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
	DownloadURL(cf swift.Connection) (string, error)
	ExpirationSec() int64
}

func Log() *logrus.Logger {
	return log
}
