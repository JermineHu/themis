package themis

import (
	"context"
	"fmt"
	"log"

	rtsp "github.com/JermineHu/themis/svc/gen/rtsp"
	"goa.design/goa/v3/security"
)

// rtsp service example implementation.
// The example methods log the requests and return zero values.
type rtspsrvc struct {
	logger *log.Logger
}

// NewRtsp returns the rtsp service implementation.
func NewRtsp(logger *log.Logger) rtsp.Service {
	return &rtspsrvc{logger}
}

// JWTAuth implements the authorization logic for service "rtsp" for the "jwt"
// security scheme.
func (s *rtspsrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	//
	// TBD: add authorization logic.
	//
	// In case of authorization failure this function should return
	// one of the generated error structs, e.g.:
	//
	//    return ctx, myservice.MakeUnauthorizedError("invalid token")
	//
	// Alternatively this function may return an instance of
	// goa.ServiceError with a Name field value that matches one of
	// the design error names, e.g:
	//
	//    return ctx, goa.PermanentError("unauthorized", "invalid token")
	//
	return ctx, fmt.Errorf("not implemented")
}

// 流的数据列表；
func (s *rtspsrvc) List(ctx context.Context, p *rtsp.ListPayload) (res *rtsp.RtspResult, view string, err error) {
	res = &rtsp.RtspResult{}
	view = "default"
	s.logger.Print("rtsp.list")
	return
}

// 创建RTSP数据
func (s *rtspsrvc) Create(ctx context.Context, p *rtsp.Rtsp) (res *rtsp.RtspResult, view string, err error) {
	res = &rtsp.RtspResult{}
	view = "default"
	s.logger.Print("rtsp.create")
	return
}

// 根据id修改数据
func (s *rtspsrvc) Update(ctx context.Context, p *rtsp.Rtsp) (res *rtsp.RtspResult, view string, err error) {
	res = &rtsp.RtspResult{}
	view = "default"
	s.logger.Print("rtsp.update")
	return
}

// 根据id删除
func (s *rtspsrvc) Delete(ctx context.Context, p *rtsp.DeletePayload) (res bool, err error) {
	s.logger.Print("rtsp.delete")
	return
}

// 根据id获取信息
func (s *rtspsrvc) Show(ctx context.Context, p *rtsp.ShowPayload) (res *rtsp.RtspResult, view string, err error) {
	res = &rtsp.RtspResult{}
	view = "default"
	s.logger.Print("rtsp.show")
	return
}
