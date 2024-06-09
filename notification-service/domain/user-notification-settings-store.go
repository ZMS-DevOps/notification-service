package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserNotificationSettingsStore interface {
	GetByUserId(id string) (*Settings, error)
	Insert(settings *Settings) (primitive.ObjectID, error)
	DeleteByUserId(id string) error
	DeleteAll()
	Update(id primitive.ObjectID, settings *Settings) error
}
