package models

import (
	"database/sql"
	"errors"
	"github.com/JermineHu/themis/svc/gen/rtsp"
	"github.com/guregu/null"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

////GetRtspWhere 根据参数获取精确查询条件
//func (qs RtspQuerySet) GetRtspWhere(where *rtsp.RtspWhere) RtspQuerySet {
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
func GetRtspList(payload *rtsp.ListPayload) (result []Rtsp, count int64, err error) {
	offsetHead := payload.OffsetHead
	OffsetTail := payload.OffsetTail

	if OffsetTail-offsetHead < 0 || OffsetTail-offsetHead > 500 {
		return result, 0, errors.New("OffsetTail must be uper than offsetHead and  OffsetTail-offsetHead must lower 200!")
	}
	qs := NewRtspQuerySet(rdb_themis)
	if err != nil {
		return
	}
	//if payload.Where != nil {
	//	qs = qs.GetRtspWhere(payload.Where)
	//}
	if payload.IsDesc {
		qs = qs.w(qs.db.Order(payload.OrderBy + " DESC"))
	} else {
		qs = qs.w(qs.db.Order(payload.OrderBy + " ASC"))
	}
	totalNum, err := qs.Count() //查询count
	count = int64(totalNum)
	list := []Rtsp{}
	err = qs.db.Offset(int(offsetHead)).Limit(int(OffsetTail - offsetHead)).Preload("Host").Find(&list).Error
	return list, count, err
}

//创建数据
func CreateRtsp(a *Rtsp) error {
	return a.Create(rdb_themis)
}

func GetRtspCount() (count int, err error) {
	qs := NewRtspQuerySet(rdb_themis)
	count, err = qs.Count() //查询count
	if err != nil {
		return
	}
	return
}

// 根据ID修改
func UpdateRtspByID(id uint64, mi *Rtsp) error {
	mi.ID = id
	qs := NewRtspQuerySet(rdb_themis)
	result := qs.db.Where("id=?", id).Save(&mi)
	if result.Error != nil {
		return rtsp.MakeBadRequest(result.Error)
	}
	if result.RowsAffected == 0 {
		return rtsp.MakeNotFound(new(rtsp.NotFound))
	}
	return nil
}

// 根据ID删除
func DeleteRtspByID(id uint64) (count int64, err error) {
	qs := NewRtspQuerySet(rdb_themis)
	db := qs.db.Unscoped().Delete(Rtsp{}, "id =?", id)
	return db.RowsAffected, db.Error
}

//获取Id获取详情
func GetRtspById(id *uint64) (result Rtsp, err error) {
	qs := NewRtspQuerySet(rdb_themis)
	pc := Rtsp{}
	err = qs.w(qs.db.Where("id =?", id).Preload("Host")).One(&pc)
	return pc, err
}
