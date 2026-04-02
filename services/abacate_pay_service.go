package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"inovare-backend/config"
	"inovare-backend/requests"
	"inovare-backend/utils"
)

type PaymentProviderService interface {
	CreatePIXBilling(ctx context.Context, req CreatePIXBillingRequest) (string, error)
}

type CreatePIXBillingRequest struct {
	ExternalID    string
	Name          string
	Description   string
	PriceInCents  int
	ReturnURL     string
	CompletionURL string
	Customer      requests.PublicCatalogPurchaseCustomer
}

type abacatePayService struct {
	client *http.Client
}

func NewAbacatePayService() PaymentProviderService {
	return &abacatePayService{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

type abacateCreateBillingRequest struct {
	Frequency     string                        `json:"frequency"`
	Methods       []string                      `json:"methods"`
	Products      []abacateCreateBillingProduct `json:"products"`
	ReturnURL     string                        `json:"returnUrl"`
	CompletionURL string                        `json:"completionUrl"`
	Customer      abacateCustomer               `json:"customer"`
}

type abacateCreateBillingProduct struct {
	ExternalID  string `json:"externalId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
}

type abacateCustomer struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Cellphone string `json:"cellphone"`
	TaxID     string `json:"taxId"`
}

type abacateCreateBillingResponse struct {
	Success bool `json:"success"`
	Data    struct {
		URL string `json:"url"`
	} `json:"data"`
}

func (s *abacatePayService) CreatePIXBilling(ctx context.Context, req CreatePIXBillingRequest) (string, error) {
	cfg := config.GetConfig()
	if strings.TrimSpace(cfg.AbacatePayToken) == "" {
		return "", utils.ErrPaymentConfigurationMissing
	}

	requestBody := abacateCreateBillingRequest{
		Frequency: "ONE_TIME",
		Methods:   []string{"PIX", "CARD"},
		Products: []abacateCreateBillingProduct{
			{
				ExternalID:  req.ExternalID,
				Name:        req.Name,
				Description: req.Description,
				Quantity:    1,
				Price:       req.PriceInCents,
			},
		},
		ReturnURL:     req.ReturnURL,
		CompletionURL: req.CompletionURL,
		Customer: abacateCustomer{
			Name:      req.Customer.Name,
			Email:     req.Customer.Email,
			Cellphone: req.Customer.Cellphone,
			TaxID:     req.Customer.TaxID,
		},
	}

	var payload bytes.Buffer
	if err := json.NewEncoder(&payload).Encode(requestBody); err != nil {
		return "", err
	}

	endpoint := strings.TrimRight(cfg.AbacatePayAPIBase, "/") + "/v1/billing/create"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &payload)
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Authorization", "Bearer "+cfg.AbacatePayToken)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("%w: %v", utils.ErrPaymentProviderUnavailable, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("%w: failed to read response", utils.ErrPaymentProviderUnavailable)
	}

	var response abacateCreateBillingResponse
	if len(responseBody) > 0 {
		_ = json.Unmarshal(responseBody, &response)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices || !response.Success || response.Data.URL == "" {
		message := extractAbacateErrorMessage(responseBody)
		if isInvalidTaxIDError(message) {
			return "", fmt.Errorf("%w: %s", utils.ErrInvalidTaxID, message)
		}
		return "", fmt.Errorf("%w: %s", utils.ErrPaymentProviderUnavailable, message)
	}

	return response.Data.URL, nil
}

func extractAbacateErrorMessage(responseBody []byte) string {
	if len(responseBody) == 0 {
		return "Abacate Pay request failed"
	}

	var payload map[string]any
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return strings.TrimSpace(string(responseBody))
	}

	if message := findStringValue(payload["error"]); message != "" {
		return message
	}
	if message := findStringValue(payload["message"]); message != "" {
		return message
	}

	return strings.TrimSpace(string(responseBody))
}

func findStringValue(value any) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case map[string]any:
		for _, nested := range v {
			if message := findStringValue(nested); message != "" {
				return message
			}
		}
	case []any:
		for _, nested := range v {
			if message := findStringValue(nested); message != "" {
				return message
			}
		}
	}

	return ""
}

func isInvalidTaxIDError(message string) bool {
	normalized := strings.ToLower(strings.TrimSpace(message))
	return strings.Contains(normalized, "invalid taxid") || strings.Contains(normalized, "invalid tax id")
}
