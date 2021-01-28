package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/ravenac95/sudolikeaboss/onepass"
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

	websocketURI := os.Getenv("SUDOLIKEABOSS_WEBSOCKET_URI")
	if websocketURI == "" {
		websocketURI = DEFAULT_WEBSOCKET_URI
	}

	websocketProtocol := os.Getenv("SUDOLIKEABOSS_WEBSOCKET_PROTOCOL")
	if websocketProtocol == "" {
		websocketProtocol = DEFAULT_WEBSOCKET_PROTOCOL
	}

	websocketOrigin := os.Getenv("SUDOLIKEABOSS_WEBSOCKET_ORIGIN")
	if websocketOrigin == "" {
		websocketOrigin = DEFAULT_WEBSOCKET_ORIGIN
	}

	stateDirectory := os.Getenv("SUDOLIKEABOSS_STATE_DIRECTORY")
	if stateDirectory == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		stateDirectory = path.Join(usr.HomeDir, ".sudolikeaboss")
	}

	return &onepass.Configuration{
		WebsocketURI:      websocketURI,
		WebsocketProtocol: websocketProtocol,
		WebsocketOrigin:   websocketOrigin,
		DefaultHost:       defaultHost,
		StateDirectory:    stateDirectory,
	}
}

func retrievePasswordFromOnepassword(configuration *onepass.Configuration, done chan bool) {
	// Load configuration from a file
	client, err := onepass.NewClientWithConfig(configuration)

	if err != nil {
		os.Exit(1)
	}

	response, err := client.Authenticate(false)

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

func registerWithOnepassword(configuration *onepass.Configuration, done chan bool) {
	// Load configuration from a file
	client, err := onepass.NewClientWithConfig(configuration)

	if err != nil {
		os.Exit(1)
	}

	_, err = client.Authenticate(true)

	if err != nil {
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Println("Congrats sudolikeaboss is registered!")

	done <- true
}

// Run the main sudolikeaboss entry point
func runSudolikeaboss() {
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

	go retrievePasswordFromOnepassword(configuration, done)

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

func runSudolikeabossRegistration() {
	done := make(chan bool)

	configuration := LoadConfiguration()

	go registerWithOnepassword(configuration, done)

	// Timeout if necessary
	select {
	case <-done:
		// Do nothing no need
		close(done)
		os.Exit(1)
	}
	// Close the app neatly
	os.Exit(0)

}
