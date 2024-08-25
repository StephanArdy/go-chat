package dto

import "time"

type FriendRequestParameter struct {
	UserID   string `json:"user_id"`
	FriendID string `json:"friend_id"`
}

type FriendRequestResponse struct {
	RequestID  string    `json:"_id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"update_at"`
}

type UpdateRequestParameter struct {
	RequestID  string `json:"request_id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Acceptance bool   `json:"acceptance"`
}

type UpdateFriendRequestResponse struct {
	RequestID  string    `json:"request_id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GetFriendListsRequest struct {
	UserID string `json:"user_id"`
}

type GetFriendListsResponse struct {
	UserID  string   `json:"user_id"`
	Friends []string `json:"friends"`
}
