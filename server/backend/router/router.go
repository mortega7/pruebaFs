package router

import (
	"log"
	"net/http"

	"github.com/mortega7/pruebaFs/server/backend/controllers"
)

func SetRoutes() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/channel", controllers.GetChannels)
	mux.HandleFunc("/api/user", controllers.GetUsers)
	mux.HandleFunc("/api/file", controllers.GetFiles)

	log.Fatal(http.ListenAndServe(controllers.API_PORT, mux))
}
