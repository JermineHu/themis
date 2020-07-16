package svc

import (
	"context"
	"github.com/JermineHu/themis/models"
	"github.com/jinzhu/copier"
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
	ctx, err := JWTCheck(ctx, token, scheme)
	if err != nil {
		return ctx, rtsp.MakeUnauthorized(err)
	}
	return ctx, err
}

// 流的数据列表；
func (s *rtspsrvc) List(ctx context.Context, p *rtsp.ListPayload) (res *rtsp.RtspList, err error) {
	res = &rtsp.RtspList{}
	if p == nil {
		lpl := rtsp.ListPayload{}
		lpl.OffsetTail = 20
		p = &lpl
	}
	list, count, err := models.GetRtspList(p)
	if err != nil {
		return
	}
	res.Count = &count
	ls := []*rtsp.RtspResult{}

	for i := range list {
		item := rtsp.RtspResult{}
		err = copier.Copy(&item, &list[i])
		if err != nil {
			return
		}
		ls = append(ls, &item)

	}
	res.PageData = ls
	return
}

// 创建RTSP数据
func (s *rtspsrvc) Create(ctx context.Context, p *rtsp.Rtsp) (res *rtsp.RtspResult, view string, err error) {
	res = &rtsp.RtspResult{}
	view = "default"
	cp := models.Rtsp{}
	err = copier.Copy(&cp, p)
	if err != nil {
		return
	}
	err = models.CreateRtsp(&cp)
	if err != nil {
		err = rtsp.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据id修改数据
func (s *rtspsrvc) Update(ctx context.Context, p *rtsp.Rtsp) (res *rtsp.RtspResult, view string, err error) {
	res = &rtsp.RtspResult{}
	view = "default"
	cp := models.Rtsp{}
	err = copier.Copy(&cp, p)
	if err != nil {
		return
	}
	err = models.UpdateRtspByID(*p.ID, &cp)
	if err != nil {
		err = rtsp.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据id删除
func (s *rtspsrvc) Delete(ctx context.Context, p *rtsp.DeletePayload) (res bool, err error) {
	count, err := models.DeleteConfigByID(p.ID)
	res = count > 0
	return
}

// 根据id获取信息
func (s *rtspsrvc) Show(ctx context.Context, p *rtsp.ShowPayload) (res *rtsp.RtspResult, view string, err error) {
	res = &rtsp.RtspResult{}
	view = "default"
	cp, err := models.GetRtspById(p.ID)
	if err != nil {
		err = rtsp.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&cp, &res)
	if err != nil {
		return
	}
	return
}
