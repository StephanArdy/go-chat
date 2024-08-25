package service

import (
	"context"
	"errors"
	"go-chat/constant"
	"go-chat/dto"
	"go-chat/model"
	"go-chat/pkg/util"
	"go-chat/repository"
	"log"
	"time"
)

type AuthService interface {
	RegisterUser(ctx context.Context, data dto.RegisterDataRequest) (resp dto.RegisterDataResponse, err error)
	CheckLogin(ctx context.Context, data dto.LoginDataRequest) (resp dto.LoginDataResponse, err error)
}

type AuthServiceImpl struct {
	authRepository repository.AuthRepository
}

func NewAuthService(authRepository repository.AuthRepository) AuthService {
	return &AuthServiceImpl{authRepository: authRepository}
}

func (a *AuthServiceImpl) RegisterUser(ctx context.Context, data dto.RegisterDataRequest) (resp dto.RegisterDataResponse, err error) {

	// check if email already exist
	userDataByEmail, err := a.authRepository.GetUserDataByEmail(ctx, data.Email)
	if err != nil {
		log.Println(err)
		return resp, err
	}

	if userDataByEmail.Email != "" {
		err = errors.New(constant.ERROR_EMAIL_EXIST)
		log.Println(err)
		return resp, err
	}

	// check if userID is already exist
	userDataByID, err := a.authRepository.GetUserDataByUserID(ctx, data.UserID)
	if err != nil {
		log.Println(err)
		return resp, err
	}

	if userDataByID.UserID != "" {
		err = errors.New(constant.ERROR_USERID_EXIST)
		log.Println(err)
		return resp, err
	}

	// hash password
	hashedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		log.Println("Fail to hash password")
		return resp, err
	}

	// define created at
	createdAt := time.Now()

	// create new user
	newUser, err := a.authRepository.CreateUser(ctx, model.User{
		UserID:    data.UserID,
		Username:  data.Username,
		Email:     data.Email,
		Password:  hashedPassword,
		CreatedAt: createdAt,
	})
	if err != nil {
		log.Println(err)
		return resp, nil
	}

	resp = dto.RegisterDataResponse{
		ID:        newUser.ID.Hex(),
		UserID:    newUser.UserID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Password:  newUser.Password,
		CreatedAt: newUser.CreatedAt,
	}

	return resp, nil
}

func (a *AuthServiceImpl) CheckLogin(ctx context.Context, data dto.LoginDataRequest) (resp dto.LoginDataResponse, err error) {

	identifier := util.CheckIdentifier(data.Identifier)

	// check account existence
	userData, err := a.authRepository.GetUserDataByParam(ctx, identifier)
	if err != nil {
		log.Println(err)
		return resp, err
	}

	if userData.UserID == "" {
		err = errors.New(constant.ERROR_LOGIN_NOT_EXIST)
		log.Println(err)
		return resp, err
	}

	// check password
	matchPass := util.CheckPassword(userData.Password, data.Password)
	if !matchPass {
		err = errors.New(constant.ERROR_PASSWORD_NOT_MATCH)
		log.Println(err)
		return resp, err
	}

	resp = dto.LoginDataResponse{
		ID:        userData.ID.Hex(),
		UserID:    userData.UserID,
		Username:  userData.Username,
		Email:     userData.Email,
		Password:  userData.Password,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Friends:   userData.Friends,
	}
	return resp, nil
}
