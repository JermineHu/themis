package apis

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
)

// This is the cellar application API design used by goa to generate
// the application code, client, tests, documentation etc.
var _ = API("themis", func() {
	Title("The themis REST API swagger document")
	Description("themis项目用于提供REST API 服务，主要包含以下几个模块的管理功能！")
	Contact(func() {
		Name("themis")
		Email("jermine.hu@qq.com")
		URL("https://jermine.vdo.pub/")
	})
	//License(func() {
	//	Name("MIT")
	//	URL("https://github.com/goadesign/goa/blob/master/LICENSE")
	//})
	Docs(func() {
		Description("themis guide")
		URL("http://10.0.52.100:10080/jermine/themis")
	})
	TermsOfService("terms") // Terms of use
	Server("themissvr", func() {
		Description("calcsvr hosts the Calculator Service.")
		// List the services hosted by this server.
		Services("health", "admin", "config", "keyboard", "host", "notice", "rtsp", "token_mgr")
		Host("development", func() {
			Description("Development hosts.")
			// Transport specific URLs, supported schemes are:
			// 'http', 'https', 'grpc' and 'grpcs' with the respective default
			// ports: 80, 443, 8080, 8443.
			URI("http://:8081/themis/v1")
			URI("grpc://:8080")

		})

		// List the Hosts and their transport URLs.
		Host("production", func() {
			Description("Production host.")
			// URIs can be parameterized using {param} notation.
			URI("https://{version}.themis.vdo.pub")
			URI("grpcs://{version}.themis.vdo.pub")

			// Variable describes a URI variable.
			Variable("version", String, "API version", func() {
				// URI parameters must have a default value and/or an
				// enum validation.
				Default("v1")
			})
		})

	})
	cors.Origin("*", func() { // Define CORS policy, may be prefixed with "*" wildcard
		cors.Methods("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS")                                                                                                                                  // One or more authorized HTTP methods
		cors.Headers("Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type,Content-Disposition,Content-Transfer-Encoding, Authorization,*") // One or more authorized headers, use "*" to authorize all
		cors.Expose("Authorization", "token", "Content-Disposition", "Content-Transfer-Encoding", "File-Name", "CheckOrigin")                                                                             // One or more headers exposed to clients
		cors.MaxAge(600)                                                                                                                                                                                  // How long to cache a preflight request response
		cors.Credentials()                                                                                                                                                                                // Sets Access-Control-Allow-Credentials header
	})
	Version("v2")
	HTTP(
		func() {
			Consumes("application/json", "application/xml", "text/xml")
			Produces("application/json", "application/xml")
		})
	//	Consumes("application/json") // Media types supported by the API
	//	Produces("application/json") // Media types generated by the API

	//Origin("https://*.ums86.com", func() {
	//	Methods("GET", "POST", "PUT", "PATCH", "DELETE")
	//	Headers("Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With, X-Tuso-Device-Token, X-Tuso-Authentication-Token, *")
	//	MaxAge(600)
	//	Credentials()
	//})

	//Origin("*", func() {
	//	Methods("GET", "POST", "PUT", "PATCH", "DELETE")
	//	Headers("Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, Authorization,*")
	//	MaxAge(600)
	//	Credentials()
	//})

	//ResponseTemplate(Created, func(pattern string) {
	//	Description("Resource created")
	//	Status(201)
	//	Headers(func() {
	//		Header("Location", String, "href to created resource", func() {
	//			Pattern(pattern)
	//		})
	//	})
	//})
})
