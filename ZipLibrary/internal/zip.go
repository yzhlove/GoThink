package internal

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type Platform int

const (
	Unix    Platform = 0
	Windows Platform = 1
)

type Packer struct {
	zipFile  *os.File
	zipWrite *zip.Writer
	prefix   string
	platform Platform
}

type ZipFunc func(packer *Packer)

func WithPlatform(platform Platform) ZipFunc {
	return func(packer *Packer) {
		packer.platform = platform
	}
}

func NewPacker(name, prefix string, options ...ZipFunc) (*Packer, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	writer := zip.NewWriter(file)
	packer := &Packer{
		zipFile:  file,
		zipWrite: writer,
		prefix:   prefix,
	}
	for _, opt := range options {
		opt(packer)
	}
	return packer, nil
}

func (p *Packer) Close() {
	p.zipWrite.Flush()
	p.zipWrite.Close()
	p.zipFile.Sync()
	p.zipFile.Close()
}

func (p *Packer) AddFile(path string, stat os.FileInfo) (err error) {
	filePath, err := filepath.Rel(p.prefix, path)
	if err != nil {
		return err
	}

	if stat == nil {
		if stat, err = os.Lstat(path); err != nil {
			return
		}
	}

	header, err := zip.FileInfoHeader(stat)
	if err != nil {
		return err
	}
	switch p.platform {
	case Unix:
		header.Name = filepath.ToSlash(filePath)
	case Windows:
		header.Name = filepath.FromSlash(filePath)
	}
	header.Name = filePath
	header.Method = zip.Deflate

	headerWriter, err := p.zipWrite.CreateHeader(header)
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(headerWriter, file)
	return err
}

func (p *Packer) AddBytes(path string, reader io.Reader) (err error) {

	file, err := p.zipWrite.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, reader)
	return
}
