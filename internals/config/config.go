package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port         string
	OpenAIKey    string
	ReplicateKey string
}

func MustLoad() Config {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	if openaiApiKey == "" {
		panic("OPENAI_API_KEY environment variable not set")
	}

	replicateToken := os.Getenv("REPLICATE_TOKEN")
	if replicateToken == "" {
		panic("REPLICATE_TOKEN environment variable not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Port:         port,
		OpenAIKey:    openaiApiKey,
		ReplicateKey: replicateToken,
	}
}
