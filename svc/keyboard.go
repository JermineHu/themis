package svc

import (
	"context"
	"github.com/JermineHu/themis/models"
	"github.com/jinzhu/copier"
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
	ctx, err := JWTCheck(ctx, token, scheme)
	if err != nil {
		return ctx, keyboard.MakeUnauthorized(err)
	}
	return ctx, err
}

// 键盘日志分页列表；
func (s *keyboardsrvc) List(ctx context.Context, p *keyboard.ListPayload) (res *keyboard.KeyboardList, err error) {
	res = &keyboard.KeyboardList{}
	if p == nil {
		lpl := keyboard.ListPayload{}
		lpl.OffsetTail = 20
		p = &lpl
	}
	list, count, err := models.GetKeyboardList(p)
	if err != nil {
		return
	}
	res.Count = &count
	ls := []*keyboard.KeyboardResult{}

	for i := range list {
		item := keyboard.KeyboardResult{}
		err = copier.Copy(&item, &list[i])
		if err != nil {
			return
		}
		ls = append(ls, &item)

	}
	res.PageData = ls
	return
}

// 创建日志数据
func (s *keyboardsrvc) Log(ctx context.Context, p *keyboard.Keyboard) (res *keyboard.KeyboardResult, err error) {
	res = &keyboard.KeyboardResult{}
	cp := models.Keyboard{}
	err = copier.Copy(&cp, p)
	if err != nil {
		return
	}
	err = models.CreateKeyboard(&cp)
	if err != nil {
		err = keyboard.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据主机ID删除，清空日志数据键盘
func (s *keyboardsrvc) Clear(ctx context.Context, p *keyboard.ClearPayload) (res bool, err error) {
	count, err := models.DeleteKeyboardByHostID(p.HostID)
	res = count > 0
	return
}
