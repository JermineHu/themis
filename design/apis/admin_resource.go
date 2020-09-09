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

var Admin = Type("Admin", func() {

	Field(1, "id", UInt64, "用户的id", func() {

	})
	Field(2, "login_name", String, "账户登陆名称", func() {
		Example("ZhangSan")
		//Meta("rpc:tag", "2")
	})
	//Field(3, "first_name", String, "用户名", func() {
	//	Example("ZhangSan")
	//	//Meta("rpc:tag", "2")
	//})
	//Field(4, "second_name", String, "用户姓", func() {
	//	Example("ZhangSan")
	//	//Meta("rpc:tag", "2")
	//})
	//Field(5, "description", String, "用户描述", func() {
	//	Example("ZhangSan")
	//	//Meta("rpc:tag", "2")
	//})
	//
	//Field(6, "mobile", String, "手机号码", func() {
	//	Example("13811123456")
	//	//Meta("rpc:tag", "3")
	//})
	Field(7, "password", String, "用户密码", func() {
		Example("123456")
		//Meta("rpc:tag", "4")
	})

	Field(8, "salt", String, "生成密码是使用的盐", func() {
		Example(RandomString(13))
		//Meta("rpc:tag", "5")
	})

	Field(9, "user_type", func() {
		Enum("USER_TYPE_NORMAL", "USER_TYPE_ADMIN")
		Example("USER_TYPE_NORMAL")
		//Meta("rpc:tag", "6")
	})

	//Field(10, "user_status", func() {
	//	Enum("USER_STATUS_NULL", "USER_STATUS_ACTIVATED", "USER_STATUS_DEACTIVATED", "USER_STATUS_CLOSED")
	//	Example("USER_STATUS_ACTIVATED")
	//	//Meta("rpc:tag", "7")
	//})

	Field(12, "created_by", String, "创建该账户的账户", func() {
		Format(FormatUUID)
		Example(NewUUIDStr())
		//Meta("rpc:tag", "9")
	})
	TokenField(13, "token", String, "JWTAuth token used to perform authorization", func() {
		//Meta("rpc:tag", "10")
	})
	//Field(14, "sex", String, "用户性别")
	//Field(15, "nickname", String, "用户昵称")
	//Field(16, "avatar_url", String, "头像")
	Field(17, "created_at", String, "创建时间")
	Field(18, "updated_at", String, "更新时间")
	//Required("id", "password", "salt", "name", "created_at", "created_by", "user_type", "user_status", "token")

})

var AdminLogin = ResultType("application/vnd.admin.login+json", func() {
	Description("管理员模型")
	Reference(Admin)
	Attributes(func() {
		Field(1, "password")
		Field(2, "login_name")
		Field(3, "token")
	})

	View("default", func() {
		Field(1, "password")
		Field(2, "login_name")
		Field(3, "token")
	})
	//Required("password", "mobile")
})

var PageModelAdmin = ResultType("application/vnd.admin_list+json", func() {
	Description("分页返回是数据模型")
	Attributes(func() {
		Attribute("count", Int64, "数据条数", func() {
			Example(200)
			Meta("rpc:tag", "1")
		})
		Attribute("page_data", CollectionOf(AdminResult), "得到的分页数据", func() {
			//Example(func() {
			//	//Value(CollectionOf(AdminResult.View("tiny")).Example(expr.NewRandom("")))
			//})
			Meta("rpc:tag", "2")
		})
	})

	View("default", func() {
		Attribute("count")
		Attribute("page_data")
	})
})

var AdminResult = ResultType("application/vnd.admin_result+json", func() {
	Description("管理员模型")
	Reference(Admin)
	Attributes(func() {
		Field(2, "id")
		Field(3, "password")
		Field(4, "salt")
		Field(5, "login_name")
		Field(6, "created_at")
		Field(7, "created_by")
		Field(8, "user_type")
		//Field(9, "user_status")
		Field(10, "token")
		//Required("id", "password", "salt", "name", "created_at", "created_by", "user_type", "user_status")
	})

	View("default", func() {
		Field(1, "id")
		Field(2, "login_name")
		Field(3, "created_at")
		Field(4, "created_by")
		Field(6, "user_type")
		//Field(7, "user_status")
		Field(9, "token")
	})

	View("tiny", func() {
		Field(1, "id")
		Field(2, "login_name")
		Field(3, "created_at")
		Field(4, "created_by")
		Field(5, "user_type")
		//Field(6, "user_status")
		Field(7, "token")

	})

	View("full", func() {
		Field(1, "id")
		Field(2, "password")
		Field(3, "salt")
		Field(4, "login_name")
		Field(5, "created_at")
		Field(6, "created_by")
		Field(8, "user_type")
		//Field(9, "user_status")
		Field(11, "token")
	})

})

// BasicAuth defines a security scheme using basic authentication. The scheme protects the "signin"
// action used to create JWTs.
//var BasicAuth = BasicAuthSecurity("BasicAuth")
var res_admin = Service("admin", func() {
	Security(JWTAuth, func() { // Use JWTAuth to auth requests to this endpoint
		Scope("api:write") // Enforce presence of "api" scope in JWTAuth claims.
	})
	HTTP(func() {
		Path("/themis/v1/admin")
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
	Error("invalid-scopes", String, "令牌授权范围错误！")
	Error("NotFound", String, "未查询到数据！")

	GRPC(func() {
		Response("Unauthorized", CodeUnauthenticated)
	})

	Method("login", func() {
		NoSecurity()
		Payload(func() {
			Field(1, "username", String, func() {
				Description("管理员账户名称")
			})

			Field(2, "password", String, func() {
				Description("管理员账户密码")
			})
			Required("username", "password")
		})

		Description("根据账户名称密码进行登陆操作！")
		Error("InvalidAccountNotFound")
		Error("InvalidAccountOrPassword")

		Result(AdminResult, func() {
			View("tiny")
		})
		HTTP(func() {
			POST("/login")
			Response(StatusOK, func() {
				Header("token:Authorization")
			})
			Response(StatusNotFound)
			Response("AuthorizedFailure", StatusForbidden)
			Response("InvalidAccountNotFound", StatusForbidden)
			Response("InvalidAccountOrPassword", StatusForbidden)
		})

		GRPC(func() {
			Response(CodeOK)
			Response(CodeNotFound)
		})
	})

	Method("logout", func() {
		Payload(func() {
			TokenField(1, "token", String, func() {
				Description("JWTAuth used for authentication")
			})
			Required("token")
		})

		Description("退出登陆（将发上来的JWT加入拒绝名单，并且设置过期时间为JWT到期时间，到期自动释放），删除本地JWT")
		Error("Unauthorized")
		Error("BadRequest")
		Error("NotFound")
		Result(Boolean)
		HTTP(func() {
			GET("/logout")
			Response(StatusOK)
			Response(StatusNotFound)
		})

		GRPC(func() {
			Response(CodeOK)
			Response(CodeNotFound)
		})
	})

	Method("list", func() {
		Description("列表数据；")
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
			Field(6, "where", Admin, "条件", func() {})
			Required("offset_head", "offset_tail")
		})
		Error("Unauthorized")
		Error("BadRequest")
		Error("NotFound")
		Result(PageModelAdmin)
		HTTP(func() {
			POST("/list")
			Response(StatusOK)
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
		Payload(Admin)
		Error("Unauthorized")
		Result(AdminResult)
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
		Description("根据id修数据")
		Payload(Admin)
		Error("Unauthorized")
		Result(AdminResult)
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
		Description("根据id信息")
		Error("Unauthorized")
		Payload(func() {
			TokenField(1, "token", String, "JWTAuth token used to perform authorization", func() {
			})
			Field(2, "id", UInt64, "id", func() {
			})

		})
		Result(AdminResult)
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
