package main

import (
	"runtime"

	config "github.com/gookit/config"
	"github.com/ptrykov/bloomy/internal"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg := readConfig()

	server := server.NewServer(cfg)
	server.Run()
}

func readConfig() *server.ServerConfig {
	config_keys := []string{"port", "channels"}
	config.WithOptions(config.ParseEnv)

	config.LoadOSEnv(config_keys)
	config.LoadFlags(config_keys)
	return &server.ServerConfig{
		Port:     config.String("port"),
		Channels: uint32(config.Int("channels")),
	}
}
