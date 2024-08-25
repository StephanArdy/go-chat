package websocket

import (
	"bytes"
	"context"
	"go-chat/model"
	"go-chat/repository"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub            *Hub
	conn           *websocket.Conn
	send           chan []byte
	chatRepository repository.ChatRepository
	roomID         string
	SenderID       string
	ReceiverID     string
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error : %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		err = c.chatRepository.SaveMessage(context.Background(), model.Message{
			MessageID:   primitive.NewObjectID(),
			SenderID:    c.SenderID,
			ReceiverID:  c.ReceiverID,
			MessageText: string(message),
			Timestamp:   time.Now(),
			ChatRoomID:  c.roomID,
		})
		if err != nil {
			log.Println(err)
		}
		c.hub.broadcast <- message
	}

}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, chatRepository repository.ChatRepository) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	roomID := r.URL.Query().Get("roomID")
	senderID := r.URL.Query().Get("userID")
	receiverID := r.URL.Query().Get("receiverID")

	client := &Client{
		hub:            hub,
		conn:           conn,
		send:           make(chan []byte, 256),
		chatRepository: chatRepository,
		roomID:         roomID,
		SenderID:       senderID,
		ReceiverID:     receiverID,
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
