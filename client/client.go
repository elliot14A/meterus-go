package client

import (
	"context"
	"fmt"
	"time"

	meter "github.com/elliot14A/meterus-go/meters/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// Client represents a client for the Meterus service.
type Client struct {
	conn   *grpc.ClientConn
	apiKey string
}

// NewMeterusClient creates a new MeterusClient with the given address and API key.
func NewMeterusClient(addr, apiKey string, opts ...grpc.DialOption) (*Client, error) {
	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticated gRPC connection: %w", err)
	}

	return &Client{
		conn:   conn,
		apiKey: apiKey,
	}, nil
}

// Close closes all client connections.
func (c *Client) Close() error {
	err := c.conn.Close()
	return err
}

// NewCloudEvent creates a new CloudEvent with the given parameters.
func NewCloudEvent(id, source, specVersion, eventType string, time time.Time, subject string, data map[string]any) (*meter.CloudEvent, error) {
	// Convert the data map to a protobuf Struct
	dataStruct, err := structpb.NewStruct(data)
	if err != nil {
		return nil, err
	}

	// Create and return the CloudEvent
	return &meter.CloudEvent{
		Id:          id,
		Source:      source,
		SpecVersion: specVersion,
		Type:        eventType,
		Time:        timestamppb.New(time),
		Subject:     subject,
		Data:        dataStruct,
	}, nil
}

func AddApiKeyAuthorizationHeader(ctx context.Context, apiKey string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "Authorization", fmt.Sprintf("Bearer %s", apiKey))
}
