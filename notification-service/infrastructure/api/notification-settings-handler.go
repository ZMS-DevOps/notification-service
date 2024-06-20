package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/mux"
	"github.com/mmmajder/zms-devops-notification-service/application"
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/request"
	"github.com/mmmajder/zms-devops-notification-service/util"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"net/http"
)

type NotificationSettingsHandler struct {
	settingsService *application.NotificationSettingsService
	traceProvider   *sdktrace.TracerProvider
	loki            promtail.Client
}

func NewNotificationSettingsHandler(settingsService *application.NotificationSettingsService, traceProvider *sdktrace.TracerProvider, loki promtail.Client) *NotificationSettingsHandler {
	return &NotificationSettingsHandler{
		settingsService: settingsService,
		traceProvider:   traceProvider,
		loki:            loki,
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
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-settings-get")
	defer func() { span.End() }()
	id := mux.Vars(r)["userId"]
	if id == "" {
		util.HttpTraceError(errors.New("user id can not be empty"), "user id can not be empty", span, handler.loki, "GetSettings", "")
		handleError(w, http.StatusBadRequest, domain.InvalidIDErrorMessage)
		return
	}

	response, err := handler.settingsService.Get(id, span, handler.loki)
	if err != nil {
		util.HttpTraceError(err, "failed to get settings", span, handler.loki, "GetSettings", "")
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.HttpTraceInfo("Settings fetched successfully", span, handler.loki, "AddRequest", "")

	writeResponse(w, http.StatusOK, response)
}

func (handler *NotificationSettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "update-settings-put")
	defer func() { span.End() }()
	id := mux.Vars(r)["userId"]
	if id == "" {
		util.HttpTraceError(errors.New("user id can not be empty"), "user id can not be empty", span, handler.loki, "UpdateSettings", "")
		handleError(w, http.StatusBadRequest, domain.InvalidIDErrorMessage)
		return
	}

	var settingsRequest request.UserNotificationSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&settingsRequest); err != nil {
		util.HttpTraceError(err, "invalid request payload", span, handler.loki, "UpdateSettings", "")
		handleError(w, http.StatusBadRequest, "Invalid user notification settings payload")
		return
	}

	if err := settingsRequest.AreValidRequestData(); err != nil {
		util.HttpTraceError(err, "invalid request date", span, handler.loki, "UpdateSettings", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.settingsService.Update(id, settingsRequest.Role, request.FromSettingsRequests(settingsRequest.Settings), span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to update settings", span, handler.loki, "UpdateSettings", "")
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.HttpTraceInfo("Settings updated successfully", span, handler.loki, "AddRequest", "")

	writeResponse(w, http.StatusAccepted, nil)
}

func (handler *NotificationSettingsHandler) DeleteSettings(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "delete-settings-delete")
	defer func() { span.End() }()
	id := mux.Vars(r)["userId"]
	if id == "" {
		util.HttpTraceError(errors.New("user id can not be empty"), "user id can not be empty", span, handler.loki, "DeleteSettings", "")
		handleError(w, http.StatusBadRequest, domain.InvalidIDErrorMessage)
		return
	}

	if err := handler.settingsService.Delete(id, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to delete settings", span, handler.loki, "DeleteSettings", "")
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.HttpTraceInfo("Settings deleted successfully", span, handler.loki, "AddRequest", "")

	writeResponse(w, http.StatusOK, nil)
}

func (handler *NotificationSettingsHandler) OnUserCreated(message *kafka.Message) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(context.TODO(), "on-user-created")
	defer func() { span.End() }()
	var userCreatedNotificationRequest request.UserCreatedNotificationRequest
	if err := json.Unmarshal(message.Value, &userCreatedNotificationRequest); err != nil {
		util.HttpTraceError(err, "failed to unmarshal data", span, handler.loki, "OnUserCreated", "")
		log.Printf("Error unmarshalling user created user request: %v", err)
	}

	handler.onCreateUserNotification(userCreatedNotificationRequest)
	util.HttpTraceInfo("On user created settings created", span, handler.loki, "AddRequest", "")

}

func (handler *NotificationSettingsHandler) onCreateUserNotification(createdUser request.UserCreatedNotificationRequest) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(context.TODO(), "on-create-user-notification")
	defer func() { span.End() }()
	if err := handler.settingsService.Insert(createdUser.UserId, createdUser.Role, span, handler.loki); err != nil {

		return
	}
}
