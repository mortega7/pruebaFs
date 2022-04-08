package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/mortega7/pruebaFs/client/controllers"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "7777"
	CONN_TYPE = "tcp"
)

func main() {
	//Se establece conexion con el servidor
	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error al conectarse al servidor: " + err.Error())
		return
	}
	defer conn.Close()

	//Se crea la carpeta del cliente para guardar los archivos
	if err := controllers.CreateFolder(conn); err != nil {
		fmt.Println("Error al crear la carpeta del cliente: " + err.Error())
		return
	}

	//Canales para recibir y enviar informacion al servidor
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
	message, _ := bufio.NewReaderSize(conn, controllers.MAX_BUFFER_CAPACITY).ReadString('|')
	message = strings.Replace(message, "|", "", 1)

	//Si el mensaje contiene "~", significa que viene la informacion de un archivo
	if strings.Contains(message, "~") {
		fileData := strings.Split(message, "~")

		if err := controllers.CopyFile(fileData[1], fileData[3], conn); err != nil {
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
	reader := bufio.NewReaderSize(os.Stdin, controllers.MAX_BUFFER_CAPACITY)
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
				file, err := controllers.DecodeFile(path)
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
