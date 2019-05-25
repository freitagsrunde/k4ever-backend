package utils

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
)

func UploadFile(file []byte, filepath string, config k4ever.Config) (string, error) {
	fullPath := string(config.FilesPath() + "/files/" + filepath)
	dir := path.Dir(fullPath)
	serverFile := string("/files/" + filepath)

	contentType := http.DetectContentType(file)

	if strings.Split(contentType, "/")[0] != "image" {
		return "", errors.New("filetype not supported")
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	if _, err := os.Stat(fullPath); err == nil {
		return "", errors.New("file already exists")
	}

	err := ioutil.WriteFile(fullPath, file, 0644)
	if err != nil {
		return "", err
	}

	return serverFile, nil
}

func DeleteFiles(filespath string, config k4ever.Config) error {
	fullPath := string(config.FilesPath() + "/files/" + filespath)

	files, err := filepath.Glob(fullPath + "*")
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
