package themis

import (
	"context"
	"fmt"
	"log"

	keyboard "github.com/JermineHu/themis/svc/gen/keyboard"
	"goa.design/goa/v3/security"
)

// keyboard service example implementation.
// The example methods log the requests and return zero values.
type keyboardsrvc struct {
	logger *log.Logger
}

// NewKeyboard returns the keyboard service implementation.
func NewKeyboard(logger *log.Logger) keyboard.Service {
	return &keyboardsrvc{logger}
}

// JWTAuth implements the authorization logic for service "keyboard" for the
// "jwt" security scheme.
func (s *keyboardsrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
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

// 键盘日志分页列表；
func (s *keyboardsrvc) List(ctx context.Context, p *keyboard.ListPayload) (res *keyboard.KeyboardResult, err error) {
	res = &keyboard.KeyboardResult{}
	s.logger.Print("keyboard.list")
	return
}

// 创建日志数据
func (s *keyboardsrvc) Log(ctx context.Context, p *keyboard.Keyboard) (res *keyboard.KeyboardResult, err error) {
	res = &keyboard.KeyboardResult{}
	s.logger.Print("keyboard.log")
	return
}

// 根据主机ID删除，清空日志数据键盘
func (s *keyboardsrvc) Clear(ctx context.Context, p *keyboard.ClearPayload) (res bool, err error) {
	s.logger.Print("keyboard.clear")
	return
}
