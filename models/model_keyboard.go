package models

import (
	"database/sql"
	"errors"
	"github.com/JermineHu/themis/svc/gen/keyboard"
	"github.com/guregu/null"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

////GetKeyboardWhere 根据参数获取精确查询条件
//func (qs KeyboardQuerySet) GetKeyboardWhere(where *keyboard.KeyboardWhere) KeyboardQuerySet {
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
func GetKeyboardList(payload *keyboard.ListPayload) (result []Keyboard, count int64, err error) {
	offsetHead := payload.OffsetHead
	OffsetTail := payload.OffsetTail

	if OffsetTail-offsetHead < 0 || OffsetTail-offsetHead > 200 {
		return result, 0, errors.New("OffsetTail must be uper than offsetHead and  OffsetTail-offsetHead must lower 200!")
	}
	qs := NewKeyboardQuerySet(rdb_themis)
	if err != nil {
		return
	}
	//if payload.Where != nil {
	//	qs = qs.GetKeyboardWhere(payload.Where)
	//}
	if payload.IsDesc {
		qs = qs.w(qs.db.Order(payload.OrderBy + " DESC"))
	} else {
		qs = qs.w(qs.db.Order(payload.OrderBy + " ASC"))
	}
	totalNum, err := qs.Count() //查询count
	count = int64(totalNum)
	list := []Keyboard{}
	err = qs.db.Offset(int(offsetHead)).Limit(int(OffsetTail - offsetHead)).Find(&list).Error
	return list, count, err
}

//创建数据
func CreateKeyboard(a *Keyboard) error {
	return a.Create(rdb_themis)
}

func GetKeyboardCount() (count int, err error) {
	qs := NewKeyboardQuerySet(rdb_themis)
	count, err = qs.Count() //查询count
	if err != nil {
		return
	}
	return
}

// 根据ID修改
func UpdateKeyboardByID(id uint64, mi *Keyboard) error {
	mi.ID = id
	qs := NewKeyboardQuerySet(rdb_themis)
	result := qs.db.Where("id=?", id).Save(&mi)
	if result.Error != nil {
		return keyboard.MakeBadRequest(result.Error)
	}
	if result.RowsAffected == 0 {
		return keyboard.MakeNotFound(new(keyboard.NotFound))
	}
	return nil
}

// 根据ID删除
func DeleteKeyboardByID(id int) (count int64, err error) {
	qs := NewKeyboardQuerySet(rdb_themis)
	db := qs.db.Unscoped().Delete(Keyboard{}, "id =?", id)
	return db.RowsAffected, db.Error
}

// 根据主机ID批量删除键盘数据
func DeleteKeyboardByHostID(host_id int) (count int64, err error) {
	qs := NewKeyboardQuerySet(rdb_themis)
	db := qs.db.Unscoped().Delete(Keyboard{}, "host_id =?", host_id)
	return db.RowsAffected, db.Error
}

//获取Id获取详情
func GetKeyboardById(id *uint64) (result Keyboard, err error) {
	qs := NewKeyboardQuerySet(rdb_themis)
	pc := Keyboard{}
	err = qs.w(qs.db.Where("id =?", id)).One(&pc)
	return pc, err
}