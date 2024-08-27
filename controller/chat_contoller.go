package controller

import (
	"context"
	"encoding/json"
	"go-chat/dto"
	"go-chat/service"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ChatController interface {
	GetMessages(w http.ResponseWriter, r *http.Request, param httprouter.Params)
	GetorCreateChatRoom(w http.ResponseWriter, r *http.Request, param httprouter.Params)
}

type ChatControllerImpl struct {
	chatService service.ChatService
}

func NewChatController(chatService service.ChatService) ChatController {
	return &ChatControllerImpl{chatService: chatService}
}

// @Summary Get messages by room ID
// @Description Retrieve a list of messages for a specific chat room
// @Tags messages
// @Accept json
// @Produce json
// @Param roomId path string true "Chat Room ID"
// @Param limit query int false "Limit the number of messages returned"
// @Success 200 {object} dto.GetMessagesResponse
// @Failure 400 {object} error
// @Router /messages/{roomID} [get]
func (c *ChatControllerImpl) GetMessages(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	roomID := param.ByName("roomId")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
	}

	messagesRequest := dto.GetMessagesRequest{
		RoomID: roomID,
		Limit:  limit,
		Offset: offset,
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&messagesRequest); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	data, err := c.chatService.GetMessages(ctx, messagesRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	resp := dto.Response{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
		return
	}
}

// @Summary Get or Create Chat Room
// @Description Retrieve an existing chat room for the specified users or create a new one if it doesn't exist.
// @Tags messages
// @Accept json
// @Produce json
// @Param user_id query string true "ID of the user in the chat room"
// @Param friend_id query string true "ID of the friend in the chat room"
// @Success 200 {object} dto.GetorCreateChatRoomResponse
// @Failure 400 {object} error
// @Router /messages/chatRoom [post]
func (c *ChatControllerImpl) GetorCreateChatRoom(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	userID1 := r.URL.Query().Get("user_id")
	userID2 := r.URL.Query().Get("friend_id")

	if userID1 == "" || userID2 == "" {
		http.Error(w, "user_id and friend_id are required", http.StatusBadRequest)
		return
	}
	ctx := context.Background()

	data, err := c.chatService.GetorCreateChatRoom(ctx, userID1, userID2)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get or create ChatRoom", http.StatusInternalServerError)
		return
	}

	resp := dto.Response{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
		return
	}
}
