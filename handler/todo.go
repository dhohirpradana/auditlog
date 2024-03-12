package handler

import (
	"auditlog/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"gopkg.in/validator.v2"
	"io"
	"net/http"
)

type TodoHandler struct {
}

func InitTodo() TodoHandler {
	return TodoHandler{}
}

func (h TodoHandler) FetchTodos(c echo.Context) (err error) {
	request, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/todos", nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	// Check the HTTP status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code", response.StatusCode)
		return
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var todos []models.Todo
	err = json.Unmarshal(body, &todos)

	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return
	}

	return c.JSON(http.StatusOK, todos)
}

func (h TodoHandler) GetTodoByID(c echo.Context) (err error) {
	request, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/todos/"+c.Param("id"), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	// Check the HTTP status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code", response.StatusCode)
		return
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var todo models.Todo
	err = json.Unmarshal(body, &todo)

	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return
	}

	return c.JSON(http.StatusOK, todo)
}

func (h TodoHandler) CreateTodo(c echo.Context) (err error) {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading request body")
	}

	//fmt.Println("Request Body:", string(body))

	// Create a Todo struct
	var todo models.Todo

	err = json.Unmarshal(body, &todo)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return
	}

	if err := validator.Validate(todo); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	//fmt.Println("Todo:", todo)

	// Post the Todo to the API using todo struct
	request, err := http.NewRequest("POST", "https://jsonplaceholder.typicode.com/todos", bytes.NewBuffer(body))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Set the Content-Type header to application/json
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	// Check the HTTP status code
	if response.StatusCode != http.StatusCreated {
		fmt.Println("Error: Unexpected status code", response.StatusCode)
		return
	}

	// Read the response body
	body, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response Body:", string(body))

	return c.JSON(http.StatusCreated, json.RawMessage(body))
}
