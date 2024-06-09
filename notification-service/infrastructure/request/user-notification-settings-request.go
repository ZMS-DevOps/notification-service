package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/mmmajder/zms-devops-notification-service/domain"
)

type NotificationSettingRequest struct {
	Type   int  `json:"type" validate:"min=0,max=4"`
	Active bool `json:"active"`
}

type UserNotificationSettingsRequest struct {
	Settings []NotificationSettingRequest `json:"settings" validate:"required"`
	Role     string                       `json:"role" validate:"required"`
}

func (request UserNotificationSettingsRequest) AreValidRequestData() error {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func FromSettingsRequests(settingsRequests []NotificationSettingRequest) []domain.NotificationSetting {
	settings := make([]domain.NotificationSetting, 0, len(settingsRequests))
	for _, setting := range settingsRequests {
		settings = append(settings, FromSettingsRequest(setting))
	}
	return settings
}

func FromSettingsRequest(settingsRequest NotificationSettingRequest) domain.NotificationSetting {
	return domain.NotificationSetting{
		Type:   domain.NotificationType(settingsRequest.Type),
		Active: settingsRequest.Active,
	}
}
