package client

import (
	"context"

	validation "github.com/elliot14A/meterus-go/validation/v1"
)

type ValidationService struct {
	client validation.ValidationServiceClient
}

func (c *Client) NewValidationService() *ValidationService {
	return &ValidationService{
		client: validation.NewValidationServiceClient(c.conn),
	}
}

func (v *ValidationService) ValidateApiKey(ctx context.Context, apiKey string, scopes []string) (bool, string, any, error) {
	ctx = AddApiKeyAuthorizationHeader(ctx, apiKey)
	res, err := v.client.ValidateApiKey(ctx, &validation.ValidateApiKeyRequest{
		RequiredScopes: scopes,
	})
	if err != nil {
		return false, "", nil, err
	}
	return true, res.Metadata.Subject, res.Metadata.AdditionalAttributes, nil
}
