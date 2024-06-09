package dto

import (
	"github.com/mmmajder/zms-devops-notification-service/domain"
)

type NotificationSettingDTO struct {
	Type   domain.NotificationType `json:"type"`
	Active bool                    `json:"active"`
}

func FromUserNotificationSettings(settings *domain.Settings) *[]NotificationSettingDTO {
	settingsDTOs := make([]NotificationSettingDTO, 0, len(settings.Settings))
	for _, setting := range settings.Settings {
		dto := FromSetting(setting)
		settingsDTOs = append(settingsDTOs, dto)
	}
	return &settingsDTOs
}

func FromSetting(setting domain.NotificationSetting) NotificationSettingDTO {
	return NotificationSettingDTO{
		Type:   setting.Type,
		Active: setting.Active,
	}
}
