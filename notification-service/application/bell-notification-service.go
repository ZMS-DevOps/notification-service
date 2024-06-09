package application

import (
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/dto"
	"net/http"
	"time"
)

type BellNotificationService struct {
	store      domain.BellNotificationStore
	HttpClient *http.Client
}

func NewBellNotificationService(store domain.BellNotificationStore, httpClient *http.Client) *BellNotificationService {
	return &BellNotificationService{
		store:      store,
		HttpClient: httpClient,
	}
}

func (service *BellNotificationService) Add(userId, message, redirectId string, shouldRedirect bool) (dto.BellNotificationDTO, error) {
	notification := &domain.BellNotification{
		UserId:         userId,
		Message:        message,
		TimeStamp:      time.Now(),
		Seen:           false,
		ShouldRedirect: shouldRedirect,
		RedirectId:     redirectId,
	}

	id, err := service.store.Insert(notification)
	if err != nil {
		return dto.BellNotificationDTO{}, err
	}

	notificationDTO := dto.FromNotification(notification)
	notificationDTO.Id = id

	return notificationDTO, nil
}

func (service *BellNotificationService) GetAllByUserId(userId string) ([]dto.BellNotificationDTO, error) {
	response, err := service.store.GetAllByUserId(userId)
	if err != nil {
		return []dto.BellNotificationDTO{}, err
	}

	return *dto.FromReviews(response), nil
}

func (service *BellNotificationService) UpdateStatus(userId string) error {
	if err := service.store.UpdateManyStatus(userId); err != nil {
		return err
	}

	return nil
}
