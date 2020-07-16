package themis

import (
	"context"
	"fmt"
	"log"

	token "github.com/JermineHu/themis/svc/gen/token"
	"goa.design/goa/v3/security"
)

// token service example implementation.
// The example methods log the requests and return zero values.
type tokensrvc struct {
	logger *log.Logger
}

// NewToken returns the token service implementation.
func NewToken(logger *log.Logger) token.Service {
	return &tokensrvc{logger}
}

// JWTAuth implements the authorization logic for service "token" for the "jwt"
// security scheme.
func (s *tokensrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
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

// token列表
func (s *tokensrvc) List(ctx context.Context, p *token.ListPayload) (res *token.TokenResult, err error) {
	res = &token.TokenResult{}
	s.logger.Print("token.list")
	return
}

// 创建数据
func (s *tokensrvc) Create(ctx context.Context, p *token.CreatePayload) (res *token.TokenResult, err error) {
	res = &token.TokenResult{}
	s.logger.Print("token.create")
	return
}

// 根据id删除
func (s *tokensrvc) Delete(ctx context.Context, p *token.DeletePayload) (res bool, err error) {
	s.logger.Print("token.delete")
	return
}
