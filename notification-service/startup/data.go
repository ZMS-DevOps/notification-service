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
		UserId:   "57325353-5469-4930-8ec9-35c003e1b967",
		Settings: notificationSettings,
	},
}
