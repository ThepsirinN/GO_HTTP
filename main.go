package main

import (
	"context"
	"database/sql"
	"fmt"
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

	"github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()
	cfg, secret := config.InitConfig(ctx)
	logger.InitLogger(cfg)
	defer logger.Sync()

	tp := tracer.InitTraceProvider(ctx, cfg.Log.Env, "test123")
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Fatal(ctx, err)
		}
	}()

	sqlDb := initDatabase(ctx, cfg, secret)

	// repositoriesV1
	repositoriesV1 := repositoriesV1.New(sqlDb)

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
		return
	}
}

func initDatabase(ctx context.Context, cfg *config.Config, secret *config.Secret) *sql.DB {
	dbCfg := mysql.Config{
		Net:       "tcp",
		Addr:      fmt.Sprintf("%v:%v", secret.Database.Host, secret.Database.Port),
		User:      secret.Database.User,
		Passwd:    secret.Database.Password,
		DBName:    cfg.Database.Name,
		ParseTime: true,
		Loc:       time.Local,
	}

	sqlDb, err := sql.Open("mysql", dbCfg.FormatDSN())
	if err != nil {
		logger.Fatal(ctx, err)
		return nil
	}

	sqlDb.SetConnMaxLifetime(cfg.Database.ConnMaxLifeTime)
	sqlDb.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDb.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	return sqlDb
}
