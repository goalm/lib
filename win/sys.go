//go:build windows

package win

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func GetUUID() string {
	cmd := exec.Command("wmic", "path", "win32_computersystemproduct", "get", "uuid")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	uuid := strings.TrimSpace(string(out[5:]))
	return uuid
}

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

func Hide(filename string) error {
	filenameW, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}
	err = syscall.SetFileAttributes(filenameW, syscall.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		return err
	}
	return nil
}

func InitializePath(path string) {
	if _, err := os.Stat(path); err == nil {
		os.RemoveAll(path)
	}
	os.MkdirAll(path, os.ModePerm)
}
