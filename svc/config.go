package svc

import (
	"context"
	"encoding/json"
	"github.com/JermineHu/themis/models"
	config "github.com/JermineHu/themis/svc/gen/config"
	"github.com/jinzhu/copier"
	"goa.design/goa/v3/security"
	"log"
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
	ctx, err := JWTCheck(ctx, token, scheme)
	if err != nil {
		return ctx, config.MakeUnauthorized(err)
	}
	return ctx, err
}

// 配置列表；
func (s *configsrvc) List(ctx context.Context, p *config.ListPayload) (res *config.ListResult, err error) {
	res = &config.ListResult{}
	if p == nil {
		lpl := config.ListPayload{}
		lpl.OffsetTail = 20
		p = &lpl
	}
	list, count, err := models.GetConfigList(p)
	if err != nil {
		return
	}
	res.Count = &count
	ls := []*config.ConfigResult{}

	for i := range list {
		item := config.ConfigResult{}
		mp := make(map[string]string)
		err = json.Unmarshal(list[i].Value, &mp)
		if err != nil {
			err = config.MakeBadRequest(err)
			return
		}
		item.Key = list[i].Key
		item.Value = mp
		ls = append(ls, &item)

	}
	res.PageData = ls
	return
}

// 创建配置
func (s *configsrvc) Create(ctx context.Context, p *config.Config) (res *config.ConfigResult, view string, err error) {
	res = &config.ConfigResult{}
	view = "default"
	cp := models.Config{}
	v, err := json.Marshal(p.Value)
	if err != nil {
		err = config.MakeBadRequest(err)
		return
	}
	cp.Value = v
	cp.Key = p.Key
	err = models.CreateConfig(&cp)
	if err != nil {
		err = config.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据key修改配置数据
func (s *configsrvc) Update(ctx context.Context, p *config.Config) (res *config.ConfigResult, view string,
	err error) {
	res = &config.ConfigResult{}
	view = "default"
	cp, err := models.GetConfigByKey(*p.Key)
	if err != nil {
		err = config.MakeBadRequest(err)
		return
	}
	v, err := json.Marshal(p.Value)
	if err != nil {
		err = config.MakeBadRequest(err)
		return
	}
	cp.Value = v
	cp.Key = p.Key
	err = models.UpdateConfigByKey(*p.Key, &cp)
	if err != nil {
		err = config.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据id删除
func (s *configsrvc) Delete(ctx context.Context, p *config.DeletePayload) (res bool, err error) {
	count, err := models.DeleteConfigByKey(p.Key)
	res = count > 0
	return
}

// 根据key获取配置数据
func (s *configsrvc) Show(ctx context.Context, p *config.ShowPayload) (res *config.ConfigResult, view string, err error) {
	res = &config.ConfigResult{}
	cp, err := models.GetConfigByKey(*p.K)
	if err != nil {
		err = config.MakeBadRequest(err)
		return
	}
	res.Key = cp.Key
	res.Value = cp.Value
	return
}
