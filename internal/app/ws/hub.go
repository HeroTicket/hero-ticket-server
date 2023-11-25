package ws

import (
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	allowOriginFunc = func(r *http.Request) bool {
		return true
	}

	upgrader = websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin:      allowOriginFunc,
	}

	hub = &Hub{
		register:   make(chan registerRequest),
		unregister: make(chan ID),
		clients:    make(map[ID]*client),
		send:       make(chan Message, 256),
	}
)

func init() {
	go hub.run()
}

type Hub struct {
	clients    map[ID]*client
	register   chan registerRequest
	unregister chan ID
	send       chan Message
}

type registerRequest struct {
	id     ID
	client *client
}

func (h *Hub) run() {
	for {
		select {
		case req := <-h.register:
			h.clients[req.id] = req.client
		case id := <-h.unregister:
			if client, ok := h.clients[id]; ok {
				close(client.send)
				delete(h.clients, id)
			}
		case msg := <-h.send:
			client, ok := h.clients[msg.ID]
			if ok {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, msg.ID)
				}
			}

		}
	}
}

func Send(msg Message) {
	hub.send <- msg
}

func Unregister(id ID) {
	hub.unregister <- id
}

// Serve godoc
//
// @Summary Serve websocket
// @Description returns websocket connection
// @Tags websocket
// @Router /ws [get]
func Serve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
			return
		}

		id := uuid.New().String()

		client := &client{
			hub:  hub,
			id:   ID(id),
			conn: conn,
			send: make(chan Message, 256),
			mu:   &sync.Mutex{},
		}

		hub.register <- registerRequest{
			id:     ID(id),
			client: client,
		}

		go client.writePump()
		go client.readPump()

		msg := Message{
			Type: IdMessage,
			ID:   ID(id),
		}

		client.send <- msg
	}
}
