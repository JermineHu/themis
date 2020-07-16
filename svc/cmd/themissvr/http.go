package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	admin "github.com/JermineHu/themis/svc/gen/admin"
	config "github.com/JermineHu/themis/svc/gen/config"
	health "github.com/JermineHu/themis/svc/gen/health"
	host "github.com/JermineHu/themis/svc/gen/host"
	adminsvr "github.com/JermineHu/themis/svc/gen/http/admin/server"
	configsvr "github.com/JermineHu/themis/svc/gen/http/config/server"
	healthsvr "github.com/JermineHu/themis/svc/gen/http/health/server"
	hostsvr "github.com/JermineHu/themis/svc/gen/http/host/server"
	keyboardsvr "github.com/JermineHu/themis/svc/gen/http/keyboard/server"
	noticesvr "github.com/JermineHu/themis/svc/gen/http/notice/server"
	rtspsvr "github.com/JermineHu/themis/svc/gen/http/rtsp/server"
	tokenmgrsvr "github.com/JermineHu/themis/svc/gen/http/token_mgr/server"
	keyboard "github.com/JermineHu/themis/svc/gen/keyboard"
	notice "github.com/JermineHu/themis/svc/gen/notice"
	rtsp "github.com/JermineHu/themis/svc/gen/rtsp"
	tokenmgr "github.com/JermineHu/themis/svc/gen/token_mgr"
	goahttp "goa.design/goa/v3/http"
	httpmdlwr "goa.design/goa/v3/http/middleware"
	"goa.design/goa/v3/middleware"
)

// handleHTTPServer starts configures and starts a HTTP server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleHTTPServer(ctx context.Context, u *url.URL, healthEndpoints *health.Endpoints, adminEndpoints *admin.Endpoints, configEndpoints *config.Endpoints, keyboardEndpoints *keyboard.Endpoints, hostEndpoints *host.Endpoints, noticeEndpoints *notice.Endpoints, rtspEndpoints *rtsp.Endpoints, tokenMgrEndpoints *tokenmgr.Endpoints, wg *sync.WaitGroup, errc chan error, logger *log.Logger, debug bool) {

	// Setup goa log adapter.
	var (
		adapter middleware.Logger
	)
	{
		adapter = middleware.NewLogger(logger)
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and configure it to serve
	// HTTP requests to the service endpoints.
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	var (
		healthServer   *healthsvr.Server
		adminServer    *adminsvr.Server
		configServer   *configsvr.Server
		keyboardServer *keyboardsvr.Server
		hostServer     *hostsvr.Server
		noticeServer   *noticesvr.Server
		rtspServer     *rtspsvr.Server
		tokenMgrServer *tokenmgrsvr.Server
	)
	{
		eh := errorHandler(logger)
		healthServer = healthsvr.New(healthEndpoints, mux, dec, enc, eh, nil)
		adminServer = adminsvr.New(adminEndpoints, mux, dec, enc, eh, nil)
		configServer = configsvr.New(configEndpoints, mux, dec, enc, eh, nil)
		keyboardServer = keyboardsvr.New(keyboardEndpoints, mux, dec, enc, eh, nil)
		hostServer = hostsvr.New(hostEndpoints, mux, dec, enc, eh, nil)
		noticeServer = noticesvr.New(noticeEndpoints, mux, dec, enc, eh, nil)
		rtspServer = rtspsvr.New(rtspEndpoints, mux, dec, enc, eh, nil)
		tokenMgrServer = tokenmgrsvr.New(tokenMgrEndpoints, mux, dec, enc, eh, nil)
		if debug {
			servers := goahttp.Servers{
				healthServer,
				adminServer,
				configServer,
				keyboardServer,
				hostServer,
				noticeServer,
				rtspServer,
				tokenMgrServer,
			}
			servers.Use(httpmdlwr.Debug(mux, os.Stdout))
		}
	}
	// Configure the mux.
	healthsvr.Mount(mux, healthServer)
	adminsvr.Mount(mux, adminServer)
	configsvr.Mount(mux, configServer)
	keyboardsvr.Mount(mux, keyboardServer)
	hostsvr.Mount(mux, hostServer)
	noticesvr.Mount(mux, noticeServer)
	rtspsvr.Mount(mux, rtspServer)
	tokenmgrsvr.Mount(mux, tokenMgrServer)

	// Wrap the multiplexer with additional middlewares. Middlewares mounted
	// here apply to all the service endpoints.
	var handler http.Handler = mux
	{
		handler = httpmdlwr.Log(adapter)(handler)
		handler = httpmdlwr.RequestID()(handler)
	}

	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: u.Host, Handler: handler}
	for _, m := range healthServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range adminServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range configServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range keyboardServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range hostServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range noticeServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range rtspServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range tokenMgrServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start HTTP server in a separate goroutine.
		go func() {
			logger.Printf("HTTP server listening on %q", u.Host)
			errc <- srv.ListenAndServe()
		}()

		<-ctx.Done()
		logger.Printf("shutting down HTTP server at %q", u.Host)

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.Shutdown(ctx)
	}()
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler(logger *log.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(middleware.RequestIDKey).(string)
		w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		logger.Printf("[%s] ERROR: %s", id, err.Error())
	}
}
