# Go To-Do Application

This is a To-Do list application.

- Server: Golang
- Client: html, css , javascript
- Database: Sqlite

## Application Requirement

### Golang server Requirement

1. Golang https://go.dev/dl/
2. Gorilla/mux package for router ```go get github.com/gorilla/mux```
3. Gorm package for connecting to database ```go get github.com/jinzhu/gorm``` 
 and for sqlite ```go get github.com/jinzhu/gorm/dialects/sqlite```
4. For logging ```go get github.com/sirupsen/logrus```
5. For tests ```go get github.com/golang/mock/gomock``` and 
	```go get github.com/stretchr/testify/assert```

## Run Tests

- Total coverage 

 From directory go-todo/internal run ```go test -cover ./...```

![cover page](/Todolist_images/coveragetest.png)

- Coverage of functions

 From directory go-todo/internal first run ```go test -coverprofile cover.out ./...```
  and run ```go cover tool -func cover.out```

![function coverpage](/Todolist_images/funccover.png)

## Start the application

- From go-todo directory, open terminal and run ```go run main.go```


## Walk through the application

Open the index.html file in any browser

### Index page 

![Index page](/Todolist_images/homepage.png)

### Create a todo list

Enter a todo list and click on add icon or enter.

![Create a list](/Todolist_images/create.png)

### The todo item is completed

Click on the tick icon for the todo item to be marked as completed.

![Completed a Todo item](/Todolist_images/completedlist.png)

### Deleting the Todo Item

Click on the delete icon for the item to be deleted

![Deleting a Todo item](/Todolist_images/Deleted.png)

### To undo a Todo Item

Click on the green tick icon again for the task to be marked Incomplete.

![Incomplete a Todo item](/Todolist_images/Incomplete.png)
