package svc

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"

	health "github.com/JermineHu/themis/svc/gen/health"
)

var (
	BuildTime = ""
	GitHash   = ""
	GitLog    = ""
	GitStatus = ""
	GoVersion = ""
	GoRuntime = ""
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
	if GitStatus == "" {
		// GitStatus 为空时，说明本地源码与最近的 commit 记录一致，无修改
		// 为它赋一个特殊值
		GitStatus = "cleanly"
	} else {
		// 将多行结果合并为一行
		GitStatus = strings.Replace(strings.Replace(GitStatus, "\r\n", " |", -1), "\n", " |", -1)
	}
	GoRuntime = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
	res.BuildTime = &BuildTime
	res.GitHash = &GitHash
	res.GitLog = &GitLog
	res.GitStatus = &GitStatus
	res.GoRuntime = &GoRuntime
	res.GoVersion = &GoVersion
	return
}
