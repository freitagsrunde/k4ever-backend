package utils

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func UploadFile(file []byte, topic string, filename string) (string, error) {
	dir := string("./files/" + topic + "/")
	path := string(dir + filename)
	serverFile := string("/files/" + topic + "/" + filename)

	contentType := http.DetectContentType(file)

	if strings.Split(contentType, "/")[0] != "image" {
		return "", errors.New("filetype not supported")
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	if _, err := os.Stat(path); err == nil {
		return "", errors.New("file already exists")
	}

	err := ioutil.WriteFile(path, file, 0644)
	if err != nil {
		return "", err
	}

	return serverFile, nil
}

func DeleteFile(topic, filename string) error {
	path := string("./files/" + topic + "/" + filename)

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
