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

var HostInfo = Type("HostInfo", func() {
	Field(1, "id", UInt64, "数据ID", func() {})
	Field(2, "ip_addr", String, "地址", func() {})
	Field(3, "host_name", String, "主机名称", func() {})
	Field(4, "mac_addr", String, "主机mac地址", func() {})
	Field(5, "mark", String, "备注信息", func() {})
	TokenField(13, "token", String, "JWTAuth token used to perform authorization", func() {
		//Meta("rpc:tag", "10")
	})

})

var PageModelHost = ResultType("application/vnd.host_list+json", func() {
	Description("分页返回是数据模型")
	Attributes(func() {
		Attribute("count", Int64, "数据条数", func() {
			Example(200)
			Meta("rpc:tag", "1")
		})
		Attribute("page_data", CollectionOf(HostResult), "得到的分页数据", func() {
			//Example(func() {
			//	//Value(CollectionOf(HostResult.View("tiny")).Example(expr.NewRandom("")))
			//})
			Meta("rpc:tag", "2")
		})
	})

	View("default", func() {
		Attribute("count")
		Attribute("page_data")
	})
})

// HostResult is the host resource media type.
var HostResult = ResultType("application/vnd.host_result+json", func() {
	Description("用户模型")
	Attributes(func() {
		Field(1, "id", UInt64, "数据ID", func() {})
		Field(2, "ip_addr", String, "地址", func() {})
		Field(3, "host_name", String, "主机名称", func() {})
		Field(4, "mac_addr", String, "主机mac地址", func() {})
		Field(5, "mark", String, "备注信息", func() {})
		Field(6, "created_at", String, "创建时间", func() {
			Format(FormatDateTime)
			//Meta("rpc:tag", "8")
		})
	})
})

// BasicAuth defines a security scheme using basic authentication. The scheme protects the "signin"
// action used to create JWTs.
//var BasicAuth = BasicAuthSecurity("BasicAuth")
var res_host = Service("host", func() {
	Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWTAuth claims.
	})
	HTTP(func() {
		Path("/themis/v1/host")
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

	GRPC(func() {
		Response("Unauthorized", CodeUnauthenticated)
	})

	Method("list", func() {
		Description("主机列表")
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
				Default("order_num")
				//Meta("rpc:tag", "4")
			})
			Field(5, "is_desc", Boolean, "是否为降序", func() {
				Example(false)
				Default(true)
				//Meta("rpc:tag", "5")
			})
			Field(6, "where", HostInfo, "条件", func() {})
			Required("offset_head", "offset_tail")
		})
		Error("Unauthorized")
		Error("BadRequest")
		Error("NotFound")
		Result(HostResult)
		HTTP(func() {
			POST("")
			Response(StatusOK, func() {

			})
			Response("Unauthorized", StatusUnauthorized)
			Response("NotFound", StatusNoContent)
		})

		GRPC(func() {
			Response(CodeOK)
			Response(CodeNotFound)
		})
	})

	Method("registry", func() {
		Description("agent注册")
		Payload(HostInfo)
		Error("Unauthorized")
		Result(HostResult)
		HTTP(func() {
			POST("")
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

	Method("update", func() {
		Description("根据id修改数据")
		Payload(HostInfo)
		Error("Unauthorized")
		Result(HostResult)
		HTTP(func() {
			PUT("/{id}")
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

	Method("delete", func() {
		Description("根据id删除")
		Error("Unauthorized")
		Payload(func() {
			TokenField(1, "token", String, "JWTAuth token used to perform authorization", func() {
			})
			Field(2, "id", Int, "要删除的id", func() {
			})
			Required("id")
		})
		Result(Boolean)
		HTTP(func() {
			DELETE("/{id}")
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
		Description("根据id获取信息")
		Error("Unauthorized")
		Payload(func() {
			TokenField(1, "token", String, "JWTAuth token used to perform authorization", func() {
			})
			Field(2, "id", UInt64, "id", func() {
			})

		})
		Result(HostResult)
		HTTP(func() {
			GET("/{id}")
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
})
