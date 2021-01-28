package onepass

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Response struct {
	Action  string          `json:"action"`
	Version string          `json:"version"`
	Payload ResponsePayload `json:"payload"`
}

type ResponsePayload struct {
	Item           ItemResponsePayload    `json:"item"`
	Algorithm      string                 `json:"alg"`
	Method         string                 `json:"method"`
	Code           string                 `json:"code"`
	Data           string                 `json:"data"`
	Hmac           string                 `json:"hmac"`
	Iv             string                 `json:"iv"`
	M3             string                 `json:"m3"`
	CS             string                 `json:"cs"`
	AssociatedData string                 `json:"adata"`
	Options        map[string]interface{} `json:"options"`
	OpenInTabMode  string                 `json:"openInTabMode"`
}

type ItemResponsePayload struct {
	Uuid           string                 `json:"uuid"`
	NakedDomains   []string               `json:"nakedDomains"`
	Overview       map[string]interface{} `json:"overview"`
	SecureContents SecureContents         `json:"secureContents"`
}

type SecureContents struct {
	HtmlForm map[string]interface{} `json:"htmlForm"`
	Fields   []map[string]string    `json:"fields"`
}

func LoadResponse(rawResponseStr string) (*Response, error) {
	rawResponseBytes := []byte(rawResponseStr)
	var response Response

	if err := json.Unmarshal(rawResponseBytes, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (response *Response) GetPassword() (string, error) {
	if response.Action != "fillItem" {
		errorMsg := fmt.Sprintf("Response action \"%s\" does not have a password", response.Action)
		return "", errors.New(errorMsg)
	}

	for _, field_obj := range response.Payload.Item.SecureContents.Fields {
		if field_obj["designation"] == "password" {
			return field_obj["value"], nil
		}
	}

	return "", errors.New("No password found in the response")
}

func getPasswordFromResponse(rawResponseStr string) (string, error) {
	rawResponseBytes := []byte(rawResponseStr)
	var response Response

	if err := json.Unmarshal(rawResponseBytes, &response); err != nil {
		return "", err
	}

	for _, field_obj := range response.Payload.Item.SecureContents.Fields {
		if field_obj["designation"] == "password" {
			return field_obj["value"], nil
		}
	}

	return "", errors.New("No password found in the response")
}
