package models

import (
	"database/sql"
	"errors"
	tokenmgr "github.com/JermineHu/themis/svc/gen/token_mgr"
	"github.com/guregu/null"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

////GetTokenWhere 根据参数获取精确查询条件
//func (qs TokenQuerySet) GetTokenWhere(where *token.TokenWhere) TokenQuerySet {
//	if where.ID != nil && *where.ID != 0 {
//		qs = qs.w(qs.db.Where("id=?", where.ID))
//	}
//	if where.IsEnable != nil {
//		qs = qs.w(qs.db.Where("is_enable=?", where.IsEnable))
//	}
//	if where.PItemStyleID != nil {
//		qs = qs.w(qs.db.Where("p_item_category_id=?", where.PItemStyleID))
//	}
//	return qs
//}

//获取数据列表
func GetTokenList(payload *tokenmgr.ListPayload) (result []Token, count int64, err error) {
	offsetHead := payload.OffsetHead
	OffsetTail := payload.OffsetTail

	if OffsetTail-offsetHead < 0 || OffsetTail-offsetHead > 200 {
		return result, 0, errors.New("OffsetTail must be uper than offsetHead and  OffsetTail-offsetHead must lower 200!")
	}
	qs := NewTokenQuerySet(rdb_themis)
	if err != nil {
		return
	}
	//if payload.Where != nil {
	//	qs = qs.GetTokenWhere(payload.Where)
	//}
	if payload.IsDesc {
		qs = qs.w(qs.db.Order(payload.OrderBy + " DESC"))
	} else {
		qs = qs.w(qs.db.Order(payload.OrderBy + " ASC"))
	}
	totalNum, err := qs.Count() //查询count
	count = int64(totalNum)
	list := []Token{}
	err = qs.db.Offset(int(offsetHead)).Limit(int(OffsetTail - offsetHead)).Find(&list).Error
	return list, count, err
}

//创建数据
func CreateToken(a *Token) error {
	return a.Create(rdb_themis)
}

func GetTokenCount() (count int, err error) {
	qs := NewTokenQuerySet(rdb_themis)
	count, err = qs.Count() //查询count
	if err != nil {
		return
	}
	return
}

// 根据ID修改
func UpdateTokenByID(id uint64, mi *Token) error {
	mi.ID = id
	qs := NewTokenQuerySet(rdb_themis)
	result := qs.db.Where("id=?", id).Save(&mi)
	if result.Error != nil {
		return tokenmgr.MakeBadRequest(result.Error)
	}
	if result.RowsAffected == 0 {
		return tokenmgr.MakeNotFound(new(tokenmgr.NotFound))
	}
	return nil
}

// 根据ID删除
func DeleteTokenByID(id int) (count int64, err error) {
	qs := NewTokenQuerySet(rdb_themis)
	db := qs.db.Unscoped().Delete(Token{}, "id =?", id)
	return db.RowsAffected, db.Error
}

//获取Id获取详情
func GetTokenById(id *uint64) (result Token, err error) {
	qs := NewTokenQuerySet(rdb_themis)
	pc := Token{}
	err = qs.w(qs.db.Where("id =?", id)).One(&pc)
	return pc, err
}

//查询token是否存在
func IsExistToken(tk string) (exist bool, err error) {
	qs := NewTokenQuerySet(rdb_themis)
	count, err := qs.w(qs.db.Where("token =?", tk)).Count()
	return count > 0, err
}
