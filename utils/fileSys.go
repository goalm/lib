package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func GetFileList(root string, suffix string) (paths []string, err error) {
	paths = make([]string, 0)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsDir() {
			return nil
		}

		if !info.IsDir() && strings.Contains(info.Name(), suffix) {
			paths = append(paths, path)
		}

		return nil
	})
	return paths, err
}

func InitializePath(path string) {
	if _, err := os.Stat(path); err == nil {
		os.RemoveAll(path)
	}
	os.MkdirAll(path, os.ModePerm)
}

func InitializePaths(path ...string) {
	for _, p := range path {
		InitializePath(p)
	}
}
