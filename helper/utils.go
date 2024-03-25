package helper

import (
	"fmt"
	"net/http"
	"strings"
)

func headersToString(headers http.Header) string {
	var headerStrings []string
	for key, values := range headers {
		for _, value := range values {
			headerStrings = append(headerStrings, fmt.Sprintf("%s: %s", key, value))
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(headerStrings, ", "))
}

func toMap(headers http.Header) map[string]string {
	headerMap := make(map[string]string)
	for key, values := range headers {
		headerMap[key] = values[0] // Assuming only one value per header key
	}
	return headerMap
}
