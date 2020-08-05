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

var PageModelToken = ResultType("application/vnd.token_list+json", func() {
	Description("分页返回是数据模型")
	Attributes(func() {
		Attribute("count", Int64, "数据条数", func() {
			Example(200)
			Meta("rpc:tag", "1")
		})
		Attribute("page_data", CollectionOf(TokenResult), "得到的分页数据", func() {
			//Example(func() {
			//	//Value(CollectionOf(TokenResult.View("tiny")).Example(expr.NewRandom("")))
			//})
			Meta("rpc:tag", "2")
		})
	})

	View("default", func() {
		Attribute("count")
		Attribute("page_data")
	})
})

// TokenResult is the token resource media type.
var TokenResult = ResultType("application/vnd.token_result+json", func() {
	Description("Token模型")
	Attributes(func() {
		Field(1, "id", UInt64, "数据ID")
		Field(2, "token", String, "Token数据")
		Field(3, "creator", UInt64, "发布者")
		Field(4, "created_at", String, "创建时间")
		Field(5, "updated_at", String, "更新时间")
		Field(6, "name", String, "名称")
		Field(7, "description", String, "描述")
	})

})

// BasicAuth defines a security scheme using basic authentication. The scheme protects the "signin"
// action used to create JWTs.
//var BasicAuth = BasicAuthSecurity("BasicAuth")
var res_token = Service("token_mgr", func() {
	Description("Token模块")
	Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
		Scope("api:write") // Enforce presence of "api" scope in JWTAuth claims.
	})
	HTTP(func() {
		Path("/themis/v1/token")
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
		Description("token列表")
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
			Required("offset_head", "offset_tail")
		})
		Error("Unauthorized")
		Error("BadRequest")
		Error("NotFound")
		Result(PageModelToken)
		HTTP(func() {
			POST("/list")
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

	Method("create", func() {
		Description("创建数据")
		Payload(func() {
			TokenField(1, "token", String, "JWTAuth token used to perform authorization", func() { //Meta("rpc:tag", "10")})
			})
			Field(2, "name", String, "名称")
			Field(3, "description", String, "描述")
		})
		Error("Unauthorized")
		Result(TokenResult)
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

})
