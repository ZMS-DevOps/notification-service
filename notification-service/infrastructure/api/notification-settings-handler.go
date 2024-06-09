package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mmmajder/zms-devops-notification-service/application"
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/request"
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
	router.HandleFunc(domain.NotificationContextPath+domain.UserIDParam, handler.GetSettings).Methods(http.MethodGet)
	router.HandleFunc(domain.NotificationContextPath+domain.UserIDParam, handler.UpdateSettings).Methods(http.MethodPut)
	router.HandleFunc(domain.NotificationContextPath+"/health/check", handler.GetHealthCheck).Methods(http.MethodGet)
}

func (handler *NotificationSettingsHandler) GetHealthCheck(w http.ResponseWriter, r *http.Request) {

	writeResponse(w, http.StatusOK, domain.HealthCheckMessage)
}

func (handler *NotificationSettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["userId"]
	if id == "" {
		handleError(w, http.StatusBadRequest, domain.InvalidIDErrorMessage)
		return
	}

	response, err := handler.settingsService.Get(id)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, response)
}

func (handler *NotificationSettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["userId"]
	if id == "" {
		handleError(w, http.StatusBadRequest, domain.InvalidIDErrorMessage)
		return
	}

	var settingsRequest request.UserNotificationSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&settingsRequest); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid user notification settings payload")
		return
	}

	if err := settingsRequest.AreValidRequestData(); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.settingsService.Update(id, settingsRequest.Role, request.FromSettingsRequests(settingsRequest.Settings)); err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusAccepted, nil)
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
