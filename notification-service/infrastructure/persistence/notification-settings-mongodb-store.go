package persistence

import (
	"context"
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "notificationdb"
	COLLECTION = "setting"
)

type NotificationSettingsMongoDBStore struct {
	settings *mongo.Collection
}

func NewNotificationSettingsMongoDBStore(client *mongo.Client) domain.UserNotificationSettingsStore {
	settings := client.Database(DATABASE).Collection(COLLECTION)
	return &NotificationSettingsMongoDBStore{
		settings: settings,
	}
}

func (store *NotificationSettingsMongoDBStore) GetByUserId(id string) (*domain.Settings, error) {
	filter := bson.M{"user_id": id}
	return store.filterOne(filter)
}

func (store *NotificationSettingsMongoDBStore) Insert(settings *domain.Settings) (primitive.ObjectID, error) {
	settings.Id = primitive.NewObjectID()
	result, err := store.settings.InsertOne(context.TODO(), settings)
	if err != nil {
		return primitive.NilObjectID, err
	}
	settings.Id = result.InsertedID.(primitive.ObjectID)
	return settings.Id, nil
}

func (store *NotificationSettingsMongoDBStore) DeleteByUserId(id string) error {
	filter := bson.M{"user_id": id}
	_, err := store.settings.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (store *NotificationSettingsMongoDBStore) DeleteAll() {
	store.settings.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *NotificationSettingsMongoDBStore) Update(id primitive.ObjectID, notificationSettings *domain.Settings) error {
	filter := bson.M{"_id": id}
	update := bson.D{
		{"$set", bson.D{
			{"settings", notificationSettings.Settings},
		}},
	}
	_, err := store.settings.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (store *NotificationSettingsMongoDBStore) filter(filter interface{}) ([]*domain.Settings, error) {
	cursor, err := store.settings.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *NotificationSettingsMongoDBStore) filterOne(filter interface{}) (notification *domain.Settings, err error) {
	result := store.settings.FindOne(context.TODO(), filter)
	err = result.Decode(&notification)
	return
}

func decode(cursor *mongo.Cursor) (notificationSettings []*domain.Settings, err error) {
	for cursor.Next(context.TODO()) {
		var userNotificationSetting domain.Settings
		err = cursor.Decode(&userNotificationSetting)
		if err != nil {
			return
		}
		notificationSettings = append(notificationSettings, &userNotificationSetting)
	}
	err = cursor.Err()
	return
}
