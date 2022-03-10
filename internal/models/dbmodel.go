package models

import "github.com/jinzhu/gorm"

type TodoItemModel struct {
	gorm.Model
	Description string
	Completed   bool
}
