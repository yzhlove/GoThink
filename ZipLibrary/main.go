package main

import (
	"bytes"
	"github.com/zip_library/internal"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func main() {

	root := "/Users/yurisa/Develop/GoWork/src/WorkSpace/GoThink/Download"
	prefix := "/Users/yurisa/Develop/GoWork/src/WorkSpace/GoThink"

	packer, err := internal.NewPacker("achive.zip", prefix)
	throwError(err)

	defer packer.Close()

	throwError(packer.AddBytes("/csv/a.txt", bytes.NewReader([]byte("hello world 1"))))
	throwError(packer.AddBytes("/csv/b.txt", bytes.NewReader([]byte("hello world 1"))))

	err = filepath.Walk(root, func(path string, detail fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if detail.IsDir() {
			return nil
		}

		if strings.HasPrefix(detail.Name(), ".") {
			return nil
		}

		if strings.HasPrefix(detail.Name(), "~") {
			return nil
		}

		return packer.AddFile(path, detail)
	})

	throwError(err)
}

func throwError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
