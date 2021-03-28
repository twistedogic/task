package fileutil

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
)

func Home() (string, error) {
	return homedir.Dir()
}

func CreateFolder(name string, mod os.FileMode) error {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return os.MkdirAll(name, mod)
	}
	return err
}

func CreateFile(name string, mod os.FileMode) error {
	_, err := os.Stat(name)
	switch {
	case os.IsNotExist(err):
		if _, err := os.Create(name); err != nil {
			return err
		}
	case err != nil:
		return err
	}
	return os.Chmod(name, mod)
}
