package themis

import (
	"context"
	"fmt"
	"log"

	host "github.com/JermineHu/themis/svc/gen/host"
	"goa.design/goa/v3/security"
)

// host service example implementation.
// The example methods log the requests and return zero values.
type hostsrvc struct {
	logger *log.Logger
}

// NewHost returns the host service implementation.
func NewHost(logger *log.Logger) host.Service {
	return &hostsrvc{logger}
}

// JWTAuth implements the authorization logic for service "host" for the "jwt"
// security scheme.
func (s *hostsrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
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

// 主机列表
func (s *hostsrvc) List(ctx context.Context, p *host.ListPayload) (res *host.HostResult, err error) {
	res = &host.HostResult{}
	s.logger.Print("host.list")
	return
}

// agent注册
func (s *hostsrvc) Registry(ctx context.Context, p *host.HostInfo) (res *host.HostResult, err error) {
	res = &host.HostResult{}
	s.logger.Print("host.registry")
	return
}

// 根据id修改数据
func (s *hostsrvc) Update(ctx context.Context, p *host.HostInfo) (res *host.HostResult, err error) {
	res = &host.HostResult{}
	s.logger.Print("host.update")
	return
}

// 根据id删除
func (s *hostsrvc) Delete(ctx context.Context, p *host.DeletePayload) (res bool, err error) {
	s.logger.Print("host.delete")
	return
}

// 根据id获取信息
func (s *hostsrvc) Show(ctx context.Context, p *host.ShowPayload) (res *host.HostResult, err error) {
	res = &host.HostResult{}
	s.logger.Print("host.show")
	return
}
