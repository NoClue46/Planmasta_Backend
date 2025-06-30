package service

import (
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"planmasta.com/internals/dto"
)

type ReplicateService struct {
	logger *slog.Logger
	apiKey string
	client *http.Client
}

func NewReplicateService(apiKey string, logger *slog.Logger) *ReplicateService {
	return &ReplicateService{
		logger: logger,
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (s *ReplicateService) SendRequest(body dto.GenerateRequest) (*ReplicateResponse, error) {
	start := time.Now()

	s.logger.Info("Sending request to replicate")

	input := ReplicateRequest{
		Input: Input{
			Prompt: body.Prompt,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(input)
	if err != nil {
		return nil, err
	}

	url := "https://api.replicate.com/v1/models/black-forest-labs/flux-kontext-max/predictions"
	if body.Quality == dto.LOW {
		url = "https://api.replicate.com/v1/models/ideogram-ai/ideogram-v3-turbo/predictions"
	} else if body.Quality == dto.MEDIUM {
		url = "https://api.replicate.com/v1/models/ideogram-ai/ideogram-v3-balanced/predictions"
	}

	slog.Info(url)

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "wait")

	resp, err := s.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var response ReplicateResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	s.logger.Info("Replicate request finished", slog.Int64("duration", time.Since(start).Milliseconds()))

	return &response, nil
}

type ReplicateResponse struct {
	Output string `json:"output"`
}

type Input struct {
	Prompt string `json:"prompt"`
}

type ReplicateRequest struct {
	Input Input `json:"input"`
}
