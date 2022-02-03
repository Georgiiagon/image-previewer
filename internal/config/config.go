package config

import (
	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	App      AppConf      `toml:"app"`
	Logger   LoggerConf   `toml:"logger"`
	Database DatabaseConf `toml:"database"`
}

type AppConf struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type LoggerConf struct {
	Level string `toml:"level"`
}

func NewConfig(filePath string) Config {
	config := Config{}

	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		panic(err)
	}

	return config
}
