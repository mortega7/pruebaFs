package controllers

import (
	"github.com/gofiber/fiber/v2"
)

//Lista todos los canales
func GetChannels(c *fiber.Ctx) error {
	return c.Status(200).JSON(Channels)
}

//Lista todos los usuarios
func GetUsers(c *fiber.Ctx) error {
	return c.Status(200).JSON(Users)
}

//Lista todos los archivos
func GetFiles(c *fiber.Ctx) error {
	return c.Status(200).JSON(Files)
}
