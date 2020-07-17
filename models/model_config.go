package models

import (
	"database/sql"
	"errors"
	"github.com/JermineHu/themis/svc/gen/config"
	"github.com/guregu/null"
	"strings"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

//GetConfigWhere 根据参数获取精确查询条件
func (qs ConfigQuerySet) GetConfigWhere(where *config.Config) ConfigQuerySet {
	if where.Key != nil && !strings.EqualFold(*where.Key, "") {
		qs = qs.w(qs.db.Where("key=?", where.Key))
	}
	return qs
}

//获取数据列表
func GetConfigList(payload *config.ListPayload) (result []Config, count int64, err error) {
	offsetHead := payload.OffsetHead
	OffsetTail := payload.OffsetTail

	if OffsetTail-offsetHead < 0 || OffsetTail-offsetHead > 200 {
		return result, 0, errors.New("OffsetTail must be uper than offsetHead and  OffsetTail-offsetHead must lower 200!")
	}
	qs := NewConfigQuerySet(rdb_themis)
	if err != nil {
		return
	}
	//if payload.Where != nil {
	//	qs = qs.GetConfigWhere(payload.Where)
	//}
	if payload.IsDesc {
		qs = qs.w(qs.db.Order(payload.OrderBy + " DESC"))
	} else {
		qs = qs.w(qs.db.Order(payload.OrderBy + " ASC"))
	}
	totalNum, err := qs.Count() //查询count
	count = int64(totalNum)
	list := []Config{}
	err = qs.db.Offset(int(offsetHead)).Limit(int(OffsetTail - offsetHead)).Find(&list).Error
	return list, count, err
}

//创建数据
func CreateConfig(a *Config) error {
	return a.Create(rdb_themis)
}

func GetConfigCount() (count int, err error) {
	qs := NewConfigQuerySet(rdb_themis)
	count, err = qs.Count() //查询count
	if err != nil {
		return
	}
	return
}

// 根据ID修改
func UpdateConfigByID(id uint64, mi *Config) error {
	mi.ID = id
	qs := NewConfigQuerySet(rdb_themis)
	result := qs.db.Where("id=?", id).Save(&mi)
	if result.Error != nil {
		return config.MakeBadRequest(result.Error)
	}
	if result.RowsAffected == 0 {
		return config.MakeNotFound(new(config.NotFound))
	}
	return nil
}

// 根据Key修改
func UpdateConfigByKey(key string, mi *Config) error {
	qs := NewConfigQuerySet(rdb_themis)
	result := qs.db.Where("key=?", key).Save(&mi)
	if result.Error != nil {
		return config.MakeBadRequest(result.Error)
	}
	if result.RowsAffected == 0 {
		return config.MakeNotFound(new(config.NotFound))
	}
	return nil
}

// 根据ID删除
func DeleteConfigByID(id int) (count int64, err error) {
	qs := NewConfigQuerySet(rdb_themis)
	db := qs.db.Unscoped().Delete(Config{}, "id =?", id)
	return db.RowsAffected, db.Error
}

// 根据Key删除
func DeleteConfigByKey(key string) (count int64, err error) {
	qs := NewConfigQuerySet(rdb_themis)
	db := qs.db.Unscoped().Delete(Config{}, "key =?", key)
	return db.RowsAffected, db.Error
}

//根据Id获取详情
func GetConfigById(id *uint64) (result Config, err error) {
	qs := NewConfigQuerySet(rdb_themis)
	pc := Config{}
	err = qs.w(qs.db.Where("id =?", id)).One(&pc)
	return pc, err
}

//根据Key获取详情
func GetConfigByKey(key string) (result Config, err error) {
	qs := NewConfigQuerySet(rdb_themis)
	pc := Config{}
	err = qs.w(qs.db.Where("key =?", key)).One(&pc)
	return pc, err
}
