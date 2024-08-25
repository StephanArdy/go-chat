package repository

import (
	"context"
	"go-chat/model"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUserDataByParam(ctx context.Context, param map[string]interface{}) (model.User, error)
	GetUserDataByEmail(ctx context.Context, email string) (user model.User, err error)
	GetUserDataByUserID(ctx context.Context, userID string) (user model.User, err error)
}

type AuthRepositoryImpl struct {
	mongo *mongo.Client
}

func NewAuthRepository(mongo *mongo.Client) AuthRepository {
	return &AuthRepositoryImpl{mongo: mongo}
}

func (a *AuthRepositoryImpl) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	collection := a.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("Users")
	res, err := collection.InsertOne(ctx, bson.M{
		"user_id":    user.UserID,
		"username":   user.Username,
		"email":      user.Email,
		"password":   user.Password,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"friends":    user.Friends,
	})
	if err != nil {
		return model.User{}, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (a *AuthRepositoryImpl) GetUserDataByParam(ctx context.Context, param map[string]interface{}) (user model.User, err error) {
	collection := a.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("Users")

	err = collection.FindOne(ctx, param).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return model.User{}, nil
	} else if err != nil {
		log.Println(err)
		return model.User{}, err
	}

	return user, nil
}

func (a *AuthRepositoryImpl) GetUserDataByEmail(ctx context.Context, email string) (user model.User, err error) {
	collection := a.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("Users")

	filter := bson.M{"email": email}
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return model.User{}, nil
	} else if err != nil {
		log.Println(err)
		return model.User{}, err
	}
	return user, nil
}

func (a *AuthRepositoryImpl) GetUserDataByUserID(ctx context.Context, userID string) (user model.User, err error) {
	collection := a.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("Users")

	filter := bson.M{"user_id": userID}
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return model.User{}, nil
	} else if err != nil {
		log.Println(err)
		return model.User{}, err
	}
	return user, nil
}
