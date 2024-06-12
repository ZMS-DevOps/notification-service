package application

import (
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/dto"
	"log"
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

func (service *NotificationSettingsService) Insert(userId, role string) error {
	var notificationSettings []domain.NotificationSetting
	log.Printf("userId: %s, role: %s", userId, role)
	if role == domain.RoleGuest {
		notificationSettings = []domain.NotificationSetting{
			{
				Type:   4,
				Active: true,
			},
		}
	} else if role == domain.HostRole {
		notificationSettings = []domain.NotificationSetting{
			{
				Type:   0,
				Active: true,
			},
			{
				Type:   1,
				Active: true,
			},
			{
				Type:   2,
				Active: true,
			},
			{
				Type:   3,
				Active: true,
			},
		}
	}

	settings := domain.Settings{
		UserId:   userId,
		Settings: notificationSettings,
	}

	if _, err := service.store.Insert(&settings); err != nil {
		return err
	}

	log.Printf("userId: %s, role: %s", userId, role)
	return nil
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

func (service *NotificationSettingsService) UserIsSubscribedToNotificationType(userId string, notificationType domain.NotificationType) bool {
	settings, err := service.store.GetByUserId(userId)
	if err != nil {
		return false
	}

	return service.findIfUserHasSpecificActiveNotification(settings, notificationType)
}

func (service *NotificationSettingsService) findIfUserHasSpecificActiveNotification(settings *domain.Settings, notificationType domain.NotificationType) bool {
	for _, setting := range settings.Settings {
		if setting.Active && setting.Type == notificationType {
			return true
		}
	}
	return false
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
