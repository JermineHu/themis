/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-09-06 17:39:24
 * @LastEditTime: 2019-09-06 17:48:01
 * @LastEditors: Please set LastEditors
 */
package apis

import (
	. "goa.design/goa/v3/dsl"
)

var Keycode = Type("keycode", func() {
	Field(1, "text", String, "键盘的值")
	Field(2, "key_code", String, "键盘的编号")
})

var Keyboard = Type("keyboard", func() {
	TokenField(1, "token", String, "JWTAuth token used to perform authorization", func() {
		//Meta("rpc:tag", "10")
	})
	Field(2, "id", UInt64, "数据ID", func() {})
	Field(3, "host_id", UInt64, "主机ID")
	Field(4, "keys", ArrayOf(Keycode), "回车后发送的事件")
	Field(5, "color", String, "颜色")
	Field(6, "e_name", String, "事件名称")
	Field(7, "e_type", String, "事件类型")
	Field(8, "count", UInt, "事件次数")
})

var KeyboardEvent = Type("KeyboardEvent", func() {
	Field(1, "type", String, "键盘的值", func() {
		Enum("heartbeat", "keyboard")
	})
	Field(2, "keyboard_info", Keyboard, "键盘信息")
	Field(3, "ext", String, "其他信息")
	Field(4, "time", String, "时间")
	Required("type")
})

var PageModelKeyboard = ResultType("application/vnd.keyboard_list+json", func() {
	Description("分页返回是数据模型")
	Attributes(func() {
		Attribute("count", Int64, "数据条数", func() {
			Example(200)
			Meta("rpc:tag", "1")
		})
		Attribute("page_data", CollectionOf(KeyboardResult), "得到的分页数据", func() {
			//Example(func() {
			//	//Value(CollectionOf(KeyboardResult.View("tiny")).Example(expr.NewRandom("")))
			//})
			Meta("rpc:tag", "2")
		})
	})

	View("default", func() {
		Attribute("count")
		Attribute("page_data")
	})
})

// KeyboardResult is the keyboard resource media type.
var KeyboardResult = ResultType("application/vnd.keyboard_result+json", func() {
	Description("键盘数据模型")
	Reference(Keyboard)
	Attributes(func() {
		Field(1, "id")
		Field(2, "host_id")
		Field(3, "keys", Any)
		Field(4, "value")
		Field(5, "created_at")
	})
})

var res_keyboard = Service("keyboard", func() {
	Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
		Scope("api:write") // Enforce presence of "api" scope in JWTAuth claims.
	})
	HTTP(func() {
		Path("/themis/v1/keyboard")
		Headers(func() {
			//Header("Authorization", String, "JWTAuth token", func() {
			//})
			//Required("Authorization")
		})
	})

	Error("Unauthorized", String, "未授权！")
	Error("InvalidAccountNotFound", String, "用户不存在，请重试！")
	Error("InvalidAccountOrPassword", String, "用户名或密码错误！")
	Error("AuthorizedFailure", String, "授权失败！")
	Error("NotFound", String, "未查询到数据！")

	GRPC(func() {
		Response("Unauthorized", CodeUnauthenticated)
	})

	Method("list", func() {
		Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
			Scope("api:read") // Enforce presence of "api" scope in JWTAuth claims.
		})
		Description("键盘日志分页列表；")
		Payload(func() {
			Field(1, "offset_head", UInt64, "从多少条开始", func() {
				Example(0)
				Minimum(0)
			})
			Field(2, "offset_tail", UInt64, "到多少条结束", func() {
				Example(20)
				Minimum(1)
			})
			TokenField(3, "token", String, "JWTAuth token used to perform authorization", func() {
			})
			Field(4, "order_by", String, "排序字段", func() {
				Example("id")
				Default("id")
				//Meta("rpc:tag", "4")
			})
			Field(5, "is_desc", Boolean, "是否为降序", func() {
				Example(false)
				Default(true)
				//Meta("rpc:tag", "5")
			})
			Field(6, "where", Keyboard, "条件", func() {})
			Required("offset_head", "offset_tail")
		})
		Error("Unauthorized")
		Error("BadRequest")
		Error("NotFound")
		Result(PageModelKeyboard)
		HTTP(func() {
			POST("/logs")
			Response(StatusOK, func() {

			})
			Response("Unauthorized", StatusUnauthorized)
			Response("NotFound", StatusNoContent)
		})
	})

	Method("list_by_host_id", func() {
		Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
			Scope("api:read") // Enforce presence of "api" scope in JWTAuth claims.
		})
		Description("键盘日志分页列表；")
		Payload(func() {
			TokenField(1, "token", String, "JWTAuth token used to perform authorization", func() {
			})
			Field(2, "host_id", UInt64, "主机ID", func() {
			})
		})
		Error("Unauthorized")
		Error("BadRequest")
		Error("NotFound")
		Result(CollectionOf(KeyboardResult))
		HTTP(func() {
			GET("/logs/{host_id}")
			Response(StatusOK, func() {

			})
			Response("Unauthorized", StatusUnauthorized)
			Response("NotFound", StatusNoContent)
		})

	})

	Method("log", func() {
		Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
			Scope("api:access") // Enforce presence of "api" scope in JWTAuth claims.
		})
		Description("创建日志数据")
		Payload(Keyboard)
		Error("Unauthorized")
		Result(KeyboardResult)
		HTTP(func() {
			POST("/log/{host_id}")
			Response(StatusOK, func() {
			})
			Response(StatusNotFound)
			Response("Unauthorized", StatusUnauthorized)
		})
	})

	Method("clear", func() {
		Description("根据主机ID删除，清空日志数据键盘")
		Error("Unauthorized")
		Payload(func() {
			TokenField(1, "token", String, "JWTAuth token used to perform authorization", func() {
			})
			Field(2, "host_id", UInt64, "要删除的host_id", func() {
			})
			Required("host_id")
		})
		Result(Boolean)
		HTTP(func() {
			DELETE("/{host_id}")
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

	Method("broker", func() {
		//Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
		//	Scope("api:read") // Enforce presence of "api" scope in JWTAuth claims.
		//})
		Description("用于建立广播消息的服务")
		NoSecurity()
		Payload(func() {
			Field(1, "token", String, "JWTAuth token used to perform authorization", func() {
			})
			Field(2, "host_id", UInt64, "对应的host_id", func() {
			})
			Required("host_id")
		})

		StreamingPayload(KeyboardEvent)
		StreamingResult(KeyboardEvent)
		Error("Unauthorized")
		Error("BadRequest")
		Error("NotFound")
		HTTP(func() {
			GET("/broker/{host_id}")
			Response(StatusOK)
			Header("token:Authorization", String, "Auth token", func() {
				//Pattern("^Bearer [^ ]+$")
			})
			//Cookie("token:Authorization", String, "Auth token", func() {
			//	Pattern("^Bearer [^ ]+$")
			//})
			Response("Unauthorized", StatusUnauthorized)
			Response("BadRequest", StatusBadRequest)
			Response("NotFound", StatusBadRequest)
		})

		GRPC(func() {
			Response(CodeOK)
			Response("Unauthorized", CodeUnauthenticated)
			Response("BadRequest", CodeFailedPrecondition)
			Response("NotFound", CodeFailedPrecondition)
		})
	})

	Method("statistics", func() {
		Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
			Scope("api:read") // Enforce presence of "api" scope in JWTAuth claims.
		})
		Description("根据主机ID获取统计数据")
		Payload(func() {
			TokenField(1, "token", String, "JWTAuth token used to perform authorization", func() {
			})
			Field(2, "host_id", UInt64, "对应的host_id", func() {
			})
			Required("host_id")
		})
		Result(MapOf(String, Keyboard))
		Error("Unauthorized")
		Error("BadRequest")
		Error("NotFound")
		HTTP(func() {
			GET("/statistics/{host_id}")
			Response(StatusOK)
			Response("Unauthorized", StatusUnauthorized)
			Response("BadRequest", StatusBadRequest)
			Response("NotFound", StatusBadRequest)
		})

		GRPC(func() {
			Response(CodeOK)
			Response("Unauthorized", CodeUnauthenticated)
			Response("BadRequest", CodeFailedPrecondition)
			Response("NotFound", CodeFailedPrecondition)
		})
	})

	Method("broker_for_hosts", func() {
		//Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
		//	Scope("api:read") // Enforce presence of "api" scope in JWTAuth claims.
		//})
		Description("用于建立广播消息的服务")
		NoSecurity()
		Payload(func() {
			Field(1, "token", String, "JWTAuth token used to perform authorization", func() {
			})
			Field(2, "host_ids", ArrayOf(UInt64), "对应的host_id", func() {
			})
		})
		StreamingResult(KeyboardEvent)
		Error("Unauthorized")
		Error("BadRequest")
		Error("NotFound")
		HTTP(func() {
			GET("/broker_for_hosts")
			Params(func() {
				Param("host_ids:ids")
			})
			Response(StatusOK)
			Header("token:Authorization", String, "Auth token", func() {
				//Pattern("^Bearer [^ ]+$")
			})
			Response("Unauthorized", StatusUnauthorized)
			Response("BadRequest", StatusBadRequest)
			Response("NotFound", StatusBadRequest)
		})

		GRPC(func() {
			Response(CodeOK)
			Response("Unauthorized", CodeUnauthenticated)
			Response("BadRequest", CodeFailedPrecondition)
			Response("NotFound", CodeFailedPrecondition)
		})
	})

})
