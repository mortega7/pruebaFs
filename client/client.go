package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/mortega7/pruebaFs/client/models"
)

const (
	CONN_HOST           = "localhost"
	CONN_PORT           = "7777"
	CONN_TYPE           = "tcp"
	CHANNELS_FOLDER     = "/home/mortega/channelsFolder"
	MAX_BUFFER_CAPACITY = 10 * 1024 * 1024
)

func main() {
	//Se establece conexion con el servidor
	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error al conectarse al servidor: " + err.Error())
		return
	}
	defer conn.Close()
	if err := createFolder(conn); err != nil {
		fmt.Println("Error al crear la carpeta del cliente: " + err.Error())
		return
	}

	receiveFromServer := make(chan string)
	sendToServer := make(chan bool)

	for {
		go handleReceiveFromServer(conn, receiveFromServer)
		go handleSendToServer(conn, sendToServer)

		select {
		case res := <-receiveFromServer:
			if res != "" {
				//Se imprime la respuesta del servidor
				fmt.Println(res)
			} else {
				//Si se recibe vacio, es porque el servidor se cerró. Se cierran los clientes
				os.Exit(0)
			}
		case <-sendToServer:
		}
	}
}

//Goroutine que se encarga de recibir la informacion del servidor
func handleReceiveFromServer(conn net.Conn, chIn chan string) {
	message, _ := bufio.NewReaderSize(conn, MAX_BUFFER_CAPACITY).ReadString('|')
	message = strings.Replace(message, "|", "", 1)

	//Si el mensaje contiene "~", significa que viene la informacion de un archivo
	if strings.Contains(message, "~") {
		fileData := strings.Split(message, "~")

		if err := copyFile(fileData[1], fileData[3], conn); err != nil {
			message = ""
			fmt.Println("Error al copiar el archivo en la carpeta: ", err.Error())
		} else {
			message = fileData[0]
		}
	}

	chIn <- message
}

//Goroutine que se encarga de enviar la informacion al servidor
func handleSendToServer(conn net.Conn, chOut chan bool) {
	reader := bufio.NewReaderSize(os.Stdin, MAX_BUFFER_CAPACITY)
	fmt.Print(">> ")
	text, _ := reader.ReadString('\n')

	if text != "" {
		//Se desconecta del servidor
		if strings.TrimSpace(string(text)) == "exit" {
			fmt.Println("$$ Adiós, vuelve pronto!")
			os.Exit(0)
		} else {
			commandParts := strings.Split(strings.TrimSpace(string(text)), " ")
			if commandParts[0] == "send" {
				//Se lee y codifica el archivo en base64
				path := strings.Join(commandParts[1:], " ")
				file, err := decodeFile(path)
				if err != "" {
					fmt.Println("Error al leer el archivo: " + err)
				} else {
					//Envia el comando al servidor
					text = commandParts[0] + " " + file.Name + " " + file.Data
					fmt.Fprintf(conn, text+"\n")
				}
			} else {
				//Envia el comando al servidor
				fmt.Fprintf(conn, text+"\n")
			}
		}
	}

	chOut <- true
}

//Funcion que obtiene la carpeta del cliente
func getFolder(conn net.Conn) string {
	addressData := strings.Split(conn.LocalAddr().String(), ":")
	path := CHANNELS_FOLDER + "/" + addressData[1]
	return path
}

//Funcion que crea la carpeta del usuario al registrarse
func createFolder(conn net.Conn) error {
	path := getFolder(conn)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

//Funcion que copia un archivo del canal a la carpeta del cliente
func copyFile(fileName, fileData string, conn net.Conn) error {
	path := getFolder(conn) + "/" + fileName
	sz := len(fileData)
	dec := make([]byte, sz*sz/base64.StdEncoding.DecodedLen(sz))
	_, err := base64.StdEncoding.Decode(dec, []byte(fileData))
	//dec, err := base64.StdEncoding.DecodeString(fileData)
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
func decodeFile(path string) (models.File, string) {
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
