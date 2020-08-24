package models

import (
	"database/sql"
	"errors"
	"github.com/JermineHu/themis/svc/gen/host"
	"github.com/guregu/null"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

////GetHostWhere 根据参数获取精确查询条件
//func (qs HostQuerySet) GetHostWhere(where *host.Host) HostQuerySet {
//	if where.ID != nil && *where.ID != 0 {
//		qs = qs.w(qs.db.Where("id=?", where.ID))
//	}
//	return qs
//}

//获取数据列表
func GetHostList(payload *host.ListPayload) (result []Host, count int64, err error) {
	offsetHead := payload.OffsetHead
	OffsetTail := payload.OffsetTail

	if OffsetTail-offsetHead < 0 || OffsetTail-offsetHead > 200 {
		return result, 0, errors.New("OffsetTail must be uper than offsetHead and  OffsetTail-offsetHead must lower 200!")
	}
	qs := NewHostQuerySet(rdb_themis)
	if err != nil {
		return
	}
	//if payload.Where != nil {
	//	qs = qs.GetHostWhere(payload.Where)
	//}
	if payload.IsDesc {
		qs = qs.w(qs.db.Order(payload.OrderBy + " DESC"))
	} else {
		qs = qs.w(qs.db.Order(payload.OrderBy + " ASC"))
	}
	totalNum, err := qs.Count() //查询count
	count = int64(totalNum)
	list := []Host{}
	err = qs.db.Offset(int(offsetHead)).Limit(int(OffsetTail - offsetHead)).Find(&list).Error
	return list, count, err
}

//获取所有主机ID
func GetAllHostID() (result []uint64, err error) {
	qs := NewHostQuerySet(rdb_themis)
	if err != nil {
		return
	}
	list := []uint64{}
	err = qs.db.Pluck("id", &list).Error
	return list, err
}

//创建数据
func CreateHost(a *Host) error {
	return a.Create(rdb_themis)
}

func GetHostCount() (count int, err error) {
	qs := NewHostQuerySet(rdb_themis)
	count, err = qs.Count() //查询count
	if err != nil {
		return
	}
	return
}

// 根据ID修改
func UpdateHostByID(id uint64, mi *Host) error {
	mi.ID = id
	qs := NewHostQuerySet(rdb_themis)
	result := qs.db.Where("id=?", id).Save(&mi)
	if result.Error != nil {
		return host.MakeBadRequest(result.Error)
	}
	if result.RowsAffected == 0 {
		return host.MakeNotFound(new(host.NotFound))
	}
	return nil
}

// 根据ID删除
func DeleteHostByID(id int) (count int64, err error) {
	qs := NewHostQuerySet(rdb_themis)
	db := qs.db.Unscoped().Delete(Host{}, "id =?", id)
	return db.RowsAffected, db.Error
}

//获取Id获取详情
func GetHostById(id *uint64) (result Host, err error) {
	qs := NewHostQuerySet(rdb_themis)
	pc := Host{}
	err = qs.w(qs.db.Where("id =?", id)).One(&pc)
	return pc, err
}
