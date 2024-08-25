package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	MessageID   primitive.ObjectID `bson:"_id, omitempty"`
	SenderID    string             `bson:"sender_id"`
	ReceiverID  string             `bson:"receiver_id"`
	MessageText string             `bson:"message_text"`
	Timestamp   time.Time          `bson:"timestamp"`
	ChatRoomID  string             `bson:"chat_room_id"`
}

type ChatRoom struct {
	ChatRoomID primitive.ObjectID `bson:"_id, omitempty"`
	UserIDs    []string           `bson:"user_ids"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

type Notification struct {
	NotificationID   primitive.ObjectID `bson:"_id, omitempty"`
	UserID           string             `bson:"user_id"`
	NotificationType string             `bson:"notification_type"`
	Content          string             `bson:"content"`
	Timestamp        time.Time          `bson:"timestamp"`
	Read             bool               `bson:"read"`
}
