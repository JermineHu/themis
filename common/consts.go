package common

// 公共配置
const (
	DOMAIN   = "DOMAIN"
	APP_PORT = "APP_PORT"
)

// 用户类型
const (
	USER_TYPE_NORMAL = "USER_TYPE_NORMAL" // 普通用户
	USER_TYPE_ADMIN  = "USER_TYPE_ADMIN"  // 管理员用户
)

// 设置事件类型
const (
	EVENT_TYPE_CLEAN            = "EVENT_TYPE_CLEAN"            // 清空事件类型
	EVENT_TYPE_DELETE_PRVE_DATA = "EVENT_TYPE_DELETE_PRVE_DATA" // 删除上一个事件类型
	EVENT_TYPE_NORMAL           = "EVENT_TYPE_NORMAL"           // 一般事件类型
)

// Redis ENV key
const (
	REDIS_SENTINEL_ADDR = "REDIS_SENTINEL_ADDR"
	REDIS_PWD           = "REDIS_PWD"
	REDIS_DB            = "REDIS_DB"
)

//For Http request and response
const (
	REMOTE_IP    = "REMOTE_IP"
	CTX_REQUEST  = "CTX_REQUEST"
	CTX_RESPONSE = "CTX_RESPONSE"
)

//Token的超时时间
const (
	TOKEN_TIMEOUT = "TOKEN_TIMEOUT"
)

// RDB ENV key
const (
	DB_CON_STR    = "DB_CON_STR"
	DB_TYPE       = "DB_TYPE"
	DB_IS_UPGRADE = "DB_IS_UPGRADE"
	DB_THEMIS     = "themis"
)
