package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"planmasta.com/internals/dto"
	"planmasta.com/internals/service"
	"time"
)

type ReplicateHandler struct {
	service *service.ReplicateService
	logger  *slog.Logger
}

func NewReplicateHandler(service *service.ReplicateService, logger *slog.Logger) *ReplicateHandler {
	return &ReplicateHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ReplicateHandler) Generate(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	var dto dto.GenerateRequest
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	response, err := h.service.SendRequest(dto)
	if err != nil {
		h.logger.Error("Failed to process request", slog.String("err", err.Error()))
		response, _ := json.Marshal(map[string]string{"error": "Failed to process request"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("Failed to process request", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	h.logger.Info(
		"Sent response",
		slog.Int64("duration", time.Since(startTime).Milliseconds()),
	)
}
