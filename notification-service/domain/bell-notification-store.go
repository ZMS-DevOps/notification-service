package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BellNotificationStore interface {
	GetAllByUserId(userId string) ([]*BellNotification, error)
	Insert(review *BellNotification) (primitive.ObjectID, error)
	UpdateStatus(id primitive.ObjectID, review *BellNotification) error
	UpdateManyStatus(userId string) error
}
