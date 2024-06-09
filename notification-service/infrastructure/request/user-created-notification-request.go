package request

import (
	"github.com/go-playground/validator/v10"
)

type UserCreatedNotificationRequest struct {
	UserId string `json:"user_id" validate:"required"`
	Role   string `json:"role" validate:"required"`
}

func (request UserCreatedNotificationRequest) AreValidRequestData() error {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}
