package dto

type Quality string

const (
	LOW    Quality = "low"
	MEDIUM Quality = "medium"
	HIGH   Quality = "high"
)

type GenerateRequest struct {
	Quality Quality `json:"quality"`
	Prompt  string  `json:"prompt"`
}
