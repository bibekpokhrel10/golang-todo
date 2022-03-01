package router

import (
	"std/internal/controllers"

	"github.com/gorilla/mux"
)

func RouterHandle() *mux.Router {
	router := mux.NewRouter()
	router.Use(controllers.CORS)
	router.HandleFunc("/healthz", controllers.Healthz).Methods("GET")
	router.HandleFunc("/todo-completed", controllers.GetCompletedItems).Methods("GET")
	router.HandleFunc("/todo-incomplete", controllers.GetIncompleteItems).Methods("GET")
	router.HandleFunc("/todo", controllers.CreateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", controllers.UpdateItem).Methods("PUT", "OPTIONS")
	router.HandleFunc("/todo/{id}", controllers.DeleteItem).Methods("DELETE", "OPTIONS")
	return router
}
