package startup

import (
	"github.com/mmmajder/zms-devops-notification-service/domain"
)

var notificationSettings = []domain.NotificationSetting{
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
var settings = []*domain.Settings{
	{
		UserId:   "77475f63-6c4c-4c55-96f4-9f91ef41be09",
		Settings: notificationSettings,
	},
}
