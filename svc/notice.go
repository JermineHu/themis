package svc

import (
	"context"
	"github.com/JermineHu/themis/models"
	"github.com/jinzhu/copier"
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
	ctx, err := JWTCheck(ctx, token, scheme)
	if err != nil {
		return ctx, notice.MakeUnauthorized(err)
	}
	return ctx, err
}

// 消息通知的列表
func (s *noticesrvc) List(ctx context.Context, p *notice.ListPayload) (res *notice.NoticeList, err error) {
	res = &notice.NoticeList{}
	if p == nil {
		lpl := notice.ListPayload{}
		lpl.OffsetTail = 20
		p = &lpl
	}
	list, count, err := models.GetNoticeList(p)
	if err != nil {
		return
	}
	res.Count = &count
	ls := []*notice.NoticeResult{}

	for i := range list {
		item := notice.NoticeResult{}
		err = copier.Copy(&item, &list[i])
		if err != nil {
			return
		}
		ls = append(ls, &item)

	}
	res.PageData = ls
	return
}

// 创建数据
func (s *noticesrvc) Create(ctx context.Context, p *notice.Notice) (res *notice.NoticeResult, err error) {
	res = &notice.NoticeResult{}
	cp := models.Notice{}
	err = copier.Copy(&cp, p)
	if err != nil {
		return
	}
	err = models.CreateNotice(&cp)
	if err != nil {
		err = notice.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据id修改通知
func (s *noticesrvc) Update(ctx context.Context, p *notice.Notice) (res *notice.NoticeResult, err error) {
	res = &notice.NoticeResult{}
	cp := models.Notice{}
	err = copier.Copy(&cp, p)
	if err != nil {
		return
	}
	err = models.UpdateNoticeByID(*p.ID, &cp)
	if err != nil {
		err = notice.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据id删除
func (s *noticesrvc) Delete(ctx context.Context, p *notice.DeletePayload) (res bool, err error) {
	count, err := models.DeleteConfigByID(p.ID)
	res = count > 0
	return
}

// 根据id获取信息
func (s *noticesrvc) Show(ctx context.Context, p *notice.ShowPayload) (res *notice.NoticeResult, err error) {
	res = &notice.NoticeResult{}
	cp, err := models.GetNoticeById(p.ID)
	if err != nil {
		err = notice.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&cp, &res)
	if err != nil {
		return
	}
	return
}
