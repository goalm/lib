//go:build windows

package win

import (
	"fmt"
	"os/exec"
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
