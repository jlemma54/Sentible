package routes

import (
	"github.com/gorilla/mux"
	"mux-mongo-api/controllers"
)

func CommentRoute(router *mux.Router) {
	router.HandleFunc("/comment", controllers.CreateComment()).Methods("POST")
	router.HandleFunc("/comment/{id1}", controllers.GetComment()).Methods("GET")
}
