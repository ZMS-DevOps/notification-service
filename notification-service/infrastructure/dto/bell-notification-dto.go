package dto

import (
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BellNotificationDTO struct {
	Id             primitive.ObjectID `json:"id"`
	UserId         string             `json:"userId"`
	Message        string             `json:"message"`
	TimeStamp      time.Time          `json:"timeStamp"`
	Seen           bool               `json:"seen"`
	ShouldRedirect bool               `json:"shouldRedirect"`
	RedirectId     string             `json:"redirectId"`
}

func FromReviews(notifications []*domain.BellNotification) *[]BellNotificationDTO {
	notificationDTOs := make([]BellNotificationDTO, 0, len(notifications))
	for _, notification := range notifications {
		dto := FromNotification(notification)
		notificationDTOs = append(notificationDTOs, dto)
	}
	return &notificationDTOs
}

func FromNotification(notification *domain.BellNotification) BellNotificationDTO {
	dto := BellNotificationDTO{
		Id:             notification.Id,
		UserId:         notification.UserId,
		Message:        notification.Message,
		TimeStamp:      notification.TimeStamp,
		Seen:           notification.Seen,
		ShouldRedirect: notification.ShouldRedirect,
		RedirectId:     notification.RedirectId,
	}
	return dto
}
