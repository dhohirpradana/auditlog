package helper

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/url"
	"strings"
)

func parseJSON(jsonStr string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func headerToMap(headers http.Header) map[string]string {
	headerMap := make(map[string]string)
	for key, values := range headers {
		headerMap[key] = values[0] // Assuming only one value per header key
	}
	return headerMap
}

func extractURL(originalURL string) (string, error) {
	requestURL := strings.Replace(originalURL, "/?url=", "", 1)

	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", err
	}
	return requestURL, nil
}

func modifyRequestHeaders(c *fiber.Ctx, request *http.Request) {
	for key, value := range c.GetReqHeaders() {
		request.Header[key] = value
	}
}
