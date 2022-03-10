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

type BaseHandler struct {
	td models.TodoItemModels
}

func NewBaseHandler(td models.TodoItemModels) *BaseHandler {
	return &BaseHandler{
		td: td,
	}
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is Ok")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive":true}`)
}

func (h *BaseHandler) GetCompletedItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get completed TodoItems")
	completedTodoItems, err := h.td.GetTodoItems(true)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completedTodoItems)
}

func (h *BaseHandler) GetIncompleteItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get incomplete TodoItems")
	IncompleteTodoItems, err := h.td.GetTodoItems(false)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "applicaton/json")
	json.NewEncoder(w).Encode(IncompleteTodoItems)
}

func (h *BaseHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	todo, err := h.td.CreateItemModel(description)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (h *BaseHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	completed, _ := strconv.ParseBool(r.FormValue("completed"))
	_, err := h.td.GetItemByID(id)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	err = h.td.UpdateItemModel(completed, id)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"updated":true}`)
}

func (h *BaseHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	_, err := h.td.GetItemByID(id)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	err = h.td.DeleteItemByID(id)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, `{"deleted":true}`)
}
