package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type NotificationType int

const (
	NewReservationRequest NotificationType = iota
	CancelReservation
	NewHostReview
	NewAccommodationReview
	ReviewReservation
)

type NotificationSetting struct {
	Type   NotificationType `bson:"type"`
	Active bool             `bson:"active"`
}

type Settings struct {
	Id       primitive.ObjectID    `bson:"_id"`
	UserId   string                `bson:"user_id"`
	Settings []NotificationSetting `bson:"settings"`
}

type BellNotification struct {
	Id             primitive.ObjectID `bson:"_id"`
	UserId         string             `bson:"user_id"`
	Message        string             `bson:"message"`
	TimeStamp      time.Time          `bson:"time_stamp"`
	Seen           bool               `bson:"seen"`
	ShouldRedirect bool               `bson:"should_redirect"`
	RedirectId     string             `bson:"redirect_id,omitempty"`
}
