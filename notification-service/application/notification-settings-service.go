package application

import (
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"net/http"
)

type NotificationSettingsService struct {
	store      domain.UserNotificationSettingsStore
	HttpClient *http.Client
}

func NewNotificationSettingsService(store domain.UserNotificationSettingsStore, httpClient *http.Client) *NotificationSettingsService {
	return &NotificationSettingsService{
		store:      store,
		HttpClient: httpClient,
	}
}

func (service *NotificationSettingsService) Delete(userId string) error {
	err := service.store.DeleteByUserId(userId)
	if err != nil {
		return err
	}

	return nil
}
