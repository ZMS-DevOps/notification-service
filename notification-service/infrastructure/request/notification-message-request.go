package request

import (
	"github.com/go-playground/validator/v10"
)

type NotificationMessageRequest struct {
	ReceiverId          string `json:"receiver_id" validate:"required"`
	Status              string `json:"status"`
	StartActionUserName string `json:"start_action_user_name" validate:"omitempty"`
	AccommodationId     string `json:"accommodation_id" validate:"omitempty"`
	ReservationId       string `json:"reservation_id" validate:"omitempty"`
}

func (request NotificationMessageRequest) AreValidRequestData() error {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}
