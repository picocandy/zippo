package zippo

import (
	"github.com/ncw/swift"
	"os"
)

var container = os.Getenv("SWIFT_CONTAINER")
var metaTempKey = os.Getenv("SWIFT_META_TEMP")

type Parker interface {
	String() string
	DownloadURL(cf swift.Connection) (string, error)
}
