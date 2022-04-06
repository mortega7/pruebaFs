package models

import "net"

type Message struct {
	Text    string
	Address string
	Channel ChannelRoom
}

type ChannelRoom struct {
	Name string
}

type User struct {
	Address string
	Conn    net.Conn
	Channel ChannelRoom
}

type File struct {
	Name    string
	Type    string
	Data    string
	Channel ChannelRoom
}
