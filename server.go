package main

import (
	"encoding/json"
	"errors"
	"gitHub.com/gorilla/websocket"
	"github.com/ravenac95/sudolikeaboss/onepass"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024000
)

type Envelope struct {
	Name string
	Type string
	Data []byte
}

type Service interface {
	Send(*Envelope) error
}

type WrappedCommand struct {
	ClientId string          `json:"slabClientId"`
	Command  onepass.Command `json:"command"`
}

type WrappedResponse struct {
	ClientId string           `json:"slabClientId"`
	Response onepass.Response `json:"response"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024000,
	WriteBufferSize: 1024000,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ClientConnection struct {
	// The websocket Connection.
	handler        *WebsocketHandler
	ws             *websocket.Conn
	sendBuffer     chan *Envelope
	receiveBuffer  chan *Envelope
	SendHandler    func(*Envelope, *ClientConnection) error
	ReceiveHandler func([]byte, *ClientConnection) error
	Uuid           string
}

func NewClientConnection(ws *websocket.Conn, handler *WebsocketHandler, sendHandler func(*Envelope, *ClientConnection) error, receiveHandler func([]byte, *ClientConnection) error) *ClientConnection {
	id := uuid.NewV4()
	return &ClientConnection{
		ws:             ws,
		handler:        handler,
		SendHandler:    sendHandler,
		ReceiveHandler: receiveHandler,
		Uuid:           id.String(),
		receiveBuffer:  make(chan *Envelope, 4),
		sendBuffer:     make(chan *Envelope, 4),
	}
}

func (clientConn *ClientConnection) SendEnvelope(event *Envelope) bool {
	select {
	case clientConn.sendBuffer <- event:
		return true
	default:
		return false
	}
}

func (clientConn *ClientConnection) CanReceiveChained() bool {
	return true
}

func (clientConn *ClientConnection) readPump() {
	log.Println("Starting read pump")

	defer func() {
		clientConn.ws.Close()
	}()

	clientConn.ws.SetReadLimit(maxMessageSize)
	clientConn.ws.SetReadDeadline(time.Now().Add(pongWait))
	clientConn.ws.SetPongHandler(func(string) error {
		clientConn.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := clientConn.ws.ReadMessage()
		if err != nil {
			break
		}

		log.Println(string(message))

		err = clientConn.ReceiveHandler(message, clientConn)
		if err != nil {
			break
		}
	}
}

func (clientConn *ClientConnection) sendPump() {
	log.Println("Starting send pump")

	ticker := time.NewTicker(pingPeriod)

	// Make sure to stop the ticker on close
	defer func() {
		ticker.Stop()
		clientConn.ws.Close()
	}()

	for {
		select {
		case receivedEnvelope, ok := <-clientConn.sendBuffer:
			log.Printf("Receiving stuff! on %s", clientConn.handler.Name)
			log.Println(string(receivedEnvelope.Data))
			if !ok {
				clientConn.sendToClient(websocket.CloseMessage, []byte{})
				return
			}

			if err := clientConn.SendHandler(receivedEnvelope, clientConn); err != nil {
				return
			}
		case <-ticker.C:
			log.Printf("ticker %s", clientConn.handler.Name)

			if err := clientConn.sendToClient(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// Send the message to the connected client
func (clientConn *ClientConnection) sendToClient(messageType int, message []byte) error {
	clientConn.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return clientConn.ws.WriteMessage(messageType, message)
}

type WebsocketHandler struct {
	Name               string
	receive            chan *Envelope
	OnWebsocketConnect func(*websocket.Conn, *WebsocketHandler)
	associate          *WebsocketHandler
	clients            map[string]*ClientConnection
	lastUsed           *ClientConnection
}

func NewWebsocketHandler(name string, onWebsocketConnect func(*websocket.Conn, *WebsocketHandler)) *WebsocketHandler {
	handler := WebsocketHandler{
		Name:               name,
		receive:            make(chan *Envelope),
		clients:            make(map[string]*ClientConnection),
		OnWebsocketConnect: onWebsocketConnect,
	}
	go handler.sendPump()
	return &handler
}

func (handler *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	handler.OnWebsocketConnect(ws, handler)
}

func (handler *WebsocketHandler) Associate(associate *WebsocketHandler) {
	handler.associate = associate
}

func (handler *WebsocketHandler) Send(event *Envelope) error {
	log.Printf("Send to %s", handler.Name)
	select {
	case handler.receive <- event:
		return nil
	default:
		return errors.New("Error sending to the associate")
	}
}

func (handler *WebsocketHandler) AddClientConnection(clientConn *ClientConnection) {
	handler.clients[clientConn.Uuid] = clientConn

	log.Println(handler.clients)

	handler.lastUsed = clientConn

	log.Println(handler.lastUsed)
}

func (handler *WebsocketHandler) GetClientConnection(id string) (*ClientConnection, error) {
	log.Printf("Retrieving %s on %s", id, handler.Name)
	clientConn, ok := handler.clients[id]
	log.Println(handler.clients)
	log.Println(handler.lastUsed)
	if !ok {
		if handler.lastUsed != nil {
			log.Println("Using last one")
			return handler.lastUsed, nil
		} else {
			return nil, errors.New("No client to connect to")
		}
	}
	return clientConn, nil
}

func (handler *WebsocketHandler) sendPump() {
	log.Printf("Starting handler send pump %s", handler.Name)

	for {
		select {
		case event, ok := <-handler.receive:
			log.Printf("Sending for %s", handler.Name)
			if !ok {
				//clientConn.sendToClient(websocket.CloseMessage, []byte{})
				return
			}
			clientConn, err := handler.GetClientConnection(event.Name)
			log.Println(clientConn)
			if err != nil {
				log.Println("error :(")
				log.Println(err)
				return
			}
			if !clientConn.SendEnvelope(event) {
				return
			}
		}
	}
}

func runServer() {
	// The slab-client handler
	slabClientHandler := NewWebsocketHandler("slab", func(ws *websocket.Conn, handler *WebsocketHandler) {
		log.Println("Connected to the sudolikeaboss client")

		clientConn := NewClientConnection(
			ws,
			handler,
			func(envelope *Envelope, clientConn *ClientConnection) error {
				log.Println("Envelope being handled by slab connection")

				if envelope.Name == clientConn.Uuid {
					if err := clientConn.sendToClient(websocket.TextMessage, envelope.Data); err != nil {
						return err
					}
				}
				return nil
			},
			// The receive handler
			func(data []byte, clientConn *ClientConnection) error {
				log.Println("Request being handled by slab connection")

				var command onepass.Command
				var wrappedCommand WrappedCommand
				if err := json.Unmarshal(data, &command); err != nil {
					return err
				}
				wrappedCommand.Command = command
				wrappedCommand.ClientId = clientConn.Uuid

				wrappedCommandJsonBytes, err := json.Marshal(wrappedCommand)

				if err != nil {
					return err
				}

				log.Println(string(wrappedCommandJsonBytes))

				envelope := Envelope{Name: clientConn.Uuid, Type: "Command", Data: wrappedCommandJsonBytes}

				clientConn.handler.associate.Send(&envelope)

				return nil
			},
		)

		handler.AddClientConnection(clientConn)

		go clientConn.sendPump()
		clientConn.readPump()
	})

	browserHandler := NewWebsocketHandler("browser", func(ws *websocket.Conn, handler *WebsocketHandler) {

		log.Println("Connected to the browser")

		clientConn := NewClientConnection(
			ws,
			handler,
			func(envelope *Envelope, clientConn *ClientConnection) error {
				log.Println("Envelope being handled by browser connection")

				if err := clientConn.sendToClient(websocket.TextMessage, envelope.Data); err != nil {
					return err
				}
				return nil
			},
			func(data []byte, clientConn *ClientConnection) error {
				log.Println("Request being handled by browser connection")

				// Unwrap the response so we can send it to the right place
				var wrappedResponse WrappedResponse
				if err := json.Unmarshal(data, &wrappedResponse); err != nil {
					return err
				}

				responseJsonBytes, err := json.Marshal(wrappedResponse.Response)

				if err != nil {
					return err
				}

				envelope := Envelope{Name: wrappedResponse.ClientId, Type: "Response", Data: responseJsonBytes}

				clientConn.handler.associate.Send(&envelope)

				return nil
			},
		)

		handler.AddClientConnection(clientConn)

		go clientConn.sendPump()
		clientConn.readPump()
	})

	slabClientHandler.Associate(browserHandler)
	browserHandler.Associate(slabClientHandler)

	// The browser handler
	http.Handle("/slab", slabClientHandler)
	http.Handle("/browser", browserHandler)
	err := http.ListenAndServe("127.0.0.1:16263", nil)
	if err != nil {
		panic(errors.New("Some error occured initializing the server"))
	}
}
