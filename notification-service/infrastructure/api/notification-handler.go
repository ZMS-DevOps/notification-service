package api

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mmmajder/zms-devops-notification-service/application"
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/request"
	"log"
	"net/http"
)

type NotificationHandler struct {
	notificationService *application.BellNotificationService
	settingsService     *application.NotificationSettingsService
	upgrades            websocket.Upgrader
	connections         []*websocket.Conn
}

func NewNotificationHandler(bellService *application.BellNotificationService, settingsService *application.NotificationSettingsService) *NotificationHandler {
	handler := NotificationHandler{
		notificationService: bellService,
		settingsService:     settingsService,
		upgrades: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		connections: []*websocket.Conn{},
	}

	return &handler
}

func (handler *NotificationHandler) Init(router *mux.Router) {
	router.HandleFunc(domain.BellNotificationContextPath+domain.UserIDParam, handler.GetAllByUserId).Methods(http.MethodGet)
	router.HandleFunc(domain.BellNotificationContextPath+domain.UserIDParam+"/seen", handler.UpdateStatus).Methods(http.MethodPut)
	router.HandleFunc("/ws", handler.WebSocketHandler)
}

func (handler *NotificationHandler) OnNewReservationRequestCreated(message *kafka.Message) {
	notificationRequest := handler.getNotificationRequest(message)
	var textMessage = "You have a new reservation request."
	if notificationRequest.Status == "automatic" {
		textMessage = "Reservation request #" + notificationRequest.ReservationId + " is automatically accepted."
	}

	handler.onCreateNewNotification(notificationRequest, domain.ReservationRedirectUrlStart+notificationRequest.ReservationId, true, domain.NotificationType(0), textMessage)
}

func (handler *NotificationHandler) OnReservationCancellation(message *kafka.Message) {
	notificationRequest := handler.getNotificationRequest(message)
	handler.onCreateNewNotification(notificationRequest, domain.ReservationRedirectUrlStart+notificationRequest.ReservationId, true, domain.NotificationType(1), "A reservation #"+notificationRequest.ReservationId+" has been cancelled.")
}

func (handler *NotificationHandler) OnHostRated(message *kafka.Message) {
	notificationRequest := handler.getNotificationRequest(message)

	handler.onCreateNewNotification(notificationRequest, "auth/view-profile/"+notificationRequest.ReceiverId, true, domain.NotificationType(2), notificationRequest.StartActionUserName+" has reviewed your profile. Check for more details.")
}

func (handler *NotificationHandler) OnAccommodationRated(message *kafka.Message) {
	notificationRequest := handler.getNotificationRequest(message)

	handler.onCreateNewNotification(notificationRequest, "accommodation/"+notificationRequest.ReceiverId, true, domain.NotificationType(3), notificationRequest.StartActionUserName+" has reviewed your accommodation. Check for more details.")
}

func (handler *NotificationHandler) OnHostRespondedToReservationRequest(message *kafka.Message) {
	notificationRequest := handler.getNotificationRequest(message)
	var textMessage = "The host has been canceled your reservation #" + notificationRequest.ReservationId
	if notificationRequest.Status == "accept-request" {
		textMessage = "The host has been confirmed your reservation #" + notificationRequest.ReservationId
	}

	handler.onCreateNewNotification(notificationRequest, domain.ReservationRedirectUrlStart+notificationRequest.ReservationId, true, domain.NotificationType(4), textMessage)
}

func (handler *NotificationHandler) getNotificationRequest(message *kafka.Message) request.NotificationMessageRequest {
	var notificationRequest request.NotificationMessageRequest
	if err := json.Unmarshal(message.Value, &notificationRequest); err != nil {
		log.Printf("Failed to unmarshal notification message for host review: %s", err)
		return request.NotificationMessageRequest{}
	}

	return notificationRequest
}

func (handler *NotificationHandler) onCreateNewNotification(notification request.NotificationMessageRequest, redirectId string, shouldRedirect bool, notificationType domain.NotificationType, message string) {

	subscribed := handler.settingsService.UserIsSubscribedToNotificationType(notification.ReceiverId, notificationType)
	if subscribed {
		notificationDTO, err := handler.notificationService.Add(notification.ReceiverId, message, redirectId, shouldRedirect)
		if err != nil {
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
	id := mux.Vars(r)["userId"]
	if id == "" {
		handleError(w, http.StatusBadRequest, domain.InvalidIDErrorMessage)
		return
	}

	response, err := handler.notificationService.GetAllByUserId(id)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, response)
}

func (handler *NotificationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["userId"]
	if id == "" {
		handleError(w, http.StatusBadRequest, domain.InvalidIDErrorMessage)
		return
	}

	if err := handler.notificationService.UpdateStatus(id); err != nil {
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
