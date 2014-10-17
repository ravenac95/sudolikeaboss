package main

import (
	"fmt"
	"github.com/ravenac95/termpass/onepass"
	"os"
	"strconv"
	"time"
)

const DEFAULT_TIMEOUT_STRING_SECONDS = "30"
const DEFAULT_HOST = "termpass://local"
const WEBSOCKET_URI = "ws://127.0.0.1:6263/4"
const WEBSOCKET_PROTOCOL = ""
const WEBSOCKET_ORIGIN = "resource://onepassword4-at-agilebits-dot-com"

func LoadConfiguration() *onepass.Configuration {
	defaultHost := os.Getenv("TERMPASS_DEFAULT_HOST")
	if defaultHost == "" {
		defaultHost = DEFAULT_HOST
	}
	return &onepass.Configuration{
		WebsocketUri:      WEBSOCKET_URI,
		WebsocketProtocol: WEBSOCKET_PROTOCOL,
		WebsocketOrigin:   WEBSOCKET_ORIGIN,
		DefaultHost:       defaultHost,
	}
}

func RunTermpass(configuration *onepass.Configuration, done chan bool) {
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

func main() {
	done := make(chan bool)

	configuration := LoadConfiguration()

	timeoutString := os.Getenv("TERMPASS_TIMEOUT_SECS")
	if timeoutString == "" {
		timeoutString = DEFAULT_TIMEOUT_STRING_SECONDS
	}

	timeout, err := strconv.ParseInt(timeoutString, 10, 16)

	if err != nil {
		os.Exit(1)
	}

	go RunTermpass(configuration, done)

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
