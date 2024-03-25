package helper

import (
	"auditlog/entity"
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"net/url"
)

type HttpHandler struct {
}

func InitHttp() HttpHandler {
	return HttpHandler{}
}

func (h HttpHandler) HTTP(c *fiber.Ctx) (err error) {
	_, err = url.ParseRequestURI(c.Query("url"))
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	urlParam := c.Query("url")

	method := c.Method()
	body := c.Body()

	request, err := http.NewRequest(method, urlParam, bytes.NewReader(body))
	if err != nil {
		fmt.Println(err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for key, value := range c.GetReqHeaders() {
		request.Header[key] = value
	}

	client := &http.Client{}

	response, err := client.Do(request)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if response == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "response is nil")
	}

	c.Status(response.StatusCode)

	contentType := response.Header.Get("Content-Type")
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var requestL entity.Request
	var responseL entity.Response
	var auditlog Auditlog

	requestL.Header = toMap(request.Header)
	requestL.Body = string(c.Body())

	responseL.Data = string(resBody)
	responseL.Code = response.StatusCode
	responseL.Header = toMap(response.Header)

	auditlog.Method = c.Method()
	auditlog.Url = urlParam
	auditlog.Request = requestL
	auditlog.Response = responseL

	go auditlog.StoreToES()

	c.Set("Content-Type", contentType)

	return c.Send(resBody)
}
