package svc

import (
	"encoding/json"
	"github.com/JermineHu/themis/common"
	"github.com/JermineHu/themis/models"
	"os"
	"strconv"
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
	initDB()     // 初始化数据库
	initStream() // 初始化流
	loadKeyMapSetting()
}

// 初始化流
func initStream() {
	RTSPConfig = loadConfig()
}

// 加载设置好的键盘事件
func loadKeyMapSetting() {
	cf, err := models.GetConfigByKey("keyboard.mean")
	if err != nil {
		panic(err)
	}
	kM := []KeyboardSetting{}
	err = json.Unmarshal(cf.Value, &kM)
	if err != nil {
		panic(err)
	}
	for k := range kM {
		ks := []string{}
		for i := range kM[k].Keycodes {
			ks = append(ks, strconv.FormatInt(int64(kM[k].Keycodes[i].KeyCode), 10))
		}
		kstr := strings.Join(ks, "-")
		kbMap[kstr] = kM[k]
	}

}
