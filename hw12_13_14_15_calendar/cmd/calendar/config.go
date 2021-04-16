package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	storageSql      = "sql"
	storageInMemory = "in-memory"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf  `json:"logger"`
	Server  ServerConf  `json:"server"`
	Storage StorageConf `json:"storage"`
}

type LoggerConf struct {
	Level    int8   `json:"level"`
	FilePath string `json:"file_path"`
}

type ServerConf struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type StorageConf struct {
	Inmemory bool
	Database DBConf `json:"database"`
	Type     string `json:"type"`
}

type DBConf struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DB       string `json:"db"`
}

func NewConfig(configFile string) (Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, fmt.Errorf("failed open config file: %w", err)
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("failed read config: %w", err)
	}

	switch config.Storage.Type {
	case storageInMemory:
		config.Storage.Inmemory = true
	case storageSql:
		config.Storage.Inmemory = false
	default:
		// По умолчанию создаем хранилище в памяти
		config.Storage.Inmemory = true
	}

	return config, nil
}

// Собирает строку DSN.
func DSN(config DBConf) string {
	var c strings.Builder
	c.WriteString(config.Username)
	c.WriteString(":")
	c.WriteString(config.Password)
	c.WriteString("@")
	c.WriteString("tcp(")
	c.WriteString(config.Host)

	if config.Port != "" {
		c.WriteString(":")
		c.WriteString(config.Port)
	}
	c.WriteString(")/")
	c.WriteString(config.DB)

	return c.String()
}
