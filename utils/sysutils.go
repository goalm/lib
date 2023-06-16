package utils

import (
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
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

// embeded FS to release dll
type FsCloner struct {
	fs   *embed.FS
	name string
	dist string
}

func New(fs *embed.FS, name, dist string) *FsCloner {
	return &FsCloner{
		fs:   fs,
		name: name,
		dist: dist,
	}
}

func (c *FsCloner) Clone() error {
	return c.clone(c.name)
}

func (c *FsCloner) clone(name string) error {
	if err := os.MkdirAll(path.Join(c.dist, name), 0755); err != nil {
		return err
	}
	dir, err := c.fs.ReadDir(name)
	if err != nil {
		panic(err)
	}

	for _, entry := range dir {
		p := path.Join(name, entry.Name())
		if entry.IsDir() {
			if err := c.clone(p); err != nil {
				return err
			}
			continue
		}
		file, err := c.fs.Open(p)
		if err != nil {
			panic(err)
		}
		newFile, err := os.OpenFile(path.Join(c.dist, p), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			file.Close()
			return err
		}

		_, err = io.Copy(newFile, file)
		file.Close()
		newFile.Close()
		if err != nil {
			panic(err)
		}
	}
	return nil
}
