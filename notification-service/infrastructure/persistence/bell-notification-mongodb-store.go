package persistence

import (
	"context"
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BellNotificationMongoDBStore struct {
	notifications *mongo.Collection
}

func NewBellNotificationMongoDBStore(client *mongo.Client) domain.BellNotificationStore {
	notifications := client.Database("notificationdb").Collection("notifications")
	return &BellNotificationMongoDBStore{
		notifications: notifications,
	}
}

func (store *BellNotificationMongoDBStore) GetAllByUserId(userId string) ([]*domain.BellNotification, error) {
	filter := bson.M{"user_id": userId}

	return store.filter(filter)
}

func (store *BellNotificationMongoDBStore) Insert(notification *domain.BellNotification) (primitive.ObjectID, error) {
	notification.Id = primitive.NewObjectID()
	result, err := store.notifications.InsertOne(context.TODO(), notification)
	if err != nil {
		return primitive.NilObjectID, err
	}
	notification.Id = result.InsertedID.(primitive.ObjectID)
	return notification.Id, nil
}

func (store *BellNotificationMongoDBStore) UpdateStatus(id primitive.ObjectID, notification *domain.BellNotification) error {
	filter := bson.M{"_id": id}
	update := bson.D{
		{"$set", bson.D{
			{"seen", notification.Seen},
		}},
	}
	_, err := store.notifications.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (store *BellNotificationMongoDBStore) UpdateManyStatus(userId string) error {
	filter := bson.M{"user_id": userId}
	update := bson.D{
		{"$set", bson.D{
			{"seen", true},
		}},
	}
	_, err := store.notifications.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (store *BellNotificationMongoDBStore) filter(filter interface{}) ([]*domain.BellNotification, error) {
	cursor, err := store.notifications.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return store.decode(cursor)
}

func (store *BellNotificationMongoDBStore) filterOne(filter interface{}) (notification *domain.BellNotification, err error) {
	result := store.notifications.FindOne(context.TODO(), filter)
	err = result.Decode(&notification)
	return
}

func (store *BellNotificationMongoDBStore) decode(cursor *mongo.Cursor) (notifications []*domain.BellNotification, err error) {
	for cursor.Next(context.TODO()) {
		var notification domain.BellNotification
		err = cursor.Decode(&notification)
		if err != nil {
			return
		}
		notifications = append(notifications, &notification)
	}
	err = cursor.Err()
	return
}
