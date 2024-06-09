package application

import (
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/dto"
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

func (service *NotificationSettingsService) Update(userId string, userRole string, settingsRequest []domain.NotificationSetting) error {
	settings, err := service.store.GetByUserId(userId)
	if err != nil {
		return err
	}
	if userRole == domain.HostRole {
		settings.Settings = service.excludeNotificationSettings(settingsRequest, domain.ReviewReservation)
	} else {
		settings.Settings = service.filterNotificationSettings(settingsRequest, domain.ReviewReservation)
	}
	err = service.store.Update(settings.Id, settings)
	if err != nil {
		return err
	}

	return nil
}

func (service *NotificationSettingsService) Get(userId string) (*[]dto.NotificationSettingDTO, error) {
	settings, err := service.store.GetByUserId(userId)
	if err != nil {
		return nil, err
	}

	return dto.FromUserNotificationSettings(settings), nil
}

func (service *NotificationSettingsService) Delete(userId string) error {
	err := service.store.DeleteByUserId(userId)
	if err != nil {
		return err
	}

	return nil
}

func (service *NotificationSettingsService) filterNotificationSettings(settings []domain.NotificationSetting, notificationType domain.NotificationType) []domain.NotificationSetting {
	var filteredSettings []domain.NotificationSetting

	for _, setting := range settings {
		if setting.Type == notificationType {
			filteredSettings = append(filteredSettings, setting)
		}
	}

	return filteredSettings
}

func (service *NotificationSettingsService) excludeNotificationSettings(settings []domain.NotificationSetting, notificationType domain.NotificationType) []domain.NotificationSetting {
	var filteredSettings []domain.NotificationSetting

	for _, setting := range settings {
		if setting.Type != notificationType {
			filteredSettings = append(filteredSettings, setting)
		}
	}

	return filteredSettings
}
