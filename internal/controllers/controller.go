package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"std/internal/models"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	models.DropTable()
}
func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is Ok")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive":true}`)
}

func GetCompletedItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get completed TodoItems")
	completedTodoItems, err := models.GetTodoItems(true)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(completedTodoItems)
}

func GetIncompleteItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get incomplete TodoItems")
	IncompleteTodoItems, err := models.GetTodoItems(false)
	w.Header().Set("Content-Type", "applicaton/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(IncompleteTodoItems)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	todo, err := models.CreateItemModel(description)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	completed, _ := strconv.ParseBool(r.FormValue("completed"))
	err := models.UpdateItemModel(completed, id)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"updated":true}`)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := models.DeleteItemByID(id)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"deleted":true}`)
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
