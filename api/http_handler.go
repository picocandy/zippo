package zippo

import (
	"github.com/ncw/swift"
)

type Handler struct {
	cf swift.Connection
}

func NewHandler(cf swift.Connection) *Handler {
	return &Handler{cf: cf}
}
