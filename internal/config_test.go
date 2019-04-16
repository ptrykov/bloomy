package server

import (
	"os"
	"testing"
)

func TestReadConfig_should_return_env_variables(t *testing.T) {
	os.Setenv("Port", "2222")
	cfg := ReadConfig()
	if cfg.Port != "2222" {
		t.Fatalf("Env variable is not read")
	}
}

func TestReadConfig_reads_variables_from_config_file(t *testing.T) {
	os.Unsetenv("Port")
	cfg := ReadConfig()
	if cfg.Port != "3333" {
		t.Fatalf("Config file is not read")
	}
	os.Chdir("internal")
}
