package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

type TestServer struct {
	DB    *gorm.DB
	Mock  sqlmock.Sqlmock
	store TodoItemModels
}

var server TestServer

func TestMain(m *testing.M) {
	Database()
	os.Exit(m.Run())
}

func Database() {
	var err error
	var db *sql.DB
	_ = os.Setenv("TestDbDriver", "postgres")
	TestDbDriver := os.Getenv("TestDbDriver")
	db, server.Mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}
	server.DB, err = gorm.Open(TestDbDriver, db)
	if err != nil {
		fmt.Printf("Cannot connect to mock %s database\n", TestDbDriver)
		log.Fatal("This is the error:", err)
	}
	server.store = InitializeTodoItemModels(server.DB)
}

func TestConnectDB(t *testing.T) {
	_, err := ConnectDB()
	if err != nil {
		t.Errorf("The error is connecting the db: %v\n", err)
		return
	}
}

func TestNewRepoUser(t *testing.T) {
	db := NewUserRepo(nil)
	if db == nil {
		t.Errorf("The error is : %v\n", db)
		return
	}
}

func TestCreateItemModel(t *testing.T) {
	Item := &TodoItemModel{
		Description: "Create TodoItem",
	}
	data := server.Mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT`)).WillReturnRows(sqlmock.NewRows([]string{"Id"}).AddRow(1))
	fmt.Println(data)
	server.Mock.ExpectCommit()
	server.Mock.ExpectBegin()
	server.Mock.MatchExpectationsInOrder(false)
	saved, err := server.store.CreateItemModel(Item.Description)
	fmt.Println(saved)
	if err != nil {
		t.Errorf("The error is creating the item: %v\n", err)
		return
	}
	_, err = server.store.CreateItemModel(Item.Description)
	if err != nil {
		assert.Error(t, err)
		return
	}
	assert.Equal(t, saved.ID, Item.ID)
	assert.Equal(t, saved.Description, Item.Description)
	assert.Equal(t, saved.Completed, Item.Completed)
}

func TestGetItemById(t *testing.T) {
	Item := &TodoItemModel{
		Description: "Hello",
		Completed:   false,
	}
	Item.ID = 1
	server.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(sqlmock.NewRows([]string{"id", "description", "completed"}).
		AddRow(1, Item.Description, Item.Completed))
	server.Mock.ExpectCommit()
	server.Mock.MatchExpectationsInOrder(false)
	saved, err := server.store.GetItemByID(int(Item.ID))
	if err != nil {
		t.Errorf("this is the error getting the item: %v\n", err)
		return
	}
	_, err = server.store.GetItemByID(int(Item.ID))
	if err != nil {
		assert.Error(t, err)
		return
	}
	assert.Equal(t, saved.ID, Item.ID)
	assert.Equal(t, saved.Description, Item.Description)
	assert.Equal(t, saved.Completed, Item.Completed)
}

func TestUpdateItemModel(t *testing.T) {
	Item := &TodoItemModel{
		Description: "Updated",
		Completed:   false,
	}
	Item.ID = 1
	server.Mock.ExpectBegin()
	server.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(
		sqlmock.NewRows([]string{"id", "description", "completed"}).AddRow(Item.ID, Item.Description, Item.Completed))
	server.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).WillReturnResult(sqlmock.NewResult(0, 1))
	server.Mock.ExpectCommit()
	server.Mock.MatchExpectationsInOrder(false)
	err := server.store.UpdateItemModel(Item.Completed, int(Item.ID))
	if err != nil {
		t.Errorf("This is error updating the item: %v", err)
		return
	}
	err = server.store.UpdateItemModel(Item.Completed, int(Item.ID))
	if err != nil {
		assert.Error(t, err)
		return
	}
}

func TestDeleteItemById(t *testing.T) {
	Item := &TodoItemModel{
		Description: "Deleted",
		Completed:   true,
	}
	Item.ID = 1
	server.Mock.ExpectBegin()
	server.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(
		sqlmock.NewRows([]string{"id", "description", "completed"}).AddRow(Item.ID, Item.Description, Item.Completed))
	server.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).WillReturnResult(sqlmock.NewResult(0, 1))
	server.Mock.ExpectCommit()
	server.Mock.MatchExpectationsInOrder(false)
	err := server.store.DeleteItemByID(int(Item.ID))
	if err != nil {
		t.Errorf("This is error deleting the item:  %v", err)
		return
	}
	err = server.store.DeleteItemByID(int(Item.ID))
	if err != nil {
		assert.Error(t, err)
		return
	}
}

func TestGetTodoItems(t *testing.T) {
	Item := &TodoItemModel{
		Description: "Hello",
		Completed:   false,
	}
	Item.ID = 1
	server.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(sqlmock.NewRows([]string{"id", "description", "completed"}).
		AddRow(Item.ID, Item.Description, Item.Completed))
	server.Mock.ExpectCommit()
	server.Mock.MatchExpectationsInOrder(false)
	saved, err := server.store.GetTodoItems(Item.Completed)
	if err != nil {
		t.Errorf("This is error getting all items: %v", err)
		return
	}
	_, err = server.store.GetTodoItems(Item.Completed)
	if err != nil {
		assert.Error(t, err)
		return
	}
	assert.NotNil(t, saved)
	assert.NoError(t, err)
}
