package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"planmasta.com/internals/service"
	"strconv"
	"time"
)

type OpenAIHandler struct {
	service *service.OpenAIService
	logger  *slog.Logger
}

func NewOpenAIHandler(srv *service.OpenAIService, logger *slog.Logger) *OpenAIHandler {
	return &OpenAIHandler{
		service: srv,
		logger:  logger,
	}
}

func (h *OpenAIHandler) Chat(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	response, statusCode, err := h.service.SendRequest(r.Body)
	if err != nil {
		h.logger.Error("Failed to process request", slog.String("err", err.Error()))
		response, _ := json.Marshal(map[string]string{"error": "Failed to process request"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(response)
		return
	}

	w.WriteHeader(statusCode)
	bytesCopied, err := w.Write(response)
	if err != nil {
		h.logger.Error("Failed to write response", slog.String("err", err.Error()))
		return
	}

	h.logger.Info(
		"Sent response",
		slog.String("size", strconv.Itoa(bytesCopied)),
		slog.Int64("duration", time.Since(startTime).Milliseconds()),
	)
}
