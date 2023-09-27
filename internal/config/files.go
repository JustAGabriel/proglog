package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	CAFile               = configFile("ca.pem")
	ServerCertFile       = configFile("server.pem")
	ServerKeyFile        = configFile("server-key.pem")
	RootClientCertFile   = configFile("root-client.pem")
	RootClientKeyFile    = configFile("root-client-key.pem")
	NobodyClientCertFile = configFile("nobody-client.pem")
	NobodyClientKeyFile  = configFile("nobody-client-key.pem")
)

func configFile(filename string) string {
	if dir := os.Getenv("CONFIG-DIR"); dir != "" {
		return filepath.Join(dir, filename)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	projectRoot := "proglog"
	splits := strings.SplitAfter(cwd, projectRoot)
	if len(splits) == 0 {
		panic(fmt.Errorf("could not extract project root %q from cwd: %q", projectRoot, cwd))
	}

	absProjectRoot := splits[0]
	return filepath.Join(absProjectRoot, "test", "cert", filename)
}
