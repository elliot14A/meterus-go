package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/elliot14A/meterus-go/meters/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// MeterusClient represents a client for the Meterus service.
type MeterusClient struct {
	conn   *grpc.ClientConn         // The gRPC client connection
	client pb.MeteringServiceClient // The Metering service client
	apiKey string                   // Meterus API KEY
}

// authInterceptor creates a gRPC interceptor that adds the API key to the outgoing context.
func authInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Add the API key to the outgoing context as a Bearer token
		ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", fmt.Sprintf("Bearer %s", apiKey))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// NewMeterusClient creates a new MeterusClient with the given address and API key.
func NewMeterusClient(addr, apiKey string, opts ...grpc.DialOption) (*MeterusClient, error) {
	// Add the auth interceptor to the dial options
	opts = append(opts, grpc.WithUnaryInterceptor(authInterceptor(apiKey)))

	// Establish a new gRPC client connection
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, err
	}

	// Create and return the MeterusClient
	return &MeterusClient{
		conn:   conn,
		client: pb.NewMeteringServiceClient(conn),
	}, nil
}

// Close closes the client connection.
func (c *MeterusClient) Close() error {
	return c.conn.Close()
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
