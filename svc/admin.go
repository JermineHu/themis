package themis

import (
	"context"
	"fmt"
	"log"

	admin "github.com/JermineHu/themis/svc/gen/admin"
	"goa.design/goa/v3/security"
)

// admin service example implementation.
// The example methods log the requests and return zero values.
type adminsrvc struct {
	logger *log.Logger
}

// NewAdmin returns the admin service implementation.
func NewAdmin(logger *log.Logger) admin.Service {
	return &adminsrvc{logger}
}

// JWTAuth implements the authorization logic for service "admin" for the "jwt"
// security scheme.
func (s *adminsrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
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

// 根据账户名称密码进行登陆操作！
func (s *adminsrvc) Login(ctx context.Context, p *admin.LoginPayload) (res *admin.AdminResult, err error) {
	res = &admin.AdminResult{}
	s.logger.Print("admin.login")
	return
}

// 退出登陆（将发上来的JWT加入拒绝名单，并且设置过期时间为JWT到期时间，到期自动释放），删除本地JWT
func (s *adminsrvc) Logout(ctx context.Context, p *admin.LogoutPayload) (res bool, err error) {
	s.logger.Print("admin.logout")
	return
}

// 列表数据；
func (s *adminsrvc) List(ctx context.Context, p *admin.ListPayload) (res *admin.AdminResult, err error) {
	res = &admin.AdminResult{}
	s.logger.Print("admin.list")
	return
}

// 创建数据
func (s *adminsrvc) Create(ctx context.Context, p *admin.Admin) (res *admin.AdminResult, err error) {
	res = &admin.AdminResult{}
	s.logger.Print("admin.create")
	return
}

// 根据id修数据
func (s *adminsrvc) Update(ctx context.Context, p *admin.Admin) (res *admin.AdminResult, err error) {
	res = &admin.AdminResult{}
	s.logger.Print("admin.update")
	return
}

// 根据id删除
func (s *adminsrvc) Delete(ctx context.Context, p *admin.DeletePayload) (res bool, err error) {
	s.logger.Print("admin.delete")
	return
}

// 根据id信息
func (s *adminsrvc) Show(ctx context.Context, p *admin.ShowPayload) (res *admin.AdminResult, err error) {
	res = &admin.AdminResult{}
	s.logger.Print("admin.show")
	return
}
