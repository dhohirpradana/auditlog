package helper

import (
	"encoding/json"
	"net/http"
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
