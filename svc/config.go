package themis

import (
	"context"
	"fmt"
	"log"

	config "github.com/JermineHu/themis/svc/gen/config"
	"goa.design/goa/v3/security"
)

// config service example implementation.
// The example methods log the requests and return zero values.
type configsrvc struct {
	logger *log.Logger
}

// NewConfig returns the config service implementation.
func NewConfig(logger *log.Logger) config.Service {
	return &configsrvc{logger}
}

// JWTAuth implements the authorization logic for service "config" for the
// "jwt" security scheme.
func (s *configsrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
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

// 配置列表；
func (s *configsrvc) List(ctx context.Context, p *config.ListPayload) (res *config.ListResult, err error) {
	res = &config.ListResult{}
	s.logger.Print("config.list")
	return
}

// 创建配置
func (s *configsrvc) Create(ctx context.Context, p *config.Config1) (res *config.Config, view string, err error) {
	res = &config.Config{}
	view = "default"
	s.logger.Print("config.create")
	return
}

// 根据id修改配置数据
func (s *configsrvc) Update(ctx context.Context, p *config.Config1) (res *config.Config, view string, err error) {
	res = &config.Config{}
	view = "default"
	s.logger.Print("config.update")
	return
}

// 根据id删除
func (s *configsrvc) Delete(ctx context.Context, p *config.DeletePayload) (res bool, err error) {
	s.logger.Print("config.delete")
	return
}

// 根据key获取配置数据
func (s *configsrvc) Show(ctx context.Context, p *config.ShowPayload) (res *config.Config, view string, err error) {
	res = &config.Config{}
	view = "default"
	s.logger.Print("config.show")
	return
}
