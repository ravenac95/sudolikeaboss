package termpass

import (
	"encoding/json"
	"errors"
)

type Response struct {
	Action  string          `json:"action"`
	Version string          `json:"version"`
	Payload ResponsePayload `json:"payload"`
}

type ResponsePayload struct {
	Item ItemResponsePayload `json:"item"`
}

type ItemResponsePayload struct {
	SecureContents SecureContents `json:"secureContents"`
}

type SecureContents struct {
	Fields []map[string]string `json:"fields"`
}

func getPasswordFromResponse(rawResponseStr string) (string, error) {
	rawResponseBytes := []byte(rawResponseStr)
	var response Response

	if err := json.Unmarshal(rawResponseBytes, &response); err != nil {
		return "", err
	}

	for _, field_obj := range response.Payload.Item.SecureContents.Fields {
		if field_obj["type"] == "P" {
			return field_obj["value"], nil
		}
	}

	return "", errors.New("No password found in the response")
}
