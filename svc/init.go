package svc

import (
	"github.com/JermineHu/themis/common"
	"github.com/JermineHu/themis/etc/config"
	"github.com/JermineHu/themis/models"
	"os"
	"strings"
)

// 初始化数据库
func initDB() {
	rdb_user_pay_service := models.GetRDBUserPayService(os.Getenv(common.DB_TYPE), os.Getenv(common.DB_CON_STR)) // 根据数据库类类型和字符串初始化链接操作
	if strings.EqualFold(os.Getenv(common.DB_IS_UPGRADE), "true") {                                              //数据库升级的开关
		models.SetUpRDBUserPayService(rdb_user_pay_service) // 升级数据库
	}

	rdb_activity := models.GetRDBActivity(os.Getenv(common.DB_TYPE), os.Getenv(common.DB_CON_STR)) // 根据数据库类类型和字符串初始化链接操作
	if strings.EqualFold(os.Getenv(common.DB_IS_UPGRADE), "true") {                                //数据库升级的开关
		models.SetUpRDBActivity(rdb_activity) // 升级数据库
	}

	rdb_zeus := models.GetRDBCrius(os.Getenv(common.DB_TYPE), os.Getenv(common.DB_CON_STR)) // 根据数据库类类型和字符串初始化链接操作
	if strings.EqualFold(os.Getenv(common.DB_IS_UPGRADE), "true") {                         //数据库升级的开关
		models.SetUpRDBZeus(rdb_zeus) // 升级数据库
	}

	rdb_auctionbid := models.GetRDBAuctionbid(os.Getenv(common.DB_TYPE), os.Getenv(common.DB_CON_STR)) // 根据数据库类类型和字符串初始化链接操作
	if strings.EqualFold(os.Getenv(common.DB_IS_UPGRADE), "true") {                                    //数据库升级的开关
		models.SetUpRDBAuctionbid(rdb_auctionbid) // 升级数据库
	}

	rdb_product_center := models.GetRDBrdb_ProductCenter(os.Getenv(common.DB_TYPE), os.Getenv(common.DB_CON_STR)) // 根据数据库类类型和字符串初始化链接操作
	if strings.EqualFold(os.Getenv(common.DB_IS_UPGRADE), "true") {                                               //数据库升级的开关
		models.SetUpRDBProductCenter(rdb_product_center) // 升级数据库
	}
}

// 初始化
func init() {
	initDB()          // 初始化数据库
	config.InitConf() // 初始化配置文件
}
