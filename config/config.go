package config

import (
	"context"
	"go_http_barko/constant"
	"log"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	App      App
	Log      Log
	Database DatabaseConfig
}

type Secret struct {
	Database DatabaseSecret
}

var config Config
var secret Secret
var syncOnce Once

func InitConfig(ctx context.Context) (*Config, *Secret) {
	syncOnce.Do(initConfig, ctx)
	return &config, &secret
}

func initConfig(ctx context.Context) {
	configPath, _ := filepath.Abs(constant.CONFIG_PATH + constant.CONFIG_FILE)
	secretPath, _ := filepath.Abs(constant.SECRET_PATH + constant.SECRET_FILE)
	if err := godotenv.Load(configPath, secretPath); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := env.Parse(&config); err != nil {
		log.Fatalf("%+v\n", err)
	}

	if err := env.Parse(&secret); err != nil {
		log.Fatalf("%+v\n", err)
	}

}

type Once struct {
	done atomic.Uint32
	m    sync.Mutex
}

func (o *Once) Do(f func(context.Context), ctx context.Context) {
	if o.done.Load() == 0 {
		o.doSlow(f, ctx)
	}
}

func (o *Once) doSlow(f func(context.Context), ctx context.Context) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done.Load() == 0 {
		defer o.done.Store(1)
		f(ctx)
	}
}
