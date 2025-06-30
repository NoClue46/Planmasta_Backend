package service

import (
	"io"
	"log/slog"
	"net/http"
)

type OpenAIService struct {
	apiKey string
	logger *slog.Logger
	client *http.Client
}

func NewOpenAIService(apiKey string, logger *slog.Logger) *OpenAIService {
	return &OpenAIService{
		apiKey: apiKey,
		logger: logger,
		client: &http.Client{},
	}
}

func (s *OpenAIService) SendRequest(body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", body)
	if err != nil {
		s.logger.Error("failed to create request", slog.String("error", err.Error()))
		return nil, http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("failed to send request", slog.String("error", err.Error()))
		return nil, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		s.logger.Info("Received response from OpenAI", slog.String("status", resp.Status))
	} else {
		s.logger.Warn("Received response from OpenAI", slog.String("status", resp.Status))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("failed to read response body", slog.String("error", err.Error()))
		return nil, http.StatusInternalServerError, err
	}

	return responseBody, resp.StatusCode, nil
}
