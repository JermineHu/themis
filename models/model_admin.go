package models

import (
	"database/sql"
	"errors"
	"github.com/JermineHu/themis/svc/gen/admin"
	"github.com/guregu/null"
	"strings"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

//GetAdminWhere 根据参数获取精确查询条件
func (qs AdminQuerySet) GetAdminWhere(where *admin.Admin) AdminQuerySet {
	if where.ID != nil && *where.ID != 0 {
		qs = qs.w(qs.db.Where("id=?", where.ID))
	}
	return qs
}

//获取数据列表
func GetAdminList(payload *admin.ListPayload) (result []Admin, count int64, err error) {
	offsetHead := payload.OffsetHead
	OffsetTail := payload.OffsetTail

	if OffsetTail-offsetHead < 0 || OffsetTail-offsetHead > 200 {
		return result, 0, errors.New("OffsetTail must be uper than offsetHead and  OffsetTail-offsetHead must lower 200!")
	}
	qs := NewAdminQuerySet(rdb_themis)
	if err != nil {
		return
	}
	if payload.Where != nil {
		qs = qs.GetAdminWhere(payload.Where)
	}
	if payload.IsDesc {
		qs = qs.w(qs.db.Order(payload.OrderBy + " DESC"))
	} else {
		qs = qs.w(qs.db.Order(payload.OrderBy + " ASC"))
	}
	totalNum, err := qs.Count() //查询count
	count = int64(totalNum)
	list := []Admin{}
	err = qs.db.Offset(int(offsetHead)).Limit(int(OffsetTail - offsetHead)).Find(&list).Error
	return list, count, err
}

//创建数据
func CreateAdmin(a *Admin) error {
	return a.Create(rdb_themis)
}

func GetAdminCount() (count int, err error) {
	qs := NewAdminQuerySet(rdb_themis)
	count, err = qs.Count() //查询count
	if err != nil {
		return
	}
	return
}

// 根据ID修改
func UpdateAdminByID(id uint64, mi *Admin) error {
	mi.ID = id
	qs := NewAdminQuerySet(rdb_themis)
	result := qs.db.Where("id=?", id).Save(&mi)
	if result.Error != nil {
		return admin.MakeBadRequest(result.Error)
	}
	if result.RowsAffected == 0 {
		return admin.MakeNotFound(new(admin.NotFound))
	}
	return nil
}

// 根据ID删除
func DeleteAdminByID(id int) (count int64, err error) {
	qs := NewAdminQuerySet(rdb_themis)
	db := qs.db.Unscoped().Delete(Admin{}, "id =?", id)
	return db.RowsAffected, db.Error
}

//获取Id获取详情
func GetAdminById(id *uint64) (result Admin, err error) {
	qs := NewAdminQuerySet(rdb_themis)
	pc := Admin{}
	err = qs.w(qs.db.Where("id =?", id)).One(&pc)
	return pc, err
}

// 根据用户名称获取单个用户
func GetAdminUserByName(name string) (account *Admin, err error) {
	if len(strings.TrimSpace(name)) == 0 {
		return account, errors.New("The account name was needed!")
	}
	qs := NewAdminQuerySet(rdb_themis)
	acco := Admin{}
	err = qs.w(qs.db.Where("login_name=?", name)).One(&acco)
	account = &acco
	return
}
