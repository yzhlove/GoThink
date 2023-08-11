package helper

import (
	"io"
	"os"
	"path/filepath"
)

func CheckExists(path string) bool {
	_, exits := os.Stat(path)
	if exits == nil {
		return true
	}

	if os.IsNotExist(exits) {
		return false
	}
	return true
}

func Mkdir(path string) error {
	if !CheckExists(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func WriteFile(path string, reader io.Reader) error {
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		file.Sync()
		file.Close()
	}()
	_, err = io.Copy(file, reader)
	return err
}
