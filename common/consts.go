package common

// 公共配置
const (
	DOMAIN = "DOMAIN"
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
