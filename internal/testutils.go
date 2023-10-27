package internal

import (
	"crypto/rand"
	"encoding/hex"
	"net"
	"os"
	"path"
	"testing"
)

func randDirName() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func testRootDir() string {
	return path.Join("/", "tmp")
}

// GetTempFile creates a file in the temporary 'dir'.
// If 'dir' is "", a random dir name will be assigned.
func GetTempFile(t *testing.T, dir, filename string) *os.File {
	// make sure that a random dir is generated if none specified
	if dir == "" {
		dirName := randDirName()
		dir = GetTempDir(t, dirName)
	} else {
		dir = path.Join(testRootDir(), dir)
	}

	f, err := os.CreateTemp(dir, filename)
	if err != nil {
		panic(err)
	}

	t.Logf("created file %q", f.Name())
	return f
}

// GetTempDir creates a temporary directory 'dir'.
// if the 'dir' string includes '*', those will be replaced with a random string.
func GetTempDir(t *testing.T, dir string) string {
	dirPath, err := os.MkdirTemp(testRootDir(), dir)
	if err != nil {
		panic(err)
	}

	t.Logf("testutil: created dir %q", dirPath)
	return dirPath
}

// FreePort returns an available port.
func FreePort(t *testing.T) int {
	for i := 0; i < 10; i++ {
		l, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			t.Logf("could not listen on free port: %v", err)
			continue
		}

		err = l.Close()
		if err != nil {
			t.Logf("could not close listener: %v", err)
			continue
		}

		port := l.Addr().(*net.TCPAddr).Port
		t.Logf("testutil: returned free TCP port '%d'", port)
		return port
	}

	t.Error("could not determine a free port")
	return -1
}
