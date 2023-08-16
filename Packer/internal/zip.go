package internal

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type ZipWriter struct {
	file  *os.File
	write *zip.Writer
}

func NewZipWrite(name string) (*ZipWriter, error) {
	w := &ZipWriter{}

	var err error
	if w.file, err = os.Create(name); err != nil {
		return nil, err
	}

	w.write = zip.NewWriter(w.file)
	return w, nil
}

func (w *ZipWriter) Close() {
	w.write.Flush()
	w.write.Close()
	w.file.Sync()
	w.file.Close()
}

func (w *ZipWriter) WriteFile(path, name string, stat os.FileInfo) (err error) {

	if stat == nil {
		if stat, err = os.Lstat(path); err != nil {
			return err
		}
	}

	header, err := zip.FileInfoHeader(stat)
	if err != nil {
		return err
	}

	header.Name = filepath.ToSlash(name)
	header.Method = zip.Deflate
	write, err := w.write.CreateHeader(header)
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(write, file)
	return
}

func (w *ZipWriter) Write(name string, reader io.Reader) (err error) {

	file, err := w.write.Create(name)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, reader)
	return
}
