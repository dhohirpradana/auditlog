package main

import (
	httpHandler "auditlog/handler"
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	article := httpHandler.InitArticle()
	todo := httpHandler.InitTodo()
	echoServer := echo.New()

	// Fetch Articles
	echoServer.GET("/articles", article.FetchArticles)

	// Fetch Todos
	echoServer.GET("/todos", todo.FetchTodos)

	// Fetch Todo by ID
	echoServer.GET("/todo/:id", todo.GetTodoByID)

	// Create a new Todo
	echoServer.POST("/todo", todo.CreateTodo)

	// Start the server
	err := echoServer.Start(":9090")
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
