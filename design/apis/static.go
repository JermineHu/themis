package apis

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("statics", func() {
	Description("静态文件")
	HTTP(func() {
		Path("/")
	})

	Files("/{*path}", "ui", func() {
		Description("JSON document containing the API swagger definition")
	})

})
