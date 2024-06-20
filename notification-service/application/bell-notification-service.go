package application

import (
	"github.com/afiskon/promtail-client/promtail"
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"github.com/mmmajder/zms-devops-notification-service/infrastructure/dto"
	"github.com/mmmajder/zms-devops-notification-service/util"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"
)

type BellNotificationService struct {
	store      domain.BellNotificationStore
	HttpClient *http.Client
	loki       promtail.Client
}

func NewBellNotificationService(store domain.BellNotificationStore, httpClient *http.Client, loki promtail.Client) *BellNotificationService {
	return &BellNotificationService{
		store:      store,
		HttpClient: httpClient,
		loki:       loki,
	}
}

func (service *BellNotificationService) Add(userId, message, redirectId string, shouldRedirect bool, span trace.Span, loki promtail.Client) (dto.BellNotificationDTO, error) {
	notification := &domain.BellNotification{
		UserId:         userId,
		Message:        message,
		TimeStamp:      time.Now(),
		Seen:           false,
		ShouldRedirect: shouldRedirect,
		RedirectId:     redirectId,
	}
	util.HttpTraceInfo("Inserting review...", span, loki, "Add", "")

	id, err := service.store.Insert(notification)
	if err != nil {
		return dto.BellNotificationDTO{}, err
	}

	notificationDTO := dto.FromNotification(notification)
	notificationDTO.Id = id

	return notificationDTO, nil
}

func (service *BellNotificationService) GetAllByUserId(userId string, span trace.Span, loki promtail.Client) ([]dto.BellNotificationDTO, error) {
	util.HttpTraceInfo("Inserting review...", span, loki, "GetAllByUserId", "")
	response, err := service.store.GetAllByUserId(userId)
	if err != nil {
		return []dto.BellNotificationDTO{}, err
	}

	return *dto.FromReviews(response), nil
}

func (service *BellNotificationService) UpdateStatus(userId string, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Inserting review...", span, loki, "UpdateStatus", "")
	if err := service.store.UpdateManyStatus(userId); err != nil {
		return err
	}

	return nil
}
