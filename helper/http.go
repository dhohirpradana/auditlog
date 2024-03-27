package helper

import (
	"auditlog/entity"
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
)

type HttpHandler struct {
}

func InitHttp() HttpHandler {
	return HttpHandler{}
}

func (h HttpHandler) HTTP(c *fiber.Ctx) (err error) {
	elasticIndex := "FALSE"
	err = godotenv.Load(".env")
	if err == nil {
		elasticIndex = os.Getenv("STORE_TO_ELASTIC")
	}

	method := c.Method()
	body := c.Body()
	originalURL := c.OriginalURL()

	requestURL, err := extractURL(originalURL)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	request, err := http.NewRequest(method, requestURL, bytes.NewReader(body))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	modifyRequestHeaders(c, request)

	client := &http.Client{}

	response, err := client.Do(request)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
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
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var requestL entity.Request
	var responseL entity.Response
	var auditlog Auditlog

	requestL.IP = c.IP()

	requestL.Header = headerToMap(request.Header)

	requestL.Body, _ = parseJSON(string(c.Body()))
	if requestL.Body == nil {
		requestL.Body = string(c.Body())
	}

	responseL.Body, _ = parseJSON(string(resBody))
	if responseL.Body == nil {
		responseL.Body = string(resBody)
	}

	responseL.Code = response.StatusCode
	responseL.Header = headerToMap(response.Header)

	auditlog.Method = c.Method()
	auditlog.Url = requestURL
	auditlog.Request = requestL
	auditlog.Response = responseL

	if elasticIndex == "TRUE" {
		go func() {
			err := auditlog.StoreToES()
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	}

	c.Set("Content-Type", contentType)

	return c.Send(resBody)
}
