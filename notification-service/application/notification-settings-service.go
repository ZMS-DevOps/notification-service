package application

import (
	"github.com/afiskon/promtail-client/promtail"
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/dto"
	"github.com/mmmajder/zms-devops-notification-service/util"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
)

type NotificationSettingsService struct {
	store      domain.UserNotificationSettingsStore
	HttpClient *http.Client
	loki       promtail.Client
}

func NewNotificationSettingsService(store domain.UserNotificationSettingsStore, httpClient *http.Client, loki promtail.Client) *NotificationSettingsService {
	return &NotificationSettingsService{
		store:      store,
		HttpClient: httpClient,
		loki:       loki,
	}
}

func (service *NotificationSettingsService) Insert(userId, role string, span trace.Span, loki promtail.Client) error {
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

	util.HttpTraceInfo("Inserting review...", span, loki, "Insert", "")
	if _, err := service.store.Insert(&settings); err != nil {
		return err
	}

	log.Printf("userId: %s, role: %s", userId, role)
	return nil
}

func (service *NotificationSettingsService) Update(userId string, userRole string, settingsRequest []domain.NotificationSetting, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Inserting review...", span, loki, "Update", "")
	settings, err := service.store.GetByUserId(userId)
	if err != nil {
		return err
	}
	if userRole == domain.HostRole {
		settings.Settings = service.excludeNotificationSettings(settingsRequest, domain.ReviewReservation)
	} else {
		settings.Settings = service.filterNotificationSettings(settingsRequest, domain.ReviewReservation)
	}
	util.HttpTraceInfo("Inserting review...", span, loki, "Update", "")
	err = service.store.Update(settings.Id, settings)
	if err != nil {
		return err
	}

	return nil
}

func (service *NotificationSettingsService) Get(userId string, span trace.Span, loki promtail.Client) (*[]dto.NotificationSettingDTO, error) {
	util.HttpTraceInfo("Inserting review...", span, loki, "Get", "")
	settings, err := service.store.GetByUserId(userId)
	if err != nil {
		return nil, err
	}

	return dto.FromUserNotificationSettings(settings), nil
}

func (service *NotificationSettingsService) Delete(userId string, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Inserting review...", span, loki, "Delete", "")
	err := service.store.DeleteByUserId(userId)
	if err != nil {
		return err
	}

	return nil
}

func (service *NotificationSettingsService) UserIsSubscribedToNotificationType(userId string, notificationType domain.NotificationType, span trace.Span, loki promtail.Client) bool {
	util.HttpTraceInfo("Inserting review...", span, loki, "UserIsSubscribedToNotificationType", "")
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
