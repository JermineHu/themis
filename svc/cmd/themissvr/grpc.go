package main

import (
	"context"
	"log"
	"net"
	"net/url"
	"sync"

	admin "github.com/JermineHu/themis/svc/gen/admin"
	config "github.com/JermineHu/themis/svc/gen/config"
	adminpb "github.com/JermineHu/themis/svc/gen/grpc/admin/pb"
	adminsvr "github.com/JermineHu/themis/svc/gen/grpc/admin/server"
	configpb "github.com/JermineHu/themis/svc/gen/grpc/config/pb"
	configsvr "github.com/JermineHu/themis/svc/gen/grpc/config/server"
	healthpb "github.com/JermineHu/themis/svc/gen/grpc/health/pb"
	healthsvr "github.com/JermineHu/themis/svc/gen/grpc/health/server"
	hostpb "github.com/JermineHu/themis/svc/gen/grpc/host/pb"
	hostsvr "github.com/JermineHu/themis/svc/gen/grpc/host/server"
	keyboardpb "github.com/JermineHu/themis/svc/gen/grpc/keyboard/pb"
	keyboardsvr "github.com/JermineHu/themis/svc/gen/grpc/keyboard/server"
	noticepb "github.com/JermineHu/themis/svc/gen/grpc/notice/pb"
	noticesvr "github.com/JermineHu/themis/svc/gen/grpc/notice/server"
	rtsppb "github.com/JermineHu/themis/svc/gen/grpc/rtsp/pb"
	rtspsvr "github.com/JermineHu/themis/svc/gen/grpc/rtsp/server"
	tokenpb "github.com/JermineHu/themis/svc/gen/grpc/token/pb"
	tokensvr "github.com/JermineHu/themis/svc/gen/grpc/token/server"
	health "github.com/JermineHu/themis/svc/gen/health"
	host "github.com/JermineHu/themis/svc/gen/host"
	keyboard "github.com/JermineHu/themis/svc/gen/keyboard"
	notice "github.com/JermineHu/themis/svc/gen/notice"
	rtsp "github.com/JermineHu/themis/svc/gen/rtsp"
	token "github.com/JermineHu/themis/svc/gen/token"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcmdlwr "goa.design/goa/v3/grpc/middleware"
	"goa.design/goa/v3/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// handleGRPCServer starts configures and starts a gRPC server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleGRPCServer(ctx context.Context, u *url.URL, healthEndpoints *health.Endpoints, adminEndpoints *admin.Endpoints, configEndpoints *config.Endpoints, keyboardEndpoints *keyboard.Endpoints, hostEndpoints *host.Endpoints, noticeEndpoints *notice.Endpoints, rtspEndpoints *rtsp.Endpoints, tokenEndpoints *token.Endpoints, wg *sync.WaitGroup, errc chan error, logger *log.Logger, debug bool) {

	// Setup goa log adapter.
	var (
		adapter middleware.Logger
	)
	{
		adapter = middleware.NewLogger(logger)
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to gRPC requests and
	// responses.
	var (
		healthServer   *healthsvr.Server
		adminServer    *adminsvr.Server
		configServer   *configsvr.Server
		keyboardServer *keyboardsvr.Server
		hostServer     *hostsvr.Server
		noticeServer   *noticesvr.Server
		rtspServer     *rtspsvr.Server
		tokenServer    *tokensvr.Server
	)
	{
		healthServer = healthsvr.New(healthEndpoints, nil)
		adminServer = adminsvr.New(adminEndpoints, nil)
		configServer = configsvr.New(configEndpoints, nil)
		keyboardServer = keyboardsvr.New(keyboardEndpoints, nil)
		hostServer = hostsvr.New(hostEndpoints, nil)
		noticeServer = noticesvr.New(noticeEndpoints, nil)
		rtspServer = rtspsvr.New(rtspEndpoints, nil)
		tokenServer = tokensvr.New(tokenEndpoints, nil)
	}

	// Initialize gRPC server with the middleware.
	srv := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			grpcmdlwr.UnaryRequestID(),
			grpcmdlwr.UnaryServerLog(adapter),
		),
	)

	// Register the servers.
	healthpb.RegisterHealthServer(srv, healthServer)
	adminpb.RegisterAdminServer(srv, adminServer)
	configpb.RegisterConfigServer(srv, configServer)
	keyboardpb.RegisterKeyboardServer(srv, keyboardServer)
	hostpb.RegisterHostServer(srv, hostServer)
	noticepb.RegisterNoticeServer(srv, noticeServer)
	rtsppb.RegisterRtspServer(srv, rtspServer)
	tokenpb.RegisterTokenServer(srv, tokenServer)

	for svc, info := range srv.GetServiceInfo() {
		for _, m := range info.Methods {
			logger.Printf("serving gRPC method %s", svc+"/"+m.Name)
		}
	}

	// Register the server reflection service on the server.
	// See https://grpc.github.io/grpc/core/md_doc_server-reflection.html.
	reflection.Register(srv)

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start gRPC server in a separate goroutine.
		go func() {
			lis, err := net.Listen("tcp", u.Host)
			if err != nil {
				errc <- err
			}
			logger.Printf("gRPC server listening on %q", u.Host)
			errc <- srv.Serve(lis)
		}()

		<-ctx.Done()
		logger.Printf("shutting down gRPC server at %q", u.Host)
		srv.Stop()
	}()
}
