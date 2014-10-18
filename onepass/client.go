package onepass

import (
	"encoding/json"
	"github.com/ravenac95/sudolikeaboss/websocketclient"
)

type Command struct {
	Action   string  `json:"action"`
	Number   int     `json:"number"`
	Version  string  `json:"version"`
	BundleId string  `json:"bundleId"`
	Payload  Payload `json:"payload"`
}

type Payload struct {
	Version      string            `json:"version,omitempty"`
	Capabilities []string          `json:"capabilities,omitempty"`
	Url          string            `json:"url,omitempty"`
	Options      map[string]string `json:"options,omitempty"`
}

type WebsocketClient interface {
	Connect() error
	Receive(v interface{}) error
	Send(v interface{}) error
}

type Configuration struct {
	WebsocketUri      string `json:"websocketUri"`
	WebsocketProtocol string `json:"websocketProtocol"`
	WebsocketOrigin   string `json:"websocketOrigin"`
	DefaultHost       string `json:"defaultHost"`
}

type OnePasswordClient struct {
	DefaultHost     string
	websocketClient WebsocketClient
	number          int
}

func NewClientWithConfig(configuration *Configuration) (*OnePasswordClient, error) {
	return NewClient(configuration.WebsocketUri, configuration.WebsocketProtocol, configuration.WebsocketOrigin, configuration.DefaultHost)
}

func NewClient(websocketUri string, websocketProtocol string, websocketOrigin string, defaultHost string) (*OnePasswordClient, error) {
	websocketClient := websocketclient.NewClient(websocketUri, websocketProtocol, websocketOrigin)

	return NewCustomClient(websocketClient, defaultHost)
}

func NewCustomClient(websocketClient WebsocketClient, defaultHost string) (*OnePasswordClient, error) {
	client := OnePasswordClient{
		websocketClient: websocketClient,
		DefaultHost:     defaultHost,
	}

	err := client.Connect()

	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (client *OnePasswordClient) Connect() error {
	err := client.websocketClient.Connect()

	return err
}

func (client *OnePasswordClient) SendShowPopupCommand() (*Response, error) {
	payload := Payload{
		Url:     client.DefaultHost,
		Options: map[string]string{"source": "toolbar-button"},
	}

	command := client.createCommand("showPopup", payload)

	response, err := client.SendCommand(command)

	if err != nil {
		return nil, err
	}
	return response, nil
}

func (client *OnePasswordClient) createCommand(action string, payload Payload) *Command {
	command := Command{
		Action:   action,
		Number:   client.number,
		Version:  "4",
		BundleId: "com.googlecode.iterm2",
		Payload:  payload,
	}

	// Increment the number (it's a 1password thing that I saw whilst listening
	// to their commands
	client.number += 1
	return &command
}

func (client *OnePasswordClient) SendHelloCommand() (*Response, error) {
	capabilities := make([]string, 0)

	payload := Payload{
		Version:      "0.0.1",
		Capabilities: capabilities,
	}

	command := client.createCommand("hello", payload)

	response, err := client.SendCommand(command)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (client *OnePasswordClient) SendCommand(command *Command) (*Response, error) {
	jsonStr, err := json.Marshal(command)

	if err != nil {
		return nil, err
	}

	err = client.websocketClient.Send(jsonStr)

	if err != nil {
		return nil, err
	}

	var rawResponseStr string

	err = client.websocketClient.Receive(&rawResponseStr)

	if err != nil {
		return nil, err
	}

	response, err := LoadResponse(rawResponseStr)

	if err != nil {
		return nil, err
	}

	return response, nil
}
