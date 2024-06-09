package api

import (
	"encoding/json"
	"github.com/mmmajder/zms-devops-notification-service/domain"
	"log"
	"net/http"
)

func handleError(w http.ResponseWriter, httpStatus int, message string) {
	w.WriteHeader(httpStatus)
	if _, err := w.Write([]byte(message)); err != nil {
		log.Printf("error writing response: %v", err)
	}
}

func writeResponse(w http.ResponseWriter, httpStatus int, data interface{}) {
	w.Header().Set(domain.ContentType, domain.JsonContentType)
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
	}
}
