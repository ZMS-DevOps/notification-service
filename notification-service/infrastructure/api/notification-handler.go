package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mmmajder/zms-devops-notification-service/application"
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/request"
	"github.com/mmmajder/zms-devops-notification-service/util"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"net/http"
)

type NotificationHandler struct {
	notificationService *application.BellNotificationService
	settingsService     *application.NotificationSettingsService
	upgrades            websocket.Upgrader
	connections         []*websocket.Conn
	traceProvider       *sdktrace.TracerProvider
	loki                promtail.Client
}

func NewNotificationHandler(bellService *application.BellNotificationService, settingsService *application.NotificationSettingsService, provider *sdktrace.TracerProvider, loki promtail.Client) *NotificationHandler {
	handler := NotificationHandler{
		notificationService: bellService,
		settingsService:     settingsService,
		upgrades: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		connections:   []*websocket.Conn{},
		traceProvider: provider,
		loki:          loki,
	}

	return &handler
}

func (handler *NotificationHandler) Init(router *mux.Router) {
	router.HandleFunc(domain.BellNotificationContextPath+domain.UserIDParam, handler.GetAllByUserId).Methods(http.MethodGet)
	router.HandleFunc(domain.BellNotificationContextPath+domain.UserIDParam+"/seen", handler.UpdateStatus).Methods(http.MethodPut)
	router.HandleFunc("/ws", handler.WebSocketHandler)
}

func (handler *NotificationHandler) OnNewReservationRequestCreated(message *kafka.Message) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(context.TODO(), "on-reservation-request-created")
	defer func() { span.End() }()
	notificationRequest := handler.getNotificationRequest(message)
	var textMessage = "You have a new reservation request."
	if notificationRequest.Status == "automatic" {
		textMessage = "Reservation request #" + notificationRequest.ReservationId + " is automatically accepted."
	}

	handler.onCreateNewNotification(notificationRequest, domain.ReservationRedirectUrlStart+notificationRequest.ReservationId, true, domain.NotificationType(0), textMessage)
}

func (handler *NotificationHandler) OnReservationCancellation(message *kafka.Message) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(context.TODO(), "on-reservation-cancellation")
	defer func() { span.End() }()
	notificationRequest := handler.getNotificationRequest(message)
	handler.onCreateNewNotification(notificationRequest, domain.ReservationRedirectUrlStart+notificationRequest.ReservationId, true, domain.NotificationType(1), "A reservation #"+notificationRequest.ReservationId+" has been cancelled.")
}

func (handler *NotificationHandler) OnHostRated(message *kafka.Message) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(context.TODO(), "on-host-rated")
	defer func() { span.End() }()
	notificationRequest := handler.getNotificationRequest(message)

	handler.onCreateNewNotification(notificationRequest, "auth/view-profile/"+notificationRequest.ReceiverId, true, domain.NotificationType(2), notificationRequest.StartActionUserName+" has reviewed your profile. Check for more details.")
}

func (handler *NotificationHandler) OnAccommodationRated(message *kafka.Message) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(context.TODO(), "on-accommodation-rated")
	defer func() { span.End() }()
	notificationRequest := handler.getNotificationRequest(message)

	handler.onCreateNewNotification(notificationRequest, "accommodation/"+notificationRequest.AccommodationId, true, domain.NotificationType(3), notificationRequest.StartActionUserName+" has reviewed your accommodation. Check for more details.")
}

func (handler *NotificationHandler) OnHostRespondedToReservationRequest(message *kafka.Message) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(context.TODO(), "on-host-responded-to-reservation-request")
	defer func() { span.End() }()
	notificationRequest := handler.getNotificationRequest(message)
	var textMessage = "The host has been canceled your reservation #" + notificationRequest.ReservationId
	if notificationRequest.Status == "accept-request" {
		textMessage = "The host has been confirmed your reservation #" + notificationRequest.ReservationId
	}

	handler.onCreateNewNotification(notificationRequest, domain.ReservationRedirectUrlStart+notificationRequest.ReservationId, true, domain.NotificationType(4), textMessage)
}

func (handler *NotificationHandler) getNotificationRequest(message *kafka.Message) request.NotificationMessageRequest {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(context.TODO(), "get-notification-request")
	defer func() { span.End() }()
	var notificationRequest request.NotificationMessageRequest
	if err := json.Unmarshal(message.Value, &notificationRequest); err != nil {
		util.HttpTraceError(err, "failed to unmarshal data", span, handler.loki, "getNotificationRequest", "")
		log.Printf("Failed to unmarshal notification message for host review: %s", err)
		return request.NotificationMessageRequest{}
	}

	return notificationRequest
}

func (handler *NotificationHandler) onCreateNewNotification(notification request.NotificationMessageRequest, redirectId string, shouldRedirect bool, notificationType domain.NotificationType, message string) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(context.TODO(), "on-create-new-notification")
	defer func() { span.End() }()
	subscribed := handler.settingsService.UserIsSubscribedToNotificationType(notification.ReceiverId, notificationType, span, handler.loki)
	if subscribed {
		notificationDTO, err := handler.notificationService.Add(notification.ReceiverId, message, redirectId, shouldRedirect, span, handler.loki)
		if err != nil {
			util.HttpTraceError(err, "failed to add notification", span, handler.loki, "onCreateNewNotification", "")
			return
		}
		jsonMessage, _ := json.Marshal(notificationDTO)
		handler.sendWebSocketMessage(jsonMessage)
	}
}

func (handler *NotificationHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := handler.upgrades.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	handler.connections = append(handler.connections, conn)
}

func (handler *NotificationHandler) GetAllByUserId(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-all-by-user-id-")
	defer func() { span.End() }()
	id := mux.Vars(r)["userId"]
	if id == "" {
		util.HttpTraceError(errors.New("user id can not be empty"), "user id can not be empty", span, handler.loki, "GetAllByUserId", "")
		handleError(w, http.StatusBadRequest, domain.InvalidIDErrorMessage)
		return
	}

	response, err := handler.notificationService.GetAllByUserId(id, span, handler.loki)
	if err != nil {
		util.HttpTraceError(err, "failed to get all notifications by id", span, handler.loki, "GetAllByUserId", "")
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, response)
}

func (handler *NotificationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "update-status-put")
	defer func() { span.End() }()
	id := mux.Vars(r)["userId"]
	if id == "" {
		util.HttpTraceError(errors.New("user id can not be empty"), "user id can not be empty", span, handler.loki, "UpdateStatus", "")
		handleError(w, http.StatusBadRequest, domain.InvalidIDErrorMessage)
		return
	}

	if err := handler.notificationService.UpdateStatus(id, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to update notification status", span, handler.loki, "UpdateStatus", "")
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusAccepted, nil)
}

func (handler *NotificationHandler) sendWebSocketMessage(jsonMessage []byte) {
	for _, conn := range handler.connections {
		err := conn.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			conn.Close()
		}
	}
}

func (handler *NotificationHandler) GetHealthCheck(w http.ResponseWriter, r *http.Request) {

	writeResponse(w, http.StatusOK, domain.HealthCheckMessage)
}
