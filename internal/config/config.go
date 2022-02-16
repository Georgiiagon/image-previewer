package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	App   AppConf
	Cache CacheConf
}

type AppConf struct {
	Host string
	Port string
}

type CacheConf struct {
	Length int
}

func New() Config {
	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file")
	}

	CacheLength, err := strconv.Atoi(os.Getenv("CACHE_LENGTH"))
	if err != nil {
		log.Panic("Error cache length env")
	}

	config := Config{
		App: AppConf{
			Host: os.Getenv("HOST"),
			Port: os.Getenv("PORT"),
		},
		Cache: CacheConf{
			Length: CacheLength,
		},
	}

	return config
}
