package startup

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mmmajder/zms-devops-notification-service/application"
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/api"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/persistence"
	"github.com/mmmajder/zms-devops-notification-service/startup/config"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type Server struct {
	config          *config.Config
	router          *mux.Router
	SettingsHandler *api.NotificationSettingsHandler
}

func NewServer(config *config.Config) *Server {
	handler := &Server{
		config: config,
		router: mux.NewRouter(),
	}

	settingsHandler := handler.setupHandlers()
	handler.SettingsHandler = settingsHandler
	return handler
}

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.router))
}

func (server *Server) setupHandlers() *api.NotificationSettingsHandler {
	mongoClient := server.initMongoClient()
	settingsStore := server.initNotificationSettingsStore(mongoClient)
	settingsService := server.initSettingsService(settingsStore)

	settingsHandler := server.initSettingsHandler(settingsService)
	settingsHandler.Init(server.router)

	return settingsHandler
}

func (server *Server) initSettingsService(store domain.UserNotificationSettingsStore) *application.NotificationSettingsService {

	return application.NewNotificationSettingsService(store, &http.Client{})
}

func (server *Server) initSettingsHandler(settingsService *application.NotificationSettingsService) *api.NotificationSettingsHandler {
	return api.NewNotificationSettingsHandler(settingsService)
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
	store.DeleteAll()
	for _, setting := range settings {
		_, _ = store.Insert(setting)
	}
	return store
}
