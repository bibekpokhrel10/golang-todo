package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
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

func TestCreateItem(t *testing.T) {
	newDescription := []struct {
		newexpected string
		descrip     string
	}{
		{"{\"Id\":1,\"Description\":\"Hello\",\"Completed\":false}\n", "?description=Hello"},
		{"{\"Id\":2,\"Description\":\"Hi\",\"Completed\":false}\n", "?description=Hi"},
		{"{\"Id\":3,\"Description\":\"Nepal\",\"Completed\":false}\n", "?description=Nepal"},
	}

	for _, des := range newDescription {
		req, err := http.NewRequest("POST", "/todo"+des.descrip, nil)
		if err != nil {
			t.Fatal("err")
		}
		res := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateItem)
		handler.ServeHTTP(res, req)
		if status := res.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: go %v want %v", status, http.StatusOK)
		}
		expected := des.newexpected
		if res.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
		}
	}
}

func TestGetIncompletedItems(t *testing.T) {
	req, err := http.NewRequest("GET", "/todo-incomplete", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(GetIncompleteItems)
	handler.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: go %v want %v", status, http.StatusOK)
	}
	expected := "[{\"Id\":1,\"Description\":\"Hello\",\"Completed\":false},{\"Id\":2,\"Description\":\"Hi\",\"Completed\":false},{\"Id\":3,\"Description\":\"Nepal\",\"Completed\":false}]\n"
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}
}

func TestGetCompletedItems(t *testing.T) {
	req, err := http.NewRequest("GET", "/todo-completed", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(GetCompletedItems)
	handler.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status: got %v want %v", status, http.StatusOK)
	}
	expected := "[]\n"
	if res.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}
}

func TestUpdateItem(t *testing.T) {
	req, err := http.NewRequest("PUT", "/todo/2?Completed=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateItem)
	handler.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("Handler return wrong status: go %v expected %v", status, http.StatusOK)
	}
	expected := `record not found`
	if res.Body.String() != expected {
		t.Errorf("Unexpected body: got %v expected %v", res.Body.String(), expected)
	}
}

func TestDeleteItem(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/todo/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteItem)
	handler.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status: got %v expected %v", status, http.StatusOK)
	}
	expected := "record not found"
	if res.Body.String() != expected {
		t.Errorf("Handler retured unexpected body: got %v expected %v", res.Body.String(), expected)
	}
}
