package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/elliot14A/meterus-go/meters/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func authInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Add the API key to the outgoing context as a Bearer token
		ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", fmt.Sprintf("Bearer %s", apiKey))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// MeterusClient represents a client for the Meterus service.
type MeterusClient struct {
	conn         *grpc.ClientConn
	client       pb.MeteringServiceClient
	addr         string
	unauthConn   *grpc.ClientConn
	unauthClient pb.MeteringServiceClient
}

// NewMeterusClient creates a new MeterusClient with the given address and API key.
func NewMeterusClient(addr, apiKey string, opts ...grpc.DialOption) (*MeterusClient, error) {
	opts = append(opts,
		grpc.WithUnaryInterceptor(authInterceptor(apiKey)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticated gRPC connection: %w", err)
	}

	unauthConn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		conn.Close() // Close the first connection if the second fails
		return nil, fmt.Errorf("failed to create unauthenticated gRPC connection: %w", err)
	}

	return &MeterusClient{
		conn:         conn,
		client:       pb.NewMeteringServiceClient(conn),
		addr:         addr,
		unauthConn:   unauthConn,
		unauthClient: pb.NewMeteringServiceClient(unauthConn),
	}, nil
}

// Close closes all client connections.
func (c *MeterusClient) Close() error {
	err1 := c.conn.Close()
	err2 := c.unauthConn.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

// Ingest sends a cloud event to the Meterus service for ingestion.
func (c *MeterusClient) Ingest(ctx context.Context, event *pb.CloudEvent) error {
	_, err := c.client.Ingest(ctx, event)
	return err
}

// ListMeters retrieves a list of meters from the Meterus service.
func (c *MeterusClient) ListMeters(ctx context.Context, limit, page int32) (*pb.ListMetersResponse, error) {
	return c.client.ListMeters(ctx, &pb.ListMetersRequest{
		Limit: limit,
		Page:  page,
	})
}

// GetMeter retrieves a specific meter from the Meterus service.
func (c *MeterusClient) GetMeter(ctx context.Context, meterIDOrSlug string) (*pb.Meter, error) {
	return c.client.GetMeter(ctx, &pb.MeterId{MeterIdOrSlug: meterIDOrSlug})
}

// CreateMeter creates a new meter in the Meterus service.
func (c *MeterusClient) CreateMeter(ctx context.Context, req *pb.CreateMeterRequest) (*pb.Meter, error) {
	return c.client.CreateMeter(ctx, req)
}

// DeleteMeter deletes a specific meter from the Meterus service.
func (c *MeterusClient) DeleteMeter(ctx context.Context, meterIDOrSlug string) error {
	_, err := c.client.DeleteMeter(ctx, &pb.MeterId{MeterIdOrSlug: meterIDOrSlug})
	return err
}

// QueryMeter queries a specific meter in the MeteringService.
func (c *MeterusClient) QueryMeter(ctx context.Context, req *pb.QueryMeterRequest) (*pb.QueryMeterResponse, error) {
	return c.client.QueryMeter(ctx, req)
}

// ListMeterSubjects retrieves a list of subjects for a specific meter from the MeteringService.
func (c *MeterusClient) ListMeterSubjects(ctx context.Context, meterIDOrSlug string) (*pb.ListMeterSubjectsResponse, error) {
	return c.client.ListMeterSubjects(ctx, &pb.ListMeterSubjectsRequest{MeterIdOrSlug: meterIDOrSlug})
}

func (c *MeterusClient) ValidateApiKey(ctx context.Context, apiKey string, scopes []string) (*pb.ValidateApiKeyResponse, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", fmt.Sprintf("Bearer %s", apiKey))
	return c.unauthClient.ValidateApiKey(ctx, &pb.ValidateApiKeyRequest{RequiredScopes: scopes})
}

// NewCloudEvent creates a new CloudEvent with the given parameters.
func NewCloudEvent(id, source, specVersion, eventType string, time time.Time, subject string, data map[string]any) (*pb.CloudEvent, error) {
	// Convert the data map to a protobuf Struct
	dataStruct, err := structpb.NewStruct(data)
	if err != nil {
		return nil, err
	}

	// Create and return the CloudEvent
	return &pb.CloudEvent{
		Id:          id,
		Source:      source,
		SpecVersion: specVersion,
		Type:        eventType,
		Time:        timestamppb.New(time),
		Subject:     subject,
		Data:        dataStruct,
	}, nil
}
