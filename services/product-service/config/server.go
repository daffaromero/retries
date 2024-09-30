package config

import (
	"fmt"
	"log"

	"github.com/daffaromero/retries/services/common/utils"
)

var EndpointPrefix = utils.GetEnv("ENDPOINT_PREFIX")

type ServerConfig struct {
	URI  string
	Port string
	Host string
}

func NewServerConfig() ServerConfig {
	uri := utils.GetEnv("SERVER_URI")
	if uri == "" {
		log.Fatal("SERVER_URI environment variable is not set")
	}
	port := utils.GetEnv("SERVER_PORT")
	if port == "" {
		log.Fatal("SERVER_PORT environment variable is not set")
	}
	return ServerConfig{
		URI:  uri,
		Port: port,
		Host: fmt.Sprintf("%s:%s", uri, port),
	}
}
