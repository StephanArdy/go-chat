package service

import (
	"context"
	"go-chat/dto"
	"go-chat/repository"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService interface {
	GetMessages(ctx context.Context, data dto.GetMessagesRequest) (resp []dto.GetMessagesResponse, err error)
	GetorCreateChatRoom(ctx context.Context, userID1 string, userID2 string) (resp dto.GetorCreateChatRoomResponse, err error)
}

type ChatServiceImpl struct {
	chatRepository repository.ChatRepository
}

func NewChatService(chatRepository repository.ChatRepository) ChatService {
	return &ChatServiceImpl{chatRepository: chatRepository}
}

func (c *ChatServiceImpl) GetMessages(ctx context.Context, data dto.GetMessagesRequest) (resp []dto.GetMessagesResponse, err error) {
	messages, err := c.chatRepository.GetMessages(ctx, data.RoomID, int64(data.Limit))
	if err != nil {
		log.Println(err)
		return []dto.GetMessagesResponse{}, err
	}

	for _, message := range messages {
		resp = append(resp, dto.GetMessagesResponse{
			MessageID:   message.MessageID.Hex(),
			SenderID:    message.SenderID,
			ReceiverID:  message.ReceiverID,
			MessageText: message.MessageText,
			Timestamp:   message.Timestamp,
			ChatRoomID:  message.ChatRoomID,
		})
	}

	return resp, nil
}

func (c *ChatServiceImpl) GetorCreateChatRoom(ctx context.Context, userID1 string, userID2 string) (resp dto.GetorCreateChatRoomResponse, err error) {
	getResp, err := c.chatRepository.GetChatRoom(ctx, userID1, userID2)
	if err != nil {
		log.Println(err)
		return resp, err
	}

	if getResp.ChatRoomID == primitive.NilObjectID {
		createResp, err := c.chatRepository.CreateChatRoom(ctx, userID1, userID2)
		if err != nil {
			log.Println(err)
			return resp, err
		}

		resp = dto.GetorCreateChatRoomResponse{
			ChatRoomID: createResp.ChatRoomID.Hex(),
			UserIDs:    createResp.UserIDs,
		}
	} else {
		resp = dto.GetorCreateChatRoomResponse{
			ChatRoomID: getResp.ChatRoomID.Hex(),
			UserIDs:    getResp.UserIDs,
		}
	}

	return
}
