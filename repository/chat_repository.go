package repository

import (
	"context"
	"go-chat/model"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, message model.Message) (err error)
	GetMessages(ctx context.Context, roomID string, limit int64, offset int64) (messages []model.Message, err error)
	CreateChatRoom(ctx context.Context, userID1 string, userID2 string) (chatRoom model.ChatRoom, err error)
	GetChatRoom(ctx context.Context, userID1 string, userID2 string) (chatRoom model.ChatRoom, err error)
	// Create Notification
}

type ChatRepositoryImpl struct {
	mongo *mongo.Client
}

func NewChatRepository(mongo *mongo.Client) ChatRepository {
	return &ChatRepositoryImpl{
		mongo: mongo,
	}
}

func (c *ChatRepositoryImpl) SaveMessage(ctx context.Context, message model.Message) (err error) {
	collection := c.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("Messages")
	res, err := collection.InsertOne(ctx, bson.M{
		"sender_id":    message.SenderID,
		"receiver_id":  message.ReceiverID,
		"message_text": message.MessageText,
		"timestamp":    message.Timestamp,
		"chat_room_id": message.ChatRoomID,
	})
	if err != nil {
		return err
	}
	message.MessageID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (c *ChatRepositoryImpl) GetMessages(ctx context.Context, roomID string, limit int64, offset int64) (messages []model.Message, err error) {
	collection := c.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("Messages")

	opts := options.FindOptions{
		Limit: &limit,
		Skip: &offset,
		Sort:  bson.D{{"timestamp", -1}},
	}

	cur, err := collection.Find(ctx, bson.D{{"chat_room_id", roomID}}, &opts)
	if err != nil {
		log.Println(err)
		return []model.Message{}, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var message model.Message
		err := cur.Decode(&message)
		if err != nil {
			log.Println("fail to decode")
			return []model.Message{}, err
		}
		messages = append(messages, message)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return messages, cur.Err()
}

func (c *ChatRepositoryImpl) CreateChatRoom(ctx context.Context, userID1 string, userID2 string) (chatRoom model.ChatRoom, err error) {
	collection := c.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("ChatRoom")

	res, err := collection.InsertOne(ctx, bson.M{
		"user_ids":   []string{userID1, userID2},
		"created_at": time.Now(),
		"updated_at": time.Time{},
	})
	if err != nil {
		return chatRoom, err
	}

	id := res.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": id}

	err = collection.FindOne(ctx, filter).Decode(&chatRoom)
	if err == mongo.ErrNoDocuments {
		return chatRoom, nil
	} else if err != nil {
		log.Println(err)
		return chatRoom, err
	}

	return chatRoom, nil
}

func (c *ChatRepositoryImpl) GetChatRoom(ctx context.Context, userID1 string, userID2 string) (chatRoom model.ChatRoom, err error) {
	collection := c.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("ChatRoom")

	filter := bson.M{"user_ids": bson.M{"$all": []string{userID1, userID2}}}

	err = collection.FindOne(ctx, filter).Decode(&chatRoom)
	if err == mongo.ErrNoDocuments {
		return chatRoom, nil
	} else if err != nil {
		log.Println(err)
		return chatRoom, err
	}

	return chatRoom, nil
}
