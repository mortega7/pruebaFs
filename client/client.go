package main

import (
	"bufio"
	"fmt"
	"log"
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
	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Fatal("Error connecting to server: " + err.Error())
		return
	}
	defer conn.Close()

	if err := controllers.CreateFolder(conn); err != nil {
		log.Fatal("Error creating client folder: " + err.Error())
		return
	}

	receiveFromServer := make(chan string)
	sendToServer := make(chan bool)

	for {
		go handleReceiveFromServer(conn, receiveFromServer)
		go handleSendToServer(conn, sendToServer)

		select {
		case res := <-receiveFromServer:
			//Prints the server response
			fmt.Println(res)
		}
	}
}

func handleReceiveFromServer(conn net.Conn, chIn chan string) {
	//"|" indicates the end of the message
	message, err := bufio.NewReaderSize(conn, controllers.MAX_BUFFER_CAPACITY).ReadString('|')
	if err != nil {
		exitClient(false)
	}

	//If the message contains "~", it means that the information of a file is coming
	message = strings.Replace(message, "|", "", 1)
	if strings.Contains(message, "~") {
		fileData := strings.Split(message, "~")
		message = fileData[0]

		err := controllers.CopyFile(fileData[1], fileData[3], conn)
		if err != nil {
			message = err.Error()
		}
	}

	chIn <- message
}

func handleSendToServer(conn net.Conn, chOut chan bool) {
	reader := bufio.NewReaderSize(os.Stdin, controllers.MAX_BUFFER_CAPACITY)
	fmt.Print(">> ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading data from client: " + err.Error())
	}

	text = strings.TrimSpace(text)
	if text != "" {
		if text == "exit" {
			exitClient(true)
		}

		if commandParts := strings.Split(text, " "); commandParts[0] == "send" {
			filePath := strings.Join(commandParts[1:], " ")
			file, err := controllers.DecodeFile(filePath)
			text = commandParts[0] + " " + file.Name + " " + file.Data
			if err != "" {
				text = "image-wrcomm " + "Error reading file: " + err
			}
		}
	} else {
		text = "wrcomm"
	}

	fmt.Fprintf(conn, text+"\n")
	chOut <- true
}

func exitClient(userExit bool) {
	message := "$$ Goodbye!"
	if !userExit {
		message = "\n$$ The server has closed the connection"
	}
	fmt.Println(message)
	os.Exit(0)
}
