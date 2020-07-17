package apis

import (
	. "goa.design/goa/v3/dsl"
)

var MapType = MapOf(String, String, func() {
	//Key(func() {
	//	Minimum(1)           // Validates keys of the map
	//})
	//Elem(func() {
	//	Pattern("[a-zA-Z]+") // Validates values of the map
	//})
})

var Config = Type("config", func() {
	Description("业务配置数据 config")
	Field(1, "id", UInt64, "ID")
	Field(2, "key", String, "配置名称")
	Field(3, "value", Any, "值")
	Field(4, "created_at", String, "创建时间")
	Field(6, "updated_at", String, "修改时间")
	TokenField(7, "token", String, "JWTAuth token used to perform authorization", func() {
	})
})

var ConfigResult = ResultType("application/vnd.config_result", func() {
	Description("业务配置数据返回对象 config")
	Attributes(func() {
		Field(1, "key", String, "配置名称")
		Field(2, "value", Any, "值")
	})

	View("default", func() {
		Field(1, "key")
		Field(2, "value")
	})

	View("full", func() {
		Field(1, "key")
		Field(2, "value")
	})
})

var res_config = Service("config", func() {
	Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWTAuth claims.
	})
	HTTP(func() {
		Path("/themis/v1/config")
	})
	Error("IllegalUserFailure", String, "非法用户！")
	Error("UserNotFund", String, "用户不存在！")
	Error("Unauthorized", String, "授权失败！")
	Error("NotFound", String, "未查询到数据！")

	Method("list", func() {
		Description("配置列表；")
		Payload(func() {
			Field(1, "offset_head", Int64, "从多少条开始", func() {
				Example(0)
				Minimum(0)
			})
			Field(2, "offset_tail", Int64, "到多少条结束", func() {
				Example(20)
				Minimum(1)
			})
			TokenField(3, "token", String, "JWTAuth token used to perform authorization", func() {
			})
			Field(4, "order_by", String, "排序字段", func() {
				Example("id")
				Default("updated_at")
				//Meta("rpc:tag", "4")
			})
			Field(5, "is_desc", Boolean, "是否为降序", func() {
				Example(false)
				Default(true)
				//Meta("rpc:tag", "5")
			})
			Required("offset_head", "offset_tail")
		})
		Error("Unauthorized")
		Error("BadRequest")
		Error("NotFound")
		//Result(func() {
		//	Field(1, "page_data", CollectionOf(ConfigResult), "得到的分类列表", func() {
		//		View("default")
		//	})
		//})

		Result(func() {
			Field(1, "count", Int64, "数据条数", func() {
				Example(200)
			})
			Field(2, "page_data", CollectionOf(ConfigResult), "得到的分页数据", func() {
				View("default")
			})
		})
		HTTP(func() {
			POST("/list")
			Response(StatusOK, func() {
			})
			Response("Unauthorized", StatusUnauthorized)
			Response("NotFound", StatusNoContent)
		})
	})
	Method("create", func() {
		Description("创建配置")
		Payload(Config)
		Error("Unauthorized")
		Result(ConfigResult)
		HTTP(func() {
			POST("")
			Response(StatusOK, func() {
			})
			Response(StatusNotFound)
			Response("Unauthorized", StatusUnauthorized)
		})
	})
	Method("update", func() {
		Description("根据key修改配置数据")
		Payload(Config)
		Error("Unauthorized")
		Result(ConfigResult)
		HTTP(func() {
			PUT("/{key}")
			Response(StatusOK, func() {
			})
			Response(StatusNotFound)
			Response("Unauthorized", StatusUnauthorized)
		})
	})

	Method("delete", func() {
		Description("根据id删除")
		Error("Unauthorized")
		Payload(func() {
			TokenField(1, "token", String, "JWTAuth token used to perform authorization", func() {
			})
			Field(2, "key", String, "要删除的key", func() {
			})
			Required("key")
		})
		Result(Boolean)
		HTTP(func() {
			DELETE("/{key}")
			Response(StatusOK, func() {
			})
			Response(StatusNotFound)
			Response("Unauthorized", StatusUnauthorized)
		})
		GRPC(func() {
			Response(CodeOK)
			Response(CodeNotFound)
		})
	})
	Method("show", func() {
		Description("根据key获取配置数据")
		Error("Unauthorized")
		Payload(func() {
			TokenField(1, "token", String, "JWTAuth token used to perform authorization", func() {
			})
			Field(2, "k", String, "对应资源的key", func() {
			})
		})
		Result(ConfigResult)
		HTTP(func() {
			GET("/show/{k}")
			Response(StatusOK, func() {
			})
			Response(StatusNotFound)
			Response("Unauthorized", StatusUnauthorized)
		})
	})
})
