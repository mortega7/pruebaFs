package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mortega7/pruebaFs/server/backend/controllers"
)

func SetRoutes() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/channel", controllers.GetChannels)
	mux.HandleFunc("/api/user", controllers.GetUsers)
	mux.HandleFunc("/api/file", controllers.GetFiles)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", controllers.API_PORT), mux))
}
