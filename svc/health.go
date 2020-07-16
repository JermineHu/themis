package themis

import (
	"context"
	"log"

	health "github.com/JermineHu/themis/svc/gen/health"
)

// health service example implementation.
// The example methods log the requests and return zero values.
type healthsrvc struct {
	logger *log.Logger
}

// NewHealth returns the health service implementation.
func NewHealth(logger *log.Logger) health.Service {
	return &healthsrvc{logger}
}

// Perform health check.
func (s *healthsrvc) Health(ctx context.Context) (err error) {
	s.logger.Print("health.health")
	return
}

// version info.
func (s *healthsrvc) BuildInfo(ctx context.Context) (res *health.BuildInfoResult, err error) {
	res = &health.BuildInfoResult{}
	s.logger.Print("health.build_info")
	return
}
