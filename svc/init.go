package svc

import (
	"encoding/json"
	"github.com/JermineHu/themis/common"
	"github.com/JermineHu/themis/models"
	"github.com/JermineHu/themis/utils"
	"github.com/jinzhu/gorm"
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
	initDB()            // 初始化数据库
	initStream()        // 初始化流
	loadKeyMapSetting() // 加载设置好的键盘事件
}

// 初始化流
func initStream() {
	RTSPConfig = loadConfig()
	DataInit() // 初始化admin
}

// 加载设置好的键盘事件
func loadKeyMapSetting() {
	cf, err := models.GetConfigByKey("keyboard.mean")
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		panic(err)
	}
	if err == nil {
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
}

type Settings struct {
	DataInit bool `json:"data_init"`
}

// 初始化admin
func DataInit() {
	key := "settings"
	_, err := models.GetConfigByKey(key)
	if gorm.IsRecordNotFoundError(err) {
		cf := models.Config{}
		cf.Key = &key
		sets := Settings{DataInit: true}
		data, _ := json.Marshal(sets)
		cf.Value = data
		err = models.CreateConfig(&cf)
		if err != nil {
			panic(err)
		}
		cp := models.Admin{}
		ln := "admin"
		cp.LoginName = &ln
		str := utils.GetRandomString(8)
		cp.Salt = &str
		pwd := utils.Md5("CZ1lHMQvZVdBe5fn" + str)
		cp.Password = &pwd
		err = models.CreateAdmin(&cp)
		if err != nil {
			panic(err)
		}
	}
}
