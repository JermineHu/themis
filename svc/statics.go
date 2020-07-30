package svc

import (
	"log"

	statics "github.com/JermineHu/themis/svc/gen/statics"
)

// statics service example implementation.
// The example methods log the requests and return zero values.
type staticssrvc struct {
	logger *log.Logger
}

// NewStatics returns the statics service implementation.
func NewStatics(logger *log.Logger) statics.Service {
	return &staticssrvc{logger}
}
