package fileutil

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
)

func Home() (string, error) {
	return homedir.Dir()
}

func CreateFolder(name string, mod os.FileMode) error {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(name, mod); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func CreateFile(name string, mod os.FileMode) error {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(name); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return os.Chmod(name, mod)
}
