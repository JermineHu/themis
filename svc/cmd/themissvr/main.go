package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/JermineHu/themis/common"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	themis "github.com/JermineHu/themis/svc"
	admin "github.com/JermineHu/themis/svc/gen/admin"
	config "github.com/JermineHu/themis/svc/gen/config"
	health "github.com/JermineHu/themis/svc/gen/health"
	host "github.com/JermineHu/themis/svc/gen/host"
	keyboard "github.com/JermineHu/themis/svc/gen/keyboard"
	notice "github.com/JermineHu/themis/svc/gen/notice"
	rtsp "github.com/JermineHu/themis/svc/gen/rtsp"
	tokenmgr "github.com/JermineHu/themis/svc/gen/token_mgr"
)

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF     = flag.String("host", "development", "Server host (valid values: development, production)")
		domainF   = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		httpPortF = flag.String("http-port", os.Getenv(common.APP_PORT), "HTTP port (overrides host HTTP port specified in service design)")
		grpcPortF = flag.String("grpc-port", "7080", "gRPC port (overrides host gRPC port specified in service design)")
		versionF  = flag.String("version", "v1", "API version")
		secureF   = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF      = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	// Setup logger. Replace logger with your own log package of choice.
	var (
		logger *log.Logger
	)
	{
		logger = log.New(os.Stderr, "[themis] ", log.Ltime)
	}

	// Initialize the services.
	var (
		healthSvc   health.Service
		adminSvc    admin.Service
		configSvc   config.Service
		keyboardSvc keyboard.Service
		hostSvc     host.Service
		noticeSvc   notice.Service
		rtspSvc     rtsp.Service
		tokenMgrSvc tokenmgr.Service
	)
	{
		healthSvc = themis.NewHealth(logger)
		adminSvc = themis.NewAdmin(logger)
		configSvc = themis.NewConfig(logger)
		keyboardSvc = themis.NewKeyboard(logger)
		hostSvc = themis.NewHost(logger)
		noticeSvc = themis.NewNotice(logger)
		rtspSvc = themis.NewRtsp(logger)
		tokenMgrSvc = themis.NewTokenMgr(logger)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		healthEndpoints   *health.Endpoints
		adminEndpoints    *admin.Endpoints
		configEndpoints   *config.Endpoints
		keyboardEndpoints *keyboard.Endpoints
		hostEndpoints     *host.Endpoints
		noticeEndpoints   *notice.Endpoints
		rtspEndpoints     *rtsp.Endpoints
		tokenMgrEndpoints *tokenmgr.Endpoints
	)
	{
		healthEndpoints = health.NewEndpoints(healthSvc)
		adminEndpoints = admin.NewEndpoints(adminSvc)
		configEndpoints = config.NewEndpoints(configSvc)
		keyboardEndpoints = keyboard.NewEndpoints(keyboardSvc)
		hostEndpoints = host.NewEndpoints(hostSvc)
		noticeEndpoints = notice.NewEndpoints(noticeSvc)
		rtspEndpoints = rtsp.NewEndpoints(rtspSvc)
		tokenMgrEndpoints = tokenmgr.NewEndpoints(tokenMgrSvc)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "development":
		{
			addr := "http://:8081/themis/v1"
			u, err := url.Parse(addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid URL %#v: %s\n", addr, err)
				os.Exit(1)
			}
			if *secureF {
				u.Scheme = "https"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *httpPortF != "" {
				h := strings.Split(u.Host, ":")[0]
				u.Host = h + ":" + *httpPortF
			} else if u.Port() == "" {
				u.Host += ":80"
			}
			handleHTTPServer(ctx, u, healthEndpoints, adminEndpoints, configEndpoints, keyboardEndpoints, hostEndpoints, noticeEndpoints, rtspEndpoints, tokenMgrEndpoints, &wg, errc, logger, *dbgF)
		}

		{
			addr := "grpc://:8080"
			u, err := url.Parse(addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid URL %#v: %s\n", addr, err)
				os.Exit(1)
			}
			if *secureF {
				u.Scheme = "grpcs"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *grpcPortF != "" {
				h := strings.Split(u.Host, ":")[0]
				u.Host = h + ":" + *grpcPortF
			} else if u.Port() == "" {
				u.Host += ":8080"
			}
			handleGRPCServer(ctx, u, healthEndpoints, adminEndpoints, configEndpoints, keyboardEndpoints, hostEndpoints, noticeEndpoints, rtspEndpoints, tokenMgrEndpoints, &wg, errc, logger, *dbgF)
		}

	case "production":
		{
			addr := "https://{version}.themis.vdo.pub"
			addr = strings.Replace(addr, "{version}", *versionF, -1)
			u, err := url.Parse(addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid URL %#v: %s\n", addr, err)
				os.Exit(1)
			}
			if *secureF {
				u.Scheme = "https"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *httpPortF != "" {
				h := strings.Split(u.Host, ":")[0]
				u.Host = h + ":" + *httpPortF
			} else if u.Port() == "" {
				u.Host += ":443"
			}
			handleHTTPServer(ctx, u, healthEndpoints, adminEndpoints, configEndpoints, keyboardEndpoints, hostEndpoints, noticeEndpoints, rtspEndpoints, tokenMgrEndpoints, &wg, errc, logger, *dbgF)
		}

		{
			addr := "grpcs://{version}.themis.vdo.pub"
			addr = strings.Replace(addr, "{version}", *versionF, -1)
			u, err := url.Parse(addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid URL %#v: %s\n", addr, err)
				os.Exit(1)
			}
			if *secureF {
				u.Scheme = "grpcs"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *grpcPortF != "" {
				h := strings.Split(u.Host, ":")[0]
				u.Host = h + ":" + *grpcPortF
			} else if u.Port() == "" {
				u.Host += ":8443"
			}
			handleGRPCServer(ctx, u, healthEndpoints, adminEndpoints, configEndpoints, keyboardEndpoints, hostEndpoints, noticeEndpoints, rtspEndpoints, tokenMgrEndpoints, &wg, errc, logger, *dbgF)
		}

	default:
		fmt.Fprintf(os.Stderr, "invalid host argument: %q (valid hosts: development|production)\n", *hostF)
	}

	// Wait for signal.
	logger.Printf("exiting (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Println("exited")
}
