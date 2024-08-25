package service

import (
	"context"
	"errors"
	"go-chat/constant"
	"go-chat/dto"
	"go-chat/model"
	"go-chat/repository"
	"log"
)

type UserService interface {
	AddFriend(ctx context.Context, req dto.FriendRequestParameter) (resp dto.FriendRequestResponse, err error)
	UpdateFriendRequest(ctx context.Context, req dto.UpdateRequestParameter) (resp dto.UpdateFriendRequestResponse, err error)
	GetFriendLists(ctx context.Context, req string) (resp dto.GetFriendListsResponse, err error)
	GetFriendRequests(ctx context.Context, req string) (resp []dto.FriendRequestResponse, err error)
}

type UserServiceImpl struct {
	AuthRepository repository.AuthRepository
	UserRepository repository.UserRepository
}

func NewUserService(a repository.AuthRepository, u repository.UserRepository) UserService {
	return &UserServiceImpl{
		AuthRepository: a,
		UserRepository: u,
	}
}

func (u *UserServiceImpl) AddFriend(ctx context.Context, req dto.FriendRequestParameter) (resp dto.FriendRequestResponse, err error) {

	checkFriend, err := u.AuthRepository.GetUserDataByUserID(ctx, req.FriendID)
	if err != nil {
		log.Println(err)
		return dto.FriendRequestResponse{}, err
	}

	if checkFriend.ID.Hex() == "0" || checkFriend.UserID == "" {
		err = errors.New(constant.ERROR_FRIEND_NOT_EXIST)
		log.Println(err)
		return dto.FriendRequestResponse{}, err
	}

	friendRequest, err := u.UserRepository.CreateFriendRequest(ctx, req.UserID, req.FriendID)
	if err != nil {
		log.Println(err)
		return dto.FriendRequestResponse{}, err
	}

	resp = dto.FriendRequestResponse{
		RequestID:  friendRequest.RequestID.Hex(),
		SenderID:   friendRequest.SenderID,
		ReceiverID: friendRequest.ReceiverID,
		Status:     friendRequest.Status,
		CreatedAt:  friendRequest.CreatedAt,
		UpdatedAt:  friendRequest.UpdatedAt,
	}
	return resp, nil
}

func (u *UserServiceImpl) UpdateFriendRequest(ctx context.Context, req dto.UpdateRequestParameter) (resp dto.UpdateFriendRequestResponse, err error) {
	var (
		updatedRequest model.FriendRequests
	)

	friendRequest, err := u.UserRepository.GetFriendRequestByRequestID(ctx, req.RequestID)
	if err != nil {
		log.Println(err)
		return resp, err
	}

	if friendRequest.RequestID.Hex() == "0" {
		err = errors.New(constant.ERROR_FRIEND_REQUEST_NOT_EXIST)
		log.Println(err)
		return resp, err
	}

	if req.Acceptance == true {
		updatedRequest, err = u.UserRepository.UpdateFriendRequest(ctx, friendRequest, constant.REQUEST_ACCEPTED_STATUS)
		if err != nil {
			err = errors.New(constant.ERROR_UPDATE_REQUEST_STATUS)
			log.Println(err)
			return resp, err
		}

		err = u.UserRepository.UpdateFriendList(ctx, req.SenderID, req.ReceiverID)
		if err != nil {
			log.Println(err)
			return resp, err
		}
	} else {
		updatedRequest, err = u.UserRepository.UpdateFriendRequest(ctx, friendRequest, constant.REQUEST_DENIED_STATUS)
		if err != nil {
			err = errors.New(constant.ERROR_UPDATE_REQUEST_STATUS)
			log.Println(err)
			return resp, err
		}
	}

	resp = dto.UpdateFriendRequestResponse{
		RequestID:  updatedRequest.RequestID.Hex(),
		SenderID:   updatedRequest.SenderID,
		ReceiverID: updatedRequest.ReceiverID,
		Status:     updatedRequest.Status,
		CreatedAt:  updatedRequest.CreatedAt,
		UpdatedAt:  updatedRequest.UpdatedAt,
	}
	return resp, nil

}

func (u *UserServiceImpl) GetFriendLists(ctx context.Context, req string) (resp dto.GetFriendListsResponse, err error) {
	friends, err := u.UserRepository.GetFriendLists(ctx, req)
	if err != nil {
		log.Println(err)
		return resp, err
	}

	resp = dto.GetFriendListsResponse{
		UserID:  req,
		Friends: friends,
	}
	return resp, nil
}

func (u *UserServiceImpl) GetFriendRequests(ctx context.Context, req string) (resp []dto.FriendRequestResponse, err error) {
	requests, err := u.UserRepository.GetFriendRequests(ctx, req)
	if err != nil {
		log.Println(err)
		return resp, err
	}

	for _, request := range requests {
		requestResponse := dto.FriendRequestResponse{
			RequestID:  request.RequestID.Hex(),
			SenderID:   request.SenderID,
			ReceiverID: request.ReceiverID,
			Status:     request.Status,
			CreatedAt:  request.CreatedAt,
			UpdatedAt:  request.UpdatedAt,
		}

		resp = append(resp, requestResponse)
	}

	return resp, nil
}
