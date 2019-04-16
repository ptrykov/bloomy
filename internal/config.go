package server

import (
	config "github.com/JeremyLoy/config"
)

type ServerConfig struct {
	Port string
}

type fileFromEnv struct {
	BloomyConfig string
}

func ReadConfig() (c *ServerConfig) {
	c = &ServerConfig{}
	fe := &fileFromEnv{}

	config.FromEnv().To(fe)

	if fe.BloomyConfig == "" {
		config.From("config/dev.config").FromEnv().To(c)
	} else {
		config.From(fe.BloomyConfig).FromEnv().To(c)
	}

	return
}
