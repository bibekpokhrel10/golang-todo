package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	//"github.com/rs/cors"

	log "github.com/sirupsen/logrus"
)

//var db, _ = gorm.Open("mysql", "root:root@/todolist?charset=utf8&parseTime=True&loc=Local")

var db, _ = gorm.Open("sqlite3", "gorm.db")

//var db, _ = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

type TodoItemModel struct {
	Id          int `grom:"primary_key"`
	Description string
	Completed   bool
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	log.WithFields(log.Fields{"description": description}).Info("Add new TodoItem. Saving to database")
	todo := &TodoItemModel{Description: description, Completed: false}
	db.Create(&todo)
	result := db.Last(&todo)
	w.Header().Set("Content-Type", "application/jason")
	json.NewEncoder(w).Encode(result.Value)

}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	todo, err := GetItemByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, err.Error())
		return

	}
	completed, _ := strconv.ParseBool(r.FormValue("completed"))
	log.WithFields(log.Fields{"Id": id, "completed": completed}).Info("Updating TodoItem")
	db.First(&todo, id)
	todo.Completed = completed
	db.Save(&todo)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"updated":true}`)

}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	log.Info("item id: ", id)
	err := DeleteItemByID(id)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"deleted":true}`)
}
func DeleteItemByID(Id int) error {
	todo := TodoItemModel{}
	err := db.Model(&todo).Where("id=?", Id).Delete(&todo).Error
	if err != nil {
		log.Warn("TodoItem not found in database")
		return err
	}
	return nil
}

func GetItemByID(Id int) (TodoItemModel, error) {
	todo := TodoItemModel{}
	err := db.First(&todo, Id).Take(&todo).Error
	if err != nil {
		log.Warn("TodoItem not found in database")
		return todo, err
	}
	return todo, nil
}

func GetCompletedItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get completed TodoItems")
	completedTodoItems, err := GetTodoItems(true)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(completedTodoItems)
}

func GetIncompleteItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get incomplete TodoItems")
	IncompleteTodoItems, err := GetTodoItems(false)
	w.Header().Set("Content-Type", "applicaton/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(IncompleteTodoItems)
}

func GetTodoItems(completed bool) ([]TodoItemModel, error) {
	var todos []TodoItemModel
	err := db.Where("Completed = ?", completed).Find(&todos).Error
	if err != nil {
		return todos, err
	}
	return todos, nil
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-type", "application/json")
	io.WriteString(w, `{"alive":true}`)
}
func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)

			return
		}

		fmt.Println("ok")

		next.ServeHTTP(w, r)
		return
	})
}

func main() {
	defer db.Close()

	db.Debug().DropTableIfExists(&TodoItemModel{})
	db.Debug().AutoMigrate(&TodoItemModel{})

	log.Info("Starting Todolist API Server")
	router := mux.NewRouter()
	router.Use(CORS)
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/todo-completed", GetCompletedItems).Methods("GET")
	router.HandleFunc("/todo-incomplete", GetIncompleteItems).Methods("GET")
	router.HandleFunc("/todo", CreateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", UpdateItem).Methods("PUT", "OPTIONS")
	router.HandleFunc("/todo/{id}", DeleteItem).Methods("DELETE", "OPTIONS")

	http.ListenAndServe(":8000", router)

}
