package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

var db, _ = gorm.Open("sqlite3", "gorm.db")

type TodoItemModel struct {
	Id          int `grom:"primary_key"`
	Description string
	Completed   bool
}

func CreateItemModel(d string) (TodoItemModel, error) {
	todo := TodoItemModel{Description: d, Completed: false}
	err := db.Model(&todo).Create(&todo).Error
	if err != nil {
		log.Error("Failed to create todolist", err)
		return todo, err
	}
	log.Info("Created Todoitem")
	return todo, nil
}

func GetTodoItems(completed bool) ([]TodoItemModel, error) {
	var todo []TodoItemModel
	err := db.Model(&todo).Where("Completed = ?", completed).Find(&todo).Error
	if err != nil {
		return todo, err
	}
	return todo, nil
}

func GetItemByID(Id int) (TodoItemModel, error) {
	todo := TodoItemModel{}
	err := db.Model(&todo).First(&todo, Id).Take(&todo).Error
	if err != nil {
		return todo, err
	}
	return todo, nil
}

func UpdateItemModel(completed bool, id int) error {
	todo, err := GetItemByID(id)
	if err != nil {
		log.Error(err)
		return err
	}
	err = db.Model(&todo).First(&todo, id).Error
	if err != nil {
		return err
	}
	todo.Completed = completed
	err = db.Model(&todo).Save(&todo).Error
	if err != nil {
		log.Error("Failed to save the todoitem", err)
		return err
	}
	log.Info("updated Todoitem")
	return nil
}

func DeleteItemByID(Id int) error {
	todo, err := GetItemByID(Id)
	if err != nil {
		log.Error(err)
		return err
	}
	err = db.Model(&todo).Where("id=?", Id).Delete(&todo).Error
	if err != nil {
		log.Error("Todoitem not found in database", err)
		return err
	}
	log.Info("Deleted todoitem")
	return nil
}

func DropTable() {
	db.Debug().DropTableIfExists(&TodoItemModel{})
	db.Debug().AutoMigrate(&TodoItemModel{})

}
