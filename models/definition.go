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
	Key   *string         `json:"key" gorm:"type:varchar(50);unique_index;not null;"`
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
	UserType  *string `gorm:"default:'USER_TYPE_ADMIN'" json:"user_type"`
}

// 主机信息
//gen:qs
type Host struct {
	CommonDBModel
	// 主机名称
	//HostName *string `json:"host_name" gorm:"unique;not null;type:varchar(100)"`
	HostName *string `json:"host_name"`
	// 备注信息
	Mark *string `json:"mark"`
	// 发布者
	CreatorID *uint64 `json:"creator_id"`
	Creator   *Admin  `json:"creator"`
	// 网卡信息
	Interfaces json.RawMessage `json:"interfaces" gorm:"type:json"`
}

type Keycode struct {
	// 键盘的值
	Text *string `json:"text"`
	// 键盘的编号
	KeyCode *string `json:"key_code"`
}

// 键盘的事件
//gen:qs
// Keyboard is the payload type of the keyboard service log method.
type Keyboard struct {
	CommonDBModel
	// 主机ID
	HostID *uint64 `json:"host_id"`
	// 回车后发送的事件
	Keys json.RawMessage `json:"keys" gorm:"type:json"`
	// 事件码
	EventCode *string `gorm:"index;type:varchar(100)" json:"event_code"`
	Count     uint    `gorm:"-" json:"count"`
}

type InterfaceInfo struct {
	// 地址
	IPAddrs []string `json:"ip_addrs"`
	// 网卡名称
	Name *string `json:"name"`
	// 主机mac地址
	MacAddr *string `json:"mac_addr"`
}

// 通知信息
//gen:qs
type Notice struct {
	CommonDBModel
	// 消息消息提醒内容
	Notice *string `json:"notice"`
	// 发布者
	CreatorID *uint64 `json:"creator_id"`
	Creator   *Admin  `json:"creator"`
}

// RTSP的数据
//gen:qs
type Rtsp struct {
	CommonDBModel
	// 发布者
	CreatorID *uint64 `json:"creator_id"`
	Creator   *Admin  `json:"creator"`
	// rtsp的地址
	RtspURL *string `json:"rtsp_url"`
	Name    *string `json:"name"`
	// 主机ID
	HostID *string `json:"host_id"`
	Host   Host    `json:"host"`
	// 其他的扩展属性，采用json字符串展示
	Ext *string `json:"ext"`
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
