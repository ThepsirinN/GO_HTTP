package main

import (
	"context"
	"fmt"
	"go_http_barko/config"
	logger "go_http_barko/utility/logger"
	handlersV1 "go_http_barko/v1/handlers"
	repositoriesV1 "go_http_barko/v1/repositories"
	servicesV1 "go_http_barko/v1/services"
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

	// repositoriesV1
	repositoriesV1 := repositoriesV1.New()

	// servicesV1
	servicesV1 := servicesV1.New(repositoriesV1)

	// handlersV1
	handlersV1 := handlersV1.New(ctx, servicesV1)

	router := initRounter(handlersV1)
	server := &http.Server{
		Addr: cfg.App.Port,
	}

	go httpServe(ctx, server, router)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)
	<-quit

	fmt.Println("recieve signal: shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}

}

func httpServe(ctx context.Context, server *http.Server, router *http.ServeMux) {
	server.Handler = router
	logger.Info(ctx, "========== Server is starting ==========")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
