package startup

import (
	"fmt"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/gorilla/mux"
	"github.com/mmmajder/zms-devops-notification-service/application"
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/api"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/persistence"
	"github.com/mmmajder/zms-devops-notification-service/startup/config"
	"go.mongodb.org/mongo-driver/mongo"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"net/http"
)

type Server struct {
	config              *config.Config
	router              *mux.Router
	SettingsHandler     *api.NotificationSettingsHandler
	NotificationHandler *api.NotificationHandler
	traceProvider       *sdktrace.TracerProvider
	loki                promtail.Client
}

func NewServer(config *config.Config, traceProvider *sdktrace.TracerProvider, loki promtail.Client) *Server {
	handler := &Server{
		config:        config,
		router:        mux.NewRouter(),
		traceProvider: traceProvider,
		loki:          loki,
	}

	notificationHandler, settingsHandler := handler.setupHandlers()
	handler.NotificationHandler = notificationHandler
	handler.SettingsHandler = settingsHandler
	return handler
}

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.router))
}

func (server *Server) setupHandlers() (*api.NotificationHandler, *api.NotificationSettingsHandler) {
	mongoClient := server.initMongoClient()
	settingsStore := server.initNotificationSettingsStore(mongoClient)
	settingsService := server.initSettingsService(settingsStore)

	bellNotificationStore := server.initBellNotificationStore(mongoClient)
	bellNotificationService := server.initBellNotificationService(bellNotificationStore)
	notificationHandler := server.initNotificationHandler(bellNotificationService, settingsService)

	settingsHandler := server.initSettingsHandler(settingsService)
	settingsHandler.Init(server.router)
	notificationHandler.Init(server.router)

	return notificationHandler, settingsHandler
}

func (server *Server) initSettingsService(store domain.UserNotificationSettingsStore) *application.NotificationSettingsService {

	return application.NewNotificationSettingsService(store, &http.Client{}, server.loki)
}

func (server *Server) initSettingsHandler(settingsService *application.NotificationSettingsService) *api.NotificationSettingsHandler {
	return api.NewNotificationSettingsHandler(settingsService, server.traceProvider, server.loki)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.DBUsername, server.config.DBPassword, server.config.DBHost, server.config.DBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initBellNotificationStore(client *mongo.Client) domain.BellNotificationStore {
	return persistence.NewBellNotificationMongoDBStore(client)
}

func (server *Server) initNotificationSettingsStore(client *mongo.Client) domain.UserNotificationSettingsStore {
	store := persistence.NewNotificationSettingsMongoDBStore(client)
	for _, setting := range settings {
		_, _ = store.Insert(setting)
	}
	return store
}

func (server *Server) initBellNotificationService(store domain.BellNotificationStore) *application.BellNotificationService {
	return application.NewBellNotificationService(store, &http.Client{}, server.loki)
}

func (server *Server) initNotificationHandler(service *application.BellNotificationService, settingsService *application.NotificationSettingsService) *api.NotificationHandler {
	return api.NewNotificationHandler(service, settingsService, server.traceProvider, server.loki)
}
