//go:build windows

package win

import (
	"embed"
	"io"
	"os"
	"path"
)

var fs embed.FS

func ReleaseVault(fs *embed.FS) (path string) {
	appFolder, err := os.UserCacheDir()
	path = appFolder + "/mt/vault"
	if err != nil {
		panic(err)
	}

	cloner := New(fs, "vault", appFolder+"/mt") // vault is the folder name in embed.FS
	if err := cloner.Clone(); err != nil {
		panic(err)
	}
	if err := Hide(path); err != nil {
		panic(err)
	}
	return
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
