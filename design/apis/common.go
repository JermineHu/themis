/*
 * @Description: In User Settings Edit
 * @Author: Jermine Hu
 * @Date: 2019-09-06 17:38:25
 * @LastEditTime: 2019-09-06 17:45:05
 * @LastEditors: Please set LastEditors
 */
/**
Create by jermine

Date  19-8-14-下午4:54

**/
package apis

import (
	"flag"
	"github.com/gofrs/uuid"
	. "goa.design/goa/v3/dsl"
	"math/rand"
	"time"
)

func NewUUIDStr() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return uuid.String()
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandomString(n int) string {
	b := make([]rune, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

var RequestParams = Type("RequestParams", func() {

	Field(1, "user_id", String, "用户的id", func() {
		Format(FormatUUID)
		Example(NewUUIDStr())
		//Meta("rpc:tag", "1")
	})
	Field(2, "start_time", String, "开始时间", func() {
		Example(time.Now().Format(time.RFC3339))
		Format(FormatDateTime)
		//	Meta("rpc:tag", "2")
	})

	Field(3, "end_time", String, "结束时间", func() {
		Example(time.Now().Format(time.RFC3339))
		Format(FormatDateTime)
		//Meta("rpc:tag", "3")
	})
	Field(4, "order_by", String, "排序字段", func() {
		Example("id")
		//Meta("rpc:tag", "4")
	})
	Field(5, "is_desc", Boolean, "是否为降序", func() {
		Example(false)
		//Meta("rpc:tag", "5")
	})
	Field(6, "offset_head", Int64, "从多少条开始", func() {
		Example(0)
	})
	Field(7, "offset_tail", Int64, "到多少条结束", func() {
		Example(20)
	})
	Field(8, "keywords", String, "关键字", func() {
		Example("")
		//	Meta("rpc:tag", "8")
	})
	//Field(9, "where", Bytes, "条件查询对象", func() {
	//	//	Meta("rpc:tag", "9")
	//})
	TokenField(10, "token", String, "JWTAuth token used to perform authorization", func() {
		//Meta("rpc:tag", "10")
	})
	Required("offset_head", "offset_tail")
})

// PageModelAccount is the PageModelAccount resource media type.
var PageModel = ResultType("application/vnd.page_mode+json", func() {
	Description("分页返回是数据模型")
	Attributes(func() {
		Attribute("count", Int64, "数据条数", func() {
			Example(200)
			Meta("rpc:tag", "1")
		})
		Attribute("page_data", Bytes, "得到的分页数据", func() {
			Meta("rpc:tag", "2")
		})

	})

	View("default", func() {
		Field(1, "count")
		Field(2, "page_data")
	})
})

var AnyType = ResultType("application/vnd.any_mode+json", func() {
	Description("通用数据模型")
	Attributes(func() {
		Attribute("data", Bytes, "二进制数据", func() {
			Meta("rpc:tag", "1")
		})
	})

	View("default", func() {
		Field(1, "data")
	})
})

// JWTAuth defines a security scheme that uses JWT tokens.
var JWTAuth = JWTSecurity("jwt", func() {
	//Header("Authorization")
	Scope("api:access", "API access") // Define "api:access" scope
	Scope("api:read", "Read-only access")
	Scope("api:write", "Read and write access")
})

// JWTAuth defines a security scheme that uses JWT tokens.
var AdminJWTAuth = JWTSecurity("adminjwt", func() {
	//Header("Authorization")
	Scope("api:access", "API access") // Define "api:access" scope
	Scope("api:read", "Read-only access")
	Scope("api:write", "Read and write access")
})

// APIKeyAuth defines a security scheme that uses API keys.
var APIKeyAuth = APIKeySecurity("api_key", func() {
	Description("Secures endpoint by requiring an API key.")
})

var BasicAuth = BasicAuthSecurity("basic_auth", func() {
	Description("基本身份认证")
	Scope("api:read", "Read-only access")

})

var (
	hostF = flag.String("host", "development", "Server host (valid values: development, production)")
)

var CustomErrorType = ResultType("application/vnd.goa.error", func() {
	Attributes(func() {
		Attribute("message", String, "Error returned.", func() {
			//Meta("struct:error:name")
		})
		Attribute("occurred_at", String, "Time error occurred.", func() {
			Format(FormatDateTime)
		})
	})
})

var SortData = Type("SortData", func() {
	Field(1, "id", UInt64, "数据编号", func() {
		Example(20)
	})
	Field(2, "ordinal", Int, "排序号", func() {
		Example(20)
	})
})
