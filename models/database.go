// This file "init.go" is created by lincan at 11/19/15.
// Copyright © 2015 - lincan. All rights reserved

package models

import (
	"database/sql"
	"fmt"
	"github.com/JermineHu/themis/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	rdb_themis *gorm.DB
)

func GetDB(dbType, cstr string) *gorm.DB {
	db, err := gorm.Open(dbType, cstr)

	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db is error")
	}
	db.LogMode(true)
	return db
}

func GetRDBThemis(dbType, cstr string) *gorm.DB {
	if rdb_themis == nil {
		rdb_themis = GetDB(dbType, fmt.Sprintf(cstr, common.DB_THEMIS))
	}
	return rdb_themis
}

//
//func GetTransaction() *gorm.DB {
//	return BeginRDB(rdb_valuation)
//}

func BeginRDB(DB *gorm.DB) *gorm.DB {
	txn := DB.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	return txn
}

func CommitRDB(txn *gorm.DB) {
	if txn == nil {
		return
	}
	txn.Commit()
	err := txn.Error
	if err != nil && err != sql.ErrTxDone {
		panic(err)
	}
}

func RollbackRDB(txn *gorm.DB) {
	if txn == nil {
		return
	}
	txn.Rollback()
	if err := txn.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
}

func SetUpRDBThemis(DB *gorm.DB) {
	if err := DB.AutoMigrate(&Admin{}).Error; err != nil {
		panic(err)
	}
	if err := DB.AutoMigrate(&Config{}).Error; err != nil {
		panic(err)
	}
	if err := DB.AutoMigrate(&Host{}).Error; err != nil {
		panic(err)
	}
	if err := DB.AutoMigrate(&Keyboard{}).Error; err != nil {
		panic(err)
	}
	if err := DB.AutoMigrate(&Notice{}).Error; err != nil {
		panic(err)
	}
	if err := DB.AutoMigrate(&Rtsp{}).Error; err != nil {
		panic(err)
	}
	if err := DB.AutoMigrate(&Token{}).Error; err != nil {
		panic(err)
	}
}
