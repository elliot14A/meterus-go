package client

import (
	"context"

	meter "github.com/elliot14A/meterus-go/meters/v1"
)

type MeteringService struct {
	client meter.MeteringServiceClient
	apiKey string
}

func (c *Client) NewMeteringService() *MeteringService {
	return &MeteringService{
		client: meter.NewMeteringServiceClient(c.conn),
		apiKey: c.apiKey,
	}
}

// Ingest sends a cloud event to the Meterus service for ingestion.
func (m *MeteringService) Ingest(ctx context.Context, event *meter.CloudEvent) error {
	ctx = AddApiKeyAuthorizationHeader(ctx, m.apiKey)
	_, err := m.client.Ingest(ctx, event)
	return err
}

// ListMeters retrieves a list of meters from the Meterus service.
func (m *MeteringService) ListMeters(ctx context.Context, limit, page int32) (*meter.ListMetersResponse, error) {
	ctx = AddApiKeyAuthorizationHeader(ctx, m.apiKey)
	return m.client.ListMeters(ctx, &meter.ListMetersRequest{
		Limit: limit,
		Page:  page,
	})
}

// GetMeter retrieves a specific meter from the Meterus service.
func (m *MeteringService) GetMeter(ctx context.Context, meterIDOrSlug string) (*meter.Meter, error) {
	ctx = AddApiKeyAuthorizationHeader(ctx, m.apiKey)
	return m.client.GetMeter(ctx, &meter.MeterId{MeterIdOrSlug: meterIDOrSlug})
}

// CreateMeter creates a new meter in the Meterus service.
func (m *MeteringService) CreateMeter(ctx context.Context, req *meter.CreateMeterRequest) (*meter.Meter, error) {
	ctx = AddApiKeyAuthorizationHeader(ctx, m.apiKey)
	return m.client.CreateMeter(ctx, req)
}

// DeleteMeter deletes a specific meter from the Meterus service.
func (m *MeteringService) DeleteMeter(ctx context.Context, meterIDOrSlug string) error {
	ctx = AddApiKeyAuthorizationHeader(ctx, m.apiKey)
	_, err := m.client.DeleteMeter(ctx, &meter.MeterId{MeterIdOrSlug: meterIDOrSlug})
	return err
}

// QueryMeter queries a specific meter in the MeteringService.
func (m *MeteringService) QueryMeter(ctx context.Context, req *meter.QueryMeterRequest) (*meter.QueryMeterResponse, error) {
	ctx = AddApiKeyAuthorizationHeader(ctx, m.apiKey)
	return m.client.QueryMeter(ctx, req)
}

// ListMeterSubjects retrieves a list of subjects for a specific meter from the MeteringService.
func (m *MeteringService) ListMeterSubjects(ctx context.Context, meterIDOrSlug string) (*meter.ListMeterSubjectsResponse, error) {
	ctx = AddApiKeyAuthorizationHeader(ctx, m.apiKey)
	return m.client.ListMeterSubjects(ctx, &meter.ListMeterSubjectsRequest{MeterIdOrSlug: meterIDOrSlug})
}
