package controllers

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/mortega7/pruebaFs/client/models"
)

const (
	CHANNELS_FOLDER     = "/home/mortega/channelsFolder"
	MAX_BUFFER_CAPACITY = 10 * 1024 * 1024
)

//Funcion que obtiene la carpeta del cliente
func GetFolder(conn net.Conn) string {
	addressData := strings.Split(conn.LocalAddr().String(), ":")
	path := CHANNELS_FOLDER + "/" + addressData[1]
	return path
}

//Funcion que crea la carpeta del usuario al registrarse
func CreateFolder(conn net.Conn) error {
	path := GetFolder(conn)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

//Funcion que copia un archivo del canal a la carpeta del cliente
func CopyFile(fileName, fileData string, conn net.Conn) error {
	path := GetFolder(conn) + "/" + fileName
	size := len(fileData)
	dec := make([]byte, size*size/base64.StdEncoding.DecodedLen(size))
	_, err := base64.StdEncoding.Decode(dec, []byte(fileData))
	if err != nil {
		return err
	}

	fc, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fc.Close()

	if _, err := fc.Write(dec); err != nil {
		return err
	}

	if err := fc.Sync(); err != nil {
		return err
	}

	return nil
}

//Obtiene el contenido del archivo en base64 (devuelve el nombre del archivo, su tipo y la cadena en base64)
func DecodeFile(path string) (models.File, string) {
	pathData := strings.Split(path, "/")
	fileName := strings.Replace(pathData[len(pathData)-1], " ", "", -1) //Se quitan los espacios del nombre para transmitir al servidor

	//Busca el archivo en la ruta indicada
	f, err := os.Open(path)
	if err != nil {
		return models.File{}, "La ruta especificada no existe"
	}

	//Lee el contenido del archivo
	reader := bufio.NewReaderSize(f, MAX_BUFFER_CAPACITY)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return models.File{}, err.Error()
	}

	//Obtiene el tipo del archivo y su contenido en base64
	mimeType := http.DetectContentType(content)
	encoded := base64.StdEncoding.EncodeToString(content)

	file := models.File{
		Name: fileName,
		Type: mimeType,
		Data: encoded,
	}

	return file, ""
}
