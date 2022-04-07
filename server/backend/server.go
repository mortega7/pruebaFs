package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/mortega7/pruebaFs/server/backend/controllers"
	"github.com/mortega7/pruebaFs/server/backend/models"
	"github.com/mortega7/pruebaFs/server/backend/router"
)

const (
	CONN_HOST           = "localhost"
	CONN_PORT           = "7777"
	CONN_TYPE           = "tcp"
	MAX_BUFFER_CAPACITY = 10 * 1024 * 1024
)

func main() {
	listen, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error en net.Listen:", err)
	}

	//Canales por defecto
	controllers.Channels = []models.ChannelRoom{
		{Name: "channel-1"},
		{Name: "channel-2"},
		{Name: "channel-3"},
	}
	fmt.Println("TCP server started")

	go broadcaster()
	go apiServer()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error en listen.Accept:", err)
			continue
		}

		go handle(conn)
	}
}

//Goroutine que se encarga de manejar las conexiones de los clientes
func handle(conn net.Conn) {
	defer conn.Close()

	//Se crea el usuario
	newUser := models.User{
		Address: conn.RemoteAddr().String(),
		Conn:    conn,
	}
	controllers.Users = append(controllers.Users, newUser)
	fmt.Println("Nuevo Usuario:", newUser.Address)

	//Se leen los comandos enviados por los clientes, se ajusta el buffer para recibir mas informacion
	input := bufio.NewScanner(conn)
	buf := make([]byte, MAX_BUFFER_CAPACITY)
	input.Buffer(buf, MAX_BUFFER_CAPACITY)
	for input.Scan() {
		if input.Text() != "" {
			messageToOwnUser, messageToOtherUsers := controllers.DecodeCommand(input.Text(), newUser.Address)
			user := controllers.FindUserByAddress(newUser.Address)

			if messageToOwnUser != "" {
				controllers.UserMessages <- newMessage(messageToOwnUser, *user)
			}
			if messageToOtherUsers != "" {
				controllers.Messages <- newMessage(messageToOtherUsers, *user)
			}
		}
	}

	//Se quita el usuario desconectado
	for i, u := range controllers.Users {
		if u.Address == conn.RemoteAddr().String() {
			fmt.Println("Usuario Desconectado:", u.Address)
			controllers.Users = append(controllers.Users[:i], controllers.Users[i+1:]...)
			break
		}
	}
}

//Goroutine que envia el mensaje a los otros usuarios que esten en el mismo canal
func broadcaster() {
	for {
		select {
		case msg := <-controllers.UserMessages:
			//Mensajes para el mismo usuario
			for _, u := range controllers.Users {
				if msg.Address == u.Conn.RemoteAddr().String() {
					fmt.Fprintln(u.Conn, msg.Text+"|")
				}
			}
		case msg := <-controllers.Messages:
			for _, u := range controllers.Users {
				//Mismo usuario, no envia el mensaje
				if msg.Address == u.Conn.RemoteAddr().String() {
					continue
				}

				//Usuarios del mismo canal, envia el mensaje
				if msg.Channel.Name == u.Channel.Name {
					fmt.Fprintln(u.Conn, "\n$$ Mensaje del canal: "+msg.Text+"|")
				}
			}
		}
	}
}

//Goroutine para controlar la conexion al api
func apiServer() {
	router.SetRoutes()
}

//Se crea un nuevo mensaje
func newMessage(msg string, user models.User) models.Message {
	message := models.Message{
		Text:    msg,
		Address: user.Conn.RemoteAddr().String(),
		Channel: user.Channel,
	}
	return message
}
