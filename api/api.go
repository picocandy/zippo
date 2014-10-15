package zippo

import "os"

var container = os.Getenv("SWIFT_CONTAINER")
var metaTempKey = os.Getenv("SWIFT_META_TEMP")
