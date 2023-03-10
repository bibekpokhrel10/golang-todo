package router

import (
	"net/http"
	"std/internal/controllers"
	"std/internal/models"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8000/resources/", http.StatusMovedPermanently)
}

func RouterHandle() *mux.Router {
	db, err := models.ConnectDB()
	if err != nil {
		log.Error(err)
		return nil
	}
	userRepo := models.NewUserRepo(db)
	h := controllers.NewBaseHandler(userRepo)
	router := mux.NewRouter()
	router.Use(CORS)
	router.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources")))).Methods("GET")
	router.HandleFunc("/", redirect)
	router.HandleFunc("/todo/healthz", controllers.Healthz).Methods("GET")
	router.HandleFunc("/todo-completed", h.GetCompletedItems).Methods("GET")
	router.HandleFunc("/todo-incomplete", h.GetIncompleteItems).Methods("GET")
	router.HandleFunc("/todo", h.CreateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", h.UpdateItem).Methods("PUT", "OPTIONS")
	router.HandleFunc("/todo/{id}", h.DeleteItem).Methods("DELETE", "OPTIONS")
	return router
}
