package config

import (
	"context"
	"log"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	App App
	Log Log
}

var config Config
var syncOnce Once

func InitConfig(ctx context.Context) *Config {
	syncOnce.Do(initConfig, ctx)
	return &config
}

func initConfig(ctx context.Context) {
	absPath, _ := filepath.Abs("./env/config.env")
	if err := godotenv.Load(absPath); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := env.Parse(&config); err != nil {
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
