package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

const HELLO_COMMAND = `
{
	"action":"hello",
	"number":2,
	"version":"4",
	"bundleId":"com.googlecode.iterm2",
	"payload": {
		"version":"0.0.1",
		"capabilities": []
	}
}
`

type Command struct {
	Action   string `json:"action"`
	Number   int    `json:"number"`
	Version  string `json:"version"`
	BundleId string `json:"bundleId"`
}

type HelloPayload struct {
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities"`
}

type ShowPopupPayload struct {
	Url     string            `json:"url"`
	Options map[string]string `json:"options"`
}

type HelloCommand struct {
	Command
	Payload HelloPayload `json:"payload"`
}

type ShowPopupCommand struct {
	Command
	Payload ShowPopupPayload `json:"payload"`
}

func main() {
	x := make([]string, 0)

	helloMsg := &HelloCommand{Command{"hello", 0, "4", "com.googlecode.iterm2"}, HelloPayload{"0.0.1", x}}

	helloJson, err := json.Marshal(helloMsg)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(helloJson))

	conn, err := websocket.Dial("ws://127.0.0.1:6263/4", "", "resource://onepassword4-at-agilebits-dot-com")

	if err != nil {
		fmt.Println("Fatal Error ", err.Error())
		os.Exit(1)
	}

	err = websocket.Message.Send(conn, helloJson)
	if err != nil {
		fmt.Println("Message failed to send")
		os.Exit(1)
	}

	var msg string

	err = websocket.Message.Receive(conn, &msg)
	if err != nil {
		fmt.Println("Failed to receive")
		os.Exit(1)
	}
	fmt.Println(msg)

	showPopupMsg := ShowPopupCommand{
		Command{"showPopup", 1, "4", "com.google.iterm2"},
		ShowPopupPayload{"https://somewheressh.dev", map[string]string{"source": "toolbar-button"}},
	}

	showPopupJson, err := json.Marshal(showPopupMsg)

	fmt.Println(string(showPopupJson))

	err = websocket.Message.Send(conn, showPopupJson)
	if err != nil {
		fmt.Println("Failed to send")
		os.Exit(1)
	}

	err = websocket.Message.Receive(conn, &msg)
	if err != nil {
		if err == io.EOF {
			// graceful shutdown by server
			return
		}
		fmt.Println("Failed to receive anything")
		os.Exit(1)
	}
	fmt.Println(msg)

	//tooLate := make(chan bool)
	done := make(chan bool)

	go func() {
		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			if err == io.EOF {
				// graceful shutdown by server
				done <- true
			}
			fmt.Println("Failed to receive anything")
			os.Exit(1)
		}
	}()

	fmt.Println("hello there it is")

	select {
	case <-done:
		fmt.Println("I'm done?!")
	case <-time.After(10 * time.Second):
		fmt.Println("too late")
		close(done)
	}
}
