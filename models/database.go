// This file "init.go" is created by lincan at 11/19/15.
// Copyright © 2015 - lincan. All rights reserved

package models

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	rdb_valuation        *gorm.DB
	rdb_user_pay_service *gorm.DB
	rdb_activity         *gorm.DB
	rdb_crius            *gorm.DB
	rdb_product_center   *gorm.DB
	rdb_auctionbid       *gorm.DB
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

func GetRDBCrius(dbType, cstr string) *gorm.DB {
	if rdb_crius == nil {
		rdb_crius = GetDB(dbType, fmt.Sprintf(cstr, common.DB_CRIUS))
	}
	return rdb_crius
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

func SetUpRDBZeus(DB *gorm.DB) {
	if err := DB.Set("gorm:table_options", "comment '订单表'").AutoMigrate(&Order{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '礼品领取表'").AutoMigrate(&Gift{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '交易明细表'").AutoMigrate(&TransactionLog{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '广告位'").AutoMigrate(&Advert{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '广告详情'").AutoMigrate(&AdvertDetail{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '马甲号'").AutoMigrate(&Puppet{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '广告详情'").AutoMigrate(&AdvertDetail{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '参数配置'").AutoMigrate(&Config{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '开团有奖活动信息表'").AutoMigrate(&Regiment{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '开团信息表'").AutoMigrate(&RegimentUp{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '白拿活动信息表'").AutoMigrate(&Baina{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '白拿记录表'").AutoMigrate(&BainaLog{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '商品分类表'").AutoMigrate(&MallItemCategory{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '商品属性表'").AutoMigrate(&MallItemAttr{}).
		AddForeignKey("item_category_id", "mall_item_categories(id)", "RESTRICT", "RESTRICT").
		Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '商品模板表'").AutoMigrate(&MallItemTemplate{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '商品款式表'").AutoMigrate(&MallItemStyle{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '商品表'").AutoMigrate(&MallItem{}).
		AddForeignKey("cat_id", "mall_item_categories(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("style_id", "mall_item_styles(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("template_id", "mall_item_templates(id)", "RESTRICT", "RESTRICT").
		Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '商品属性关联表'").AutoMigrate(&MallItemAttrRef{}).
		AddForeignKey("attr_id", "mall_item_attrs(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("mall_item_id", "mall_items(id)", "RESTRICT", "RESTRICT").
		Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '商品轮播图表'").AutoMigrate(&MallItemBanner{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '商品详情图'").AutoMigrate(&MallItemDetailImg{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '商品售卖规格表'").AutoMigrate(&MallItemSellSpecs{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '订单表'").AutoMigrate(&MallOrder{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '订单中的商品表'").AutoMigrate(&MallOrderItem{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '订单退款记录'").AutoMigrate(&MallOrderRefund{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '商城活动类型'").AutoMigrate(&MallActivityType{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '商城活动'").AutoMigrate(&MallActivity{}).Error; err != nil {
		panic(err)
	}
	if err := DB.Set("gorm:table_options", "comment '商城活动商品'").AutoMigrate(&MallActivityItem{}).Error; err != nil {
		panic(err)
	}
}
