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

type UserRepository interface {
	CreateFriendRequest(ctx context.Context, senderID string, receiverID string) (friendReq model.FriendRequests, err error)
	GetFriendRequests(ctx context.Context, userID string) (friendrequests []model.FriendRequests, err error)
	GetFriendRequestByRequestID(ctx context.Context, requestID string) (friendReq model.FriendRequests, err error)
	UpdateFriendRequest(ctx context.Context, friendRequest model.FriendRequests, status string) (updatedFriendRequest model.FriendRequests, err error)
	UpdateFriendList(ctx context.Context, userID string, friendID string) (err error)
	GetFriendLists(ctx context.Context, userID string) (friends []string, err error)
}

type UserRepositoryImpl struct {
	mongo *mongo.Client
}

func NewUserRepository(mongo *mongo.Client) UserRepository {
	return &UserRepositoryImpl{mongo: mongo}
}

func (u *UserRepositoryImpl) CreateFriendRequest(ctx context.Context, senderID string, receiverID string) (friendReq model.FriendRequests, err error) {
	collection := u.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("FriendRequests")

	timeNow := time.Now()

	res, err := collection.InsertOne(ctx, bson.M{
		"sender_id":   senderID,
		"receiver_id": receiverID,
		"status":      "pending",
		"created_at":  timeNow,
		"updated_at":  time.Time{},
	})
	if err != nil {
		return model.FriendRequests{}, err
	}

	friendReq.RequestID = res.InsertedID.(primitive.ObjectID)
	friendReq.SenderID = senderID
	friendReq.ReceiverID = receiverID
	friendReq.Status = "pending"
	friendReq.CreatedAt = timeNow
	friendReq.UpdatedAt = time.Time{}

	return friendReq, nil
}

func (u *UserRepositoryImpl) GetFriendRequests(ctx context.Context, userID string) (friendRequests []model.FriendRequests, err error) {
	collection := u.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("FriendRequests")

	filter := bson.D{
		{"receiver_id", userID},
		{"status", "pending"},
	}

	opts := options.FindOptions{
		Sort: bson.D{{"created_at", 1}},
	}

	cur, err := collection.Find(ctx, filter, &opts)
	if err != nil {
		log.Println(err)
		return []model.FriendRequests{}, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var friendReq model.FriendRequests
		err := cur.Decode(&friendReq)
		if err != nil {
			log.Println("fail to decode")
			return []model.FriendRequests{}, err
		}
		friendRequests = append(friendRequests, friendReq)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return friendRequests, cur.Err()
}

func (u *UserRepositoryImpl) GetFriendRequestByRequestID(ctx context.Context, requestID string) (friendReq model.FriendRequests, err error) {
	collection := u.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("FriendRequests")

	objectID, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		log.Println(err)
		return model.FriendRequests{}, err
	}

	filter := bson.D{{"_id", objectID}}

	err = collection.FindOne(ctx, filter).Decode(&friendReq)
	if err == mongo.ErrNoDocuments {
		return model.FriendRequests{}, nil
	} else if err != nil {
		log.Println(err)
		return model.FriendRequests{}, err
	}

	return friendReq, nil
}

func (u *UserRepositoryImpl) UpdateFriendRequest(ctx context.Context, friendRequest model.FriendRequests, status string) (updatedFriendRequest model.FriendRequests, err error) {
	collection := u.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("FriendRequests")

	filter := bson.D{{"_id", friendRequest.RequestID}}
	update := bson.D{{"$set", bson.D{{"status", status}, {"updated_at", time.Now()}}}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Failed to update data in database: ", err)
		return model.FriendRequests{}, err
	}

	if result.MatchedCount == 0 {
		return model.FriendRequests{}, mongo.ErrNoDocuments
	}

	err = collection.FindOne(ctx, filter).Decode(&updatedFriendRequest)
	if err != nil {
		log.Println("Failed to retrieve updated request: ", err)
		return model.FriendRequests{}, err
	}

	return updatedFriendRequest, nil
}

func (u *UserRepositoryImpl) UpdateFriendList(ctx context.Context, userID string, friendID string) (err error) {
	collection := u.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("Users")

	// update sender friend list
	initFilter := bson.D{
		{"user_id", userID},
		{"$or", bson.A{
			bson.D{{"friends", bson.D{{"$exists", false}}}},
			bson.D{{"friends", nil}},
		}},
	}
	initUpdate := bson.D{
		{"$set", bson.D{{"friends", bson.A{}}}},
	}

	_, err = collection.UpdateOne(ctx, initFilter, initUpdate)
	if err != nil {
		log.Println(err)
		return err
	}

	filter := bson.D{{"user_id", userID}}
	update := bson.D{
		{"$push", bson.D{{"friends", friendID}}},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Failed to update data in database: ", err)
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	if result.ModifiedCount == 0 {
		log.Printf("Document matched but not modified for userID: %s\n", userID)
	}

	// update receiver friend list
	initReceiverFilter := bson.D{
		{"user_id", friendID},
		{"$or", bson.A{
			bson.D{{"friends", bson.D{{"$exists", false}}}},
			bson.D{{"friends", nil}},
		}},
	}
	initReceiverUpdate := bson.D{
		{"$set", bson.D{{"friends", bson.A{}}}},
	}

	_, err = collection.UpdateOne(ctx, initReceiverFilter, initReceiverUpdate)
	if err != nil {
		log.Println(err)
		return err
	}

	receiverFilter := bson.D{{"user_id", friendID}}
	receiverUpdate := bson.D{
		{"$push", bson.D{{"friends", userID}}},
	}

	updateReceiverresult, err := collection.UpdateOne(ctx, receiverFilter, receiverUpdate)
	if err != nil {
		log.Println("Failed to update data in database: ", err)
		return err
	}

	if updateReceiverresult.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	if updateReceiverresult.ModifiedCount == 0 {
		log.Printf("Document matched but not modified for userID: %s\n", friendID)
	}

	return nil
}

func (u *UserRepositoryImpl) GetFriendLists(ctx context.Context, userID string) (friends []string, err error) {
	collection := u.mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("Users")

	filter := bson.D{{"user_id", userID}}
	opts := options.FindOne().SetProjection(bson.D{{"friends", 1}})

	result := model.FriendList{}

	err = collection.FindOne(ctx, filter, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return []string{}, nil
	} else if err != nil {
		log.Println(err)
		return []string{}, err
	}

	friends = result.Friends

	return friends, nil
}
