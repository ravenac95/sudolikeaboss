package websocketclient

import (
	ws "golang.org/x/net/websocket"
)

type Codec interface {
	Receive(*ws.Conn, interface{}) error
	Send(*ws.Conn, interface{}) error
}

// A websocket client meant to be used with the 1password
type Client struct {
	WebsocketUri      string
	WebsocketProtocol string
	WebsocketOrigin   string
	conn              *ws.Conn
	dial              func(string, string, string) (*ws.Conn, error)
	codec             Codec
}

func NewClient(websocketUri string, websocketProtocol string, websocketOrigin string) *Client {
	return NewCustomClient(websocketUri, websocketProtocol, websocketOrigin, ws.Dial, ws.Message)
}

func NewCustomClient(websocketUri string, websocketProtocol string, websocketOrigin string,
	dial func(string, string, string) (*ws.Conn, error), codec Codec) *Client {

	client := Client{
		WebsocketUri:      websocketUri,
		WebsocketProtocol: websocketProtocol,
		WebsocketOrigin:   websocketOrigin,
		dial:              dial,
		codec:             codec,
	}

	return &client
}

func (client *Client) Connect() error {
	conn, err := client.dial(client.WebsocketUri, client.WebsocketProtocol, client.WebsocketOrigin)

	if err != nil {
		return err
	}

	client.conn = conn

	return nil
}

func (client *Client) Receive(v interface{}) error {
	return client.codec.Receive(client.conn, v)
}

func (client *Client) Send(v interface{}) error {
	return client.codec.Send(client.conn, v)
}
