package models

import (
	"database/sql"
	"encoding/json"
	"github.com/guregu/null"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

//go:generate goqueryset -in definition.go
type Model struct {
	ID        *string    `gorm:"primary_key;unique;not null;type:varchar(100);comment:'数据编号'" json:"id"`
	CreatedAt time.Time  `json:"created_at" gorm:"comment:'数据创建时间'"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"comment:'数据更新时间'"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at" gorm:"comment:'数据删除时间'"`
}

type CommonDBModel struct {
	ID        uint64     `gorm:"primary_key;not null;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time  `json:"created_at" gorm:"index"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// 字典配置
//gen:qs
type Config struct {
	CommonDBModel
	// 配置名称
	Key   null.String     `json:"key" gorm:"type:varchar(50);unique_index;not null;"`
	Value json.RawMessage `json:"value" gorm:"type:json"`
}

// 管理员
//gen:qs
type Admin struct {
	CommonDBModel
	// 账户登陆名称
	LoginName *string `json:"login_name" gorm:"unique"`
	Password  *string ` json:"password"`
	// 生成密码是使用的盐
	Salt *string `json:"salt"`
	// 创建该账户的账户
	CreatedBy *uint64 `json:"created_by"`
}

// 主机信息
//gen:qs
type Host struct {
	CommonDBModel
	// 地址
	IPAddr *string
	// 主机名称
	HostName *string
	// 主机mac地址
	MacAddr *string
	// 备注信息
	Mark *string
	// 发布者
	CreatorID *uint64
	Creator   *Admin
}

// 键盘的事件
//gen:qs
type Keyboard struct {
	CommonDBModel
	// 主机ID
	HostID *string
	// 键盘的编号
	Key *string
	// 值
	Value *string
}

// 通知信息
//gen:qs
type Notice struct {
	CommonDBModel
	// 消息消息提醒内容
	Notice *string
	// 发布者
	CreatorID *uint64
	Creator   *Admin
}

// RTSP的数据
//gen:qs
type Rtsp struct {
	CommonDBModel
	// 发布者
	CreatorID *uint64
	Creator   *Admin
	// rtsp的地址
	RtspURL *string
	// 主机ID
	HostID *string
	Host   Host
	// 其他的扩展属性，采用json字符串展示
	Ext *string
}

// 后台生成token
//gen:qs
type Token struct {
	CommonDBModel
	// Token数据
	Token *string `gorm:"unique;not null;type:varchar(100)"`
	// 发布者
	CreatorID *uint64
	Creator   *Admin
	// 名称
	Name *string
	// 描述
	Description *string
}
