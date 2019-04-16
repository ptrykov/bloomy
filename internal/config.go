package server

import (
	config "github.com/JeremyLoy/config"
)

type ServerConfig struct {
	Port string
}

func ReadConfig() (c *ServerConfig) {
	c = &ServerConfig{}
	config.From("config/dev.config").FromEnv().To(c)
	return
}
