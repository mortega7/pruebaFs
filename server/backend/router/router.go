package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mortega7/pruebaFs/server/backend/controllers"
)

func SetRoutes() {
	app := fiber.New()
	app.Use(cors.New()) //CORS (Cross-Origin Resource Sharing)

	//Canales
	app.Get("/api/channel", controllers.GetChannels)

	//Usuarios
	app.Get("/api/user", controllers.GetUsers)

	//Archivos
	app.Get("/api/file", controllers.GetFiles)

	log.Fatal(app.Listen(":3000"))
}
