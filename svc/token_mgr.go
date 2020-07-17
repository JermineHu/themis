package svc

import (
	"context"
	"github.com/JermineHu/themis/models"
	"github.com/google/uuid"
	"log"
	"time"

	tokenmgr "github.com/JermineHu/themis/svc/gen/token_mgr"
	"goa.design/goa/v3/security"
)

// token_mgr service example implementation.
// The example methods log the requests and return zero values.
type tokenMgrsrvc struct {
	logger *log.Logger
}

// NewTokenMgr returns the token_mgr service implementation.
func NewTokenMgr(logger *log.Logger) tokenmgr.Service {
	return &tokenMgrsrvc{logger}
}

// JWTAuth implements the authorization logic for service "token_mgr" for the
// "jwt" security scheme.
func (s *tokenMgrsrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	ctx, err := JWTCheck(ctx, token, scheme)
	if err != nil {
		return ctx, tokenmgr.MakeUnauthorized(err)
	}
	return ctx, err
}

// token列表
func (s *tokenMgrsrvc) List(ctx context.Context, p *tokenmgr.ListPayload) (res *tokenmgr.TokenList, err error) {
	usr_id, err := GetUserIDByJWT(*p.Token)
	if err != nil {
		return
	}
	res = &tokenmgr.TokenList{}
	if p == nil {
		lpl := tokenmgr.ListPayload{}
		lpl.OffsetTail = 20
		p = &lpl
	}
	list, count, err := models.GetTokenList(p)
	res.Count = &count
	ls := []*tokenmgr.TokenResult{}
	for i := range list {
		temp := tokenmgr.TokenResult{}
		temp.ID = &list[i].ID
		temp.Name = list[i].Name
		temp.Token = list[i].Token
		temp.Description = list[i].Description
		ct := list[i].CreatedAt.Format(time.RFC3339)
		ut := list[i].UpdatedAt.Format(time.RFC3339)
		temp.CreatedAt = &ct
		temp.UpdatedAt = &ut
		temp.Creator = usr_id
		ls = append(ls, &temp)
	}
	res.PageData = ls
	return
}

// 创建数据
func (s *tokenMgrsrvc) Create(ctx context.Context, p *tokenmgr.CreatePayload) (res *tokenmgr.TokenResult, err error) {
	usr_id, err := GetUserIDByJWT(*p.Token)
	if err != nil {
		return
	}
	res = &tokenmgr.TokenResult{}
	uuid := uuid.New().String()
	cp := models.Token{
		Token:       &uuid,
		CreatorID:   usr_id,
		Name:        p.Name,
		Description: p.Description,
	}
	err = models.CreateToken(&cp)
	if err != nil {
		err = tokenmgr.MakeBadRequest(err)
		return
	}
	res.ID = &cp.ID
	res.Token = cp.Token
	res.Name = cp.Name
	return
}

// 根据id删除
func (s *tokenMgrsrvc) Delete(ctx context.Context, p *tokenmgr.DeletePayload) (res bool, err error) {
	count, err := models.DeleteTokenByID(p.ID)
	res = count > 0
	return
}
