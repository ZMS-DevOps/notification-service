package api

import (
	"github.com/gorilla/mux"
	"github.com/mmmajder/zms-devops-notification-service/application"
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"net/http"
)

type NotificationSettingsHandler struct {
	settingsService *application.NotificationSettingsService
}

func NewNotificationSettingsHandler(settingsService *application.NotificationSettingsService) *NotificationSettingsHandler {
	return &NotificationSettingsHandler{
		settingsService: settingsService,
	}
}

func (handler *NotificationSettingsHandler) Init(router *mux.Router) {
	router.HandleFunc(domain.NotificationContextPath+"/health/check", handler.GetHealthCheck).Methods(http.MethodGet)
}

func (handler *NotificationSettingsHandler) GetHealthCheck(w http.ResponseWriter, r *http.Request) {

	writeResponse(w, http.StatusOK, domain.HealthCheckMessage)
}

func (handler *NotificationSettingsHandler) DeleteSettings(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["userId"]
	if id == "" {
		handleError(w, http.StatusBadRequest, domain.InvalidIDErrorMessage)
		return
	}

	if err := handler.settingsService.Delete(id); err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, nil)
}
