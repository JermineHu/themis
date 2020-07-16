package apis

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("health", func() {
	HTTP(func() {
		Path("/themis/_ah")
	})
	Method("health", func() {
		Description("Perform health check.")
		HTTP(func() {
			GET("/health")
			Response(StatusOK)
		})
		GRPC(func() {
			Response(CodeOK)
		})

	})
	Method("build_info", func() {
		Description("version info.")
		Result(func() {
			Field(1, "git_hash", String, "编译时所在的tag版本", func() {

			})
			Field(2, "git_log", String, "编译时git的提交日志", func() {

			})
			Field(3, "git_status", String, "编译时git当前状态", func() {

			})
			Field(4, "build_time", String, "编译时间", func() {

			})
			Field(5, "go_version", String, "编译时所使用的go版本", func() {

			})
			Field(6, "go_runtime", String, "go的运行时", func() {
			})
		})
		HTTP(func() {
			GET("/build_info")
			Response(StatusOK)
		})
		GRPC(func() {
			Response(CodeOK)
		})

	})
})
