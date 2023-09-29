package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const targetDir string = "proglog"

func TestGetDirPath(t *testing.T) {
	scenarios := map[string]func(t *testing.T){
		"returns most nested dir path possible":               testGetMostNextedDirPath,
		"returns based on one and only dir occurence in path": testExtractSimpleSubpath,
		"returns error if dir name not in path":               testReturnsErrorIfDirNameNotInPath,
		"does ignore file named as dir":                       testIgnoresFilename,
	}

	for title, test := range scenarios {
		t.Run(title, func(t *testing.T) {
			test(t)
		})
	}
}

func testGetMostNextedDirPath(t *testing.T) {
	// arrange
	const path string = "/home/runner/work/proglog/proglog/"
	const expectedPath string = "/home/runner/work/proglog/proglog/"

	// act
	result, err := extractDirPath(path, targetDir)

	//assert
	require.NoError(t, err)
	require.Equal(t, expectedPath, result)
}

func testExtractSimpleSubpath(t *testing.T) {
	// arrange
	const path string = "/home/runner/work/proglog/test/"
	const expectedPath string = "/home/runner/work/proglog/"

	// act
	result, err := extractDirPath(path, targetDir)

	//assert
	require.NoError(t, err)
	require.Equal(t, expectedPath, result)
}

func testReturnsErrorIfDirNameNotInPath(t *testing.T) {
	// arrange
	const path string = "/home/runner/work/no-proglog-here/"

	// act
	result, err := extractDirPath(path, targetDir)

	//assert
	require.Error(t, err)
	require.Equal(t, "", result)
}

func testIgnoresFilename(t *testing.T) {
	// arrange
	const path string = "/home/runner/work/proglog/proglog"
	const expectedPath string = "/home/runner/work/proglog/"

	// act
	result, err := extractDirPath(path, targetDir)

	//assert
	require.NoError(t, err)
	require.Equal(t, expectedPath, result)
}
