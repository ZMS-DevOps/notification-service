package domain

const (
	NotificationContextPath     string = "/notification"
	BellNotificationContextPath string = "/notification/bell"
	ContentType                 string = "Content-Type"
	JsonContentType             string = "application/json"
	HealthCheckMessage          string = "NOTIFICATION SERVICE IS HEALTH"
	InvalidIDErrorMessage       string = "Invalid user ID"
	HostRole                    string = "host"
	UserIDParam                 string = "/{userId}"
	ReservationRedirectUrlStart string = "reservation/view/"
	RoleGuest                   string = "guest"
	RoleHost                    string = "host"
)
