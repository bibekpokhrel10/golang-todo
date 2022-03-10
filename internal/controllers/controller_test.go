package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"std/internal/controllers/mockdb"
	"std/internal/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHealthz(t *testing.T) {
	url := "/healthz"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(Healthz)
	handler.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := `{"alive":true}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}
}

func TestNewBaseHandler(t *testing.T) {
	db := NewBaseHandler(nil)
	if db == nil {
		t.Errorf("The error is : %v", db)
		return
	}
}

func TestCreateItem(t *testing.T) {
	Item := models.TodoItemModel{
		Description: "Good",
		Completed:   false,
	}
	dummyError := errors.New("Dummy error")
	testCases := []struct {
		name          string
		setCreateItem func(cstore *mockdb.MockTodoItemModels)
	}{
		{
			name: "CreateItem",
			setCreateItem: func(cstore *mockdb.MockTodoItemModels) {
				cstore.EXPECT().CreateItemModel(gomock.Any()).Return(Item, nil).Times(1)
			},
		},
		{
			name: "CreateItemError",
			setCreateItem: func(cstore *mockdb.MockTodoItemModels) {
				cstore.EXPECT().CreateItemModel(gomock.Any()).Return(Item, dummyError).Times(1)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			cstore := mockdb.NewMockTodoItemModels(mockCtrl)
			tc.setCreateItem(cstore)
			TestCreateItem := &BaseHandler{td: cstore}
			req, err := http.NewRequest("POST", "/todo?description=Hello", nil)
			assert.NoError(t, err)
			res := httptest.NewRecorder()
			TestCreateItem.CreateItem(res, req)
		})
	}

}

func TestGetCompletedItems(t *testing.T) {
	Item := []models.TodoItemModel{}
	dummyError := errors.New("Dummy error")
	testCases := []struct {
		name                string
		setGetCompletedItem func(gcstore *mockdb.MockTodoItemModels)
	}{
		{
			name: "GetCompletedItem",
			setGetCompletedItem: func(gcstore *mockdb.MockTodoItemModels) {
				gcstore.EXPECT().GetTodoItems(gomock.Any()).Return(Item, nil).Times(1)
			},
		},
		{
			name: "GetCompletedItemError",
			setGetCompletedItem: func(gcstore *mockdb.MockTodoItemModels) {
				gcstore.EXPECT().GetTodoItems(gomock.Any()).Return(Item, dummyError).Times(1)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			gcstore := mockdb.NewMockTodoItemModels(mockCtrl)
			tc.setGetCompletedItem(gcstore)
			TestGetCompletedItems := &BaseHandler{td: gcstore}
			req, err := http.NewRequest("GET", "/todo-completed", nil)
			assert.NoError(t, err)
			res := httptest.NewRecorder()
			TestGetCompletedItems.GetCompletedItems(res, req)
		})
	}
}

func TestGetInompleteItemsError(t *testing.T) {
	Item := []models.TodoItemModel{}
	dummyError := errors.New("Dummy error")
	testCases := []struct {
		name                 string
		setGetIncompleteItem func(gicstore *mockdb.MockTodoItemModels)
	}{
		{
			name: "GetIncompleteItem",
			setGetIncompleteItem: func(gicstore *mockdb.MockTodoItemModels) {
				gicstore.EXPECT().GetTodoItems(gomock.Any()).Return(Item, nil).Times(1)
			},
		},
		{
			name: "GetInompleteItemError",
			setGetIncompleteItem: func(gicstore *mockdb.MockTodoItemModels) {
				gicstore.EXPECT().GetTodoItems(gomock.Any()).Return(Item, dummyError).Times(1)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			gicstore := mockdb.NewMockTodoItemModels(mockCtrl)
			tc.setGetIncompleteItem(gicstore)
			TestGetIncompleteItems := &BaseHandler{td: gicstore}
			req, err := http.NewRequest("GET", "/todo-incomplete", nil)
			assert.NoError(t, err)
			res := httptest.NewRecorder()
			TestGetIncompleteItems.GetIncompleteItems(res, req)
		})
	}
}

func TestDeleteItem(t *testing.T) {
	Item := models.TodoItemModel{
		Description: "Good",
		Completed:   false,
	}
	dummyError := errors.New("Dummy error")
	testCases := []struct {
		name          string
		setDeleteItem func(dstore *mockdb.MockTodoItemModels)
	}{
		{
			name: "DeleteItem",
			setDeleteItem: func(dstore *mockdb.MockTodoItemModels) {
				dstore.EXPECT().GetItemByID(gomock.Any()).Return(Item, nil).Times(1)
				dstore.EXPECT().DeleteItemByID(gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name: "DeleteItemByIdError",
			setDeleteItem: func(dstore *mockdb.MockTodoItemModels) {
				dstore.EXPECT().GetItemByID(gomock.Any()).Return(Item, nil).Times(1)
				dstore.EXPECT().DeleteItemByID(gomock.Any()).Return(dummyError).Times(1)
			},
		},
		{
			name: "GetItemError",
			setDeleteItem: func(dstore *mockdb.MockTodoItemModels) {
				dstore.EXPECT().GetItemByID(gomock.Any()).Return(Item, dummyError).Times(1)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			dstore := mockdb.NewMockTodoItemModels(mockCtrl)
			tc.setDeleteItem(dstore)
			TestDeleteItem := &BaseHandler{td: dstore}
			req, err := http.NewRequest("DELETE", "/todo/1", nil)
			assert.NoError(t, err)
			res := httptest.NewRecorder()
			TestDeleteItem.DeleteItem(res, req)
		})
	}
}

func TestUpdateItem(t *testing.T) {
	Item := models.TodoItemModel{
		Description: "Good",
		Completed:   false,
	}
	dummyError := errors.New("Dummy Error")
	testCases := []struct {
		name          string
		setUpdateItem func(ustore *mockdb.MockTodoItemModels)
	}{
		{
			name: "UdateItem",
			setUpdateItem: func(ustore *mockdb.MockTodoItemModels) {
				ustore.EXPECT().GetItemByID(gomock.Any()).Return(Item, nil).Times(1)
				ustore.EXPECT().UpdateItemModel(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name: "UpdateItemError",
			setUpdateItem: func(ustore *mockdb.MockTodoItemModels) {
				ustore.EXPECT().GetItemByID(gomock.Any()).Return(Item, nil).Times(1)
				ustore.EXPECT().UpdateItemModel(gomock.Any(), gomock.Any()).Return(dummyError).Times(1)
			},
		},
		{
			name: "GetItemError",
			setUpdateItem: func(ustore *mockdb.MockTodoItemModels) {
				ustore.EXPECT().GetItemByID(gomock.Any()).Return(Item, dummyError).Times(1)

			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ustore := mockdb.NewMockTodoItemModels(mockCtrl)
			tc.setUpdateItem(ustore)
			TestUpdateItem := &BaseHandler{td: ustore}
			req, err := http.NewRequest("PUT", "/todo/1?completed=true", nil)
			assert.NoError(t, err)
			res := httptest.NewRecorder()
			TestUpdateItem.UpdateItem(res, req)
		})
	}
}
