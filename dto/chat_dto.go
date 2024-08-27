package dto

import "time"

type GetMessagesRequest struct {
	RoomID string `json:"chat_room_id"`
	Limit  int64    `json:"limit"`
	Offset int64    `json:"offset"`
}

type GetMessagesResponse struct {
	MessageID   string    `json:"_id"`
	SenderID    string    `json:"sender_id"`
	ReceiverID  string    `json:"receiver_id"`
	MessageText string    `json:"message_text"`
	Timestamp   time.Time `json:"timestamp"`
	ChatRoomID  string    `json:"chat_room_id"`
}

type GetorCreateChatRoomResponse struct {
	ChatRoomID string   `json:"chat_room_id"`
	UserIDs    []string `json:"user_ids"`
}
