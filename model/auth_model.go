package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sessions struct {
	ID           primitive.ObjectID `bson:"_id, omitempty"`
	UserID       string             `bson:"user_id"`
	SessionToken string             `bson:"session_token"`
	CreatedAt    time.Time          `bson:"created_at"`
	ExpiresAt    time.Time          `bson:"expires_at"`
}
