package models

import (
	"database/sql"
	"errors"
	"github.com/JermineHu/themis/svc/gen/notice"
	"github.com/guregu/null"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

////GetNoticeWhere 根据参数获取精确查询条件
//func (qs NoticeQuerySet) GetNoticeWhere(where *notice.NoticeWhere) NoticeQuerySet {
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
func GetNoticeList(payload *notice.ListPayload) (result []Notice, count int64, err error) {
	offsetHead := payload.OffsetHead
	OffsetTail := payload.OffsetTail

	if OffsetTail-offsetHead < 0 || OffsetTail-offsetHead > 200 {
		return result, 0, errors.New("OffsetTail must be uper than offsetHead and  OffsetTail-offsetHead must lower 200!")
	}
	qs := NewNoticeQuerySet(rdb_themis)
	if err != nil {
		return
	}
	//if payload.Where != nil {
	//	qs = qs.GetNoticeWhere(payload.Where)
	//}
	if payload.IsDesc {
		qs = qs.w(qs.db.Order(payload.OrderBy + " DESC"))
	} else {
		qs = qs.w(qs.db.Order(payload.OrderBy + " ASC"))
	}
	totalNum, err := qs.Count() //查询count
	count = int64(totalNum)
	list := []Notice{}
	err = qs.db.Offset(int(offsetHead)).Limit(int(OffsetTail - offsetHead)).Find(&list).Error
	return list, count, err
}

//创建数据
func CreateNotice(a *Notice) error {
	return a.Create(rdb_themis)
}

func GetNoticeCount() (count int, err error) {
	qs := NewNoticeQuerySet(rdb_themis)
	count, err = qs.Count() //查询count
	if err != nil {
		return
	}
	return
}

// 根据ID修改
func UpdateNoticeByID(id uint64, mi *Notice) error {
	mi.ID = id
	qs := NewNoticeQuerySet(rdb_themis)
	result := qs.db.Where("id=?", id).Save(&mi)
	if result.Error != nil {
		return notice.MakeBadRequest(result.Error)
	}
	if result.RowsAffected == 0 {
		return notice.MakeNotFound(new(notice.NotFound))
	}
	return nil
}

// 根据ID删除
func DeleteNoticeByID(id int) (count int64, err error) {
	qs := NewNoticeQuerySet(rdb_themis)
	db := qs.db.Unscoped().Delete(Notice{}, "id =?", id)
	return db.RowsAffected, db.Error
}

//获取Id获取详情
func GetNoticeById(id *uint64) (result Notice, err error) {
	qs := NewNoticeQuerySet(rdb_themis)
	pc := Notice{}
	err = qs.w(qs.db.Where("id =?", id)).One(&pc)
	return pc, err
}
