package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	CAFile               = configFile("cert", "ca.pem")
	ServerCertFile       = configFile("cert", "server.pem")
	ServerKeyFile        = configFile("cert", "server-key.pem")
	RootClientCertFile   = configFile("cert", "root-client.pem")
	RootClientKeyFile    = configFile("cert", "root-client-key.pem")
	NobodyClientCertFile = configFile("cert", "nobody-client.pem")
	NobodyClientKeyFile  = configFile("cert", "nobody-client-key.pem")
)

type PathExtractionFailed struct {
	targetDir, providedPath string
}

func (pef PathExtractionFailed) Error() string {
	return fmt.Sprintf("could not extract path for dir %q from path: %q", pef.targetDir, pef.providedPath)
}

// extractDirPath extracts the longest/most nested path for a given dir name.
// Path's without a trailing '/' are treated as paths referencing a file.
func extractDirPath(path, dir string) (string, error) {
	if !strings.HasSuffix(path, "/") {
		path = filepath.Dir(path)
		path = path + "/"
	}

	if !strings.Contains(path, "/"+dir+"/") {
		return "", PathExtractionFailed{targetDir: dir, providedPath: path}
	}

	splits := strings.SplitAfter(path, dir)
	if len(splits) == 0 {
		return "", PathExtractionFailed{targetDir: dir, providedPath: path}
	}

	// assert that longest path possible (to handle nested 'proglog' dir's) is extracted
	pathStrings := splits[:len(splits)-1]

	extractedPath := filepath.Join(pathStrings...)
	extractedPath = extractedPath + "/"
	return extractedPath, nil
}

func configFile(path ...string) string {
	if dir := os.Getenv("CONFIG-DIR"); dir != "" {
		path = append([]string{dir}, path...)
		return filepath.Join(path...)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	projectRoot := "proglog"
	absProjectRoot, err := extractDirPath(cwd, projectRoot)
	if err != nil {
		panic(fmt.Errorf("could not extract project root %q from cwd: %q", projectRoot, cwd))
	}

	path = append([]string{absProjectRoot, "test"}, path...)
	return filepath.Join(path...)
}
