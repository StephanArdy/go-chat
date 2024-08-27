package websocket

import (
	"context"
	"encoding/json"
	"go-chat/dto"
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
	var (
		msgData     map[string]interface{}
		textMessage string
	)

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

		if err := json.Unmarshal(message, &msgData); err != nil {
			log.Println("Failed to parse message:", err)
			continue
		}

		if action, ok := msgData["action"].(string); ok && action == "get_messages" {
			roomID, ok := msgData["room_id"].(string)
			if !ok {
				log.Println("room_id is not a string")
				continue
			}

			limit, ok := msgData["limit"].(float64)
			if !ok {
				log.Println("limit is not a number")
				continue
			}

			offset, ok := msgData["offset"].(float64)
			if !ok {
				log.Println("offset is not a number")
				continue
			}

			messageRequest := dto.GetMessagesRequest{
				RoomID: roomID,
				Limit:  int64(limit),
				Offset: int64(offset),
			}

			messages, err := c.chatRepository.GetMessages(context.Background(), messageRequest.RoomID, messageRequest.Limit, messageRequest.Offset)
			if err != nil {
				log.Println("Failed to retrieve messages: ", err)
				return
			}

			if len(messages) == 0 {
				response := map[string]interface{}{
					"action": "messages",
					"data":   model.Message{},
				}
				responseBytes, _ := json.Marshal(response)
				c.send <- responseBytes
				continue
			} else {
				response := map[string]interface{}{
					"action": "messages",
					"data":   messages,
				}
				responseBytes, _ := json.Marshal(response)
				c.send <- responseBytes
				continue
			}
		}

		if msgData["action"] == "send_message" {
			textMessage = msgData["message_text"].(string)

			err = c.chatRepository.SaveMessage(context.Background(), model.Message{
				MessageID:   primitive.NewObjectID(),
				SenderID:    c.SenderID,
				ReceiverID:  c.ReceiverID,
				MessageText: textMessage,
				Timestamp:   time.Now(),
				ChatRoomID:  c.roomID,
			})
			if err != nil {
				log.Println("Failed to save message: ", err)
				return
			}
			c.hub.broadcast <- message
		}
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
		log.Println("Websocket upgrade error:", err)
		http.Error(w, "Failed to upgrade Websocket", http.StatusInternalServerError)
		return
	}
	log.Println("Websocket connection successfully upgraded")

	params := r.URL.Query()
	roomID := params.Get("roomID")
	senderID := params.Get("userID")
	receiverID := params.Get("receiverID")

	if roomID == "" || senderID == "" || receiverID == "" {
		log.Println("Missing required query parameters")
		http.Error(w, "Missing required query parameters", http.StatusBadRequest)
		return
	}

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
