package controllers

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/mortega7/pruebaFs/server/backend/models"
)

const (
	ERR_UNDEF_COMM       = "Comando no válido"
	ERR_UNDEF_MSG        = "Mensaje no especificado"
	ERR_UNDEF_CHAN       = "Canal no especificado"
	ERR_UNDEF_PATH       = "Archivo no especificado"
	ERR_UNDEF_FILEDATA   = "Datos del archivo no especificados"
	ERR_NOT_FOUND_CHAN   = "Canal no encontrado"
	ERR_NOT_SUBSCRIPTION = "Debes suscribirte primero a un canal"
)

var UserMessages = make(chan models.Message)
var Messages = make(chan models.Message)
var Users []models.User
var Channels []models.ChannelRoom
var Files []models.File

//Crea los canales por defecto
func CreateDefaultChannels() {
	Channels = []models.ChannelRoom{
		{Name: "channel-1"},
		{Name: "channel-2"},
		{Name: "channel-3"},
	}
}

//Crea un nuevo mensaje
func NewMessage(msg string, user models.User) models.Message {
	message := models.Message{
		Text:    msg,
		Address: user.Conn.RemoteAddr().String(),
		Channel: user.Channel,
	}
	return message
}

//Crea un nuevo usuario
func NewUser(conn net.Conn) models.User {
	user := models.User{
		Address: conn.RemoteAddr().String(),
		Conn:    conn,
	}
	return user
}

//Decodifica el comando enviado por el cliente y ejecuta la funcion deseada
func DecodeCommand(command string, address string) (string, string) {
	ownMessage := ""
	othersMessage := ""
	commandParts := strings.Split(command, " ")

	if len(commandParts) < 1 {
		return "Comando incorrecto o incompleto", ""
	}

	//Se ejecuta el comando
	switch commandParts[0] {
	case "list":
		ownMessage = ListAllChannels(address)
	case "create":
		if commandParts[1] != "" {
			ownMessage = CreateChannel(commandParts[1], address)
		} else {
			ownMessage = ERR_UNDEF_CHAN
		}
	case "subs":
		if commandParts[1] != "" {
			ownMessage = SubscribeToChannel(commandParts[1], address)
		} else {
			ownMessage = ERR_UNDEF_CHAN
		}
	case "broadcast":
		if commandParts[1] != "" {
			ownMessage, othersMessage = BroadcastToChannel(address, strings.Join(commandParts[1:], " "))
		} else {
			ownMessage = ERR_UNDEF_MSG
		}
	case "send":
		if commandParts[1] != "" {
			if commandParts[2] != "" {
				ownMessage, othersMessage = SendFileToChannel(address, commandParts[1:])
			} else {
				ownMessage = ERR_UNDEF_FILEDATA
			}
		} else {
			ownMessage = ERR_UNDEF_PATH
		}
	default:
		ownMessage = ERR_UNDEF_COMM
	}

	return ownMessage, othersMessage
}

//Se obtiene un usuario por su direccion (devuelve el puntero)
func FindUserByAddress(address string) *models.User {
	var user *models.User
	for i := range Users {
		if Users[i].Address == address {
			user = &Users[i]
			return user
		}
	}
	return nil
}

//Se obtiene un canal por su nombre (devuelve el puntero)
func FindChannelByName(channelName string) *models.ChannelRoom {
	var channel *models.ChannelRoom
	for i := range Channels {
		if Channels[i].Name == channelName {
			channel = &Channels[i]
			return channel
		}
	}
	return nil
}

//Crea un canal
func CreateChannel(channelName, address string) string {
	var response string
	channel := FindChannelByName(channelName)
	if channel == nil {
		user := FindUserByAddress(address)
		newChannel := models.ChannelRoom{
			Name: channelName,
		}

		Channels = append(Channels, newChannel)
		user.Channel.Name = newChannel.Name
		response = "Canal creado con éxito"
	} else {
		response = "El canal ya existe"
	}
	return response
}

//Devuelve todos los canales y notifica si el usuario esta suscrito a alguno
func ListAllChannels(address string) string {
	user := FindUserByAddress(address)
	var response = "Lista de canales:\n"
	for _, ch := range Channels {
		if user.Channel.Name == ch.Name {
			response += fmt.Sprintf("\t%s <Suscrito>\n", ch.Name)
		} else {
			response += fmt.Sprintf("\t%s\n", ch.Name)
		}
	}
	response = response[:len(response)-1]
	return response
}

//Suscribe el usuario a un canal
func SubscribeToChannel(channelName string, address string) string {
	//Se verifica que el canal exista
	var response string
	user := FindUserByAddress(address)
	channel := FindChannelByName(channelName)
	if channel == nil {
		response = ERR_NOT_FOUND_CHAN
	} else {
		user.Channel.Name = channel.Name
		response = "Suscripción exitosa"
	}
	return response
}

//Funcion para enviar un mensaje escrito al canal suscrito por el cliente
func BroadcastToChannel(address, message string) (string, string) {
	responseOwn := ""
	responseOthers := ""
	user := FindUserByAddress(address)
	if user.Channel.Name != "" {
		responseOthers = message
	} else {
		responseOwn = ERR_NOT_SUBSCRIPTION
	}
	return responseOwn, responseOthers
}

//Funcion para enviar el archivo al canal suscrito por el cliente
func SendFileToChannel(address string, commands []string) (string, string) {
	responseOwn := ""
	responseOthers := ""
	user := FindUserByAddress(address)
	if user.Channel.Name != "" {
		//Se verifica que el canal exista
		channel := FindChannelByName(user.Channel.Name)
		if channel == nil {
			responseOwn = ERR_NOT_FOUND_CHAN
		} else {
			//Se crea el archivo en el path
			file, err := CreateBase64File(*user, commands)
			if err != nil {
				return "Error al recibir el archivo: " + err.Error(), ""
			}

			responseOwn = "Archivo enviado con éxito"
			responseOthers = "Se ha recibido el archivo " + file.Name + "~" + file.Name + "~" + file.Type + "~" + file.Data
		}
	} else {
		responseOwn = ERR_NOT_SUBSCRIPTION
	}
	return responseOwn, responseOthers
}

//Se genera el archivo con base al codigo enviado (commands: filename, data)
func CreateBase64File(user models.User, commands []string) (models.File, error) {
	//Esquema de URI: data:[<media type>][;base64],<data> (por defecto es text/plain;charset=US-ASCII)
	var mediaType string
	var data string
	dataFile := strings.Split(commands[1], ",")

	if len(dataFile) > 1 {
		dataType := strings.Split(dataFile[0], ";")
		mediaType = dataType[0][strings.IndexByte(dataType[0], ':')+1:]
		data = dataFile[1]
	} else {
		mediaType = "text/plain"
		data = dataFile[0]
	}

	//Si el archivo ya existe, se renombre
	count := 1
	fileName := commands[0]
	fileValid := false
	for !fileValid {
		fileOldName := fileName
		for _, f := range Files {
			if f.Name == fileName {
				fn := strings.Split(fileName, ".")
				fileName = fn[0] + "-" + strconv.Itoa(count) + "." + fn[1]
				count = count + 1
			}
		}

		if fileOldName == fileName {
			fileValid = true
		}
	}

	//Se crea el struct con la informacion
	file := models.File{
		Name:    fileName,
		Type:    mediaType,
		Data:    data,
		Channel: user.Channel,
	}
	Files = append(Files, file)
	return file, nil
}
