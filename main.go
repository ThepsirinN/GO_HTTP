package main

import (
	"context"
	"go_http_barko/config"
	logger "go_http_barko/utility/logger"
	"go_http_barko/utility/tracer"
	handlersV1 "go_http_barko/v1/handlers"
	repositoriesV1 "go_http_barko/v1/repositories"
	servicesV1 "go_http_barko/v1/services"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	cfg := config.InitConfig(ctx)
	logger.InitLogger(cfg)
	defer logger.Sync()

	tp := tracer.InitTraceProvider(ctx, cfg.Log.Env, "test123")
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Fatal(ctx, err)
		}
	}()

	// repositoriesV1
	repositoriesV1 := repositoriesV1.New()

	// servicesV1
	servicesV1 := servicesV1.New(repositoriesV1)

	// handlersV1
	handlersV1 := handlersV1.New(servicesV1)

	router := initRounter(handlersV1)
	server := &http.Server{
		Addr:         cfg.App.Port,
		BaseContext:  func(net.Listener) context.Context { return ctx }, // very important for tracing
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go httpServe(ctx, server, router)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)
	<-quit

	logger.Info(ctx, "recieve signal: shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(ctx, err)
	}

}

func httpServe(ctx context.Context, server *http.Server, router http.Handler) {
	server.Handler = router
	logger.Info(ctx, "========== Server is starting ==========")
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal(ctx, err)
	}
}
