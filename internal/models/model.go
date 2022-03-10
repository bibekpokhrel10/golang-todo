package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "gorm.db")
	if err != nil {
		return nil, err
	}
	db.Debug().DropTableIfExists(&TodoItemModel{})
	db.Debug().AutoMigrate(&TodoItemModel{})
	return db, nil
}

func NewUserRepo(db *gorm.DB) *databases {
	return &databases{
		db: db,
	}
}

type databases struct {
	db *gorm.DB
}

type TodoItemModels interface {
	CreateItemModel(d string) (TodoItemModel, error)
	GetTodoItems(completed bool) ([]TodoItemModel, error)
	GetItemByID(Id int) (TodoItemModel, error)
	UpdateItemModel(completed bool, id int) error
	DeleteItemByID(Id int) error
}

func InitializeTodoItemModels(db *gorm.DB) TodoItemModels {
	return &databases{db}
}

func (data *databases) CreateItemModel(d string) (TodoItemModel, error) {
	todo := TodoItemModel{Description: d, Completed: false}
	err := data.db.Model(&todo).Create(&todo).Error
	if err != nil {
		log.Error(err)
		return todo, err
	}
	log.Info("Created TodoItem")
	return todo, nil
}

func (data *databases) GetTodoItems(completed bool) ([]TodoItemModel, error) {
	var todo []TodoItemModel
	err := data.db.Model(&todo).Where("completed = ?", completed).Find(&todo).Error
	if err != nil {
		log.Error(err)
		return todo, err
	}
	return todo, nil
}

func (data *databases) GetItemByID(Id int) (TodoItemModel, error) {
	todo := TodoItemModel{}
	err := data.db.Model(&todo).Where("id = ?", Id).Take(&todo).Error
	if err != nil {
		log.Error(err)
		return todo, err
	}
	return todo, nil
}

func (data *databases) UpdateItemModel(completed bool, Id int) error {
	todo, _ := data.GetItemByID(Id)
	todo.Completed = completed
	err := data.db.Model(&todo).Where("id = ?", Id).Save(&todo).Error
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Updated Todoitem")
	return nil
}

func (data *databases) DeleteItemByID(Id int) error {
	todo := TodoItemModel{}
	err := data.db.Model(&todo).Where("id=?", Id).Delete(&todo).Error
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Deleted todoitem")
	return nil
}
