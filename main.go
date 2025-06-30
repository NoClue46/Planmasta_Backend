package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"planmasta.com/internals/config"
	"planmasta.com/internals/handlers"
	"planmasta.com/internals/service"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	openaiService := service.NewOpenAIService(cfg.OpenAIKey, log)
	replicateService := service.NewReplicateService(cfg.ReplicateKey, log)

	openaiHandler := handlers.NewOpenAIHandler(openaiService, log)
	replicateHandler := handlers.NewReplicateHandler(replicateService, log)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/chat", openaiHandler.Chat)
	r.Post("/replicate", replicateHandler.Generate)

	log.Info("starting server", slog.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		panic(err)
	}
}
