package svc

import (
	"context"
	"github.com/JermineHu/themis/models"
	"github.com/JermineHu/themis/utils"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"log"
	"strings"
	"time"

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
	ctx, err := JWTCheck(ctx, token, scheme)
	if err != nil {
		return ctx, admin.MakeUnauthorized(err)
	}
	return ctx, err
}

// 根据账户名称密码进行登陆操作！
func (s *adminsrvc) Login(ctx context.Context, p *admin.LoginPayload) (res *admin.AdminResult, err error) {
	res = &admin.AdminResult{}
	adm, err := models.GetAdminUserByName(p.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = admin.MakeNotFound(new(admin.InvalidAccountNotFound))
			return
		}
		return
	}
	if adm == nil {
		err = admin.MakeNotFound(new(admin.InvalidAccountNotFound))
		return
	}
	if !strings.EqualFold(utils.Md5(p.Password+*adm.Salt), *adm.Password) {
		err = admin.MakeInvalidAccountOrPassword(new(admin.InvalidAccountOrPassword))
		return
	}

	res.LoginName = adm.LoginName
	res.ID = &adm.ID
	ad := admin.Admin{}
	err = copier.Copy(&ad, &res)
	if err != nil {
		return
	}
	tk, err := makeJWTWithAdmin(ad)
	if err != nil {
		return res, err
	}
	res.Token = tk
	res.UserType = adm.UserType
	return
}

// 退出登陆（将发上来的JWT加入拒绝名单，并且设置过期时间为JWT到期时间，到期自动释放），删除本地JWT
func (s *adminsrvc) Logout(ctx context.Context, p *admin.LogoutPayload) (res bool, err error) {
	s.logger.Print("admin.logout")
	return
}

// 列表数据；
func (s *adminsrvc) List(ctx context.Context, p *admin.ListPayload) (res *admin.AdminList, err error) {
	res = &admin.AdminList{}
	if p == nil {
		lpl := admin.ListPayload{}
		lpl.OffsetTail = 20
		p = &lpl
	}
	list, count, err := models.GetAdminList(p)
	if err != nil {
		return
	}
	res.Count = &count
	ls := []*admin.AdminResult{}

	for i := range list {
		item := admin.AdminResult{}
		err = copier.Copy(&item, &list[i])
		if err != nil {
			return
		}
		ct := list[i].CreatedAt.Format(time.RFC3339)
		item.CreatedAt = &ct
		//ut := list[i].UpdatedAt.Format(time.RFC3339)
		//item.UpdatedAt = &ut
		ls = append(ls, &item)

	}
	res.PageData = ls
	return
}

// 创建数据
func (s *adminsrvc) Create(ctx context.Context, p *admin.Admin) (res *admin.AdminResult, err error) {
	res = &admin.AdminResult{}
	cp := models.Admin{}
	err = copier.Copy(&cp, p)
	if err != nil {
		return
	}
	str := utils.GetRandomString(8)
	cp.Salt = &str
	pwd := utils.Md5(*p.Password + str)
	cp.Password = &pwd
	err = models.CreateAdmin(&cp)
	if err != nil {
		err = admin.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据id修数据
func (s *adminsrvc) Update(ctx context.Context, p *admin.Admin) (res *admin.AdminResult, err error) {
	res = &admin.AdminResult{}
	cp := models.Admin{}
	err = copier.Copy(&cp, p)
	if err != nil {
		return
	}
	str := utils.GetRandomString(8)
	cp.Salt = &str
	pwd := utils.Md5(*p.Password + str)
	cp.Password = &pwd
	err = models.UpdateAdminByID(*p.ID, &cp)
	if err != nil {
		err = admin.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据id删除
func (s *adminsrvc) Delete(ctx context.Context, p *admin.DeletePayload) (res bool, err error) {
	count, err := models.DeleteAdminByID(p.ID)
	res = count > 0
	return
}

// 根据id信息
func (s *adminsrvc) Show(ctx context.Context, p *admin.ShowPayload) (res *admin.AdminResult, err error) {
	res = &admin.AdminResult{}
	cp, err := models.GetAdminById(p.ID)
	if err != nil {
		err = admin.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&cp, &res)
	if err != nil {
		return
	}
	return
}
