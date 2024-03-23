package config

type App struct {
	Port string `env:"CONFIG.APP.PORT"`
}
