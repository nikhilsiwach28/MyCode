package config

import (
	"log"
	"os"
	"strconv"
)

type ServerConfig struct {
	address      string
	port         string
	readTimeout  int
	writeTimeout int
}

func (c *ServerConfig) GetPort() string {
	return c.port
}

func (c ServerConfig) GetAddress() string {
	return c.address
}

func (c ServerConfig) GetReadTimeout() int {
	return c.readTimeout
}

func (c ServerConfig) GetWriteTimeout() int {
	return c.writeTimeout
}

func NewServerConfig() ServerConfig {
	rTimeout, errR := strconv.Atoi(getEnvWithDefault("READ_TIMEOUT", "5"))
	wTimeout, errW := strconv.Atoi(getEnvWithDefault("WRITE_TIMEOUT", "5"))

	if errR != nil || errW != nil {
		log.Fatal(errR, errW)
	}

	return ServerConfig{
		address:      getEnvWithDefault("SERVER_ADDRESS", "localhost"),
		port:         getEnvWithDefault("SERVER_PORT", "8080"),
		readTimeout:  rTimeout,
		writeTimeout: wTimeout,
	}
}

func getEnvWithDefault(key string, defaultValue string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	} else {
		return defaultValue
	}
}

func AuthEnabled() bool {
	return getEnvWithDefault("AUTH_ENABLED", "false") != "false"
}
