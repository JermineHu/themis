package svc

import (
	"github.com/JermineHu/themis/common"
	"github.com/JermineHu/themis/models"
	"os"
	"strings"
)

// 初始化数据库
func initDB() {
	rdb_themis := models.GetRDBThemis(os.Getenv(common.DB_TYPE), os.Getenv(common.DB_CON_STR)) // 根据数据库类类型和字符串初始化链接操作
	if strings.EqualFold(os.Getenv(common.DB_IS_UPGRADE), "true") {                            //数据库升级的开关
		models.SetUpRDBThemis(rdb_themis) // 升级数据库
	}
}

// 初始化
func init() {
	initDB() // 初始化数据库
}
