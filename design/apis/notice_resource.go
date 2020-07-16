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

var Notice = Type("Notice", func() {
	TokenField(1, "token", String, "JWTAuth token used to perform authorization", func() {
		//Meta("rpc:tag", "10")
	})
	Field(2, "id", UInt64, "数据ID", func() {})
	Field(3, "notice", String, "消息消息提醒内容")
})

var PageModelNotice = ResultType("application/vnd.notice_list+json", func() {
	Description("分页返回是数据模型")
	Attributes(func() {
		Attribute("count", Int64, "数据条数", func() {
			Example(200)
			Meta("rpc:tag", "1")
		})
		Attribute("page_data", CollectionOf(NoticeResult), "得到的分页数据", func() {
			//Example(func() {
			//	//Value(CollectionOf(NoticeResult.View("tiny")).Example(expr.NewRandom("")))
			//})
			Meta("rpc:tag", "2")
		})
	})

	View("default", func() {
		Attribute("count")
		Attribute("page_data")
	})
})

// NoticeResult is the notice resource media type.
var NoticeResult = ResultType("application/vnd.notice_result+json", func() {
	Description("用户模型")
	Attributes(func() {
		Field(2, "id", UInt64, "数据ID")
		Field(3, "notice", String, "消息消息提醒内容")
		Field(4, "creator", UInt64, "发布者")
		Field(5, "created_at", String, "创建时间")
		Field(6, "updated_at", String, "更新时间")
	})

})

// BasicAuth defines a security scheme using basic authentication. The scheme protects the "signin"
// action used to create JWTs.
//var BasicAuth = BasicAuthSecurity("BasicAuth")
var res_notice = Service("notice", func() {
	Description("消息通知模块")
	Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWTAuth claims.
	})
	HTTP(func() {
		Path("/themis/v1/notice")
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
		Description("消息通知的列表")
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
			Field(6, "where", Notice, "条件", func() {})
			Required("offset_head", "offset_tail")
		})
		Error("Unauthorized")
		Error("BadRequest")
		Error("NotFound")
		Result(NoticeResult)
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

	Method("create", func() {
		Description("创建数据")
		Payload(Notice)
		Error("Unauthorized")
		Result(NoticeResult)
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
		Description("根据id修改通知")
		Payload(Notice)
		Error("Unauthorized")
		Result(NoticeResult)
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
		Result(NoticeResult)
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
