package services

import (
	"context"

	"github.com/pchchv/hcms/models"
)

// BitrixClient defines the interface for sending leads to Bitrix24.
type BitrixClient interface {
	SendLead(ctx context.Context, lead models.Lead, webhookURL string) error
}

// HTTPBitrixClient implements BitrixClient using net/http.
type HTTPBitrixClient struct{}

type bitrixEmail struct {
	Value     string `json:"VALUE"`
	ValueType string `json:"VALUE_TYPE"`
}

type bitrixPhone struct {
	Value     string `json:"VALUE"`
	ValueType string `json:"VALUE_TYPE"`
}
