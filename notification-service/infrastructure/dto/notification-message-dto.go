package dto

type NotificationMessageDTO struct {
	Message         string `json:"message"`
	UserId          string `json:"userId"`
	Status          string `json:"status"`
	AccommodationId string `json:"accommodationId"`
	Name            string `json:"name"`
}
