package controllers

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mortega7/pruebaFs/client/models"
)

const (
	CHANNELS_FOLDER     = "/home/mortega/channelsFolder"
	MAX_MEGABYTES       = 20
	MAX_BUFFER_CAPACITY = MAX_MEGABYTES * 1024 * 1024
)

func GetFolder(conn net.Conn) string {
	addressData := strings.Split(conn.LocalAddr().String(), ":")
	path := CHANNELS_FOLDER + "/" + addressData[1]
	return path
}

func CreateFolder(conn net.Conn) error {
	path := GetFolder(conn)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func CopyFile(fileName, fileData string, conn net.Conn) error {
	size := len(fileData)
	decoded := make([]byte, size*size/base64.StdEncoding.DecodedLen(size))
	_, err := base64.StdEncoding.Decode(decoded, []byte(fileData))
	if err != nil {
		return err
	}

	filePath := GetFolder(conn) + "/" + fileName
	fileCreate, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fileCreate.Close()

	if _, err := fileCreate.Write(decoded); err != nil {
		return err
	}

	if err := fileCreate.Sync(); err != nil {
		return err
	}

	return nil
}

func DecodeFile(filePath string) (models.File, string) {
	fileOpen, err := os.Open(filePath)
	if err != nil {
		return models.File{}, "The specified path does not exist"
	}
	defer fileOpen.Close()

	reader := bufio.NewReaderSize(fileOpen, MAX_BUFFER_CAPACITY)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return models.File{}, err.Error()
	}

	if len(content) > MAX_BUFFER_CAPACITY {
		return models.File{}, "The file size exceeds the maximum allowed of " + strconv.Itoa(MAX_MEGABYTES) + " MB"
	}

	pathData := strings.Split(filePath, "/")
	fileName := strings.Replace(pathData[len(pathData)-1], " ", "", -1)
	mimeType := http.DetectContentType(content)
	encoded := base64.StdEncoding.EncodeToString(content)

	file := models.File{
		Name: fileName,
		Type: mimeType,
		Data: encoded,
	}

	return file, ""
}
