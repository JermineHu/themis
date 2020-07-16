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

var Rtsp = Type("Rtsp", func() {
	TokenField(1, "token", String, "JWTAuth token used to perform authorization", func() {
		//Meta("rpc:tag", "10")
	})
	Field(2, "rtsp_url", String, "rtsp的地址")
	Field(3, "host_id", String, "主机ID")
	Field(4, "ext", String, "其他的扩展属性，采用json字符串展示")
	Field(5, "created_at", String, "创建时间")
	Field(6, "updated_at", String, "更新时间")
	Field(7, "id", UInt64, "数据ID", func() {})
})

var PageModelRtsp = ResultType("application/vnd.rtsp_list+json", func() {
	Description("分页返回是数据模型")
	Attributes(func() {
		Attribute("count", Int64, "数据条数", func() {
			Example(200)
			Meta("rpc:tag", "1")
		})
		Attribute("page_data", CollectionOf(RtspResult), "得到的分页数据", func() {
			//Example(func() {
			//	//Value(CollectionOf(RtspResult.View("tiny")).Example(expr.NewRandom("")))
			//})
			Meta("rpc:tag", "2")
		})
	})

	View("default", func() {
		Attribute("count")
		Attribute("page_data")
	})
})

// RtspResult is the rtsp resource media type.
var RtspResult = ResultType("application/vnd.rtsp_result+json", func() {
	Description("RTSP数据信息")
	Attributes(func() {
		Field(1, "id", UInt64, "数据ID")
		Field(2, "rtsp_url", String, "rtsp地址")
		Field(3, "host_id", String, "主机ID")
		Field(4, "ext", String, "扩展信息")
		Field(5, "created_at", String, "创建时间")
		Field(6, "updated_at", String, "更新时间")
		Field(7, "host", HostInfo, "主机信息")
		Field(8, "play_url", String, "播放的地址")
	})

	View("default", func() {
		Field(1, "id")
		Field(2, "rtsp_url")
		Field(3, "host_id")
		Field(4, "ext")
		Field(5, "created_at")
		Field(6, "updated_at")
		Field(7, "host")
	})

	View("tiny", func() {
		Field(1, "id")
		Field(2, "rtsp_url")
		Field(3, "host_id")
		Field(5, "created_at")
		Field(6, "updated_at")
	})

	View("full", func() {
		Field(1, "id")
		Field(2, "rtsp_url")
		Field(3, "host_id")
		Field(4, "ext")
		Field(5, "created_at")
		Field(6, "updated_at")
		Field(7, "host")
	})

})

// BasicAuth defines a security scheme using basic authentication. The scheme protects the "signin"
// action used to create JWTs.
//var BasicAuth = BasicAuthSecurity("BasicAuth")
var res_rtsp = Service("rtsp", func() {
	Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWTAuth claims.
	})
	HTTP(func() {
		Path("/themis/v1/rtsp")
		Headers(func() {
			//Header("Authorization", String, "JWTAuth token", func() {
			//})
			//Required("Authorization")
		})
	})

	Error("NotFound", String, "未查询到数据！")
	Error("Unauthorized", String, "未授权！")
	Error("InvalidAccountNotFound", String, "用户不存在，请重试！")
	Error("InvalidAccountOrPassword", String, "用户名或密码错误！")
	Error("AuthorizedFailure", String, "授权失败！")

	GRPC(func() {
		Response("Unauthorized", CodeUnauthenticated)
	})

	Method("list", func() {
		Description("流的数据列表；")
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
			Field(6, "where", Rtsp, "条件", func() {})
			Required("offset_head", "offset_tail")
		})
		Error("Unauthorized")
		Error("BadRequest")
		Error("NotFound")
		Result(PageModelRtsp)
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
		Description("创建RTSP数据")
		Payload(Rtsp)
		Error("Unauthorized")
		Result(RtspResult)
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
		Payload(Rtsp)
		Error("Unauthorized")
		Result(RtspResult)
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
		Result(RtspResult)
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
