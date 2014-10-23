package main

import (
	"fmt"
	"github.com/ravenac95/sudolikeaboss/onepass"
	"os"
	"strconv"
	"time"
)

const DEFAULT_TIMEOUT_STRING_SECONDS = "30"
const DEFAULT_HOST = "sudolikeaboss://local"
const DEFAULT_WEBSOCKET_URI = "ws://127.0.0.1:6263/4"
const DEFAULT_WEBSOCKET_PROTOCOL = ""
const DEFAULT_WEBSOCKET_ORIGIN = "chrome-extension://aomjjhallfgjeglblehebfpbcfeobpgk"

func LoadConfiguration() *onepass.Configuration {
	defaultHost := os.Getenv("SUDOLIKEABOSS_DEFAULT_HOST")
	if defaultHost == "" {
		defaultHost = DEFAULT_HOST
	}

	websocketUri := os.Getenv("SUDOLIKEABOSS_WEBSOCKET_URI")
	if websocketUri == "" {
		websocketUri = DEFAULT_WEBSOCKET_URI
	}

	websocketProtocol := os.Getenv("SUDOLIKEABOSS_WEBSOCKET_PROTOCOL")
	if websocketProtocol == "" {
		websocketProtocol = DEFAULT_WEBSOCKET_PROTOCOL
	}

	websocketOrigin := os.Getenv("SUDOLIKEABOSS_WEBSOCKET_ORIGIN")
	if websocketOrigin == "" {
		websocketOrigin = DEFAULT_WEBSOCKET_ORIGIN
	}

	return &onepass.Configuration{
		WebsocketUri:      websocketUri,
		WebsocketProtocol: websocketProtocol,
		WebsocketOrigin:   websocketOrigin,
		DefaultHost:       defaultHost,
	}
}

func RunSudolikeaboss(configuration *onepass.Configuration, done chan bool) {
	// Load configuration from a file
	client, err := onepass.NewClientWithConfig(configuration)

	if err != nil {
		os.Exit(1)
	}

	response, err := client.SendHelloCommand()

	if err != nil {
		os.Exit(1)
	}

	response, err = client.SendShowPopupCommand()

	if err != nil {
		os.Exit(1)
	}

	password, err := response.GetPassword()
	fmt.Println(password)

	done <- true
}

func runGetPassword() {
	done := make(chan bool)

	configuration := LoadConfiguration()

	timeoutString := os.Getenv("SUDOLIKEABOSS_TIMEOUT_SECS")
	if timeoutString == "" {
		timeoutString = DEFAULT_TIMEOUT_STRING_SECONDS
	}

	timeout, err := strconv.ParseInt(timeoutString, 10, 16)

	if err != nil {
		os.Exit(1)
	}

	go RunSudolikeaboss(configuration, done)

	// Timeout if necessary
	select {
	case <-done:
		// Do nothing no need
	case <-time.After(time.Duration(timeout) * time.Second):
		close(done)
		os.Exit(1)
	}
	// Close the app neatly
	os.Exit(0)
}
