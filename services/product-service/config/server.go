package config

import (
	"fmt"

	"github.com/daffaromero/retries/services/common/utils"
)

var EndpointPrefix = utils.GetEnv("ENDPOINT_PREFIX")

type serverConfig struct {
	URI  string
	Port string
	Host string
}

func NewServerConfig() serverConfig {
	return serverConfig{
		URI:  utils.GetEnv("SERVER_URI"),
		Port: utils.GetEnv("SERVER_PORT"),
		Host: fmt.Sprintf("%s:%s", utils.GetEnv("SERVER_URI"), utils.GetEnv("SERVER_PORT")),
	}
}
