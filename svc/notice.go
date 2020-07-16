package themis

import (
	"context"
	"fmt"
	"log"

	notice "github.com/JermineHu/themis/svc/gen/notice"
	"goa.design/goa/v3/security"
)

// notice service example implementation.
// The example methods log the requests and return zero values.
type noticesrvc struct {
	logger *log.Logger
}

// NewNotice returns the notice service implementation.
func NewNotice(logger *log.Logger) notice.Service {
	return &noticesrvc{logger}
}

// JWTAuth implements the authorization logic for service "notice" for the
// "jwt" security scheme.
func (s *noticesrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
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

// 消息通知的列表
func (s *noticesrvc) List(ctx context.Context, p *notice.ListPayload) (res *notice.NoticeResult, err error) {
	res = &notice.NoticeResult{}
	s.logger.Print("notice.list")
	return
}

// 创建数据
func (s *noticesrvc) Create(ctx context.Context, p *notice.NoticeType) (res *notice.NoticeResult, err error) {
	res = &notice.NoticeResult{}
	s.logger.Print("notice.create")
	return
}

// 根据id修改通知
func (s *noticesrvc) Update(ctx context.Context, p *notice.NoticeType) (res *notice.NoticeResult, err error) {
	res = &notice.NoticeResult{}
	s.logger.Print("notice.update")
	return
}

// 根据id删除
func (s *noticesrvc) Delete(ctx context.Context, p *notice.DeletePayload) (res bool, err error) {
	s.logger.Print("notice.delete")
	return
}

// 根据id获取信息
func (s *noticesrvc) Show(ctx context.Context, p *notice.ShowPayload) (res *notice.NoticeResult, err error) {
	res = &notice.NoticeResult{}
	s.logger.Print("notice.show")
	return
}
