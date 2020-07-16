package svc

import (
	"context"
	"github.com/JermineHu/themis/models"
	"github.com/jinzhu/copier"
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
	ctx, err := JWTCheck(ctx, token, scheme)
	if err != nil {
		return ctx, host.MakeUnauthorized(err)
	}
	return ctx, err
}

// 主机列表
func (s *hostsrvc) List(ctx context.Context, p *host.ListPayload) (res *host.HostList, err error) {
	res = &host.HostList{}
	if p == nil {
		lpl := host.ListPayload{}
		lpl.OffsetTail = 20
		p = &lpl
	}
	list, count, err := models.GetHostList(p)
	if err != nil {
		return
	}
	res.Count = &count
	ls := []*host.HostResult{}

	for i := range list {
		item := host.HostResult{}
		err = copier.Copy(&item, &list[i])
		if err != nil {
			return
		}
		ls = append(ls, &item)

	}
	res.PageData = ls
	return
}

// agent注册
func (s *hostsrvc) Registry(ctx context.Context, p *host.HostInfo) (res *host.HostResult, err error) {
	res = &host.HostResult{}
	cp := models.Host{}
	err = copier.Copy(&cp, p)
	if err != nil {
		return
	}
	err = models.CreateHost(&cp)
	if err != nil {
		err = host.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据id修改数据
func (s *hostsrvc) Update(ctx context.Context, p *host.HostInfo) (res *host.HostResult, err error) {
	res = &host.HostResult{}
	cp := models.Host{}
	err = copier.Copy(&cp, p)
	if err != nil {
		return
	}
	err = models.UpdateHostByID(*p.ID, &cp)
	if err != nil {
		err = host.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据id删除
func (s *hostsrvc) Delete(ctx context.Context, p *host.DeletePayload) (res bool, err error) {
	count, err := models.DeleteHostByID(p.ID)
	res = count > 0
	return
}

// 根据id获取信息
func (s *hostsrvc) Show(ctx context.Context, p *host.ShowPayload) (res *host.HostResult, err error) {
	res = &host.HostResult{}
	cp, err := models.GetHostById(p.ID)
	if err != nil {
		err = host.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&cp, &res)
	if err != nil {
		return
	}
	return
}
