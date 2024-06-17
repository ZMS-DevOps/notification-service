package startup

import (
	"github.com/mmmajder/zms-devops-notification-service/domain"
)

var notificationSettingsForHost = []domain.NotificationSetting{
	{
		Type:   1,
		Active: true,
	},
	{
		Type:   0,
		Active: true,
	},
	{
		Type:   2,
		Active: false,
	},
	{
		Type:   3,
		Active: true,
	},
}

var notificationSettingsForGuest = []domain.NotificationSetting{
	{
		Type:   4,
		Active: true,
	},
}
var settings = []*domain.Settings{
	{
		UserId:   "3f92c83e-966d-41e6-8bb5-c076737d89ee",
		Settings: notificationSettingsForHost,
	},
	{
		UserId:   "f3c0120b-39f3-45cf-a771-e062c6932ce2",
		Settings: notificationSettingsForGuest,
	},
}
